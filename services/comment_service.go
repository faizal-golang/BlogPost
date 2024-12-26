package services

import (
	"blog-post/models"
	"blog-post/repositories"
	"fmt"
)

type CommentService interface {
	GetCommentsByArticleID(articleID uint) ([]models.Comment, error)
	PostComment(comment *models.Comment) (*models.Comment, error)
	CreateCommentOnComments(parentID uint, articleID uint, nickname string, content string) (*models.Comment, error)
}

type commentService struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) GetCommentsByArticleID(articleID uint) ([]models.Comment, error) {
	return s.repo.FindByArticleID(articleID)
}

func (s *commentService) PostComment(comment *models.Comment) (*models.Comment, error) {
	// Check if the article exists
	// var article models.Article
	if _, err := s.repo.FindByArticleID(comment.ArticleID); err != nil {
		return nil, fmt.Errorf("article not found: %v", err)
	}

	// Create a new comment
	comment = &models.Comment{
		ArticleID: comment.ArticleID,
		Nickname:  comment.Nickname,
		Content:   comment.Nickname,
	}
	fmt.Println("enterr commentservice")

	// Create the comment in the database
	err, _ := s.repo.Create(comment)
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %v", err)
	}

	com, err := s.repo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	return com, nil
}

func (s *commentService) CreateCommentOnComments(parentID uint, articleID uint, nickname string, content string) (*models.Comment, error) {
	return s.repo.CreateCommentOnComment(parentID, articleID, nickname, content)
}
