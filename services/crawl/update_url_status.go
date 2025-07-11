package crawl

import (
	"fmt"
	"log"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
)

func (ct *CrawlTask) UpdateUrlStatus(status string) error {
	if err := ct.urlRepo.UpdateStatus(ct.crawlJob.URLID, status); err != nil {
		log.Printf("Failed to update URL status to %s: %v", status, err)
		crawl_manager.BroadcastError(ct.crawlJob, fmt.Sprintf("Failed to update status: %v", err))
		return err
	}
	return nil
}
