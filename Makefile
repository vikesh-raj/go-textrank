.PHONY: test bench clean generate-stopwords examples

stopwords: ./cmd/stopwords/main.go
	go build ./cmd/stopwords

generate-stopwords: stopwords
	./stopwords textrank/stopwords/stopwords_ar.txt textrank/stopwords/stopwords_ar_gen.go stopwords arabic
	./stopwords textrank/stopwords/stopwords_da.txt textrank/stopwords/stopwords_da_gen.go stopwords danish
	./stopwords textrank/stopwords/stopwords_de.txt textrank/stopwords/stopwords_de_gen.go stopwords german
	./stopwords textrank/stopwords/stopwords_en.txt textrank/stopwords/stopwords_en_gen.go stopwords english
	./stopwords textrank/stopwords/stopwords_es.txt textrank/stopwords/stopwords_es_gen.go stopwords spanish
	./stopwords textrank/stopwords/stopwords_it.txt textrank/stopwords/stopwords_it_gen.go stopwords italian
	./stopwords textrank/stopwords/stopwords_pl.txt textrank/stopwords/stopwords_pl_gen.go stopwords polish
	./stopwords textrank/stopwords/stopwords_pt.txt textrank/stopwords/stopwords_pt_gen.go stopwords portuguese
	./stopwords textrank/stopwords/stopwords_sv.txt textrank/stopwords/stopwords_sv_gen.go stopwords swedish

test:
	go test -cover -coverprofile=c.out ./textrank && go tool cover -html=c.out -o coverage.html

bench:
	go test -benchmem ./textrank -bench Benchmark.*

clean:
	rm -f *.out coverage.html stopwords

examples:
	(cd examples && go run main.go)