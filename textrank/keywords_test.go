package textrank

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeywords(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/mihalcea_tarau.txt")
	require.NoError(t, err)
	text := string(content)

	scores, err := Keywords(text, Options{Language: "english"}, 0.2, 0)
	require.NoError(t, err)

	keywords := ScoreSentenceToText(scores)
	expectedKeywords := []string{"hurricane gilbert", "dominican coast", "southeast winds", "storm", "puerto rico", "miles south", "heavy rains", "said residents", "mph", "sunday", "saturday", "weather", "alerted", "province", "flood", "barahona", "reported", "north"}
	assert.Equal(t, expectedKeywords, keywords)
}
