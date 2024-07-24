package wld

import (
	"fmt"
)

// ReadAscii reads the ascii file at path
func (wld *Wld) ReadAscii(path string) error {
	wld.mu.Lock()
	defer wld.mu.Unlock()

	asciiReader, err := LoadAsciiFile(path, wld)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	err = asciiReader.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, asciiReader.lineNumber, err)
	}
	return nil
}
