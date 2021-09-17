package textrank

import (
	"fmt"
	"sort"
)

func Summarize(text string, options Options) ([]ScoreSentence, error) {
	sentences := SplitSentences(text)
	return SummarizeSentences(sentences, options)
}

func SummarizeSentences(sentencetexts []string, options Options) ([]ScoreSentence, error) {
	sentences := PreProcessSentences(sentencetexts, options.Language, options.AdditionalStopWords)

	g := NewGraphForSentences(sentences)
	if len(g.edges) == 0 {
		return nil, fmt.Errorf("all the sentences have no similarity")
	}

	prunedGraph, mapOldtoNew := g.PruneGraph()
	if prunedGraph.nodeCount == 0 {
		return nil, fmt.Errorf("all the sentences have no relation")
	}

	scores, err := pageRank(prunedGraph.nodeCount, prunedGraph.edges)
	if err != nil {
		return nil, err
	}

	output := make([]ScoreSentence, len(sentences))
	for i := range sentences {
		score := 0.0
		newIndex, ok := mapOldtoNew[i]
		if ok {
			score = scores[newIndex]
		}
		output[i] = ScoreSentence{
			Text:  sentences[i].OriginalText,
			Score: score,
			Index: sentences[i].Index,
		}
	}
	return output, nil
}

func PickTopSentencesByRatio(scores []ScoreSentence, ratio float64) []string {
	length := int(float64(len(scores)) * ratio)
	if length > len(scores) {
		length = len(scores)
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	output := make([]string, length)
	for i := range output {
		output[i] = scores[i].Text
	}

	// Sort the top items by Index again.
	sort.Slice(output, func(i, j int) bool {
		return scores[i].Index < scores[j].Index
	})
	return output
}

func PickTopSentence(scores []ScoreSentence, count int) string {
	max := scores[0].Score
	maxIndex := 0
	for i, score := range scores {
		if score.Score > max {
			max = score.Score
			maxIndex = i
		}
	}
	return scores[maxIndex].Text
}
