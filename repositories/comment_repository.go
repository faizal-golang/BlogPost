package repositories

import (
	"blog-post/db"
	"blog-post/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CommentRepository interface {
	FindByArticleID(articleID uint) ([]models.Comment, error)
	Create(comment *models.Comment) (error, *models.Comment)
	FindByID(commentID uint) (*models.Comment, error)
	CreateCommentOnComment(parentID uint, articleID uint, nickname string, content string) (*models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{db: db.DB}
}

func (r *commentRepository) FindByArticleID(articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Preload("Replies").Where("article_id = ?", articleID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Create(comment *models.Comment) (error, *models.Comment) {
	// Attempt to create the comment in the database
	err := r.db.Create(comment).Error
	if err != nil {
		return err, nil
	}

	// After successful creation, return the created comment object
	return nil, comment
}

func (r *commentRepository) FindByID(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.First(&comment, commentID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) CreateCommentOnComment(parentID uint, articleID uint, nickname string, content string) (*models.Comment, error) {
	if parentID == 0 {
		return nil, fmt.Errorf("invalid parent comment ID")
	}

	var parentComment models.Comment
	if err := r.db.Where("id = ?", parentID).First(&parentComment).Error; err != nil {
		return nil, fmt.Errorf("parent comment does not exist: %v", err)
	}
	comment := &models.Reply{
		CommentID: parentComment.ID,
		Content:   content,
		Nickname:  nickname,
		CreatedAt: time.Now(),
	}
	parentComment.Replies = append(parentComment.Replies, *comment)

	if err := r.db.Save(&parentComment).Error; err != nil {
		return nil, fmt.Errorf("failed to update parent comment: %v", err)
	}
	return &parentComment, nil
}
