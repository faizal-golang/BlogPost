package handlers

import (
	"blog-post/config"
	"blog-post/db"
	models "blog-post/models"
	"blog-post/repositories"
	"blog-post/services"
	"blog-post/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetAllArticles(t *testing.T) {
	config.LoadConfigGForMockDB()

	// Setup test DB connection
	DB, err = setupTestDB()
	assert.NoError(t, err)

	// // Ensure that DB is closed after the test is complete
	// d, _ := db.DB()
	// defer d.Close()

	// Initialize router and repositories
	r := mux.NewRouter()
	db.InitDB()
	articleRepo := repositories.NewArticleRepository()
	assert.NotNil(t, articleRepo, "Article repository should not be nil")

	// Initialize the service with the repository
	articleService := services.NewArticleService(articleRepo)
	assert.NotNil(t, articleService, "Article service should not be nil")

	// Initialize handlers with the service
	articleHandler := NewArticleHandler(articleService)
	assert.NotNil(t, articleHandler, "Article handler should not be nil")
	// Define the route
	r.HandleFunc("/articles", articleHandler.ListArticles).Methods("GET")

	// Create the HTTP request
	req, err := http.NewRequest("GET", "/articles", nil)
	assert.NoError(t, err)

	// Record the HTTP response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert that the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body into the Article model
	var response []models.Article
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check if the response contains articles or if it's empty (can modify based on expected behavior)
	if len(response) == 0 {
		t.Errorf("Expected non-empty response, got empty")
	}

	// Optionally, print the response to help with debugging
	fmt.Println("Response:", response)
}

func TestPostArticle(t *testing.T) {
	config.LoadConfigGForMockDB()
	utils.ConnectRabbitMQ()

	// Setup test DB connection
	DB, err = setupTestDB()
	assert.NoError(t, err)

	r := mux.NewRouter()
	db.InitDB()
	articleRepo := repositories.NewArticleRepository()
	assert.NotNil(t, articleRepo, "Article repository should not be nil")

	// Initialize the service with the repository
	articleService := services.NewArticleService(articleRepo)
	assert.NotNil(t, articleService, "Article service should not be nil")

	// Initialize handlers with the service
	articleHandler := NewArticleHandler(articleService)

	PostArticle := models.Article{
		Title:     "Sample Article",
		Content:   "This is a sample article content.",
		Nickname:  "JohnDoe",
		CreatedAt: time.Now(),
	}
	jsonData, err := json.Marshal(PostArticle)
	assert.NoError(t, err)
	r.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")
	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Record the HTTP response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Article
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.ID)
	assert.Equal(t, PostArticle.Title, response.Title)
	assert.Equal(t, PostArticle.Content, response.Content)
	assert.Equal(t, PostArticle.Nickname, response.Nickname)
}

func TestGetArticle(t *testing.T) {
	config.LoadConfigGForMockDB()

	// Setup test DB connection
	DB, err = setupTestDB()
	assert.NoError(t, err)

	// // Ensure that DB is closed after the test is complete
	// d, _ := db.DB()
	// defer d.Close()

	// Initialize router and repositories
	r := mux.NewRouter()
	db.InitDB()
	articleRepo := repositories.NewArticleRepository()
	assert.NotNil(t, articleRepo, "Article repository should not be nil")

	// Initialize the service with the repository
	articleService := services.NewArticleService(articleRepo)
	assert.NotNil(t, articleService, "Article service should not be nil")

	// Initialize handlers with the service
	articleHandler := NewArticleHandler(articleService)
	assert.NotNil(t, articleHandler, "Article handler should not be nil")
	// Define the route
	r.HandleFunc("/articles/{id}", articleHandler.ListArticles).Methods("GET")

	// Create the HTTP request
	req, err := http.NewRequest("GET", "/articles/1", nil)
	assert.NoError(t, err)

	// Record the HTTP response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.Article
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	if reflect.DeepEqual(models.Article{}, response) {
		t.Errorf("Error on response")
	}
}
