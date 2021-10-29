package textrank

type Options struct {
	Debug               bool
	Language            string
	AdditionalStopWords []string
	DeAccent            bool
}

type ScoreSentence struct {
	Score float64
	Text  string
	Index int
}

type sentence struct {
	OriginalText string
	Text         string
	Index        int
	Words        stringSet
	Len          int
}

type word struct {
	Text  string
	Lemma string
	Tag   string
	Index int
}

func ScoreSentenceToText(scores []ScoreSentence) []string {
	output := make([]string, len(scores))
	for i, score := range scores {
		output[i] = score.Text
	}
	return output
}
