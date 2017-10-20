package certstream

import (
	"time"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"github.com/jmoiron/jsonq"
)

func CertStreamEventStream(skipHeartbeats bool) chan jsonq.JsonQuery {
	outputStream := make(chan jsonq.JsonQuery)

	var log = logging.MustGetLogger("CertStream")

	go func() {
		for {
			c, _, err := websocket.DefaultDialer.Dial("wss://certstream.calidog.io", nil)
			if err != nil {
				log.Error("Error connecting to certstream! Sleeping a few seconds and reconnecting... ", err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer c.Close()
			defer close(outputStream)

			for {
				var v interface{}
				err = c.ReadJSON(&v)
				if err != nil {
					log.Fatalf("Error decoding json frame!", err)
				}

				log.Info(v)

				jq := jsonq.NewQuery(v)

				res, err := jq.String("message_type")
				if err != nil {
					log.Fatalf("Error creating jq!", err)
				}

				if (skipHeartbeats && res == "heartbeat"){
					continue
				}

				outputStream <- *jq
			}
		}
	}()

	return outputStream
}

