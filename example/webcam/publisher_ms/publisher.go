package main

import (
	"github.com/l1va/roms"
	"github.com/l1va/roms/example/webcam"

	//"gocv.io/x/gocv"
)

func main() {
	c := roms.BinaryConnection(webcam.ServerURL)
	defer c.Close()

	//commented for successful build on travis
	// (opencv4 should be installed, but it is not 5 min,
	// you can help with it)

	/*rate := roms.NewRate(1)

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
