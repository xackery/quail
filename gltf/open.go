package gltf

import "github.com/qmuntal/gltf"

func Open(path string) (*gltf.Document, error) {
	return gltf.Open(path)
}
