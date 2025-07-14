package crawl

import (
	"fmt"
	"log"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/utils/crawl/crawl_manager"
)

func (ct *CrawlTask) UpdateURLRecord(url models.URL) error {
	urlRecord, err := ct.urlRepo.GetByID(ct.CrawlJob.URLID)
	if err != nil {
		log.Printf("Failed to get URL record: %v", err)
		crawl_manager.BroadcastError(ct.CrawlJob, fmt.Sprintf("Failed to get URL record: %v", err))
		return err
	}

	urlRecord.Title = url.Title
	urlRecord.StatusCode = url.StatusCode
	urlRecord.HTMLVersion = url.HTMLVersion
	urlRecord.LoginForm = url.LoginForm
	urlRecord.Tags = url.Tags
	urlRecord.Links = url.Links
	urlRecord.JobId = fmt.Sprint(ct.CrawlJob.ID)
	urlRecord.Status = "done"

	if err := ct.urlRepo.Update(urlRecord); err != nil {
		log.Printf("Failed to update URL record: %v", err)
		crawl_manager.BroadcastError(ct.CrawlJob, fmt.Sprintf("Failed to save crawl results: %v", err))

		ct.urlRepo.UpdateStatus(ct.CrawlJob.URLID, "error")
		return err
	}

	return nil
}
