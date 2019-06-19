package reqresp

type AddTwoIntsRequest struct {
	A int
	B int
}

type SumResponse struct {
	C       int
	Comment string
}

const SumServiceTopic = "sum_service"
