package main

import limitedbasedstream "go-streaming/limited-based-stream"

func main() {
	// nats.PublisherNat()
	// natsstreaming.JsPublisher()
	limitedbasedstream.JsPublisherLimit()
	// retentionpolicy.RetentionPolicy()
	// pullconsumer.JsPullPublisher()
}
