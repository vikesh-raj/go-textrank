package textrank

import (
	"github.com/kljensen/snowball/english"
	"github.com/kljensen/snowball/french"
	"github.com/kljensen/snowball/norwegian"
	"github.com/kljensen/snowball/russian"
	"github.com/kljensen/snowball/spanish"
	"github.com/kljensen/snowball/swedish"
)

func getStemmer(language string) func(string, bool) string {
	switch language {
	case "", "english":
		return english.Stem
	case "spanish":
		return spanish.Stem
	case "french":
		return french.Stem
	case "russian":
		return russian.Stem
	case "swedish":
		return swedish.Stem
	case "norwegian":
		return norwegian.Stem
	default:
		return nil
	}
}

func applyStemmer(words []string, stemmer func(string, bool) string) []string {
	output := make([]string, len(words))
	for i, word := range words {
		output[i] = stemmer(word, false)
	}
	return output
}
