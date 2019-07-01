package main

import (
	"fmt"
	"github.com/l1va/rhombus/gorhom"
	"time"
)

type JointState struct {
	Name     []string
	Position []float64
	Velocity []float64
}

func main() {
	c := gorhom.LocalJSONConnection()
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
