package main

import (
	"fmt"
	"log"

	qr_encode "github.com/psxzz/go-qr/pkg/qr-encode"
)

func main() {
	encoder := qr_encode.NewEncoder(qr_encode.L)

	fmt.Println("Start...")
	grid, err := encoder.Encode2D("Hello world!")
	if err != nil {
		log.Fatalf("encoder: %v", err)
	}

	fmt.Printf("grid: %v\n", grid)

}
