package repositories

import (
	"blog-post/db"
	"blog-post/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	FindAll(page, limit int) ([]models.Article, error)
	FindByID(id uint) (*models.Article, error)
	Create(article *models.Article) error
}

type ArticleRepositories struct {
	Db *gorm.DB
}

func NewArticleRepository() ArticleRepository {
	return &ArticleRepositories{Db: db.DB}
}

func (r *ArticleRepositories) FindAll(page, limit int) ([]models.Article, error) {
	var articles []models.Article
	offset := (page - 1) * limit

	// Preload the comments to fetch them along with the articles
	if err := r.Db.Preload("Comments").Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepositories) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	if err := r.Db.Preload("Comments").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepositories) Create(article *models.Article) error {
	return r.Db.Create(article).Error
}
