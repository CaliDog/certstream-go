package main

import (
	"github.com/op/go-logging"
	"github.com/CaliDog/certstream-go"
)

var log = logging.MustGetLogger("example")

func main() {
	stream, errStream := certstream.CertStreamEventStream(false)
	for {
		select {
			case jq := <-stream:
				message_type, err := jq.String("message_type")

				if err != nil{
					log.Fatal("Error decoding jq string")
				}

				log.Info("Message type -> ", message_type)
				log.Info("recv: ", jq)

			case err := <-errStream:
				log.Error(err)

		}
	}
}
