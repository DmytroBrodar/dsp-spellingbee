package manager

import (
	"spellingbee/internal/game"
	"sync"
)

// singleton pattern. we keep a single manager with one game inside
// we can access it by manager.Get() everywhere on server
type Manager struct {
	Game *game.Game
}

// for singleton use sync.Once
var (
	instance *Manager
	once     sync.Once
)

// Get returns only one Manager, firstly it create the instance and next times it return the same instance
func Get() *Manager {
	once.Do(func() {
		instance = &Manager{}
	})
	return instance
}
