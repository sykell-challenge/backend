package crawl_manager

import (
	"fmt"
	"sykell-challenge/backend/db"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/repositories"
	"sykell-challenge/backend/utils"
	crawlUtils "sykell-challenge/backend/utils/crawl"

	"github.com/gocolly/colly"
)

type CrawlManager struct {
	data        *models.URL
	currentHost string
	collector   *colly.Collector
	linksFound  []string
	urlRepo     *repositories.URLRepository
	jobRepo     *repositories.CrawlJobRepository
}

func InitializeCrawlManager(url string) *CrawlManager {
	var data models.URL
	data.URL = url
	data.Links = models.Links{}
	data.Tags = models.Tags{}

	db := db.GetDB()
	urlRepo := repositories.NewURLRepository(db)
	jobRepo := repositories.NewCrawlJobRepository(db)

	cm := &CrawlManager{
		data:        &data,
		currentHost: utils.GetHostFromURL(url),
		linksFound:  []string{},
		urlRepo:     urlRepo,
		jobRepo:     jobRepo,
	}

	cm.initCrawler()

	return cm
}

func (cm *CrawlManager) Crawl() crawlUtils.CrawlData {
	cm.collector.Visit(cm.data.URL)

	cm.collector.Wait()

	if err := cm.urlRepo.Update(cm.data); err != nil {
		fmt.Println("crawl_manager.go:52 Tried to update URL record:", cm.data)
		fmt.Println("crawl_manager.go:53 Failed to update URL record: ", err)
		BroadcastError(models.CrawlJob{
			URL:   cm.data.URL,
			URLID: cm.data.ID,
		}, "Failed to update URL record: "+err.Error())
		return crawlUtils.CrawlData{}
	}

	fmt.Printf("✅ Successfully updated URL record:\n%+v\n", cm.data)

	cm.jobRepo.UpdateStatus(cm.data.JobId, "running")

	fmt.Printf("✅ Successfully updated job status to 'running' for Job ID: %s\n", cm.data.JobId)
	cm.jobRepo.UpdateProgress(cm.data.JobId, 75)
	fmt.Printf("✅ Successfully updated job progress to 75%% for Job ID: %s\n", cm.data.JobId)

	currentJob, err := cm.jobRepo.GetByID(cm.data.JobId)
	if err != nil {
		BroadcastError(models.CrawlJob{
			URL:   cm.data.URL,
			URLID: cm.data.ID,
		}, "Failed to retrieve current job: "+err.Error())
		return crawlUtils.CrawlData{}
	}

	fmt.Printf("Successfully retrieved current job: %+v\n", currentJob)

	BroadcastHalfCompleted(*currentJob, crawlUtils.CrawlData{
		MainData:  *cm.data,
		LinkCount: len(cm.data.Links),
	})

	cm.processLinks()

	return crawlUtils.CrawlData{
		MainData:  *cm.data,
		LinkCount: len(cm.data.Links),
	}
}
