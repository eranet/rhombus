package rog

import "github.com/nats-io/go-nats"

// Handler is a specific callback used for Subscribe. It is generalized to
// an interface{}, but we will discover its format and arguments at runtime
// and perform the correct callback, including de-marshaling JSON strings
// back into the appropriate struct based on the signature of the Handler.
type Handler interface{}

type RogConnection struct {
	c *nats.EncodedConn
}

func Connection() *RogConnection {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	return &RogConnection{c}
}

func (rc *RogConnection) Subscribe(subj string, cb Handler) error {
	_, err := rc.c.Subscribe(subj, cb)
	return err
}

func (rc *RogConnection) Publish(subj string, v interface{}) error {
	return rc.c.Publish(subj, v)
}

func (rc *RogConnection) Spin() {
	doneChan := make(chan struct{})
	<-doneChan
}

func (rc *RogConnection) Close() {
	rc.c.Close()
}
