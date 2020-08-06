package rhomgo

import (
	"errors"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
)

// Handler is a specific callback used for Subscribe. It is generalized to
// an interface{}, but we will discover its format and arguments at runtime
// and perform the correct callback, including de-marshaling JSON strings
// back into the appropriate struct based on the signature of the Handler.
type Handler interface{}

type RhombusConnection struct {
	c        *nats.EncodedConn
	doneChan chan struct{}
}

func getConnection(serverURL string, encoder string) *RhombusConnection {
	nc, _ := nats.Connect(serverURL)
	c, _ := nats.NewEncodedConn(nc, encoder)
	return &RhombusConnection{c, make(chan struct{})}
}

// LocalJSONConnection create a new connection to the local messaging server with json encoder
func LocalJSONConnection() *RhombusConnection {
	return getConnection(nats.DefaultURL, nats.JSON_ENCODER)
}

// LocalBinaryConnection create a new connection to the local messaging server with binary encoder
func LocalBinaryConnection() *RhombusConnection {
	return getConnection(nats.DefaultURL, nats.DEFAULT_ENCODER)
}

// JSONConnection create a new connection to the messaging server with json encoder
func JSONConnection(serverURL string) *RhombusConnection {
	return getConnection(serverURL, nats.JSON_ENCODER)
}

// BinaryConnection for binary data
func BinaryConnection(serverURL string) *RhombusConnection {
	return getConnection(serverURL, nats.DEFAULT_ENCODER)
}

// Subscribe to the given topic
// subj: topic
// cb: callback function, must receive one struct parameter
func (rc *RhombusConnection) Subscribe(subj string, cb Handler) error {
	_, err := rc.c.Subscribe(subj, cb)
	return err
}

// Publish a struct to given topic
// subj: topic
// v: struct
func (rc *RhombusConnection) Publish(subj string, v interface{}) error {
	return rc.c.Publish(subj, v)
}

// Publish a request to given responder_ms and wait
// subj: topic
// v: request struct
// vPtr: reply struct, must be pointer
// timeout: timeout duration
func (rc *RhombusConnection) Request(subj string, v interface{},
	vPtr interface{}, timeout time.Duration) error {
	err := rc.c.Request(subj, v, vPtr, timeout)
	return err
}

// Advertise a responder_ms
// subj: topic
// cb: callback function, must receive 1 struct parameter
// and return 1 struct output
func (rc *RhombusConnection) Service(subj string, cb Handler) error {
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
func (rc *RhombusConnection) Spin() {
	<-rc.doneChan
}

// method for waiting until the process is terminated
func (rc *RhombusConnection) SpinDone() {
	rc.doneChan <- struct{}{}
}

// close connection to the messaging server
func (rc *RhombusConnection) Close() {
	rc.c.Close()
}
