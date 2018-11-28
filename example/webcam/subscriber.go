package main

import (
	"io/ioutil"
	"rog"
)

func main() {
	cnx := rog.BinaryConnection()
	defer cnx.Close()

	cnx.Subscribe("webcam", func(data []byte) {
		ioutil.WriteFile("recieved.jpg", data, 0644)
	})

	cnx.Spin()
}
