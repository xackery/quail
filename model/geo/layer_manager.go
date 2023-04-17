package geo

import (
	"bufio"
	"fmt"
	"os"
)

// LayerManager is a layer manager
type LayerManager struct {
	layers []*Layer
}

// Layer returns a layer by name
func (e *LayerManager) Layer(name string) (*Layer, bool) {
	for _, layer := range e.layers {
		if layer.Name == name {
			return layer, true
		}
	}
	return nil, false
}

// BlenderExport writes all layers to a directory
func (e *LayerManager) BlenderExport(dir string) error {
	if len(e.layers) > 0 {
		mw, err := os.Create(fmt.Sprintf("%s/layer.txt", dir))
		if err != nil {
			return fmt.Errorf("create file %s/layer.txt: %w", dir, err)
		}
		defer mw.Close()
		layer := &Layer{}
		err = layer.WriteHeader(mw)
		if err != nil {
			return fmt.Errorf("write header: %w", err)
		}
		for _, layer = range e.layers {
			err = layer.Write(mw)
			if err != nil {
				return fmt.Errorf("write layer: %w", err)
			}
		}
	}

	return nil
}

// ReadFile reads a material file
func (e *LayerManager) ReadFile(dir string) error {
	var err error
	layerPath := fmt.Sprintf("%s/layer.txt", dir)
	err = e.layerRead(layerPath)
	if err != nil {
		return fmt.Errorf("read layer: %w", err)
	}

	return nil
}

// Inspect prints out the layer manager
func (e *LayerManager) Inspect() {
	fmt.Println(len(e.layers), "layers:")
	for _, layer := range e.layers {
		fmt.Printf("	%s\n", layer.Name)
	}
}

// Add adds a layer to the layer manager
func (e *LayerManager) Add(layer *Layer) {
	e.layers = append(e.layers, layer)
}

func (e *LayerManager) layerRead(path string) error {
	r, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		e.layers = append(e.layers, &Layer{Name: line})
	}
	return nil
}

// Layers returns all layers
func (e *LayerManager) Layers() []*Layer {
	return e.layers
}

// Count returns the number of layers
func (e *LayerManager) Count() int {
	return len(e.layers)
}
