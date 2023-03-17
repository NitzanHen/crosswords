package crossword

import "regexp"

type Word string

type Corpus struct {
	words []Word
	cache map[string][]Word
}

func NewCorpus(words []Word) Corpus {
	return Corpus{words, make(map[string][]Word)}
}

func (c *Corpus) Filter(regex regexp.Regexp) []Word {
	str := regex.String()

	if cached, ok := c.cache[str]; ok {
		return cached
	}

	matches := Filter(
		c.words,
		func(w Word) bool { return regex.MatchString(string(w)) },
	)

	c.cache[str] = matches

	return matches
}
