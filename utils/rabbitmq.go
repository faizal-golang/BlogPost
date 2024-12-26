package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Declare a global variable for the RabbitMQ channel
var Channel *amqp.Channel
var Connection *amqp.Connection

// ConnectRabbitMQ establishes a connection with RabbitMQ server
func ConnectRabbitMQ() {
	var err error
	Connection, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	Channel, err = Connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare a queue for articles and comments
	_, err = Channel.QueueDeclare(
		"article_queue", // Queue name
		true,            // Durable
		false,           // Auto-Delete
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ and declared queues.")
}

// PublishArticle publishes a message to the RabbitMQ queue
func PublishArticle(articleID uint, message string) error {
	// Prepare the message
	msg := fmt.Sprintf("Article ID: %d, Message: %s", articleID, message)
	err := Channel.Publish(
		"",              // Default exchange
		"article_queue", // Routing key (queue name)
		false,           // Mandatory
		false,           // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to publish a message: %v", err)
	}

	log.Printf("Message published to RabbitMQ: %s", msg)
	return nil
}

// ConsumeArticleMessages listens for messages from the RabbitMQ queue
func ConsumeArticleMessages() {
	messages, err := Channel.Consume(
		"article_queue", // Queue name
		"",              // Consumer name
		true,            // Auto-acknowledge
		false,           // Exclusive
		false,           // No local
		false,           // No wait
		nil,             // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// Process incoming messages
	for message := range messages {
		log.Printf("Received a message: %s", message.Body)
		// Here you can process the message, for example:
		// - Update the database
		// - Send an email
		// - Notify other services
	}
}
