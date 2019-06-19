package main

import (
	"github.com/l1va/rhombus/gorhom"
	//"github.com/l1va/rhombus/gorhom/example/webcam"

	//"gocv.io/x/gocv"
)

func main() {
	c := gorhom.LocalBinaryConnection()
	defer c.Close()

	//commented for successful build on travis
	// (opencv4 should be installed,
	// you can help with commands for travis)

	/*rate := rhombus.NewRate(1)

	cam, _ := gocv.OpenVideoCapture(0)
	defer cam.Close()

	var img gocv.Mat
	defer img.Close()

	for {
		cam.Read(&img)
		data, _ := gocv.IMEncode(".jpg", img)
		c.Publish(webcam.WebcamTopic, data)

		rate.Sleep()
	}*/
}
