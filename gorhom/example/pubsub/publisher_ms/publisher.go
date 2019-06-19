package main

import (
	"github.com/l1va/rhombus/gorhom"
	"github.com/l1va/rhombus/gorhom/example/pubsub"
)

func main() {
	c := gorhom.LocalJSONConnection()
	defer c.Close()

	r := gorhom.NewRate(10) // 10 hz

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
