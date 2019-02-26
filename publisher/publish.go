package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/publishTask", handlePublishTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/\n", port, port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

func createPubSubClient() (context.Context, *pubsub.Client) {
	ctx := context.Background()

	projectID := "cloud-functions-talk"

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return ctx, client
}

func getTopic(ctx context.Context, client *pubsub.Client, topicID string) *pubsub.Topic {
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

func handlePublishTask(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/publishTask" {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("X-Appengine-Cron") == "" {
		log.Println("X-Appengine-Cron header not present")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	message := "Testing123"

	ctx, client := createPubSubClient()

	msg := &pubsub.Message{
		Data: []byte(message),
	}

	topic := getTopic(ctx, client, "pubsub-example")

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		log.Fatalf("Could not publish message: %v", err)
		return
	}

	log.Printf("Message Recieved: %v", message)

	w.WriteHeader(http.StatusOK)
}
