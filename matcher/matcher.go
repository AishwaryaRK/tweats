package matcher

import (
	"container/heap"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/AishwaryaRK/tweats/datamodel"
)

// Constants.
const (
	AnyInterest    = "ANYTHING"
	AnyInterestStr = "ANYTHING UNDER THE SUN 🌞"
)

// Match defines a match
type Match struct {
	MatchedTweeps   []datamodel.Tweep
	MatchedInterest string
}

// The variables within package.
var (
	_interestPriorityQueue *PriorityQueue
	_interestMapping       map[string][]datamodel.Tweep
	_wildcardTweeps        []datamodel.Tweep

	_leftOutTweeps []datamodel.Tweep
)

func init() {
	Clear()
}

// Clear clears all data in matcher
func Clear() {
	_interestMapping = make(map[string][]datamodel.Tweep)
	_interestPriorityQueue = &PriorityQueue{}
	heap.Init(_interestPriorityQueue)
	_wildcardTweeps = make([]datamodel.Tweep, 0)
	_leftOutTweeps = make([]datamodel.Tweep, 0)
}

// AddTweeps addds one or more tweep to the interest pool for matching.
func AddTweeps(tweeps ...datamodel.Tweep) (err error) {
	// firstly group the tweeps by interest.
	for _, tweep := range tweeps {
		// add tweeps to each interest
		for _, interest := range tweep.Interests {
			if strings.Contains(interest, AnyInterest) {
				_wildcardTweeps = append(_wildcardTweeps, tweep)
			} else {
				_interestMapping[interest] = append(_interestMapping[interest], tweep)
			}
		}
	}
	// secondly push the groups into a priority queue.
	for interest, tweeps := range _interestMapping {
		heap.Push(_interestPriorityQueue, &Item{
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
	for _interestPriorityQueue.Len() > 0 {
		// Let's get the interest with least tweeps inside.
		item := heap.Pop(_interestPriorityQueue).(*Item)

		// if no tweeps are in this interst group anymore.
		if len(item.tweeps) == 0 {
			continue
		}

		// if this interest only has 1 person inside.
		if len(item.tweeps) == 1 {
			// if this tweep exists in other interest also, skip this.
			if existsInOtherGroups(item.tweeps[0]) {
				continue
			}
			// otherwise we pick a random user from the widlcard interest group to and form a match
			if len(_wildcardTweeps) > 0 {
				match := Match{[]datamodel.Tweep{item.tweeps[0]}, item.interest}
				pickedTweeps, err := randomPick(&_wildcardTweeps, 1)
				if err == nil {
					match.MatchedTweeps = append(match.MatchedTweeps, pickedTweeps...)
					matches = append(matches, match)
					removeMatchedTweeps(match.MatchedTweeps)
				}
			} else {
				// try to add to leftover tweeps
				addToLeftoverTweeps(item.tweeps[0])
			}
		} else if len(item.tweeps) <= 3 {
			// if we have 2 or 3 tweeps, we will do a match for all of them.
			match := Match{item.tweeps, item.interest}
			matches = append(matches, match)
			removeMatchedTweeps(match.MatchedTweeps)
		} else {
			for count := 5; count >= 2; count-- {
				if len(item.tweeps) >= count*2 {
					fmt.Println("picking up", count, "tweeps from interst", item.interest)
					match, err := getMatchByCount(&item.tweeps, item.interest, count)
					if err == nil {
						matches = append(matches, match)
						removeMatchedTweeps(match.MatchedTweeps)
					}
					break
				}
			}
			heap.Init(_interestPriorityQueue)
			// once we are done with creating one match, we will push this back to the priority queue for next round if there are still more than 1 tweep left.
			if len(item.tweeps) > 0 {
				heap.Push(_interestPriorityQueue, item)
			}
		}

		fmt.Println("The priority queue after getting a match:")

		for _, item := range *_interestPriorityQueue {
			fmt.Printf("interest %s: has %d tweeps\n", item.interest, len(item.tweeps))
		}
	}
	matches, err = handleLeftovers(matches)
	return
}

func handleLeftovers(currentMatches []Match) (matches []Match, err error) {
	fmt.Println("current left over tweeps:", len(_leftOutTweeps))
	// if leftout is not 0, it means we don't have wildcards left.
	if len(_leftOutTweeps) > 0 {
		// for each left over tweeps, we try to insert into one of the existing groups based on intersection of interests.
		for _, match := range currentMatches {
			for i := range _leftOutTweeps {
				if containsStr(_leftOutTweeps[i].Interests, match.MatchedInterest) {
					match.MatchedTweeps = append(match.MatchedTweeps, _leftOutTweeps[i])
					break
				}
			}
		}
	}

	fmt.Println("current wild card tweeps:", len(_wildcardTweeps))
	if len(_wildcardTweeps) > 0 {
		// send out emails to these guys, and let them freestyle.
		match := Match{_wildcardTweeps, AnyInterestStr}
		matches = append(currentMatches, match)
	} else {
		matches = currentMatches
	}
	return
}

func existsInOtherGroups(tweep datamodel.Tweep) bool {
	for i := range *_interestPriorityQueue {
		if containTweeps((*_interestPriorityQueue)[i].tweeps, tweep) {
			return true
		}
	}
	if containTweeps(_wildcardTweeps, tweep) {
		return true
	}
	return false
}

func addToLeftoverTweeps(tweep datamodel.Tweep) {
	for i := range *_interestPriorityQueue {
		if containTweeps((*_interestPriorityQueue)[i].tweeps, tweep) {
			return
		}
	}
	if containTweeps(_wildcardTweeps, tweep) {
		return
	}
	if containTweeps(_leftOutTweeps, tweep) {
		return
	}
	_leftOutTweeps = append(_leftOutTweeps, tweep)
}

func removeMatchedTweeps(tweeps []datamodel.Tweep) {
	for i := range *_interestPriorityQueue {
		for _, tweep := range tweeps {
			for j := range (*_interestPriorityQueue)[i].tweeps {
				if (*_interestPriorityQueue)[i].tweeps[j].LDAP == tweep.LDAP {
					(*_interestPriorityQueue)[i].tweeps = append((*_interestPriorityQueue)[i].tweeps[:j], (*_interestPriorityQueue)[i].tweeps[j+1:]...)
					break
				}
			}
		}
	}

	for _, tweep := range tweeps {
		for i := range _wildcardTweeps {
			if _wildcardTweeps[i].LDAP == tweep.LDAP {
				_wildcardTweeps = append(_wildcardTweeps[:i], _wildcardTweeps[i+1:]...)
				break
			}
		}
	}

	for _, tweep := range tweeps {
		for i := range _leftOutTweeps {
			if _leftOutTweeps[i].LDAP == tweep.LDAP {
				_leftOutTweeps = append(_leftOutTweeps[:i], _leftOutTweeps[i+1:]...)
				break
			}
		}
	}
}

func containTweeps(tweeps []datamodel.Tweep, tweep datamodel.Tweep) bool {
	for _, tp := range tweeps {
		if tp.LDAP == tweep.LDAP {
			return true
		}
	}
	return false
}

func getMatchByCount(tweeps *[]datamodel.Tweep, interest string, count int) (match Match, err error) {
	pickedTweeps, err := randomPick(tweeps, count)
	if err == nil {
		match = Match{pickedTweeps, interest}
	}
	return
}

func containsStr(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
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
		index := rand.Intn(len((*tweeps)))
		randomTweeps = append(randomTweeps, (*tweeps)[index])
		(*tweeps) = append((*tweeps)[:index], (*tweeps)[index+1:]...)
	}
	return
}
