package example

type Person struct {
	Name    string
	Address string
	Age     int
	Cnt     int
}

const HelloTopic = "hello"

type AddTwoInts struct {
	A int
	B int
}

type AddTwoIntsResponse struct {
	C int
	Comment string
}

const ServiceTopic = "service"
