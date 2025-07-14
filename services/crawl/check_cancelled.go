package crawl

import (
	"context"
	"fmt"
	"log"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
)

func (ct *CrawlTask) CheckCancelled(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Printf("Crawl task cancelled before starting: %s", ct.CrawlJob.URL)
		crawl_manager.BroadcastCancelled(ct.CrawlJob)
		ct.urlRepo.UpdateStatus(ct.CrawlJob.URLID, "cancelled")
		ct.jobRepo.UpdateStatus(fmt.Sprint(ct.CrawlJob.ID), "cancelled")
		return ctx.Err()
	default:
		return nil
	}
}
