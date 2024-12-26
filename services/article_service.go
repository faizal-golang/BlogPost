package services

import (
	"blog-post/models"
	"blog-post/repositories"
)

type ArticleService interface {
	GetAllArticles(page, limit int) ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	PostArticle(article *models.Article) error
}

type articleService struct {
	repo repositories.ArticleRepository
}

func NewArticleService(repo repositories.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) GetAllArticles(page, limit int) ([]models.Article, error) {
	return s.repo.FindAll(page, limit)
}

func (s *articleService) GetArticleByID(id uint) (*models.Article, error) {
	return s.repo.FindByID(id)
}

func (s *articleService) PostArticle(article *models.Article) error {
	return s.repo.Create(article)
}
