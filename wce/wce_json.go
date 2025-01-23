package wce

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReadJSON reads the json file at path
func ReadJSON(name string, path string) (*Wce, error) {

	wce := New(name)
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w, path: %s", err, path)
	}

	dec := json.NewDecoder(r)
	err = dec.Decode(wce)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return wce, nil
}

// WriteJSON writes the json file at path
func (wce *Wce) WriteJSON(path string) error {
	var err error

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer w.Close()

	enc := json.NewEncoder(w)
	err = enc.Encode(wce)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
