package vision

import (
	"context"
	"fmt"
	"log"

	vision "cloud.google.com/go/vision/apiv1"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// DetectLabels detects the labels of an image uploaded to a GCS bucket.
func DetectLabels(ctx context.Context, e GCSEvent) error {
	log.Printf("Processing file: %s", e.Name)

	// Creates a client.
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create vision client: %v", err)
		return err
	}

	img := vision.NewImageFromURI(fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name))

	labels, err := visionClient.DetectLabels(ctx, img, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
		return err
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Println(label.Description)
	}

	return nil
}
