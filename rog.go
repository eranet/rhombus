package rog

import (
	"errors"
	"reflect"
	"time"
	"github.com/nats-io/go-nats"
)

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

func (rc *RogConnection) RequestService(subj string, v interface{},
		vPtr interface{}, timeout time.Duration) error {
	err := rc.c.Request(subj, v, vPtr, timeout)
	return err
}

func (rc *RogConnection) Service(subj string, cb Handler) error {
	cbType := reflect.TypeOf(cb)

	if cbType.Kind() != reflect.Func {
		return errors.New("Handler needs to be a func")
	}
	if cbType.NumIn() != 1 {
		return errors.New("Handler needs to have 1 parameter")
	}
	if cbType.NumOut() != 1{
		return errors.New("Handler needs to have 1 output parameter")
	}
	argType := cbType.In(0)

	cbValue := reflect.ValueOf(cb)
	err := rc.Subscribe(subj, func(msg *nats.Msg) {
		var oPtr reflect.Value
		if argType.Kind() != reflect.Ptr {
			oPtr = reflect.New(argType)
		} else {
			oPtr = reflect.New(argType.Elem())
		}
		_ = rc.c.Enc.Decode(msg.Subject, msg.Data, oPtr.Interface())

		r := cbValue.Call([]reflect.Value{oPtr})
		rc.Publish(msg.Reply, r[0].Interface())
	})
	return err
}

func (rc *RogConnection) Spin() {
	doneChan := make(chan struct{})
	<-doneChan
}

func (rc *RogConnection) Close() {
	rc.c.Close()
}
