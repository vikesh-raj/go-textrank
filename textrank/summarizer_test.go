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
		ratio     float64
	}{
		{
			input:    "testdata/mihalcea_tarau.txt",
			expected: "testdata/mihalcea_tarau.summ.txt",
			language: "english",
		},
		{
			input:    "testdata/mihalcea_tarau.txt",
			expected: "testdata/mihalcea_tarau.top.txt",
			ratio:    0.02,
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
			ratio := item.ratio
			if ratio == 0.0 {
				ratio = 0.2
			}
			top := PickTopSentencesByRatio(scores, ratio)
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
	scores, err := SummarizeSentences(sentences, Options{Debug: true, Language: "english"})
	require.NoError(t, err)

	summary := PickTopSentence(scores)
	expectedSummary := "Hurricane Gilbert swept toward the Dominican Republic Sunday, and the Civil Defense alerted its heavily populated south coast to prepare for high winds, heavy rains and high seas."

	assert.Equal(t, expectedSummary, summary)
}

func TestSummarizationWordCount(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/mihalcea_tarau.txt")
	require.NoError(t, err)
	text := string(content)

	sentences := strings.Split(text, "\n")
	scores, err := SummarizeSentences(sentences, Options{Language: "english"})
	printScores(scores)
	require.NoError(t, err)

	summarySentences := PickTopSentencesByWordCount(scores, 10)
	summary := strings.Join(summarySentences, " ")
	expectedSummary := "Hurricane Gilbert swept toward the Dominican Republic Sunday, and the Civil Defense alerted its heavily populated south coast to prepare for high winds, heavy rains and high seas."

	assert.Equal(t, expectedSummary, summary)

	summarySentences = PickTopSentencesByWordCount(scores, 70)
	summary = strings.Join(summarySentences, " ")
	expectedSummary = "Hurricane Gilbert swept toward the Dominican Republic Sunday, and the Civil Defense alerted its heavily populated south coast to prepare for high winds, heavy rains and high seas. The National Hurricane Center in Miami reported its position at 2 a.m. Sunday at latitude 16.1 north, longitude 67.5 west, about 140 miles south of Ponce, Puerto Rico, and 200 miles southeast of Santo Domingo."
	assert.Equal(t, expectedSummary, summary)

	summarySentences = PickTopSentencesByWordCount(scores, 120)
	summary = strings.Join(summarySentences, " ")
	expectedSummary = "Hurricane Gilbert swept toward the Dominican Republic Sunday, and the Civil Defense alerted its heavily populated south coast to prepare for high winds, heavy rains and high seas. The National Hurricane Center in Miami reported its position at 2 a.m. Sunday at latitude 16.1 north, longitude 67.5 west, about 140 miles south of Ponce, Puerto Rico, and 200 miles southeast of Santo Domingo. The National Weather Service in San Juan, Puerto Rico, said Gilbert was moving westward at 15 mph with a ``broad area of cloudiness and heavy weather'' rotating around the center of the storm. Strong winds associated with the Gilbert brought coastal flooding, strong southeast winds and up to 12 feet feet to Puerto Rico's south coast."
	assert.Equal(t, expectedSummary, summary)
}
