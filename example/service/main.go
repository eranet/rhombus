package main

import (
	"fmt"
	"rog"
	"rog/example"
)

func main() {
	c := rog.Connection()
	defer c.Close()

	c.Service(example.ServiceTopic, func(r *example.AddTwoInts) *example.AddTwoIntsResponse {
		fmt.Printf("Received a request: %v\n", r)
		return &example.AddTwoIntsResponse{C: r.A + r.B, Comment: "This is a struct"}
	})

	c.Spin()
}
