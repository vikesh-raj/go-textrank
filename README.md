# text rank

go-textrank is go port for the [textrank](https://github.com/summanlp/textrank) python library.



Usage :

```go

import "github.com/vikesh-raj/go-textrank/textrank"

func main() {
    sentences := []string{
        "sentence 1",
        "sentence 2",
    }

    // Summarization
    scores, _ := textrank.SummarizeSentences(sentences, textrank.Options{Language: "english"})
    top := textrank.PickTopSentencesByRatio(scores, 0.3)
    summary := strings.Join(top, "\n")
    fmt.Println("---- Summary ---")
    fmt.Println(summary)

    // Keywords
    text := strings.Join(sentences, " ")
    keywordScores, _ := textrank.Keywords(text, textrank.Options{Language: "english"}, 0.0, 15)
    keywords := textrank.ScoreSentenceToText(keywordScores)
    allKeywords := strings.Join(keywords, "\n")
    fmt.Println("---- Keywords ---")
    fmt.Println(allKeywords)
}
```
