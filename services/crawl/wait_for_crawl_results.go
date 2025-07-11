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
		log.Printf("Crawl task cancelled during crawling: %s", ct.crawlJob.URL)
		crawl_manager.BroadcastCancelled(ct.crawlJob)
		_ = ct.UpdateUrlStatus("cancelled")
		return crawlUtils.CrawlData{}, ctx.Err()

	case err := <-crawlErr:
		log.Printf("Crawl task failed: %v", err)
		crawl_manager.BroadcastError(ct.crawlJob, fmt.Sprintf("Crawl failed: %v", err))
		_ = ct.UpdateUrlStatus("error")
		return crawlUtils.CrawlData{}, err

	case crawlData := <-crawlDone:
		return crawlData, nil
	}
}
