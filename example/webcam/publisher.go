package main

import (
	"rog"

	"gocv.io/x/gocv"
)

// WebcamTopic Topic name for webcam images
const WebcamTopic = "webcam"

const imgFileName = "tmp_image.jpg"

func main() {
	cnx := rog.BinaryConnection()
	defer cnx.Close()

	rate := rog.NewRate(1)

	webcam, _ := gocv.OpenVideoCapture(0)
	defer webcam.Close()

	for {
		img := gocv.NewMat()
		defer img.Close()

		webcam.Read(&img)
		data, _ := gocv.IMEncode(".jpg", img)
		cnx.Publish(WebcamTopic, data)

		rate.Sleep()
	}
}
