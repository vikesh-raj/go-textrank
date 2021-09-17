package textrank

import (
	"fmt"
	"sort"
	"strings"
)

func Keywords(text string, options Options, ratio float64, numOutputWords int) ([]ScoreSentence, error) {
	words := PreProcessWords(text, options.Language, options.DeAccent, options.AdditionalStopWords)
	lemmaMap, lemmas := makeLemmas(words)

	graph := createGraph(words, lemmaMap)
	if graph.nodeCount == 0 {
		return nil, fmt.Errorf("no useful words present in the text")
	}

	scores, err := pageRank(graph)
	if err != nil {
		return nil, err
	}
	sentenceScores := makeScores(scores, lemmas)
	topScores := sentenceScores[:calcuateNumWords(len(lemmas), ratio, numOutputWords)]
	return combineKeywords(topScores, text), nil
}

func createGraph(words []word, lemmaMap map[string]int) graph {
	edges := make([]edge, 0, len(words))
	for i := 0; i < len(words)-1; i++ {
		first := words[i].Lemma
		second := words[i+1].Lemma
		edges = append(edges, edge{lemmaMap[first], lemmaMap[second], 1.0})
	}

	return graph{
		nodeCount: len(lemmaMap),
		edges:     edges,
	}
}

func calcuateNumWords(total int, ratio float64, wordCount int) int {
	numWords := 10
	if ratio != 0.0 {
		numWords = int(ratio * float64(total))
	} else if wordCount != 0 {
		numWords = wordCount
	}
	if numWords > total {
		numWords = total
	}
	return numWords
}

func makeScores(scores []float64, lemmas []word) []ScoreSentence {
	sscores := make([]ScoreSentence, len(scores))
	for i := range sscores {
		sscores[i] = ScoreSentence{Score: scores[i], Text: lemmas[i].Text, Index: i}
	}

	sort.Slice(sscores, func(i, j int) bool {
		return sscores[i].Score > sscores[j].Score
	})

	return sscores
}

func makeLemmas(words []word) (map[string]int, []word) {
	lemmaMap := make(map[string]int, len(words))
	outputWords := make([]word, 0, len(words))
	for _, w := range words {
		_, ok := lemmaMap[w.Lemma]
		if !ok {
			lemmaMap[w.Lemma] = len(outputWords)
			outputWords = append(outputWords, w)
		}
	}
	return lemmaMap, outputWords
}

func combineKeywords(topScores []ScoreSentence, text string) []ScoreSentence {
	topScoreMap := make(map[string]ScoreSentence, len(topScores))
	for _, score := range topScores {
		topScoreMap[score.Text] = score
	}

	var combinedKeywords []ScoreSentence
	allwords := strings.Fields(text)
	for i := 0; i < len(allwords)-1; i++ {
		// If successive words are in top keywords, add them to combined list
		word := strings.ToLower(allwords[i])
		nextWord := strings.ToLower(allwords[i+1])
		ts1, ok1 := topScoreMap[word]
		ts2, ok2 := topScoreMap[nextWord]
		if ok1 && ok2 {
			combined := fmt.Sprintf("%s %s", ts1.Text, ts2.Text)
			score := ts1.Score + ts2.Score
			combinedKeywords = append(combinedKeywords, ScoreSentence{Text: combined, Score: score})
			delete(topScoreMap, ts1.Text)
			delete(topScoreMap, ts2.Text)
		}
	}
	for _, score := range topScoreMap {
		combinedKeywords = append(combinedKeywords, score)
	}

	sort.Slice(combinedKeywords, func(i, j int) bool {
		return combinedKeywords[i].Score > combinedKeywords[j].Score
	})

	// fmt.Println(combinedKeywords)
	return combinedKeywords
}
