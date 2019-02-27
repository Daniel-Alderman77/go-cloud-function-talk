package pubsub

import (
	"context"
	"log"
)

// Message is the payload of a Pub/Sub event
type Message struct {
	Data []byte `json:"data"`
}

// PrintMessage consumes a Pub/Sub and prints message.
func PrintMessage(ctx context.Context, m Message) error {
	log.Println(string(m.Data))
	return nil
}
