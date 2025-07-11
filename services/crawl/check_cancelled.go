package crawl

import (
	"context"
	"log"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
)

func (ct *CrawlTask) CheckCancelled(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Printf("Crawl task cancelled before starting: %s", ct.crawlJob.URL)
		crawl_manager.BroadcastCancelled(ct.crawlJob)
		ct.urlRepo.UpdateStatus(ct.crawlJob.URLID, "cancelled")
		return ctx.Err()
	default:
		return nil
	}
}
