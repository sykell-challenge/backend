package crawl

import (
	"context"
	"fmt"
	"log"
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/services/taskq"
	crawlUtils "sykell-challenge/backend/utils/crawl"
	"sykell-challenge/backend/utils/crawl/crawl_manager"

	"github.com/google/uuid"
)

type CrawlTask struct {
	crawlJob     crawlUtils.CrawlJob
	urlRepo      *repositories.URLRepository
	crawlManager *crawl_manager.CrawlManager
}

func (ct *CrawlTask) Do(ctx context.Context) error {
	log.Printf("Starting crawl task for URL: %s (ID: %d)", ct.crawlJob.URL, ct.crawlJob.URLID)

	jobCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	taskq.RegisterJob(ct.crawlJob.ID, cancel)
	defer taskq.UnregisterJob(ct.crawlJob.ID)

	if err := ct.CheckCancelled(jobCtx); err != nil {
		return err
	}

	crawl_manager.BroadcastJobStarted(ct.crawlJob)

	if err := ct.UpdateUrlStatus("running"); err != nil {
		return err
	}

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
		crawl_manager.BroadcastError(ct.crawlJob, fmt.Sprintf("Failed to save crawl results: %v", err))
		return err
	}

	crawl_manager.BroadcastHalfCompleted(ct.crawlJob, crawlData)

	log.Printf("Crawl task completed successfully for URL: %s", ct.crawlJob.URL)
	return nil
}

func CreateCrawlTask(url string, urlID uint) *CrawlTask {
	db := db.GetDB()
	urlRepo := repositories.NewURLRepository(db)
	crawlManager := crawl_manager.InitializeCrawlManager(url)

	return &CrawlTask{
		crawlJob: crawlUtils.CrawlJob{
			ID:    uuid.New().String(),
			URL:   url,
			URLID: urlID,
		},
		urlRepo:      urlRepo,
		crawlManager: crawlManager,
	}
}

func (ct *CrawlTask) GetJobID() string {
	return ct.crawlJob.ID
}
