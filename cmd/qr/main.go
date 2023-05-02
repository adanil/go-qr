package main

import (
	"flag"
	"fmt"
	"log"

	qr_encode "github.com/psxzz/go-qr/pkg/qr-encode"
)

var corrFlag = flag.String("corr", "L", "Correction level")
var corr qr_encode.CodeLevel

func main() {
	flag.Parse()

	switch *corrFlag {
	case "L":
		corr = qr_encode.L
	case "M":
		corr = qr_encode.M
	case "Q":
		corr = qr_encode.Q
	case "H":
		corr = qr_encode.H
	}

	encoder := qr_encode.NewEncoder(corr)

	fmt.Println("Start...")
	grid, err := encoder.Encode2D("Hello world")
	if err != nil {
		log.Fatalf("encoder: %v", err)
	}

	fmt.Printf("grid: %v\n", grid)

}
