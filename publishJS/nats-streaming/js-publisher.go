package natsstreaming

import (
	"log"

	"github.com/nats-io/nats.go"
)

func JsPublisher() {
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
		Name:     "ORDERS",
		Subjects: []string{"orders.*"},
	})

	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 5; i++ {
		msg := []byte("Order created with ID: " + string(i+48)) // ASCII conversion to keep it simple
		_, err := js.Publish("orders.created", msg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Published message: %s", msg)
	}
	
}
