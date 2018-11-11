package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"rog"
	"rog/example"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./requester a b")
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

	c := rog.Connection()
	defer c.Close()

	req := example.AddTwoInts{A: a, B: b}
	var res example.AddTwoIntsResponse
	err = c.RequestService(example.ServiceTopic, &req, &res, time.Second)

	if err != nil {
		fmt.Println("Received an error: ", err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("%d + %d = %d\n", req.A, req.B, res.C)
	}
}
