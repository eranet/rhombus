package main

import (

	"sync"

	"github.com/l1va/rhombus/gorhom"
)

// create wait group
var wg sync.WaitGroup

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
		println("Received a msg: %+v\n", js)
	})

	println("after")
	c.Spin()

}
