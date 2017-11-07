package main

import (
	"github.com/CaliDog/certstream-go"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func main() {
	stream := certstream.CertStreamEventStream(false)

	for jq := range stream {

		message_type, err := jq.String("message_type")

		if err != nil {
			log.Fatalf("Error parsing message_type: %v", err)
		}

		log.Info("Message type -> ", message_type)
		log.Info("recv: ", jq)

	}
}
