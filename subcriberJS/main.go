package main

import limitedbasedstream "sub/limited-based-stream"

func main() {
	// nats.Subcriber()
	// natsstreaming.JsSubcriber()
	// limitedbasedstream.JsSubcriberLimit()
	limitedbasedstream.MaxAckPendingSubcriber()
	// pullconsumer.JsPullSubcriber()
}
