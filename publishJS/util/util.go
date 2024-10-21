package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"golang.org/x/exp/rand"
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

func CheckInfoStreamExist(js nats.JetStreamContext, streamName string) bool {
	_, err := js.StreamInfo(streamName)
	return err == nil
}

func PrintStreamState(ctx context.Context, js nats.JetStreamContext, streamName string) {
	// Get the stream info
	info, err := js.StreamInfo(streamName)
	if err != nil {
		log.Fatalf("Error getting stream info: %v", err)
	}

	// Pretty-print the stream state
	b, err := json.MarshalIndent(info.State, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling stream state: %v", err)
	}

	fmt.Println("Inspecting stream info:")
	fmt.Println(string(b))
}

func GenerateLargePayload(size int) []byte {
	rand.Seed(uint64(time.Now().Unix()))
	payload := make([]byte, size)
	rand.Read(payload) // Fill it with random data
	return payload
}

type LargeMessage struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Payload []byte `json:"payload"` // This is where we store a large byte array
}
