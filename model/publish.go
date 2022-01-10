package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/nats-io/nats.go"
)

type order struct {
	Space string
	Point json.RawMessage
}

func main() {

	file, _ := ioutil.ReadFile("model.json")

	var i interface{}

	json.Unmarshal(file, &i)

	nc, _ := nats.Connect(nats.DefaultURL, nats.Token("mytoken"))
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	sendCH := make(chan *interface{})
	ec.BindSendChan("order", sendCH)

	sendCH <- &i

	nc.Drain()
	nc.Close()
}
