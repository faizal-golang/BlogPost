package utils

import (
	"blog-post/models"
	"encoding/json"
	"net/http"
	"time"
)

type CommentResponse struct {
	CommentID uint           `json:"commentId"`
	ArticleID uint           `json:"articleId"`
	Content   string         `json:"content"`
	Nickname  string         `json:"nickName"`
	CreatedAt time.Time      `json:"createdAt"`
	Replies   []models.Reply `json:"replies"`
}

func ConvertToResponse(comment *models.Comment) *CommentResponse {
	return &CommentResponse{
		CommentID: comment.ID,
		ArticleID: comment.ArticleID,
		Content:   comment.Content,
		Nickname:  comment.Nickname,
		CreatedAt: comment.CreatedAt,
		Replies:   comment.Replies,
	}
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}
