package main

import (
	"fmt"
	"github.com/l1va/rhombus/gorhom"
	"github.com/l1va/rhombus/gorhom/example/reqresp"
)

func sumHandler(r *reqresp.AddTwoIntsRequest) *reqresp.SumResponse {
	fmt.Printf("Received a request: %v\n", r)
	return &reqresp.SumResponse{C: r.A + r.B, Comment: "This is an answer"}
}

func main() {
	c := gorhom.LocalJSONConnection()
	defer c.Close()

	c.Service(reqresp.SumServiceTopic, sumHandler)

	c.Spin()
}
