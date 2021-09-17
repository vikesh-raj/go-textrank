package textrank

// The stringSet type is a type alias of `map[string]struct{}`
type stringSet map[string]struct{}

// Adds an word to the set
func (s stringSet) add(word string) {
	s[word] = struct{}{}
}

// Returns a boolean value describing if the word exists in the set
func (s stringSet) has(word string) bool {
	_, ok := s[word]
	return ok
}

func (s stringSet) len() int {
	return len(s)
}

func (s stringSet) addAll(words []string) {
	for _, word := range words {
		s.add(word)
	}
}

func NewStringSet(words []string) stringSet {
	w := make(stringSet)
	w.addAll(words)
	return w
}

func Intersect(set1, set2 stringSet) stringSet {
	output := make(stringSet)
	for word := range set1 {
		if set2.has(word) {
			output.add(word)
		}
	}
	return output
}

func CountCommonWords(set1, set2 stringSet) int {
	return Intersect(set1, set2).len()
}
