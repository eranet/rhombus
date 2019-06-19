package main

import (
	"fmt"
	"github.com/l1va/rhombus/gorhom"
	"github.com/l1va/rhombus/gorhom/example/reqresp"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./requester_ms a b")
		os.Exit(1)
	}
	a, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	c := gorhom.LocalJSONConnection()
	defer c.Close()

	req := reqresp.AddTwoIntsRequest{A: a, B: b}
	var res reqresp.SumResponse
	err = c.Request(reqresp.SumServiceTopic, &req, &res, time.Second)

	if err != nil {
		fmt.Println("Received an error: ", err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("%d + %d = %d\n", req.A, req.B, res.C)
	}
}
