package textrank

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func printScores(scores []ScoreSentence) {
	for _, score := range scores {
		fmt.Printf("--- %f ::: %s\n", score.Score, score.Text)
	}
}

func TestSummarizationBasic(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/mihalcea_tarau.txt")
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/mihalcea_tarau.summ.txt")
	require.NoError(t, err)

	sentences := strings.Split(string(content), "\n")
	scores, err := SummarizeSentences(sentences, Options{
		Language:            "english",
		AdditionalStopWords: nil,
	})
	require.NoError(t, err)

	top := PickTopSentencesByRatio(scores, 0.2)
	summary := strings.Join(top, "\n")

	printScores(scores)

	assert.Equal(t, string(expected), summary)
}
