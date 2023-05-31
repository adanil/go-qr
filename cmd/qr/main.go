package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/psxzz/go-qr/pkg/qr"
)

const defaultPixelSize = 40

func main() {
	encoder := qr.NewEncoder(qr.WithCorrectionLevel(qr.M))

	code, err := encoder.Encode("Hello, Alina! You are so cute.")

	if err != nil {
		panic(err)
	}
	fmt.Printf("code: %v\n", code)

	white, pink := color.RGBA{R: 255, G: 255, B: 255, A: 0xff}, color.RGBA{R: 227, G: 61, B: 148, A: 0xff} //nolint:gomnd
	img := code.GetImageWithColors(defaultPixelSize, white, pink)

	f, _ := os.Create("qr.png")
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
