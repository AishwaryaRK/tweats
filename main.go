package main

import (
	"fmt"
	"tweats/matcher"
	"tweats/tweatsreader"
)

func main() {
	tweeps, err := tweatsreader.Read()
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	matcher.AddTweeps(tweeps...)
	matches, err := matcher.GetMatches()
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	for i, match := range matches {
		fmt.Printf("\nMatch %d: \n", i+1)
		fmt.Printf("--------- \n")
		for _, tweep := range match.MatchedTweeps {
			fmt.Printf("%v \n", tweep)
		}
		fmt.Printf("Common interests: %v \n", match.MatchedInterest)
	}

	//mailsender.Send(matches[0].MatchedTweeps)
}
