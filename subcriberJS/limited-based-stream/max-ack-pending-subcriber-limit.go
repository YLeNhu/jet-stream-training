package limitedbasedstream

import (
	"fmt"
	"log"
	"sub/util"
	"time"

	"github.com/nats-io/nats.go"
)

func MaxAckPendingSubcriber() {
	log.Println("Listening for messages on 'orders.created'...")

	// connect nats
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	// subcriber
	// nats.MaxAckPending(20)
	// Create a Pull-based consumer (without MaxAckPending limit)
	sub, err := js.PullSubscribe("events.created", "worker", nats.MaxAckPending(20))
	if err != nil {
		log.Fatal(err)
	}
	count := 0

	// Keep the process alive to continue receiving messages
	for {
		// Fetch up to 10 messages
		msgs, err := sub.Fetch(10, nats.MaxWait(0))
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			continue
		}

		// Process each message
		for _, msg := range msgs {
			fmt.Printf("Received message: %s\n", string(msg.Data))

			// Simulate slow processing without acknowledging
			time.Sleep(2 * time.Second)
			fmt.Println("after 2s")

			count += 1
			// Eventually acknowledge the message
			msg.Ack()

			fmt.Printf("total message: %d\n", count)

			util.CheckPendingAcks(js, "EVENTS", "worker")
		}
	}

}
