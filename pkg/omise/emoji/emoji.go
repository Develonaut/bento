// Package emoji provides standardized emoji utilities for the Omise TUI.
package emoji

import (
	"hash/fnv"
	"math/rand"
	"time"
)

// Sushi-themed emojis for general use
var Sushi = []string{"🍣", "🍙", "🥢", "🍥"}

// Status emojis
const (
	Bento     = "🍱"
	Executing = "⏳"
	Success   = "✓"
	Failure   = "✗"
	Error     = "❌"
	Pending   = "•"
)

// GetDeterministic returns a deterministic emoji from a list based on a string key
// Uses FNV-1a hash to ensure the same key always returns the same emoji
func GetDeterministic(key string, emojis []string) string {
	if len(emojis) == 0 {
		return ""
	}
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32()
	return emojis[hash%uint32(len(emojis))]
}

// GetSushi returns a deterministic sushi emoji based on a key
func GetSushi(key string) string {
	return GetDeterministic(key, Sushi)
}

// RandomSushi returns a random sushi emoji
func RandomSushi() string {
	if len(Sushi) == 0 {
		return Bento
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Sushi[r.Intn(len(Sushi))]
}
