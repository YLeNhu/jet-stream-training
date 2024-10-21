package nats

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

func PublisherNat() {
	var wg sync.WaitGroup
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}
	defer nc.Drain() // prevent lost data or discarded
	defer nc.Close()

	// Ping the broker by requesting a reply from a simple service
	subject := "ping"

	//time.Sleep(1 * time.Second)

	// Start the ping check
	fmt.Println("Pinging NATS broker...")

	// Request a reply from the "ping" subject and wait for a response

	fmt.Println("Pinging NATS broker...")

	start := time.Now()
	latency := time.Since(start)
	channel := make(chan string, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			task := fmt.Sprintf("task %d", i)
			msg, err := nc.Request(subject, []byte(task), 2*time.Second)
			if err != nil {
				log.Fatalf("Error during ping request: %v", err)
			}
			channel <- string(msg.Data)
		}(i)
	}

	wg.Wait()
	close(channel) // lack of will run forever

	for v := range channel {
		fmt.Printf("data response: %s\n", v)
		fmt.Printf("Ping successful! Latency: %v\n", latency)

	}
}
