// What it does:
//
// This example uses the CascadeClassifier class to detect faces from url,
// and draw a rectangle around each of them, before displaying them within a Window.
//
// How to run:
//
// facedetect-from-url [image URL] [classifier XML file] [image file]
//
// 		go run ./cmd/facedetect-from-url/main.go https://raw.githubusercontent.com/hybridgroup/gocv/release/images/face.jpg data/haarcascade_frontalface_default.xml output.jpg
//
//这个不行	go run 12-facedetect-from-url.go file:///home/wzw/project/go/src/a-opencv/data/face.jpg data/haarcascade_frontalface_default.xml output.jpg
//这个行	go run 12-facedetect-from-url.go https://pics5.baidu.com/feed/30adcbef76094b36bc25a739363e3bd18f109dad.png?token=e4d0490adeac15a2312f547b061ca06f data/haarcascade_frontalface_default.xml output.jpg
//		go run 12-facedetect-from-url.go https://pics3.baidu.com/feed/0df431adcbef7609bdc516b9500de5c47dd99e27.jpeg?token=729fcc045aaed399747cd1d7e7c251ca data/haarcascade_frontalface_default.xml output.jpg
//		go run 12-facedetect-from-url.go https://pics6.baidu.com/feed/1c950a7b02087bf4f83a0c6c8e60102413dfcf99.jpeg?token=34df580b62896b91c2e5f372266635a5 data/haarcascade_frontalface_default.xml output.jpg
// +build example

package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("How to run:\n\tfacedetect-from-url [image URL] [classifier XML file] [image file]")
		return
	}

	// parse args
	imageURL := os.Args[1]
	xmlFile := os.Args[2]
	saveFile := os.Args[3]

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	// get image from URL
	res, err := http.Get(imageURL)
	if err != nil {
		log.Fatal(err)
	}

	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	img, err := gocv.IMDecode(resByte, 1)
	if err != nil {
		log.Fatal(err)
	}

	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))
	// draw a rectangle around each face on the original image,
	// along with text identifing as "Human"
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)

		size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
	}
	gocv.IMWrite(saveFile, img)
	fmt.Printf("saved to %s\n", saveFile)
}
