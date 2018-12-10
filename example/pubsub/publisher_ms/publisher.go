package main

import (
	"github.com/l1va/roms"
	"github.com/l1va/roms/example/pubsub"
)

func main() {
	c := roms.Connection(pubsub.ServerURL)
	defer c.Close()

	r := roms.NewRate(10) // 10 hz

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
