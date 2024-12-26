package main

import (
	"log"
	"net/http"

	"blog-post/config"
	database "blog-post/db"
	"blog-post/handlers"
	"blog-post/repositories"
	"blog-post/routes"
	"blog-post/services"
	"blog-post/utils"
)

func init() {
	config.LoadConfig()
	// Initialize the database connection
	database.InitDB()
	utils.ConnectRabbitMQ()
	go utils.ConsumeArticleMessages()

}
func main() {

	// Initialize repositories
	articleRepo := repositories.NewArticleRepository()
	commentRepo := repositories.NewCommentRepository()

	// Initialize services
	articleService := services.NewArticleService(articleRepo)
	commentService := services.NewCommentService(commentRepo)

	// Initialize handlers
	articleHandler := handlers.NewArticleHandler(articleService)
	commentHandler := handlers.NewCommentHandler(commentService)

	// Set up the routes
	router := routes.NewRouter(articleHandler, commentHandler)

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
