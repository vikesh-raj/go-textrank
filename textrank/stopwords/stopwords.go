package stopwords

func GetStopWords(language string) []string {
	switch language {
	case "", "english", "en":
		return english[:]
	case "danish", "da":
		return danish[:]
	case "german", "de":
		return german[:]
	case "spanish", "es":
		return spanish[:]
	case "italian", "it":
		return italian[:]
	case "polish", "pl":
		return polish[:]
	case "portuguese", "pt":
		return portuguese[:]
	case "swedish", "sv":
		return swedish[:]
	}

	return nil
}
