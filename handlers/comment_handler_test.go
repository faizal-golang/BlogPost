package handlers

import (
	"blog-post/config"
	"blog-post/db"
	models "blog-post/models"
	"blog-post/repositories"
	"blog-post/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = "3306"
	user     = "root"
	password = "root"
	dbName   = "blog_api"
)

func setupTestDB() (*gorm.DB, error) {
	var dsn string
	// Use the password if it is set
	if password != "" {
		dsn = user + ":" + password + "@(" + host + ":" + port + ")/" + dbName + "?parseTime=true"
	} else {
		// If no password, omit it from the DSN
		dsn = user + "@(" + host + ":" + port + ")/" + dbName + "?parseTime=true"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {

	}

	// Auto-migrate the necessary models
	db.AutoMigrate(&models.Article{})
	db.AutoMigrate(&models.Comment{})

	return db, nil
}

var DB *gorm.DB
var err error

func TestAddComment(t *testing.T) {
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
	commentRepo := repositories.NewCommentRepository()
	assert.NotNil(t, commentRepo, "Article repository should not be nil")

	// Initialize the service with the repository
	commentService := services.NewCommentService(commentRepo)
	assert.NotNil(t, commentService, "Article service should not be nil")

	// Initialize handlers with the service
	commentHandler := NewCommentHandler(commentService)
	assert.NotNil(t, commentHandler, "Article handler should not be nil")
	// Define the route
	addComment := models.Comment{
		ArticleID: uint(1),
		Content:   "JohnDoe",
		Nickname:  "JohnDoe",
		CreatedAt: time.Now(),
	}
	jsonData, err := json.Marshal(addComment)
	r.HandleFunc("/articles/{id}/comments", commentHandler.AddComment).Methods("POST")

	// Create the HTTP request
	req, err := http.NewRequest("POST", "/articles/1/comments", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Record the HTTP response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Comment
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.ArticleID)
	assert.Equal(t, addComment.ArticleID, response.ArticleID)
	assert.Equal(t, addComment.Content, response.Content)
	assert.Equal(t, addComment.Nickname, response.Nickname)

	var savedComment models.Comment
	err = db.DB.First(&savedComment, "article_id = ? AND nickname = ?", addComment.ArticleID, addComment.Nickname).Error
	assert.NoError(t, err)
	assert.Equal(t, addComment.Content, savedComment.Content)
}

func TestGetArticleComments(t *testing.T) {
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
	commentRepo := repositories.NewCommentRepository()
	assert.NotNil(t, commentRepo, "Article repository should not be nil")

	// Initialize the service with the repository
	commentService := services.NewCommentService(commentRepo)
	assert.NotNil(t, commentService, "Article service should not be nil")

	// Initialize handlers with the service
	commentHandler := NewCommentHandler(commentService)
	assert.NotNil(t, commentHandler, "Article handler should not be nil")
	// Define the route
	r.HandleFunc("/articles/{id}/comments", commentHandler.ListComments).Methods("GET")

	// Create the HTTP request
	req, err := http.NewRequest("GET", "/articles/1/comments", nil)
	assert.NoError(t, err)

	// Record the HTTP response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Comment
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	if reflect.DeepEqual([]models.Comment{}, response) {
		t.Errorf("Error on response")
	}
}

func TestGetCommentOnComment(t *testing.T) {
	config.LoadConfigGForMockDB()

	// Setup test DB connection
	DB, err = setupTestDB()
	assert.NoError(t, err)
	// Initialize router and repositories
	r := mux.NewRouter()
	db.InitDB()
	commentRepo := repositories.NewCommentRepository()
	assert.NotNil(t, commentRepo, "Article repository should not be nil")

	// Initialize the service with the repository
	commentService := services.NewCommentService(commentRepo)
	assert.NotNil(t, commentService, "Article service should not be nil")

	// Initialize handlers with the service
	commentHandler := NewCommentHandler(commentService)
	assert.NotNil(t, commentHandler, "Article handler should not be nil")
	// Define the route
	CommentReply := models.Reply{
		Nickname: "Faizii A",
		Content:  "Nice",
	}
	val, _ := json.Marshal(CommentReply)
	r.HandleFunc("/articles/{article_id}/comments/{comment_id}", commentHandler.CreateCommentOnComment).Methods("POST")

	// Create the HTTP request
	req, err := http.NewRequest("POST", "/articles/1/comments/1", bytes.NewBuffer(val))
	assert.NoError(t, err)

	// Record the HTTP response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Comment
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

}
