package main

import (
	"github.com/l1va/rhombus"
	"github.com/l1va/rhombus/example/pubsub"
	"fmt"
)

func main() {
	c := rhombus.LocalConnection()
	defer c.Close()

	c.Subscribe(pubsub.HelloTopic, func(p *pubsub.Person) {
		fmt.Printf("Received a person: %+v\n", p)
	})

	c.Spin()
}
