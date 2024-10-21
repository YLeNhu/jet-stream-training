package util

import (
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNat() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	return nc
}

func CreateStream(js nats.JetStreamContext, config nats.StreamConfig) {
	jsInfo, err := js.AddStream(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stream EVENTS created with limits", jsInfo)
}
