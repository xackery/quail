package wld

import (
	"fmt"
	"path/filepath"
)

// ReadAscii reads the ascii file at path
func (wld *Wld) ReadAscii(path string) error {
	wld.mu.Lock()
	defer wld.mu.Unlock()

	asciiReader, err := LoadAsciiFile(path, wld)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = asciiReader.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, asciiReader.lineNumber, err)
	}
	fmt.Println(asciiReader.TotalLineCountRead(), "total lines parsed for", filepath.Base(path))
	return nil
}
