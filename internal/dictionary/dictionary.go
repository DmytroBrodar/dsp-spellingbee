package dictionary

import (
	"encoding/json"
	"os"
	"strings"
)

// Dictionary  interface
type Dictionary interface {
	IsValid(word string) bool
}

// RealDictionary load words from words_dictionary.json
type RealDictionary struct {
	words map[string]bool
}

func NewDictionary(filename string) *RealDictionary {
	data, err := os.ReadFile(filename)
	if err != nil {
		return &RealDictionary{words: map[string]bool{}}
	}

	// data format {"key" : value,...}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return &RealDictionary{words: map[string]bool{}}
	}

	ws := make(map[string]bool, len(raw))
	for k := range raw {
		ws[strings.ToLower(k)] = true
	}
	return &RealDictionary{words: ws}
}

func (rd *RealDictionary) IsValid(word string) bool {
	_, ok := rd.words[strings.ToLower(strings.TrimSpace(word))]
	return ok
}
