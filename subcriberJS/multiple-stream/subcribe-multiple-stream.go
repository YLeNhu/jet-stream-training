package multiple_stream

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sub/util"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
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

	// create dead letter queue
	util.CreateStreamDeadLetterQueue(js)

	if err != nil {
		log.Fatal(err)
	}

	//  AckNonePolicy: delivery done and auto ack
	// AckExplicitPolicy: manual ack
	// Define consumer configuration
	consumerConfig := &nats.ConsumerConfig{
		Durable:       "metrics_subscriber",   // Durable name of the consumer
		AckPolicy:     nats.AckExplicitPolicy, // Explicit acknowledgment required
		AckWait:       3 * time.Second,        // Wait 10 seconds for an ack before retry
		MaxAckPending: 10,                     // Maximum number of unacknowledged messages
		MaxDeliver:    3,                      // Retry each message 3 times
		DeliverPolicy: nats.DeliverAllPolicy,  // Deliver all messages from the start
		FilterSubject: "metrics.>",            // Listen for all "metrics.*" subjects
	}

	_, err = js.AddConsumer("METRICS_WORLDWIDE", consumerConfig)
	if err != nil {
		log.Fatalf("Error adding consumer: %v", err)
	}

	fanInChannel := make(chan MetricMessage)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			// Recover from panic and log the error
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()

		fmt.Println("start Subcriber metrics")
		//PullSubscribe -> async function

		sub, err := js.PullSubscribe("metrics.>", "metrics_subscriber") // nats.MaxAckPending(10), nats.AckWait(5*time.Second), nats.ManualAck()

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
				err := processMessage(js, msg, fanInChannel)
				fmt.Printf("error ne %v\n", err)
				if err != nil {
					sendToDLQ(js, msg)
				}
				// Parse the message subject to extract metricType and country
				// parts := strings.Split(msg.Subject, ".")
				// if len(parts) < 3 {
				// 	log.Printf("Invalid subject format: %s", msg.Subject)
				// 	msg.Nak() // NAK (Negative Acknowledge) for invalid messages
				// 	continue
				// }

				// metricType := parts[1]
				// country := parts[2]

				// // Send the message to the fanIn channel for processing
				// fanInChannel <- MetricMessage{
				// 	Country: country,
				// 	Metric:  metricType,
				// 	Value:   string(msg.Data),
				// }

				// index, _ := strconv.Atoi(strings.Split(string(msg.Data), "-")[2])

				// if index%2 == 0 {
				// 	if err := msg.Ack(); err != nil {
				// 		log.Printf("Failed to acknowledge message: %v", err)
				// 	}
				// }
				// // Acknowledge the message after processing

				// util.CheckPendingAcks(js, "METRICS_WORLDWIDE", "metrics_subscriber")
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

func processMessage(js nats.JetStreamContext, msg *nats.Msg, fanInChannel chan MetricMessage) (err error) {
	time.Sleep(time.Second * 2)
	defer func() {
		// Recover from panic inside message processing
		if r := recover(); r != nil {
			fmt.Printf("recover bao ne %v\n", r)
			err = fmt.Errorf("panic occurred: %v", r) // Set the error to return to indicate the panic
			// Negative Acknowledge the message to retry later
			errNak := msg.Nak()
			if errNak != nil {
				err = fmt.Errorf("nak error ")
			}
		}
	}()
	// Simulate panic for specific message content
	parts := strings.Split(msg.Subject, ".")
	if len(parts) < 3 {
		log.Printf("Invalid subject format: %s", msg.Subject)
		err := msg.Nak()
		if err != nil {
			return fmt.Errorf("message nak error ")
		}
		return fmt.Errorf("invalid subject format: %s", msg.Subject)
	}

	metricType := parts[1]
	country := parts[2]

	// Send the message to the fanIn channel
	fanInChannel <- MetricMessage{
		Country: country,
		Metric:  metricType,
		Value:   string(msg.Data),
	}

	// Simulate panic on purpose
	index, err := strconv.Atoi(strings.Split(string(msg.Data), "-")[2])
	if err != nil {
		log.Printf("Invalid message data: %v", err)
		// Terminate the invalid message so it won't be retried
		if err := msg.Term(); err != nil {
			log.Printf("Failed to terminate message: %v", err)
		}
		return fmt.Errorf("invalid message data: %v", err)
	}

	// Forcing a panic when index % 5 == 0 to simulate an unexpected crash
	if index%2 != 0 {
		panic(fmt.Sprintf("Simulated crash: index %d caused a panic", index))
	}

	// Normal message acknowledgment if index % 2 == 0
	if index%2 == 0 {
		if err := msg.Ack(); err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		}
	}

	//util.CheckPendingAcks(js, "METRICS_WORLDWIDE", "metrics_subscriber")
	return nil
}

func sendToDLQ(js nats.JetStreamContext, msg *nats.Msg) {
	fmt.Println("send dead letter queue")
	dlqSubject := "dlq.metrics" // Subject for the DLQ

	// Publish the failed message to the DLQ stream
	_, err := js.Publish(dlqSubject, msg.Data)
	if err != nil {
		log.Printf("Failed to publish message to DLQ: %v", err)
	}

	// Terminate the original message so it won't be retried
	if err := msg.Term(); err != nil {
		log.Printf("Failed to terminate message: %v", err)
	}
}
