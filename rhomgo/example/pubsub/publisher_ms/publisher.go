package main

import (
	"github.com/l1va/rhombus/rhomgo"
	"github.com/l1va/rhombus/rhomgo/example/pubsub"
)

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	r := rhomgo.NewRate(10) // 10 hz

	cnt := 0
	for {
		me := &pubsub.Person{
			Name:    "derek",
			Age:     22,
			Address: "140 New Montgomery Street, San Francisco, CA",
			Cnt:     cnt}
		c.Publish(pubsub.HelloTopic, me)
		cnt += 1
		r.Sleep()
	}
}
