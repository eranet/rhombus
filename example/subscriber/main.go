package main

import (
	"rog"
	"rog/example"
	"fmt"
)

func main() {
	c := rog.Connection()
	defer c.Close()

	c.Subscribe(example.HelloTopic, func(p *example.Person) {
		fmt.Printf("Received a person: %+v\n", p)
	})

	c.Spin()
}
