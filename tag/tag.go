package tag

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

var (
	mu             sync.RWMutex
	lastPos        int64
	tags           []tag
	lastColorIndex int
	coder          encdec.Coder
)

type tag struct {
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Color   string `json:"color"`
	Caption string `json:"caption"`
}

// New creates a new tag manager
func New() {
	if flag.Lookup("test.v") == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	tags = []tag{}
}

func NewWithCoder(c encdec.Coder) {
	if flag.Lookup("test.v") == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	tags = []tag{}
	coder = c
	lastPos = coder.Pos()
}

func SetCoder(c encdec.Coder) {
	if flag.Lookup("test.v") == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	coder = c
	lastPos = coder.Pos()
}

// Mark requires coder to be set, but smartly marks from last position to current position
func Mark(color string, caption string) {
	if flag.Lookup("test.v") == nil {
		return
	}
	if coder == nil {
		panic("mark requires a coder to be set")
	}
	Add(lastPos, coder.Pos(), color, caption)
}

// MarkRand requires coder to be set, but smartly marks from last position to current position
func MarkRand(caption string) {
	if flag.Lookup("test.v") == nil {
		return
	}
	if coder == nil {
		panic("mark requires a coder to be set")
	}
	AddRand(lastPos, coder.Pos(), caption)
}

// LastPos returns the last position
func LastPos() int64 {
	if flag.Lookup("test.v") == nil {
		return 0
	}
	mu.Lock()
	defer mu.Unlock()
	return lastPos
}

// Close closes the tag manager
func Close() {
	if flag.Lookup("test.v") == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	tags = []tag{}
	coder = nil
}

// Add creates a new tag
func Add(from, to int64, color, caption string) {
	if flag.Lookup("test.v") == nil {
		return
	}
	if log.LogLevel() != 0 {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	tags = append(tags, tag{
		From:    from,
		To:      to - 1,
		Color:   color,
		Caption: caption,
	})
	lastPos = to
}

func Addf(from, to int64, color, format string, args ...interface{}) {
	if flag.Lookup("test.v") == nil {
		return
	}
	Add(from, to, color, fmt.Sprintf(format, args...))
}

// Add creates a new tag with random color
func AddRand(from, to int64, caption string) {
	if flag.Lookup("test.v") == nil {
		return
	}
	if log.LogLevel() != 0 {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	// make a string array of colors
	colors := []string{
		"red",
		"green",
		"yellow",
		"teal",
		"purple",
		"pink",
		"brown",
		"gray",
		"orange",
		"blue",
	}
	if from < 0 {
		from = 0
	}
	colorIndex := lastColorIndex + 1
	if colorIndex >= len(colors) {
		colorIndex = 0
	}
	color := colors[colorIndex]
	lastColorIndex = colorIndex

	tags = append(tags, tag{
		From:    from,
		To:      to - 1,
		Color:   color,
		Caption: caption,
	})
	lastPos = to
}

func AddRandf(from, to int64, format string, args ...interface{}) {
	if flag.Lookup("test.v") == nil {
		return
	}
	AddRand(from, to, fmt.Sprintf(format, args...))
}

// Write writes the tag file
func Write(path string) error {
	if flag.Lookup("test.v") == nil {
		return nil
	}
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
	enc.SetIndent("", "  ")
	err = enc.Encode(tags)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
