package queue

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

// CreatePubSubClient creates authenticated Pub/Sub client.
func CreatePubSubClient() (context.Context, *pubsub.Client) {
	ctx := context.Background()

	// Creates a client.
	client, err := pubsub.NewClient(ctx, "cloud-functions-talk")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return ctx, client
}

// GetTopic checks if the topic exists before returning a pointer to it.
func GetTopic(ctx context.Context, client *pubsub.Client, topicID string) *pubsub.Topic {
	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Error checking for topic: %v", err)
	}
	if !exists {
		log.Fatalf("Topic does not exist: %v", err)
	}

	return topic
}
