package main

import (
	"fmt"
	"tweats/tweats_reader"
)

func main() {
	tweeps, err := tweats_reader.Read()
	if err != nil {
		fmt.Errorf("err: %v", err)
		return
	}

	fmt.Printf("%v", tweeps)
}
