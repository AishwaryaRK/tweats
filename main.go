package main

import (
	"fmt"

	"github.com/AishwaryaRK/tweats/tweats_reader"
)

func main() {
	tweeps, err := tweats_reader.Read()
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", tweeps)
}
