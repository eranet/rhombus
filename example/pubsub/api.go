package pubsub

import "github.com/nats-io/go-nats"

const ServerURL = nats.DefaultURL

type Person struct {
	Name    string
	Address string
	Age     int
	Cnt     int
}

const HelloTopic = "hello"
