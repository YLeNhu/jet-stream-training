package nats

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

func Subcriber() {
	var wg sync.WaitGroup
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}
	defer nc.Close()

	subject := "ping"

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting subscriber...")

		print("waiting listen message")
		// Subscribe normal
		// QueueSubscribe add into group
		_, err = nc.QueueSubscribe(subject, "test", func(m *nats.Msg) {
			message := fmt.Sprintf("Received message: %s", string(m.Data))
			log.Printf(message)
			// Simply reply with "pong"
			m.Respond([]byte(message))
		})
		if err != nil {
			log.Fatalf("Error subscribing to subject %s: %v", subject, err)
		}
		// Keep the subscriber running for 30 seconds

		select {} // Block forever here to keep the goroutine alive

	}()

	wg.Wait()
}
