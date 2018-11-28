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

func getConnection(encoder string) *RogConnection {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, encoder)
	return &RogConnection{c}
}

// Create a new connection to the messaging server
func Connection() *RogConnection {
	return getConnection(nats.JSON_ENCODER)
}

// BinaryConnection for binary data
func BinaryConnection() *RogConnection {
	return getConnection(nats.DEFAULT_ENCODER)
}

// Subscribe to the given topic
// subj: topic
// cb: callback function, must receive one struct parameter
func (rc *RogConnection) Subscribe(subj string, cb Handler) error {
	_, err := rc.c.Subscribe(subj, cb)
	return err
}

// Publish a struct to given topic
// subj: topic
// v: struct
func (rc *RogConnection) Publish(subj string, v interface{}) error {
	return rc.c.Publish(subj, v)
}

// Publish a request to given service and wait
// subj: topic
// v: request struct
// vPtr: reply struct, must be pointer
// timeout: timeout duration
func (rc *RogConnection) RequestService(subj string, v interface{},
	vPtr interface{}, timeout time.Duration) error {
	err := rc.c.Request(subj, v, vPtr, timeout)
	return err
}

// Advertise a service
// subj: topic
// cb: callback function, must receive 1 struct parameter
// and return 1 struct output
func (rc *RogConnection) Service(subj string, cb Handler) error {
	// need to do a little bit of reflect magic
	// in order to be able to receive callback of any type
	// (inspired by go-nats client)
	cbType := reflect.TypeOf(cb)

	if cbType.Kind() != reflect.Func {
		return errors.New("Handler needs to be a func")
	}
	if cbType.NumIn() != 1 {
		return errors.New("Handler needs to have 1 parameter")
	}
	if cbType.NumOut() != 1 {
		return errors.New("Handler needs to have 1 output parameter")
	}
	argType := cbType.In(0) // type of the first argument in callback function

	cbValue := reflect.ValueOf(cb) // value of callback function so that we can call it later
	err := rc.Subscribe(subj, func(msg *nats.Msg) {
		// manually decoding incoming message to a suitable struct
		var oPtr reflect.Value
		if argType.Kind() != reflect.Ptr {
			oPtr = reflect.New(argType)
		} else {
			oPtr = reflect.New(argType.Elem())
		}
		_ = rc.c.Enc.Decode(msg.Subject, msg.Data, oPtr.Interface())
		// if callback argument is not a pointer
		if argType.Kind() != reflect.Ptr {
			oPtr = reflect.Indirect(oPtr)
		}

		// calling the callback
		r := cbValue.Call([]reflect.Value{oPtr})
		// publishing the reply
		rc.Publish(msg.Reply, r[0].Interface())
	})
	return err
}

// method for waiting until the process is terminated
func (rc *RogConnection) Spin() {
	doneChan := make(chan struct{})
	<-doneChan
}

// close connection to the messaging server
func (rc *RogConnection) Close() {
	rc.c.Close()
}
