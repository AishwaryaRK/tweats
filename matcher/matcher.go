package matcher

import (
	"container/heap"
	"errors"
	"math/rand"
	"strings"

	"github.com/AishwaryaRK/tweats/datamodel"
)

// Constants.
const (
	AnyInterest = "Anything"
)

// Match defines a match
type Match struct {
	MatchedTweeps   []datamodel.Tweep
	MatchedInterest string
}

// The variables within package.
var (
	_interestPriorityQueue PriorityQueue
	_interestMapping       map[string][]datamodel.Tweep
	_wildcardTweeps        []datamodel.Tweep

	_leftOutTweeps []datamodel.Tweep
)

func init() {
	_interestMapping = make(map[string][]datamodel.Tweep)
	_interestPriorityQueue = make(PriorityQueue, 0)
	_wildcardTweeps = make([]datamodel.Tweep, 0)
	_leftOutTweeps = make([]datamodel.Tweep, 0)
}

func contains(interests []string, interest string) bool {
	for _, in := range interests {
		if strings.Contains(in, AnyInterest) {
			return true
		}
	}
	return false
}

func randomPick(tweeps *[]datamodel.Tweep, count int) (randomTweeps []datamodel.Tweep, err error) {
	if len(*tweeps) < count {
		return nil, errors.New("not enough tweeps to pick from")
	}
	// randomly pick for count times.
	for i := 0; i < count; i++ {
		index := rand.Intn(count)
		randomTweeps = append(randomTweeps, (*tweeps)[index])
		(*tweeps) = append((*tweeps)[:index], (*tweeps)[index+1:]...)
	}
	return
}

// AddTweeps addds one or more tweep to the interest pool for matching.
func AddTweeps(tweeps ...datamodel.Tweep) (err error) {
	// firstly group the tweeps by interest.
	for _, tweep := range tweeps {
		// if this tweep contains a "wild card" intertest, we will put them into wild card group, which is special.
		if contains(tweep.Interests, AnyInterest) {
			_wildcardTweeps = append(_wildcardTweeps, tweep)
			continue
		}
		// otherwise, we can place them into individual interest group.
		for _, interest := range tweep.Interests {
			_interestMapping[interest] = append(_interestMapping[interest], tweep)
		}
	}
	// secondly push the groups into a priority queue.
	for interest, tweeps := range _interestMapping {
		heap.Push(&_interestPriorityQueue, &Item{
			interest, tweeps,
		})
	}
	return nil
}

// GetMatches calculates the matches
func GetMatches() (matches []Match, err error) {
	// Match mechanism is explained as follows:

	// Our logic is to firstly find out the interest that have least people and start from there.
	// After find a match, we will again check what is the interest that have least people.
	// Repeat this process.

	// To match people in the same interest:

	// If the group has 2 or 3 people, we will match them directly.
	// If the group has more than this number, we will check by the following rule:
	// If the group has X*2 people, where 2 <= X <= 5, we will try to arrange a group of X, randomly.
	// If the group has only 1 person, we will get another person from the wild card list and match them.

	// We try to exhaust all interests
	for len(_interestPriorityQueue) > 0 {
		// Let's get the interest with least tweeps inside.
		item := heap.Pop(&_interestPriorityQueue).(*Item)

		// if this interest only has 1 person inside.
		if len(item.tweeps) == 1 {
			// we randomly pick another user from the wildcard list if there are any.
			if len(_wildcardTweeps) > 0 {
				match := Match{[]datamodel.Tweep{item.tweeps[0]}, item.interest}
				pickedTweeps, err := randomPick(&_wildcardTweeps, 1)
				if err != nil {
					match.MatchedTweeps = append(match.MatchedTweeps, pickedTweeps...)
					matches = append(matches, match)
				}
			} else {
				// otherwise we will unfortunately leave this tweep out.
				_leftOutTweeps = append(_leftOutTweeps, item.tweeps[0])
			}
		} else if len(item.tweeps) <= 3 {
			// if we have 2 or 3 tweeps, we will do a match for all of them.
			match := Match{item.tweeps, item.interest}
			matches = append(matches, match)
		} else {
			for count := 5; count >= 2; count-- {
				if len(item.tweeps) >= count*2 {
					match, err := getMatchByCount(&item.tweeps, item.interest, count)
					if err != nil {
						matches = append(matches, match)
					}
					break
				}
			}
			// once we are done with creating one match, we will push this back to the priority queue for next round if there are still more than 1 tweep left.
			if len(item.tweeps) > 0 {
				heap.Push(&_interestPriorityQueue, item)
			}
		}
	}
	otherMatches, err := handleLeftovers()
	matches = append(matches, otherMatches...)
	return
}

func handleLeftovers() (matches []Match, err error) {
	// if leftout is not 0, it means we don't have wildcards left.
	if len(_leftOutTweeps) > 0 {
		// send out emails to these guys, and let them freestyle.
		match := Match{_leftOutTweeps, "Anything"}
		matches = append(matches, match)
	} else if len(_wildcardTweeps) > 0 {
		// send out emails to these guys, and let them freestyle.
		match := Match{_wildcardTweeps, "Anything"}
		matches = append(matches, match)
	}
	return
}

func getMatchByCount(tweeps *[]datamodel.Tweep, interest string, count int) (match Match, err error) {
	pickedTweeps, err := randomPick(tweeps, count)
	if err != nil {
		match = Match{pickedTweeps, interest}
	}
	return
}
