package vision

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/storage"
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

	// Create a vision client
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create vision client: %v", err)
		return err
	}

	// Load image from GCS bucket
	img := vision.NewImageFromURI(fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name))

	labels, err := visionClient.DetectLabels(ctx, img, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
		return err
	}

	props, err := visionClient.DetectSafeSearch(ctx, img, nil)
	if err != nil {
		log.Fatalf("Failed to detect safe search: %v", err)
		return err
	}

	// Split string to generate output filename
	filename := strings.Split(e.Name, ".")[0] + "_labels.txt"

	// Create a GCS client
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}

	// Create object to write output to
	outputBlob := storageClient.Bucket(e.Bucket).Object(filename)
	w := outputBlob.NewWriter(ctx)
	defer w.Close()

	fmt.Println("Labels:")
	fmt.Fprintln(w, "Labels:")
	for _, label := range labels {
		fmt.Printf("%v: %v", label.Description, label.Score)
		fmt.Fprintf(w, "%v: %v\n", label.Description, label.Score)
	}

	fmt.Println("Safe Search:")
	fmt.Fprintln(w, "Safe Search:")
	fmt.Printf("Racy: %v", props.Racy)
	fmt.Fprintf(w, "Racy: %v\n", props.Racy)

	return nil
}
