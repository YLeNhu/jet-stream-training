package multiple_stream

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"strings"
	"sub/util"
	"sync"
	"time"
)

type MetricMessage struct {
	Country string
	Metric  string
	Value   string
}

func SubcriberMultipleStream() {
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	fanInChannel := make(chan MetricMessage)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("start Subcriber metrics")
		//PullSubscribe -> async function

		sub, err := js.PullSubscribe("metrics.>", "metrics_subscriber",
			nats.MaxAckPending(10), nats.AckWait(5*time.Second), nats.ManualAck())

		if err != nil {
			log.Fatalf("Error subscribing to stream: %v", err)
		}
		for {
			msgs, err := sub.Fetch(10_000, nats.MaxWait(5*time.Second)) // fetch is sync function
			if err != nil {
				log.Printf("Error fetching messages: %v", err)
				continue
			}

			for _, msg := range msgs {
				// Parse the message subject to extract metricType and country
				parts := strings.Split(msg.Subject, ".")
				if len(parts) < 3 {
					log.Printf("Invalid subject format: %s", msg.Subject)
					msg.Nak() // NAK (Negative Acknowledge) for invalid messages
					continue
				}

				metricType := parts[1]
				country := parts[2]

				// Send the message to the fanIn channel for processing
				fanInChannel <- MetricMessage{
					Country: country,
					Metric:  metricType,
					Value:   string(msg.Data),
				}

				index, _ := strconv.Atoi(strings.Split(string(msg.Data), "-")[2])

				if index%2 == 0 {
					if err := msg.Ack(); err != nil {
						log.Printf("Failed to acknowledge message: %v", err)
					}
				}
				// Acknowledge the message after processing

				util.CheckPendingAcks(js, "METRICS_WORLDWIDE", "metrics_subscriber")
			}
		}

		// using Subscribe
		//func(msg *nats.Msg) {
		//	parts := strings.Split(msg.Subject, ".")
		//
		//	metricType := parts[1]
		//	country := parts[2]
		//
		//	fanInChannel <- MetricMessage{
		//		Country: country,
		//		Metric:  metricType,
		//		Value:   string(msg.Data),
		//	}
		//
		//	if err := msg.Ack(); err != nil {
		//		log.Printf("Failed to acknowledge message: %v", err)
		//	}
		//}
		//
		//if err != nil {
		//	log.Fatalf("Error subscribing to wildcard subject: %v", err)
		//} else {
		//	log.Println("Subscribed to all metrics.* subjects")
		//}
	}()

	go func() {
		for message := range fanInChannel {
			// Process the fanned-in message
			fmt.Printf("Received Metric from %s - %s: %s\n", message.Country, message.Metric, message.Value)
		}
	}()

	// Wait for all subscriptions to complete
	wg.Wait()

	select {}
}
