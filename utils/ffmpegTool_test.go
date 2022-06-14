package utils

import (
	"github.com/disintegration/imaging"
	"testing"
)

func TestExampleReadFrameAsJpeg(t *testing.T) {

	reader := ExampleReadFrameAsJpeg("../public/bear.mp4", 5)
	img, err := imaging.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}
	err = imaging.Save(img, "../public/bear.jpeg")
	if err != nil {
		t.Fatal(err)
	}

}
