package textrank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const epsilon = 0.001

func TestSimilarityScore(t *testing.T) {

	table := []struct {
		first    string
		second   string
		expected float64
	}{
		{
			first:    "hurrican gilbert head dominican coast",
			second:   "hurrican gilbert swept dominican republ sunday civil defens alert heavili popul south coast prepar high wind heavi rain high sea",
			expected: 2.0,
		},
	}

	for _, item := range table {
		first := PreProcessSentence(item.first, "english", nil)
		second := PreProcessSentence(item.second, "english", nil)
		similarity := getSimilarity(first, second)
		assert.InDelta(t, item.expected, similarity, epsilon, "first = %s, second = %s", item.first, item.second)
	}
}
