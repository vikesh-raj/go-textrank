package textrank

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/vikesh-raj/go-textrank/textrank/stopwords"
	"golang.org/x/text/unicode/norm"
)

var AB_ACRONYM_LETTERS = regexp.MustCompile(`([a-zA-Z])\.([a-zA-Z])\.`)
var REGEX_EMAIL = regexp.MustCompile("[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*")

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

func processText(text string) string {
	text = charRemover(isNewLine, text, false)
	text = strings.ToLower(text)
	text = charRemover(isPunc, text, true)
	text = charRemover(unicode.IsDigit, text, false)
	return text
}

func PreProcessSentence(text string, language string, additionalStopWords []string) sentence {
	s := sentence{
		OriginalText: text,
	}
	text = removeEmail(text)
	text = processText(text)
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

func isMarkNonSpacing(r rune) bool {
	return unicode.Is(unicode.Mn, r)
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

func deaccentText(text string) string {
	text = norm.NFD.String(text)
	text = charRemover(isMarkNonSpacing, text, false)
	text = norm.NFC.String(text)
	return text
}

func PreProcessWords(text, language string, deaccent bool, additionalStopWords []string) []word {
	text = removeAcronyms(text)
	text = strings.ToLower(text)
	if deaccent {
		text = deaccentText(text)
	}
	originalWords := tokenize(text)
	stopwords := stopwords.GetStopWords(language)
	stemmer := getStemmer(language)
	output := make([]word, 0, len(originalWords))

	for i, w := range originalWords {
		processed := processText(w)
		if Contains(processed, stopwords) {
			continue
		}
		if Contains(processed, additionalStopWords) {
			continue
		}
		if stemmer != nil {
			processed = stemmer(processed, false)
		}
		output = append(output, word{
			Text:  w,
			Lemma: processed,
			Index: i,
		})
	}
	return output
}

func removeAcronyms(text string) string {
	return AB_ACRONYM_LETTERS.ReplaceAllString(text, "${1}${2}")
}

func removeEmail(text string) string {
	return REGEX_EMAIL.ReplaceAllString(text, "")
}
