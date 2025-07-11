package crawl_manager

import (
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/utils"
	crawlUtils "sykell-challenge/backend/utils/crawl"

	"github.com/gocolly/colly"
)

type CrawlManager struct {
	data        *models.URL
	currentHost string
	collector   *colly.Collector
	linksFound  []string
}

func InitializeCrawlManager(url string) *CrawlManager {
	var data models.URL
	data.URL = url
	data.Links = models.Links{}
	data.Tags = models.Tags{}

	pm := &CrawlManager{
		data:        &data,
		currentHost: utils.GetHostFromURL(url),
		linksFound:  []string{},
	}

	pm.initCrawler()

	return pm
}

func (pm *CrawlManager) Crawl() crawlUtils.CrawlData {
	pm.collector.Visit(pm.data.URL)

	pm.collector.Wait()

	BroadcastHalfCompleted(crawlUtils.CrawlJob{
		ID:    pm.data.CrawlJobID,
		URL:   pm.data.URL,
		URLID: pm.data.ID,
	}, crawlUtils.CrawlData{
		MainData:  *pm.data,
		LinkCount: len(pm.data.Links),
	})

	pm.processLinks()

	return crawlUtils.CrawlData{
		MainData:  *pm.data,
		LinkCount: len(pm.data.Links),
	}
}
