package pullconsumer

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func JsPullSubcriber() {
	log.Printf("Start fetching messages")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	// PullSubscribe
	sub, err := js.PullSubscribe("events.*", "worker-1")
	if err != nil {
		log.Fatal(err)
	}

	for {

		//  nats.MaxWait(10*time.Second) -> ensuring that your consumer doesn't block indefinitely while waiting for messages
		// Fetch up to 5 messages with a maximum wait time of 10 seconds
		msgs, err := sub.Fetch(5, nats.MaxWait(10*time.Second))
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			time.Sleep(1 * time.Second) // Wait before retrying if there's an error
			continue
		}

		// Process each message
		for _, msg := range msgs {
			fmt.Printf("Received message: %s\n", string(msg.Data))

			// Acknowledge the message
			msg.Ack()
		}

		// Simulate some processing delay
		time.Sleep(2 * time.Second)
	}
}
