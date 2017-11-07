package certstream

import (
	"time"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/jsonq"
	"github.com/pkg/errors"
)

func CertStreamEventStream(skipHeartbeats bool) (chan jsonq.JsonQuery, chan error) {
	outputStream := make(chan jsonq.JsonQuery)
	errStream := make(chan error)

	go func() {
		for {
			c, _, err := websocket.DefaultDialer.Dial("wss://certstream.calidog.io", nil)

			if err != nil {
				errStream <- errors.Wrap(err, "Error connecting to certstream! Sleeping a few seconds and reconnecting... ")
				time.Sleep(5 * time.Second)
				continue
			}

			defer c.Close()
			defer close(outputStream)

			for {
				var v interface{}
				err = c.ReadJSON(&v)
				if err != nil {
					errStream <- errors.Wrap(err, "Error decoding json frame!")
				}

				jq := jsonq.NewQuery(v)

				res, err := jq.String("message_type")
				if err != nil {
					errStream <- errors.Wrap(err, "Error creating jq object!")
				}

				if skipHeartbeats && res == "heartbeat" {
					continue
				}

				outputStream <- *jq
			}
		}
	}()

	return outputStream, errStream
}

