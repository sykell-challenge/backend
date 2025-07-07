package repositories

import (
	"sykell-challenge/backend/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *TagRepository) GetByURLID(urlID uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("url_id = ?", urlID).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) GetByTagName(tagName string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("tag_name = ?", tagName).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) UpdateCount(id uint, count int) error {
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Update("count", count).Error
}

func (r *TagRepository) IncrementCount(id uint) error {
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Update("count", gorm.Expr("count + ?", 1)).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

func (r *TagRepository) DeleteByURLID(urlID uint) error {
	return r.db.Where("url_id = ?", urlID).Delete(&models.Tag{}).Error
}

func (r *TagRepository) GetTagCountSummary() (map[string]int, error) {
	var results []struct {
		TagName string
		Total   int
	}

	err := r.db.Model(&models.Tag{}).
		Select("tag_name, SUM(count) as total").
		Group("tag_name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]int)
	for _, result := range results {
		summary[result.TagName] = result.Total
	}

	return summary, nil
}

func (r *TagRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Tag{}).Count(&count).Error
	return count, err
}
