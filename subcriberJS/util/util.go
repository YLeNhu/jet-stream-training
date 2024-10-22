package util

import (
	"fmt"
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

func CheckPendingAcks(js nats.JetStreamContext, streamName string, consumerName string) {
	consumerInfo, err := js.ConsumerInfo(streamName, consumerName)
	if err != nil {
		log.Printf("Error fetching consumer info: %v", err)
		return
	}

	fmt.Printf("Pending Acks: %d / Max Pending Acks: %d\n", consumerInfo.NumAckPending, consumerInfo.Config.MaxAckPending)
}
