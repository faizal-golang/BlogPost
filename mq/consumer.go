package mq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type ArticleConsumer struct {
	channel *amqp.Channel
}

func NewArticleConsumer(channel *amqp.Channel) *ArticleConsumer {
	return &ArticleConsumer{channel: channel}
}

func (c *ArticleConsumer) StartConsuming(queueName string) {
	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var article map[string]interface{}
		if err := json.Unmarshal(msg.Body, &article); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		log.Printf("New Article Received: %+v\n", article)
		// Process the article (e.g., send to search service, cache, etc.)
	}
}
