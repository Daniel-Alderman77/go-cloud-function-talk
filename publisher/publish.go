package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/Daniel-Alderman77/go-cloud-functions-talk/publisher/queue"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/publishTask", handlePublishTask)
	mux.HandleFunc("/uploadImage", handleUploadImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s/\n", port, port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
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

	ctx, client := queue.CreatePubSubClient()

	msg := &pubsub.Message{
		Data: []byte(message),
	}

	topic := queue.GetTopic(ctx, client, "pubsub-example")

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		log.Fatalf("Could not publish message: %v", err)
		return
	}

	log.Printf("Message Recieved: %v", message)

	w.WriteHeader(http.StatusOK)
}

func handleUploadImage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/uploadImage" {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("X-Appengine-Cron") == "" {
		log.Println("X-Appengine-Cron header not present")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	file, err := ioutil.ReadFile("odie.png")
	if err != nil {
		log.Fatalf("Could not read odie image file: %v", err)
	}

	message := base64.StdEncoding.EncodeToString(file)

	req, err := http.NewRequest("POST", "https://europe-west1-cloud-functions-talk.cloudfunctions.net/image-upload", strings.NewReader(message))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
