package tag

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/xackery/quail/log"
)

var (
	mu      sync.RWMutex
	lastPos int
	tags    []tag
)

type tag struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Color   string `json:"color"`
	Caption string `json:"caption"`
}

// New creates a new tag manager
func New() {
	mu.Lock()
	defer mu.Unlock()
	tags = []tag{}
}

// LastPos returns the last position
func LastPos() int {
	mu.Lock()
	defer mu.Unlock()
	return lastPos
}

// Close closes the tag manager
func Close() {
	mu.Lock()
	defer mu.Unlock()
	tags = []tag{}
}

// Add creates a new tag
func Add(from, to int, color, caption string) {
	if log.LogLevel() != 0 {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	tags = append(tags, tag{
		From:    from,
		To:      to,
		Color:   color,
		Caption: caption,
	})
	lastPos = to
}

// Write writes the tag file
func Write(path string) error {
	mu.Lock()
	defer mu.Unlock()

	if len(tags) == 0 {
		return nil
	}

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	enc := json.NewEncoder(w)
	err = enc.Encode(tags)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
