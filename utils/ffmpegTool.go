package utils

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}
func ReadFrameAsJpeg(videoPath string, outPath string) error {
	reader := ExampleReadFrameAsJpeg(videoPath, 5)
	img, err := imaging.Decode(reader)
	if err != nil {
		PrintLogError("imaging读取失败;videoPath=", videoPath, " outPath=", outPath, " ", err)
		return err
	}
	err = imaging.Save(img, "outPath")
	if err != nil {
		PrintLogError("imaging保存失败;videoPath=", videoPath, " outPath=", outPath, " ", err)
		return err
	}
	return nil
}
