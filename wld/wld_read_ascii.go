package wld

import (
	"fmt"
	"os"
)

func (wld *Wld) ReadAscii(path string) error {
	wld.mu.Lock()
	defer wld.mu.Unlock()

	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	asciiReader := NewAsciiReader(r)
	err = asciiReader.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, asciiReader.lineNumber, err)
	}
	return nil
}
