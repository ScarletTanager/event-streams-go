// Example function-based high-level Apache Kafka consumer
package main

/**
 * Copyright 2016 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// consumer_example implements a consumer using the non-channel Poll() API
// to retrieve messages and events.

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	// if len(os.Args) < 4 {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
	// 		os.Args[0])
	// 	os.Exit(1)
	// }

	// broker := os.Args[1]
	// group := os.Args[2]
	// password := os.Args[3]
	// topics := os.Args[4:]
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "broker-0-0d8c8z8plq5d25gn.kafka.svc11.us-south.eventstreams.cloud.ibm.com:9093,broker-1-0d8c8z8plq5d25gn.kafka.svc11.us-south.eventstreams.cloud.ibm.com:9093,broker-2-0d8c8z8plq5d25gn.kafka.svc11.us-south.eventstreams.cloud.ibm.com:9093,broker-3-0d8c8z8plq5d25gn.kafka.svc11.us-south.eventstreams.cloud.ibm.com:9093,broker-4-0d8c8z8plq5d25gn.kafka.svc11.us-south.eventstreams.cloud.ibm.com:9093,broker-5-0d8c8z8plq5d25gn.kafka.svc11.us-south.eventstreams.cloud.ibm.com:9093",
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family":                 "v4",
		"group.id":                              "myGroup",
		"session.timeout.ms":                    6000,
		"auto.offset.reset":                     "earliest",
		"sasl.username":                         "token",
		"sasl.password":                         "iqL6fNnrB99Z1LUvoettHz7MyDDzGe48hKJKZbx7Nvge",
		"security.protocol":                     "sasl_ssl",
		"sasl.mechanism":                        "PLAIN",
		"ssl.endpoint.identification.algorithm": "HTTPS",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics([]string{"myTopic"}, nil)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}
