package limitedbasedstream

import (
	"fmt"
	"log"
	"sub/util"

	"github.com/nats-io/nats.go"
)

func JsSubcriberLimit() {
	// connect nats
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	// subcriber
	sub, err := js.Subscribe("events.created", func(msg *nats.Msg) {
		message := string(msg.Data)
		log.Printf("Received message: %s", message)
		msg.Ack() // Acknowledge the message
		msg.Respond([]byte(fmt.Sprintf("nhận message rồi nha message là %s", message)))
	}, nats.ManualAck(), nats.Durable("events-processor"), nats.MaxAckPending(20))
	// , nats.Durable("events-processor")

	if err != nil {
		log.Fatal(err)
	}
	// sub.Unsubscribe()
	log.Println("data subject = ", sub.Subject)
	log.Println("Listening for messages on 'orders.created'...")

	// Keep the process alive to continue receiving messages
	select {}

}
