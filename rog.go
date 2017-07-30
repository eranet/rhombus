package rog

import "github.com/nats-io/go-nats"

func Connection() *nats.EncodedConn {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	return c
}

func Listen() {
	doneChan := make(chan struct{})
	<-doneChan
}
