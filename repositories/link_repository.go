package repositories

import (
	"sykell-challenge/backend/models"

	"gorm.io/gorm"
)

type LinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(link *models.Link) error {
	return r.db.Create(link).Error
}

func (r *LinkRepository) GetByID(id uint) (*models.Link, error) {
	var link models.Link
	err := r.db.First(&link, id).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) GetAll() ([]models.Link, error) {
	var links []models.Link
	err := r.db.Find(&links).Error
	return links, err
}

func (r *LinkRepository) GetByURLID(urlID uint) ([]models.Link, error) {
	var links []models.Link
	err := r.db.Where("url_id = ?", urlID).Find(&links).Error
	return links, err
}

func (r *LinkRepository) GetByType(linkType string) ([]models.Link, error) {
	var links []models.Link
	err := r.db.Where("type = ?", linkType).Find(&links).Error
	return links, err
}

func (r *LinkRepository) GetByStatusCode(statusCode int) ([]models.Link, error) {
	var links []models.Link
	err := r.db.Where("status_code = ?", statusCode).Find(&links).Error
	return links, err
}

func (r *LinkRepository) GetByTypeAndURLID(urlID uint, linkType string) ([]models.Link, error) {
	var links []models.Link
	err := r.db.Where("url_id = ? AND type = ?", urlID, linkType).Find(&links).Error
	return links, err
}

func (r *LinkRepository) Update(link *models.Link) error {
	return r.db.Save(link).Error
}

func (r *LinkRepository) UpdateStatusCode(id uint, statusCode int) error {
	return r.db.Model(&models.Link{}).Where("id = ?", id).Update("status_code", statusCode).Error
}

func (r *LinkRepository) UpdateType(id uint, linkType string) error {
	return r.db.Model(&models.Link{}).Where("id = ?", id).Update("type", linkType).Error
}

func (r *LinkRepository) Delete(id uint) error {
	return r.db.Delete(&models.Link{}, id).Error
}

func (r *LinkRepository) DeleteByURLID(urlID uint) error {
	return r.db.Where("url_id = ?", urlID).Delete(&models.Link{}).Error
}

func (r *LinkRepository) GetLinksByLink(linkString string) ([]models.Link, error) {
	var links []models.Link
	err := r.db.Where("link = ?", linkString).Find(&links).Error
	return links, err
}

func (r *LinkRepository) GetTypeDistribution() (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	err := r.db.Model(&models.Link{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int64)
	for _, result := range results {
		distribution[result.Type] = result.Count
	}

	return distribution, nil
}

func (r *LinkRepository) GetStatusCodeDistribution() (map[int]int64, error) {
	var results []struct {
		StatusCode int
		Count      int64
	}

	err := r.db.Model(&models.Link{}).
		Select("status_code, COUNT(*) as count").
		Group("status_code").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[int]int64)
	for _, result := range results {
		distribution[result.StatusCode] = result.Count
	}

	return distribution, nil
}

func (r *LinkRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Link{}).Count(&count).Error
	return count, err
}
