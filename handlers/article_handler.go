package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"blog-post/models"
	"blog-post/services"
	"blog-post/utils"
)

type ArticleHandler struct {
	service services.ArticleService
}

func NewArticleHandler(service services.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// ListArticles handles GET /articles?page={page}&limit={limit}
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	articles, err := h.service.GetAllArticles(page, limit)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to fetch articles")
		return
	}

	utils.JSONResponse(w, http.StatusOK, articles)
}

// GetArticle handles GET /articles/{id}
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/articles/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	article, err := h.service.GetArticleByID(uint(id))
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "Article not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, article)
}

// CreateArticle handles POST /articles
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.PostArticle(&article); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create article")
		return
	}
	err = utils.PublishArticle(article.ID, "A new article has been posted.")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to publish to RabbitMQ")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, article)
}
