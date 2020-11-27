package main

import (
	"github.com/eranet/rhombus/rhomgo"
	"github.com/eranet/rhombus/rhomgo/example/webcam"
	"io/ioutil"
)

func photoSaver(data []byte) {
	ioutil.WriteFile("recieved.jpg", data, 0644)
}

func main() {
	c := rhomgo.LocalBinaryConnection()
	defer c.Close()

	c.Subscribe(webcam.WebcamTopic, photoSaver)

	c.Spin()
}
