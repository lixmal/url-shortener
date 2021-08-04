// package database provides a simple key/value store
package database

import (
	"sync"
)

type db struct {
	sync.RWMutex
	db map[string]string
}

// TODO: Using a shared global var for simplicity, add persistence later
var database = db{
	db: map[string]string{},
}

// Set sets the value for the given key
func Set(key, value string) {
	database.Lock()
	defer database.Unlock()

	database.db[key] = value
}

// Lookup retrieves the value for the given key and provides a second
// bool return value analogous to map lookups
func Lookup(key string) (string, bool) {
	database.RLock()
	defer database.RUnlock()

	v, ok := database.db[key]
	return v, ok
}
