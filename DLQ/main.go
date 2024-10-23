package main

import (
	"dead-letter-queue/util"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	fmt.Println("start subcriber dead letter queue")
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	consumerConfig := &nats.ConsumerConfig{
		Durable:       "dead_letter_queue_subscriber", // Durable name of the consumer
		AckPolicy:     nats.AckExplicitPolicy,         // Explicit acknowledgment required
		AckWait:       3 * time.Second,                // Wait 10 seconds for an ack before retry
		MaxAckPending: 10,                             // Maximum number of unacknowledged messages
		MaxDeliver:    3,                              // Retry each message 3 times
		DeliverPolicy: nats.DeliverAllPolicy,          // Deliver all messages from the start
		FilterSubject: "dlq.>",                        // Listen for all "metrics.*" subjects
	}

	_, err = js.AddConsumer("METRICS_DEAD_LETTER_QUEUE", consumerConfig)

	if err != nil {
		log.Fatalf("Error adding consumer: %v", err)
	}

	sub, err := js.PullSubscribe("dlq.>", "dead_letter_queue_subscriber")
	if err != nil {
		log.Fatalf("Error subscribing to DLQ: %v", err)
	}

	for {
		// Fetch messages in batches from DLQ
		msgs, err := sub.Fetch(10, nats.MaxWait(5*time.Second))
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			continue
		}

		for _, msg := range msgs {
			// Handle DLQ message
			fmt.Printf("Received message from DLQ: Subject: %s, Data: %s\n", msg.Subject, string(msg.Data))

			// Process the message (this is where you would handle the failure)
			// Acknowledge the message if processed successfully
			err := msg.Ack()
			if err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}
}
