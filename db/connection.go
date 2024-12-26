package db

import (
	"fmt"
	"log"
	"os"

	"blog-post/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := GetDBConnectionString() // Fetch the connection string
	fmt.Println("dgsgsg", dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Connected to the database successfully")

	err = DB.AutoMigrate(
		&models.Article{},
		&models.Comment{},
		&models.Reply{},
	)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

}

func GetDBConnectionString() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	return connectionString
}
