package main

import (
	"github.com/l1va/rhombus/gorhom"
)

type Position struct {
	Value float64
}

func main() {
	c := gorhom.LocalJSONConnection()
	defer c.Close()

	println("before")
	err := c.Publish("simple_gripper/left_finger_tip/command", Position{Value: 2.3})
	if err != nil {
		println("error pub:", err)
	}

	println("after")

}
