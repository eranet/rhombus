package main

import (
	"fmt"
	"github.com/eranet/rhombus/rhomgo"
	"time"
)

type JointState struct {
	Name     []string
	Position []float64
	Velocity []float64
}

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	println("before")
	c.Subscribe("simple_gripper/joint_states", func(js *JointState) {
		t := time.Now()
		println("Received a msg:", t.UnixNano()/1000000.0)
		fmt.Printf("%+v\n", js)

	})

	println("after")
	c.Spin()

}
