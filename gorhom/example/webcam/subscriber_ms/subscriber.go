package main

import (
	"github.com/l1va/rhombus/gorhom"
	"github.com/l1va/rhombus/gorhom/example/webcam"
	"io/ioutil"
)

func photoSaver(data []byte) {
	ioutil.WriteFile("recieved.jpg", data, 0644)
}

func main() {
	c := gorhom.LocalBinaryConnection()
	defer c.Close()

	c.Subscribe(webcam.WebcamTopic, photoSaver)

	c.Spin()
}
