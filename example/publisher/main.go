package main

import (
	"rog"
	"rog/example"
)

func main() {
	c := rog.Connection()
	defer c.Close()

	me := &example.Person{Name: "derek", Age: 22, Address: "140 New Montgomery Street, San Francisco, CA"}

	c.Publish(example.HelloTopic, me)
}
