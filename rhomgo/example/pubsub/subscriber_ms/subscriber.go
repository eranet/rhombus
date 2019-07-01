package main

import (
	"fmt"
	"github.com/l1va/rhombus/rhomgo"
	"github.com/l1va/rhombus/rhomgo/example/pubsub"
)

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	c.Subscribe(pubsub.HelloTopic, func(p *pubsub.Person) {
		fmt.Printf("Received a person: %+v\n", p)
	})

	c.Spin()
}
