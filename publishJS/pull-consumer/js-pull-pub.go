package pullconsumer

import (
	"log"

	"github.com/nats-io/nats.go"
)

func JsPullPublisher() {
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

	// Define a stream, which stores messages under subjects
	js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.*"},
	})

	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 4; i++ {
		msg := []byte("events created with ID: " + string(i+48)) // ASCII conversion to keep it simple
		_, err := js.Publish("events.created", msg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Published message: %s", msg)
	}

}
