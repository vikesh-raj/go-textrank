package textrank

import (
	"strings"
	"unicode"

	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/data"
	"github.com/vikesh-raj/go-textrank/textrank/stopwords"
)

var sentenceTokenizer *sentences.DefaultSentenceTokenizer

func init() {
	b, _ := data.Asset("data/english.json")
	training, _ := sentences.LoadTraining(b)
	sentenceTokenizer = sentences.NewSentenceTokenizer(training)
}

// SplitSentences splits a text into sentences
func SplitSentences(text string) []string {
	sentences := sentenceTokenizer.Tokenize(text)
	ret := make([]string, len(sentences))
	for i, sentence := range sentences {
		ret[i] = charRemover(isNewLine, strings.TrimSpace(sentence.Text), false)
	}
	return ret
}

func PreProcessSentences(sentences []string, language string, additionalStopWords []string) []sentence {
	output := make([]sentence, 0, len(sentences))
	for _, sentence := range sentences {
		s := PreProcessSentence(sentence, language, additionalStopWords)
		if s.Text != "" {
			s.Index = len(output)
			output = append(output, s)
		}
	}
	return output
}

func PreProcessSentence(text string, language string, additionalStopWords []string) sentence {
	s := sentence{
		OriginalText: text,
	}
	text = charRemover(isNewLine, text, false)
	text = strings.ToLower(text)
	text = charRemover(isPunc, text, true)
	text = charRemover(unicode.IsDigit, text, false)

	stopwords := stopwords.GetStopWords(language)
	tokens := strings.Fields(text)
	if len(stopwords) != 0 {
		tokens = removeStopwords(tokens, stopwords)
	}
	if len(additionalStopWords) != 0 {
		tokens = removeStopwords(tokens, additionalStopWords)
	}
	stemmer := getStemmer(language)
	if stemmer != nil {
		tokens = applyStemmer(tokens, stemmer)
	}
	text = strings.Join(tokens, " ")
	s.Words = NewStringSet(tokens)
	s.Text = text
	s.Len = len(tokens)
	return s
}

func charRemover(predicate func(rune) bool, text string, space bool) string {
	return strings.Map(func(r rune) rune {
		if predicate(r) {
			if space {
				return ' '
			}
			return -1
		}
		return r
	}, text)
}

func isNewLine(r rune) bool {
	return r == '\r' || r == '\n' || r == '\t'
}

func isPunc(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

func removeStopwords(tokens, stopwords []string) []string {
	output := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !Contains(token, stopwords) {
			output = append(output, token)
		}
	}
	return output
}

func tokenize(text string) []string {
	var b strings.Builder
	tokens := make([]string, 0)

	for _, ch := range text {
		if unicode.IsSpace(ch) {
			if b.Len() != 0 {
				tokens = append(tokens, b.String())
			}
			b.Reset()
		} else if unicode.IsLetter(ch) {
			b.WriteRune(ch)
		}
	}
	if b.Len() != 0 {
		tokens = append(tokens, b.String())
	}

	return tokens
}

// Contains checks if a string 's' is present in array of strings 'a'
func Contains(s string, a []string) bool {
	for _, item := range a {
		if s == item {
			return true
		}
	}
	return false
}
