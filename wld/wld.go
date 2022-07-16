package wld

import "github.com/xackery/quail/common"

// WLD is a wld file struct
type WLD struct {
	path               string
	s3d                common.Archiver
	name               string
	BspRegionCount     uint32
	Hash               map[int]string
	fragments          []*fragmentInfo
	materials          []*common.Material
	vertices           []*common.Vertex
	faces              []*common.Face
	files              []common.Filer
	gltfMaterialBuffer map[string]*uint32
	gltfBoneBuffer     map[int]uint32
}

type fragmentInfo struct {
	name string
	data common.WldFragmenter
}

func New(name string, path string) (*WLD, error) {
	e := &WLD{
		name: name,
		path: path,
	}
	return e, nil
}

func NewS3D(name string, s3d common.Archiver) (*WLD, error) {
	e := &WLD{
		name: name,
		s3d:  s3d,
	}
	return e, nil
}
