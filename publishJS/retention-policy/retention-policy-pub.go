package retentionpolicy

import (
	"context"
	"go-streaming/util"
	"log"

	"github.com/nats-io/nats.go"
)

func RetentionPolicy() {
	// connect nats
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	streamName := "EVENTS"

	// nats.InterestPolicy if no consumer active, and stream not save message when publish 
	cfg := nats.StreamConfig{
		Name:      "EVENTS",
		Retention: nats.InterestPolicy, // Retains based on defined limits
		Subjects:  []string{"events.>"},
		MaxMsgs:   1000,             // Retain up to 1000 messages
		MaxBytes:  10 * 1024 * 1024, // Retain up to 10MB of messages
	}

	util.CreateStream(js, cfg)

	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	log.Println("Stream created with InterestPolicy")

	go func() {
		for i := 0; i < 100; i++ {
			msg := []byte("Event created with ID: " + string(i)) // ASCII conversion to keep it simple
			_, err := js.Publish("events.created", msg)
			if err != nil {
				log.Printf("Error publishing message %d: %v", i+1, err)
			} else {
				log.Printf("Published: %s", msg)
			}
		}

	}()

	util.PrintStreamState(context.Background(), js, streamName)

	select {}
}
