package main

import (
	"fmt"

	"github.com/psxzz/go-qr/pkg/qr"
)

func main() {
	encoder := qr.NewEncoder(qr.WithCorrectionLevel(qr.M))

	code, err := encoder.Encode("https://github.com/psxzz/go-qr")

	if err != nil {
		panic(err)
	}
	fmt.Printf("code: %v\n", code)
}
