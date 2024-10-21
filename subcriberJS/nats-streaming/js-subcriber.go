package natsstreaming

import (
	"log"

	"github.com/nats-io/nats.go"
)

func JsSubcriber() {
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

	sub, err := js.Subscribe("orders.created", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
		msg.Ack() // Acknowledge the message
	}, nats.ManualAck(), nats.Durable("order-processor"))
	//get name stream : nats stream ls -> ORDERS
	//get consumer for stream: nats consumer ls ORDERS

	// Uses a durable subscription (nats.Durable("order-processor")), which ensures that if the subscriber disconnects,
	// it will pick up from where it left off when it reconnects.
	if err != nil {
		log.Fatal(err)
	}
	// sub.Unsubscribe()
	log.Println("data subject = ", sub.Subject)
	log.Println("Listening for messages on 'orders.created'...")

	// Keep the process alive to continue receiving messages
	select {}
}
