package decorate

import (
	"fmt"
	"spellingbee/internal/dictionary"
	"strings"
)

// decorator pattern wraps dictionary and prints logs
type Logging struct {
	Inner dictionary.Dictionary
}

func (d Logging) IsValid(word string) bool {
	w := strings.ToLower(strings.TrimSpace(word))
	ok := d.Inner.IsValid(w)
	fmt.Printf("[LOGGING] dictionary check: '%s' -> %v\n", w, ok)
	return ok
}

// helper
func WrapLogging(inner dictionary.Dictionary) dictionary.Dictionary {
	return Logging{Inner: inner}
}
