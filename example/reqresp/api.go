package reqresp

import "github.com/nats-io/go-nats"

const ServerURL = nats.DefaultURL

type AddTwoIntsRequest struct {
	A int
	B int
}

type SumResponse struct {
	C       int
	Comment string
}

const SumServiceTopic = "sum_service"
