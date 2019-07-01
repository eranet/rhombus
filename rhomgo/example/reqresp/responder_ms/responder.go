package main

import (
	"fmt"
	"github.com/l1va/rhombus/rhomgo"
	"github.com/l1va/rhombus/rhomgo/example/reqresp"
)

func sumHandler(r *reqresp.AddTwoIntsRequest) *reqresp.SumResponse {
	fmt.Printf("Received a request: %v\n", r)
	return &reqresp.SumResponse{C: r.A + r.B, Comment: "This is an answer"}
}

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	c.Service(reqresp.SumServiceTopic, sumHandler)

	c.Spin()
}
