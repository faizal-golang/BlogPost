package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"blog-post/models"
	"blog-post/services"
	"blog-post/utils"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	service services.CommentService
}

func NewCommentHandler(service services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// ListComments handles GET /articles/{id}/comments
func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/articles/"):]
	idStr = idStr[:len(idStr)-len("/comments")]
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	comments, err := h.service.GetCommentsByArticleID(uint(articleID))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	utils.JSONResponse(w, http.StatusOK, comments)
}

// AddComment handles POST /articles/{id}/comments
func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/articles/"):]
	idStr = idStr[:len(idStr)-len("/comments")]
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	comment.ArticleID = uint(articleID)

	comments, err := h.service.PostComment(&comment)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}
	FinalValues := utils.ConvertToResponse(comments)
	utils.JSONResponse(w, http.StatusCreated, FinalValues)
}

func (h *CommentHandler) CreateCommentOnComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Decode the request body into a new comment
	var request struct {
		Nickname string `json:"nickname"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Invalid request body")
		return
	}
	articleID, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "invalid article ID")
		return
	}
	parentID, err := strconv.Atoi(vars["comment_id"])
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Invalid Parent ID")
		return
	}

	// Call the service layer to create the comment on a comment
	comment, err := h.service.CreateCommentOnComments(uint(parentID), uint(articleID), request.Nickname, request.Content)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating comment on comment: %v", err))
		return
	}
	FinalValues := utils.ConvertToResponse(comment)
	utils.JSONResponse(w, http.StatusCreated, FinalValues)
}
