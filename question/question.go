package question

import (
	"math/rand"
	"time"
)

var QUESTION_MAP = map[string][]string{
	"Sports": []string{
		"Who is your favourite sports person?",
		"What is your favorite sport?",
		"What outdoor activities do you most enjoy? How often do you get time to participate?",
	},
	"Music and Entertainment": []string{
		"What is your favourite music/band?",
		"What is your favourite TV show/movie?",
		"Have you watched any cool shows recently?",
		"If you could have a meal with anyone/shadow someone's life for a day in music and entertainment (dead or alive), who would it be?",
	},
	"Current Affairs": []string{
		"What is the weirdest news story you’ve read recently?",
		"Are you interested in politics?",
		"Do you follow the stock market?",
		"What is your favourite podcast on Current Affairs?",
	},
	"High-Tech": []string{
		"What do you think about Elon Musk?",
		"Any new or upcoming game I should know of?",
		"Whats the best implementation of VR/AR you’ve seen?",
		"What feature do you think we can introduce to Twitter to make us interesting again?",
		"What is your favourite tech gadget you cannot live without?",
	},
	"Travel & Lifestyle": []string{
		"What is the top travel destination on your bucket list?",
		"If you could live in a different place anywhere in the world for a year, where would it be?",
		"What is your favourite travel destination in the world?",
		"Where should I go to change my currency?",
		"What’s your favorite restaurants in Singapore?",
		"What is your favourite cuisine to eat?",
	},
	"Nature Lovers": []string{
		"What’s your favourite outdoor activity?",
		"Which nature park do you think is the best in Singapore?",
		"Top three trees?",
		"Do you like Floral more or Fauna more? Why?",
		"Do you like Cat more or dog more? Why?",
	},
}

func GenRandomQuestion(category string) (question string) {
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
	questions := QUESTION_MAP[category]
	question = questions[r1.Intn(len(questions))]
	return
}