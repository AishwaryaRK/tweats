package matcher

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/AishwaryaRK/tweats/datamodel"
	"github.com/WUMUXIAN/go-common-utils/cryptowrapper"
)

// Manually generate 200 tweeps, randomly distribute them amongst 8 interests.

var (
	testInterests = []string{
		"a", "b", "c", "d", "e", "f", "ANYTHING UNDER THE SUN 🌞",
	}
)

func TestMatcher(t *testing.T) {
	tweeps := make([]datamodel.Tweep, 0)
	// Generate 200 Tweeps:
	for i := 0; i < 200; i++ {
		tweep := datamodel.Tweep{
			LDAP:      cryptowrapper.GenUUID(),
			Interests: []string{},
		}
		// For each tweep, we add some random interest to it, ranging from 1 to 5.
		randomCount := rand.Intn(5) + 1
		// fmt.Println("generateing", randomCount, "interests")
		for j := 0; j < randomCount; j++ {
			for {
				index := rand.Intn(len(testInterests))
				if !containsStr(tweep.Interests, testInterests[index]) {
					tweep.Interests = append(tweep.Interests, testInterests[index])
					break
				}
			}
		}
		tweeps = append(tweeps, tweep)
	}

	fmt.Println("Total tweeps generated:", len(tweeps))

	// Add tweeps into matcher
	AddTweeps(tweeps...)

	// At this point, let's see the statistics
	fmt.Println("The mapping:")

	for interest, tweeps := range _interestMapping {
		fmt.Printf("interest %s: has %d tweeps\n", interest, len(tweeps))
	}

	fmt.Println("wildcard tweeps:", len(_wildcardTweeps))

	fmt.Println("The priority queue:")

	for _, item := range *_interestPriorityQueue {
		fmt.Printf("interest %s: has %d tweeps\n", item.interest, len(item.tweeps))
	}

	// Perform the matching.
	matches, err := GetMatches()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Total matches generated:", len(matches))
	total := 0
	for i, match := range matches {
		total += len(match.MatchedTweeps)
		fmt.Printf("match %d: interest %s, tweeps count: %d\n", i, match.MatchedInterest, len(match.MatchedTweeps))
	}

	fmt.Println("Total tweeps getting matched:", total)
}
