package main

import (
	"io/ioutil"
	"github.com/l1va/roms"
	"github.com/l1va/roms/example/webcam"
)

func photoSaver(data []byte) {
	ioutil.WriteFile("recieved.jpg", data, 0644)
}

func main() {
	c := roms.BinaryConnection(webcam.ServerURL)
	defer c.Close()

	c.Subscribe(webcam.WebcamTopic, photoSaver)

	c.Spin()
}
