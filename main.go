package main

import (
	"fmt"

	"github.com/AishwaryaRK/tweats/tweatsreader"
)

func main() {
	tweeps, err := tweatsreader.Read()
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", tweeps)
}
