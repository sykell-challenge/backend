package crawl

import (
	"context"
	"fmt"
	"log"
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/services/taskq"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
	"time"
)

type CrawlTask struct {
	CrawlJob     models.CrawlJob
	urlRepo      *repositories.URLRepository
	jobRepo      *repositories.CrawlJobRepository
	crawlManager *crawl_manager.CrawlManager
}

func (ct *CrawlTask) Do(ctx context.Context) error {
	log.Printf("Starting crawl task for URL: %s (ID: %d)", ct.CrawlJob.URL, ct.CrawlJob.URLID)

	jobCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobId := fmt.Sprint(ct.CrawlJob.ID)

	taskq.RegisterJob(jobId, cancel)
	defer taskq.UnregisterJob(jobId)

	if err := ct.CheckCancelled(jobCtx); err != nil {
		return err
	}

	if err := ct.UpdateUrlStatus("running"); err != nil {
		return err
	}

	ct.jobRepo.UpdateStatus(jobId, "running")
	ct.CrawlJob.Status = "running"
	now := time.Now()
	ct.CrawlJob.StartedAt = &now
	ct.CrawlJob.Progress = 25

	ct.jobRepo.Update(jobId, &ct.CrawlJob)

	crawl_manager.BroadcastJobStarted(ct.CrawlJob)

	if err := ct.CheckCancelled(jobCtx); err != nil {
		return err
	}

	crawlDone, crawlErr := ct.RunCrawlAsync()

	crawlData, err := ct.WaitForCrawlResult(jobCtx, crawlDone, crawlErr)
	if err != nil {
		return err
	}

	if err := ct.UpdateURLRecord(crawlData.MainData); err != nil {
		log.Printf("Failed to update URL record: %v", err)
		crawl_manager.BroadcastError(ct.CrawlJob, fmt.Sprintf("Failed to save crawl results: %v", err))
		return err
	}

	ct.CrawlJob.Status = "completed"
	now = time.Now()
	ct.CrawlJob.CompletedAt = &now
	ct.CrawlJob.Progress = 100 // Set progress to 100% on completion

	ct.jobRepo.Update(jobId, &ct.CrawlJob)

	crawl_manager.BroadcastCompleted(ct.CrawlJob, crawlData)

	log.Printf("Crawl task completed successfully for URL: %s", ct.CrawlJob.URL)
	return nil
}

func CreateCrawlTask(url string, urlID uint) *CrawlTask {
	db := db.GetDB()
	urlRepo := repositories.NewURLRepository(db)
	jobsRepo := repositories.NewCrawlJobRepository(db)
	crawlManager := crawl_manager.InitializeCrawlManager(url)

	startedAt := time.Now()

	crawlJob := models.CrawlJob{
		URL:       url,
		URLID:     urlID,
		Status:    "queued",
		StartedAt: &startedAt,
	}

	jobsRepo.Create(&crawlJob)

	crawl_manager.BroadcastJobQueued(fmt.Sprintf("%d", crawlJob.ID), url, urlID)

	return &CrawlTask{
		CrawlJob:     crawlJob,
		urlRepo:      urlRepo,
		jobRepo:      jobsRepo,
		crawlManager: crawlManager,
	}
}
