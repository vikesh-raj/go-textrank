package textrank

import (
	"fmt"
	"sort"
	"strings"
)

func SummarizeSentences(sentencetexts []string, options Options) ([]ScoreSentence, error) {
	sentences := PreProcessSentences(sentencetexts, options.Language, options.AdditionalStopWords)

	if options.Debug {
		dumpSentences(sentences)
	}

	g := NewGraphForSentences(sentences)
	if len(g.edges) == 0 {
		return nil, fmt.Errorf("all the sentences have no similarity")
	}

	prunedGraph, mapOldtoNew := g.PruneGraph()
	if prunedGraph.nodeCount == 0 {
		return nil, fmt.Errorf("all the sentences have no relation")
	}

	scores, err := pageRank(prunedGraph)
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
	if length == 0 {
		length = 1
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

func PickTopSentencesByWordCount(scores []ScoreSentence, wordCount int) []string {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	currentWordCount := 0
	sentences := make([]string, 0, 10)
	ids := make([]int, 0, 10)
	for i, score := range scores {
		numWords := len(strings.Fields(score.Text))
		currentWordCount += numWords
		if i != 0 && wordCount < currentWordCount {
			break
		}
		sentences = append(sentences, score.Text)
		ids = append(ids, score.Index)
	}

	// Sort the top items by Index again.
	sort.Slice(sentences, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	return sentences
}

func PickTopSentence(scores []ScoreSentence) string {
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

func dumpSentences(sentences []sentence) {
	fmt.Println("pre processed sentences = ")
	for _, sentence := range sentences {
		fmt.Println(sentence.Index, " ==> ", sentence.Text)
	}
}
