package main

import multiple_stream "sub/multiple-stream"

func main() {
	// nats.Subcriber()
	// natsstreaming.JsSubcriber()
	// limitedbasedstream.JsSubcriberLimit()
	//limitedbasedstream.MaxAckPendingSubcriber()
	// pullconsumer.JsPullSubcriber()
	multiple_stream.SubcriberMultipleStream()
	//multiple_stream.SubcriberDeadLetterQueue()
}
