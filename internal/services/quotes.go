package services

import "math/rand"

var quotes = []string{
	"Become the undisputed Lord of the Cut - take away the title from CJ",
	"Be in amazing shape when getting married",
	"Inspire other people at the gym",
	"Get better than your previous all time best",
	"Be light enough to do iron cross again",
	"Be able to do muscle ups every day!",
	"Be truly obsessed",
}

func getRandomQuote() string {
	randInt := rand.Intn(len(quotes))
	return quotes[randInt]
}
