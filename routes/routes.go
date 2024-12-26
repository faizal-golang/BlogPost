package routes

import (
	"net/http"

	"blog-post/handlers"

	"github.com/gorilla/mux"
)

type Routes struct {
	ArticleHandler *handlers.ArticleHandler
	CommentHandler *handlers.CommentHandler
}

func NewRouter(articleHandler *handlers.ArticleHandler, commentHandler *handlers.CommentHandler) *mux.Router {
	r := mux.NewRouter()
	routes := Routes{
		ArticleHandler: articleHandler,
		CommentHandler: commentHandler,
	}

	// Article Routes
	r.HandleFunc("/articles", routes.ArticleHandler.ListArticles).Methods(http.MethodGet)
	r.HandleFunc("/articles/{id}", routes.ArticleHandler.GetArticle).Methods(http.MethodGet)
	r.HandleFunc("/articles", routes.ArticleHandler.CreateArticle).Methods(http.MethodPost)

	// Comment Routes
	r.HandleFunc("/articles/{id}/comments", routes.CommentHandler.ListComments).Methods(http.MethodGet)
	r.HandleFunc("/articles/{id}/comments", routes.CommentHandler.AddComment).Methods(http.MethodPost)
	r.HandleFunc("/articles/{article_id}/comments/{comment_id}", commentHandler.CreateCommentOnComment).Methods(http.MethodPost)

	return r
}
