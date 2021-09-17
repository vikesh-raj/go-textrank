package textrank

type Options struct {
	Language            string
	AdditionalStopWords []string
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
