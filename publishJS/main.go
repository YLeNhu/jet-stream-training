package main

import multiple_stream "go-streaming/multiple-stream"

func main() {
	// nats.PublisherNat()
	// natsstreaming.JsPublisher()
	//limitedbasedstream.JsPublisherLimit()
	// retentionpolicy.RetentionPolicy()
	// pullconsumer.JsPullPublisher()
	multiple_stream.PublishMultipleStream()
}
