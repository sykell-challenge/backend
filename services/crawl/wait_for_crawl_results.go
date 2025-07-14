package crawl

import (
	"context"
	"fmt"
	"log"
	crawlUtils "sykell-challenge/backend/utils/crawl"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
)

func (ct *CrawlTask) WaitForCrawlResult(ctx context.Context, crawlDone <-chan crawlUtils.CrawlData, crawlErr <-chan error) (crawlUtils.CrawlData, error) {
	select {
	case <-ctx.Done():
		log.Printf("Crawl task cancelled during crawling: %s", ct.CrawlJob.URL)
		crawl_manager.BroadcastCancelled(ct.CrawlJob)
		ct.UpdateUrlStatus("cancelled")
		ct.UpdateJobStatus("cancelled")
		return crawlUtils.CrawlData{}, ctx.Err()

	case err := <-crawlErr:
		log.Printf("Crawl task failed: %v", err)
		crawl_manager.BroadcastError(ct.CrawlJob, fmt.Sprintf("Crawl failed: %v", err))
		ct.UpdateUrlStatus("error")
		ct.UpdateJobStatus("error")
		return crawlUtils.CrawlData{}, err

	case crawlData := <-crawlDone:
		// crawl_manager.BroadcastHalfCompleted(ct.crawlJob, crawlData)
		return crawlData, nil
	}
}
