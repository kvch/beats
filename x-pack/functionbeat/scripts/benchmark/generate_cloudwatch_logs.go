// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var (
	eventCount int
)

func init() {
	flag.IntVar(&eventCount, "count", 5, "number of messages to generate")
}

func main() {
	flag.Parse()
	testEvents := generateCloudwatchLogEvents(eventCount)
	b, err := json.Marshal(testEvents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func generateCloudwatchLogEvents(n int) events.CloudwatchLogsEvent {
	logEvents := make([]events.CloudwatchLogsLogEvent, n)
	for i := 0; i < n; i++ {
		num := strconv.Itoa(i)
		logEvents[i] = events.CloudwatchLogsLogEvent{
			ID:        "1234-" + num,
			Timestamp: time.Now().Unix(),
			Message:   "hello world " + num,
		}
	}

	rawEvent := events.CloudwatchLogsData{
		Owner:     "foobar",
		LogGroup:  "foo",
		LogStream: "/var/foobar",
		LogEvents: logEvents,
	}

	b, _ := json.Marshal(&rawEvent)

	data := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, data)
	zw := gzip.NewWriter(encoder)
	zw.Write(b)
	zw.Close()
	encoder.Close()

	return events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: data.String(),
		},
	}
}
