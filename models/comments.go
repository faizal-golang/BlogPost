package models

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint      `gorm:"not null" json:"articleId"` // Foreign key for the article
	Content   string    `gorm:"not null" json:"content"`
	Nickname  string    `gorm:"not null" json:"nickname"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	Replies   []Reply   `gorm:"foreignKey:CommentID" json:"replies"` // One-to-many relationship with Reply
}

type Reply struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CommentID uint      `gorm:"not null" json:"commentId"` // Foreign key for the parent comment
	Content   string    `gorm:"not null" json:"content"`
	Nickname  string    `gorm:"not null" json:"nickname"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
