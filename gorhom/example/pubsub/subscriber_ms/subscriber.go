package main

import (
	"fmt"
	"github.com/l1va/rhombus/gorhom"
	"github.com/l1va/rhombus/gorhom/example/pubsub"
)

func main() {
	c := gorhom.LocalJSONConnection()
	defer c.Close()

	c.Subscribe(pubsub.HelloTopic, func(p *pubsub.Person) {
		fmt.Printf("Received a person: %+v\n", p)
	})

	c.Spin()
}
