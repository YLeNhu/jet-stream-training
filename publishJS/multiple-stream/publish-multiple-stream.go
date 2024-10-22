package multiple_stream

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go-streaming/util"
	"log"
	"reflect"
	"sync"
)

func PublishMultipleStream() {
	nc := util.ConnectNat()
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()

	if err != nil {
		log.Fatal(err)
	}

	streamName := "METRICS_WORLDWIDE"

	util.CreateStream(js, nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"metrics.>"},
		Storage:  nats.FileStorage,
	})

	var waitGroup sync.WaitGroup
	//{"CPU usage: 45%", "Memory usage: 45%", "Network usage: 300mbps"}
	countryMetrics := map[util.Country][]util.Metrics{
		util.Singapore: {
			{CPU: "CPU usage: 45%", Memory: "Memory usage: 45%", Network: "Network usage: 300mbps"},
		},
		util.USA:     {},
		util.Germany: {},
	}

	waitGroup.Add(1)
	go func() {
		for country, metrics := range countryMetrics {
			fmt.Printf("processing country %v\n", country)
			for _, item := range metrics {
				message := processMetric(&item, country, &js)
				print(message)
			}
		}

		defer waitGroup.Done()
	}()

	waitGroup.Wait()
	//select {}
}

func processMetric(metric *util.Metrics, country util.Country, js *nats.JetStreamContext) string {
	val := reflect.ValueOf(metric).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i) // Get the field
		message := val.Field(i).String()
		switch field.Name {
		case "CPU":
			emitMessage(js, fmt.Sprintf("metrics.cpu.%s", country), message)
			//case "Memory":
			//	emitMessage(js, fmt.Sprintf("metrics.memory.%s", country), message)
			//case "Network":
			//	emitMessage(js, fmt.Sprintf("metrics.network.%s", country), message)
			//default:
			//	fmt.Println("Unknown metric")
		}
	}

	return "send success"
}

func emitMessage(js *nats.JetStreamContext, subject string, message string) {
	for i := 0; i < 10; i++ {
		_, err := (*js).Publish(subject, []byte(fmt.Sprintf("%s-index-%d", message, i)))
		if err != nil {
			log.Printf("Error publishing message  %v", err)
		} else {
			log.Printf("Published: %s", message)
		}
	}

	return
}
