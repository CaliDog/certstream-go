package certstream

import (
	"time"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/jsonq"
)

type CertstreamError struct {
	errorCode string
	errorString string
	error *error
}

func CertStreamEventStream(skipHeartbeats bool) (chan jsonq.JsonQuery, chan CertstreamError) {
	outputStream := make(chan jsonq.JsonQuery)
	errStream := make(chan CertstreamError)

	go func() {
		for {
			c, _, err := websocket.DefaultDialer.Dial("wss://certstream.calidog.io", nil)

			if err != nil {
				errStream <- CertstreamError{
					"CONNECTION_ERROR",
					"Error connecting to certstream! Sleeping a few seconds and reconnecting... ",
					&err,
				}
				time.Sleep(5 * time.Second)
				continue
			}

			defer c.Close()
			defer close(outputStream)

			for {
				var v interface{}
				err = c.ReadJSON(&v)
				if err != nil {
					errStream <- CertstreamError{
						"JSON_DECODE_ERROR",
						"Error decoding json frame!",
						&err,
					}
				}

				jq := jsonq.NewQuery(v)

				res, err := jq.String("message_type")
				if err != nil {
					errStream <- CertstreamError{
						"JQ_ERROR",
						"Error creating jq object!",
						&err,
					}
				}

				if (skipHeartbeats && res == "heartbeat"){
					continue
				}

				outputStream <- *jq
			}
		}
	}()

	return outputStream, errStream
}

