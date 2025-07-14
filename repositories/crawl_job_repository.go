package repositories

import (
	"sykell-challenge/backend/models"
	"time"

	"gorm.io/gorm"
)

type CrawlJobRepository struct {
	db *gorm.DB
}

func NewCrawlJobRepository(db *gorm.DB) *CrawlJobRepository {
	return &CrawlJobRepository{db: db}
}

func (r *CrawlJobRepository) Create(job *models.CrawlJob) error {
	return r.db.Create(job).Error
}

func (r *CrawlJobRepository) GetByID(id string) (*models.CrawlJob, error) {
	var job models.CrawlJob
	err := r.db.Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *CrawlJobRepository) Update(jobID string, job *models.CrawlJob) error {
	// return r.db.Save(job).Error
	return r.db.Model(&models.CrawlJob{}).Where("id = ?", jobID).Updates(job).Error
}

func (r *CrawlJobRepository) UpdateStatus(jobID string, status string) error {
	return r.db.Model(&models.CrawlJob{}).Where("id = ?", jobID).Update("status", status).Error
}

func (r *CrawlJobRepository) UpdateProgress(jobID string, progress int) error {
	return r.db.Model(&models.CrawlJob{}).Where("id = ?", jobID).Update("progress", progress).Error
}

func (r *CrawlJobRepository) GetActiveJobs() ([]models.CrawlJob, error) {
	var jobs []models.CrawlJob
	err := r.db.Where("status IN ?", []string{"queued", "running"}).Find(&jobs).Error
	return jobs, err
}

func (r *CrawlJobRepository) GetJobsByStatus(status string) ([]models.CrawlJob, error) {
	var jobs []models.CrawlJob
	err := r.db.Where("status = ?", status).Find(&jobs).Error
	return jobs, err
}

func (r *CrawlJobRepository) GetJobsCreatedAfter(after time.Time) ([]models.CrawlJob, error) {
	var jobs []models.CrawlJob
	err := r.db.Where("created_at > ?", after).Find(&jobs).Error
	return jobs, err
}

func (r *CrawlJobRepository) GetJobHistory() ([]models.CrawlJob, error) {
	var jobs []models.CrawlJob
	err := r.db.Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *CrawlJobRepository) DeleteOldJobs(olderThan time.Time) error {
	return r.db.Where("created_at < ? AND status NOT IN ?", olderThan, []string{"running", "queued"}).Delete(&models.CrawlJob{}).Error
}

func (r *CrawlJobRepository) GetJobsByURLID(urlID uint) ([]models.CrawlJob, error) {
	var jobs []models.CrawlJob
	err := r.db.Where("url_id = ?", urlID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}
