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

func TestSummarization(t *testing.T) {

	table := []struct {
		run       bool
		fail      bool
		inputText string
		input     string
		notEmpty  bool
		expected  string
		language  string
		stopwords []string
	}{
		{
			input:    "testdata/mihalcea_tarau.txt",
			expected: "testdata/mihalcea_tarau.summ.txt",
			language: "english",
		},
		{
			input:     "testdata/mihalcea_tarau.txt",
			expected:  "testdata/mihalcea_tarau.summ.txt",
			language:  "english",
			stopwords: []string{"press", "strong", "people"},
		},
		{
			input:    "testdata/few_distinct_words.txt",
			fail:     true,
			language: "english",
		},
		{
			input:     "testdata/few_distinct_words.txt",
			fail:      true,
			language:  "english",
			stopwords: []string{"here", "there"},
		},
		{
			input:    "testdata/unrelated.txt",
			fail:     true,
			language: "english",
		},
		{
			// Empty text
			fail:     true,
			language: "english",
		},
		{
			inputText: "single sentence",
			fail:      true,
			language:  "english",
		},
		{
			input:    "testdata/spanish.txt",
			notEmpty: true,
			language: "spanish",
		},
		{
			input:    "testdata/polish.txt",
			notEmpty: true,
			language: "polish",
		},
		{
			input:    "testdata/arabic.txt",
			notEmpty: true,
			language: "arabic",
		},
	}

	runMode := false
	// runMode = true

	for _, item := range table {
		if runMode && !item.run {
			continue
		}

		var text string
		if item.inputText != "" {
			text = item.inputText
		} else if item.input != "" {
			content, err := ioutil.ReadFile(item.input)
			require.NoError(t, err)
			text = string(content)
		}

		sentences := strings.Split(text, "\n")
		scores, err := SummarizeSentences(sentences, Options{
			Language:            item.language,
			AdditionalStopWords: item.stopwords,
		})

		if item.fail {
			assert.NotNil(t, err)
		} else {
			require.NoError(t, err)

			top := PickTopSentencesByRatio(scores, 0.2)
			summary := strings.Join(top, "\n")

			if item.notEmpty {
				if !assert.NotEmpty(t, top) {
					printScores(scores)
				}
			} else {
				expected, err := ioutil.ReadFile(item.expected)
				require.NoError(t, err)

				if !assert.Equal(t, string(expected), summary) {
					printScores(scores)
				}
			}
		}
	}
}

func TestSummarizationTopSentence(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/mihalcea_tarau.txt")
	require.NoError(t, err)
	text := string(content)

	sentences := strings.Split(text, "\n")
	scores, err := SummarizeSentences(sentences, Options{Language: "english"})
	require.NoError(t, err)

	summary := PickTopSentence(scores)
	expectedSummary := "Hurricane Gilbert swept toward the Dominican Republic Sunday, and the Civil Defense alerted its heavily populated south coast to prepare for high winds, heavy rains and high seas."

	assert.Equal(t, expectedSummary, summary)
}
