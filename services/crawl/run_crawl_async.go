package crawl

import (
	"fmt"
	crawlUtils "sykell-challenge/backend/utils/crawl"
)

func (ct *CrawlTask) RunCrawlAsync() (chan crawlUtils.CrawlData, chan error) {
	crawlDone := make(chan crawlUtils.CrawlData, 1)
	crawlErr := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				crawlErr <- fmt.Errorf("crawl panic: %v", r)
			}
		}()
		crawlData := ct.crawlManager.Crawl()
		crawlDone <- crawlData
	}()

	return crawlDone, crawlErr
}
