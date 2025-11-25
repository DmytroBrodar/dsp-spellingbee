package game

import "strings"

// game rules (length, center letter, allowed letters, scoring)

// game keeps the 7 letters, center letter, and the running score
type Game struct {
	Letters []rune // all allowed letters
	Center  rune   // must use letters
	Score   int    // total score
}

// check if a word is valid (at least 4 letters, include center letter, used allowed letters)
func (g *Game) IsValidWord(word string) (bool, string) {
	w := strings.ToLower(strings.TrimSpace(word))

	if len(w) < 4 {
		return false, "Too short (minimum 4 letters)"
	}
	if !strings.ContainsRune(w, g.Center) {
		return false, "Word must include the center letter"
	}
	for _, ch := range w {
		if !strings.ContainsRune(string(g.Letters), ch) {
			return false, "You used a letter not from the set"
		}
	}
	return true, "ok"
}

// calculate scoring (4 letters = 1, more than 4 = word length, if guessed all 7  = +7)
func (g *Game) ScoreWord(word string) int {
	w := strings.ToLower(strings.TrimSpace(word))

	points := 1
	if len(w) > 4 {
		points = len(w)
	}
	if isPangram(w, g.Letters) {
		points += 7
	}
	g.Score += points
	return points

}

func isPangram(w string, letters []rune) bool {
	for _, l := range letters {
		if !strings.ContainsRune(w, l) {
			return false
		}
	}
	return true
}
