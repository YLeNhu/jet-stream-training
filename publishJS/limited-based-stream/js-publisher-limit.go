package limitedbasedstream

import (
	"context"
	"go-streaming/util"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

func JsPublisherLimit() {
	// connect nats
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	streamName := "EVENTS"

	if util.CheckInfoStreamExist(js, streamName) {
		log.Println("Stream already exists:", streamName)
	} else {
		// create stream
		util.CreateStream(js, nats.StreamConfig{
			Name:       "EVENTS",
			Subjects:   []string{"events.>"},
			Storage:    nats.FileStorage, // Choose between File or Memory storage
			MaxMsgs:    100,              // Maximum number of messages if have new message discard old message
			MaxBytes:   10 * 1024 * 1024, // Maximum size in bytes (10 MB)
			MaxAge:     5 * time.Minute,  // Maximum age of messages -> nano second
			MaxMsgSize: 1 * 1024 * 1024,
		})
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	// send large message -> Error publishing message 100: nats: maximum payload exceeded
	// message := util.LargeMessage{
	// 	ID:      "1234",
	// 	Name:    "TestMessage",
	// 	Payload: util.GenerateLargePayload(1 * 1024 * 1024), // 1 MB payload
	// }

	// data, err := json.Marshal(message)

	go func() {
		for i := 0; i < 16; i++ {
			msg := []byte("Event created with ID: gia bao") // ASCII conversion to keep it simple
			_, err := js.Publish("events.created", msg)
			if err != nil {
				log.Printf("Error publishing message %d: %v", i+1, err)
			} else {
				log.Printf("Published: %s", msg)
			}
		}

		defer waitGroup.Done()
	}()
	waitGroup.Wait()
	util.PrintStreamState(context.Background(), js, streamName)

	select {}
}
