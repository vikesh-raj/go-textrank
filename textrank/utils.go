package textrank

import "math"

func getSimilarity(s1, s2 sentence) float64 {
	commonWordCount := CountCommonWords(s1.Words, s2.Words)
	logS1 := math.Log10(float64(s1.Len))
	logS2 := math.Log10(float64(s2.Len))

	if logS1+logS2 == 0 {
		return 0
	}
	return float64(commonWordCount) / (logS1 + logS2)
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
