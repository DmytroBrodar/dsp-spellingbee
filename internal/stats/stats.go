package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// This struct will keep all stats in memory (in the server)
type Stats struct {
	TotalWords int
	ValidWords int
	Pangrams   int
	TotalScore int
	mu         sync.Mutex
}

// LoadStats loads stats from a file or creates new ones if file doesn't exist
func LoadStats(filename string) *Stats {
	data, err := os.ReadFile(filename)
	if err != nil {
		// file not found or can't read → return empty stats
		fmt.Println("Stats: no existing file, starting fresh:", filename)
		return &Stats{}
	}

	fmt.Println("Stats: loaded from file:", filename)

	parts := strings.Split(string(data), "\n")
	stats := &Stats{}

	// each line is like key=value
	for _, line := range parts {
		if strings.TrimSpace(line) == "" {
			continue
		}
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		num, _ := strconv.Atoi(value)

		switch key {
		case "TotalWords":
			stats.TotalWords = num
		case "ValidWords":
			stats.ValidWords = num
		case "Pangrams":
			stats.Pangrams = num
		case "TotalScore":
			stats.TotalScore = num
		}
	}
	return stats
}

// Save stats to file (NO locking here: Update already holds the lock)
func (s *Stats) Save(filename string) {
	// make sure directory exists if any
	dir := filepath.Dir(filename)
	if dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}

	content := fmt.Sprintf(
		"TotalWords=%d\nValidWords=%d\nPangrams=%d\nTotalScore=%d\n",
		s.TotalWords, s.ValidWords, s.Pangrams, s.TotalScore,
	)

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		fmt.Println("ERROR writing stats:", err)
	} else {
		fmt.Println("Stats: saved to file:", filename)
	}
}

// Update stats after every word
func (s *Stats) Update(valid bool, score int, pangram bool, filename string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("Stats: Updating:", "valid =", valid, "score =", score, "pangram =", pangram, "file =", filename)

	s.TotalWords++

	if valid {
		s.ValidWords++
	}
	if pangram {
		s.Pangrams++
	}
	s.TotalScore += score

	s.Save(filename)
}
