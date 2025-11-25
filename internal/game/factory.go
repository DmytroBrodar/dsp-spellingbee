package game

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

// factory pattern. this function create a Game from json file
// Here we pick a random pangram seed and take its unique letters
func NewGameFromFile(filename string) *Game {
	data, err := os.ReadFile(filename)
	if err != nil {
		// return empty game if file missing
		return &Game{}
	}

	var pangrams map[string]string
	if err := json.Unmarshal(data, &pangrams); err != nil {
		return &Game{}
	}
	// collect keys
	keys := make([]string, 0, len(pangrams))
	for k := range pangrams {
		keys = append(keys, k)
	}

	// pick random seed
	rand.Seed(time.Now().UnixNano())
	seed := keys[rand.Intn(len(keys))]

	// get unique letters from the seed word
	used := map[rune]bool{}
	var letters []rune
	for _, ch := range seed {
		if !used[ch] {
			used[ch] = true
			letters = append(letters, ch)
		}
	}

	// chhose random center letter
	if len(letters) == 0 {
		return &Game{}
	}
	center := letters[rand.Intn(len(letters))]

	return &Game{Letters: letters, Center: center, Score: 0}
}
