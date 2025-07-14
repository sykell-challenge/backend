package crawl

import (
	"fmt"
	"log"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
)

func (ct *CrawlTask) UpdateUrlStatus(status string) error {
	if err := ct.urlRepo.UpdateStatus(ct.CrawlJob.URLID, status); err != nil {
		log.Printf("Failed to update URL (%s) status to %s: %v", ct.CrawlJob.URL, status, err)
		crawl_manager.BroadcastError(ct.CrawlJob, fmt.Sprintf("Failed to update url status: %v", err))
		return err
	}
	return nil
}

func (ct *CrawlTask) UpdateJobStatus(status string) error {
	jobId := fmt.Sprint(ct.CrawlJob.ID)
	if err := ct.jobRepo.UpdateStatus(jobId, status); err != nil {
		log.Printf("Failed to update JOB (%s) status to %s: %v", jobId, status, err)
		crawl_manager.BroadcastError(ct.CrawlJob, fmt.Sprintf("Failed to update job status: %v", err))
		return err
	}
	return nil
}
