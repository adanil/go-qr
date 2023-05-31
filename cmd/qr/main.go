package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/psxzz/go-qr/pkg/qr"
)

var corrFlag = flag.String("corr", "L", "Correction level")
var corr qr.Correction

func main() {
	flag.Parse()

	switch *corrFlag {
	case "L":
		corr = qr.L
	case "M":
		corr = qr.M
	case "Q":
		corr = qr.Q
	case "H":
		corr = qr.H
	}

	encoder := qr.NewEncoder(
		qr.WithCorrectionLevel(corr),
		qr.WithVersionRange(0, 40),
	)

	code, err := encoder.Encode("https://github.com/psxzz/go-qr")

	if err != nil {
		log.Fatalf("encoder: %v", err)
	}

	fmt.Printf("code: %v\n", code)

}
