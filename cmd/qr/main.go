package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/psxzz/go-qr/pkg/qr"
)

const (
	imageSize = 500
)

func main() {
	encoder := qr.NewEncoder(qr.WithCorrectionLevel(qr.H))

	code, err := encoder.Encode("https://github.com/psxzz/go-qr")

	if err != nil {
		panic(err)
	}
	fmt.Printf("code: %v\n", code)

	white, pink := color.RGBA{R: 255, G: 255, B: 255, A: 0xff}, color.RGBA{R: 227, G: 61, B: 148, A: 0xff} //nolint:gomnd
	img, err := code.GetImageWithColors(imageSize, white, pink)
	if err != nil {
		log.Fatalf("couldn't create image: %v\n", err)
	}

	f, _ := os.Create("qr.png")
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
