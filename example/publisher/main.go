package main

import (
	"rog"
	"rog/example"
)

func main() {
	c := rog.Connection()
	defer c.Close()

	r := rog.NewRate(10) // 10 hz

	cnt := 0
	for {
		me := &example.Person{
			Name: "derek",
			Age: 22,
			Address: "140 New Montgomery Street, San Francisco, CA",
			Cnt: cnt}
		c.Publish(example.HelloTopic, me)
		cnt += 1
		r.Sleep()
	}
}
