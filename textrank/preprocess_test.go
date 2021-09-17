package textrank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	table := []struct {
		text   string
		tokens []string
	}{
		{"This is b pin!! %% $$11 aa farà sjätte", []string{"This", "is", "b", "pin", "aa", "farà", "sjätte"}},
	}

	for _, item := range table {
		ans := tokenize(item.text)
		assert.Equal(t, item.tokens, ans, "input = %s", item.text)
	}

}

func TestPreProcessSentence(t *testing.T) {

	table := []struct {
		text     string
		expected string
	}{
		{
			text:     "Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.!!",
			expected: "lorem ipsum industri standard dummi text unknown printer took galley scrambl make specimen book",
		},
		{
			text:     "``There is no need for alarm,'' Civil Defense Director Eugenio Cabral said in a television alert shortly before midnight Saturday.",
			expected: "need alarm civil defens director eugenio cabral said televis alert short midnight saturday",
		},
	}

	for _, item := range table {
		s := PreProcessSentence(item.text, "english", []string{"type"})
		assert.Equal(t, item.expected, s.Text, "input = %s", item.text)
	}
}
