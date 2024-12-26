package models

import (
	"time"
)

// Article represents the article schema
type Article struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	Nickname  string     `gorm:"type:varchar(100);not null"`
	Title     string     `gorm:"type:varchar(200);not null"`
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	Comments  []Comment  `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"`
}
