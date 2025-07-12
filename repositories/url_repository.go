package repositories

import (
	"errors"
	"fmt"
	"strings"
	"sykell-challenge/backend/models"

	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *URLRepository) GetByID(id uint) (*models.URL, error) {
	var url models.URL
	err := r.db.First(&url, id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) GetAll() ([]models.URL, error) {
	var urls []models.URL
	err := r.db.Find(&urls).Error
	return urls, err
}

func (r *URLRepository) GetByStatus(status string) ([]models.URL, error) {
	var urls []models.URL
	err := r.db.Where("status = ?", status).Find(&urls).Error
	return urls, err
}

func (r *URLRepository) Update(url *models.URL) error {
	return r.db.Save(url).Error
}

func (r *URLRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.URL{}).Where("id = ?", id).Update("status", status).Error
}

func (r *URLRepository) Delete(id uint) error {
	return r.db.Delete(&models.URL{}, id).Error
}

func (r *URLRepository) GetByURL(urlString string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("url = ?", urlString).First(&url).Error
	if err != nil {
		return nil, err
	}

	if url.URL != urlString {
		return nil, errors.New("Not found!")
	}
	return &url, nil
}

func (r *URLRepository) GetByJobID(jobID string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("crawl_job_id = ?", jobID).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.URL{}).Count(&count).Error
	return count, err
}

// URLQueryParams for advanced querying
type URLQueryParams struct {
	Page        int    `form:"page,default=1"`
	Limit       int    `form:"limit,default=10"`
	SortBy      string `form:"sort_by,default=id"`
	SortOrder   string `form:"sort_order,default=desc"`
	Status      string `form:"status"`
	HTMLVersion string `form:"html_version"`
	LoginForm   *bool  `form:"login_form"`
	Search      string `form:"search"`
}

// URLListResponse for paginated responses
type URLListResponse struct {
	Data       []models.URL `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}

// GetAllWithParams returns URLs with pagination, sorting, and filtering
func (r *URLRepository) GetAllWithParams(params URLQueryParams) (*URLListResponse, error) {
	var urls []models.URL
	var total int64

	query := r.db.Model(&models.URL{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.HTMLVersion != "" {
		query = query.Where("html_version = ?", params.HTMLVersion)
	}
	if params.LoginForm != nil {
		query = query.Where("login_form = ?", *params.LoginForm)
	}

	// Apply search (fuzzy search on URL field)
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(url) LIKE ?", searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	orderClause := r.buildOrderClause(params.SortBy, params.SortOrder)
	query = query.Order(orderClause)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query
	if err := query.Find(&urls).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &URLListResponse{
		Data:       urls,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// buildOrderClause creates the ORDER BY clause for sorting
func (r *URLRepository) buildOrderClause(sortBy, sortOrder string) string {
	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Map sort fields to actual database columns or computed values
	switch sortBy {
	case "url":
		return fmt.Sprintf("url %s", sortOrder)
	case "status":
		return fmt.Sprintf("status %s", sortOrder)
	case "html_version":
		return fmt.Sprintf("html_version %s", sortOrder)
	case "login_form":
		return fmt.Sprintf("login_form %s", sortOrder)
	case "created_at":
		return fmt.Sprintf("created_at %s", sortOrder)
	case "updated_at":
		return fmt.Sprintf("updated_at %s", sortOrder)
	case "internal_links":
		return fmt.Sprintf("JSON_LENGTH(JSON_EXTRACT(links, '$[*]')) %s", sortOrder)
	case "external_links":
		return fmt.Sprintf("(SELECT COUNT(*) FROM JSON_TABLE(links, '$[*]' COLUMNS (type VARCHAR(20) PATH '$.type')) AS jt WHERE jt.type = 'external') %s", sortOrder)
	case "inaccessible_links":
		return fmt.Sprintf("(SELECT COUNT(*) FROM JSON_TABLE(links, '$[*]' COLUMNS (type VARCHAR(20) PATH '$.type')) AS jt WHERE jt.type = 'inaccessible') %s", sortOrder)
	default:
		return fmt.Sprintf("id %s", sortOrder)
	}
}

// GetURLLinks returns categorized links for a URL
func (r *URLRepository) GetURLLinks(id uint) (map[string][]models.Link, error) {
	var url models.URL
	err := r.db.First(&url, id).Error
	if err != nil {
		return nil, err
	}

	result := map[string][]models.Link{
		"internal":     []models.Link{},
		"external":     []models.Link{},
		"inaccessible": []models.Link{},
	}

	for _, link := range url.Links {
		switch link.Type {
		case "internal":
			result["internal"] = append(result["internal"], link)
		case "external":
			result["external"] = append(result["external"], link)
		case "inaccessible":
			result["inaccessible"] = append(result["inaccessible"], link)
		}
	}

	return result, nil
}

// SearchURLs performs fuzzy search across multiple fields
func (r *URLRepository) SearchURLs(query string, limit int) ([]models.URL, error) {
	var urls []models.URL
	searchTerm := "%" + strings.ToLower(query) + "%"

	err := r.db.Where("LOWER(url) LIKE ? OR LOWER(html_version) LIKE ?", searchTerm, searchTerm).
		Limit(limit).
		Find(&urls).Error

	return urls, err
}

func (r *URLRepository) GetByCrawlJobID(crawlJobID string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("crawl_jobId = ?", crawlJobID).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}
