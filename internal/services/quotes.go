package services

import "math/rand"

var quotes = []string{
	"Become the undisputed Lord of the Cut - take away the title from CJ",
}

func getRandomQuote() string {
	randInt := rand.Intn(len(quotes))
	return quotes[randInt]
}
