package dictionary

import "strings"

// Proxy patters. add caching method for getting the used before word faster.
type Proxy struct {
	Inner Dictionary
	cache map[string]bool
}

func NewProxy(inner Dictionary) *Proxy {
	return &Proxy{Inner: inner, cache: make(map[string]bool)}
}

func (proxy *Proxy) IsValid(word string) bool {
	key := strings.ToLower(strings.TrimSpace(word))
	if val, ok := proxy.cache[key]; ok {
		return val // return from cache
	}
	val := proxy.Inner.IsValid(key)
	proxy.cache[key] = val
	return val
}
