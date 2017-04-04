package main

import (
	"fmt"
	"log"

	"github.com/DeedleFake/signage"
)

func main() {
	entries, err := signage.GetSigned()
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}

	for _, e := range entries {
		fmt.Printf("%#v\n", e)
	}
}
