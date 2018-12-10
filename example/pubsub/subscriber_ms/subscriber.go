package main

import (
	"github.com/l1va/roms"
	"github.com/l1va/roms/example/pubsub"
	"fmt"
)

func main() {
	c := roms.Connection(pubsub.ServerURL)
	defer c.Close()

	c.Subscribe(pubsub.HelloTopic, func(p *pubsub.Person) {
		fmt.Printf("Received a person: %+v\n", p)
	})

	c.Spin()
}
