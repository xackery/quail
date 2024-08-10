package rawfrag

import (
	"io"

	"github.com/xackery/quail/model"
)

// WldFragDmTrackDef2 is DmTrackDef2 in libeq, Mesh Animated Vertices in openzone, DMTRACKDEF in wld, MeshAnimatedVertices in lantern
type WldFragDmTrackDef2 struct {
	NameRef     int32                     `yaml:"name_ref"`
	Flags       uint32                    `yaml:"flags"`
	VertexCount uint16                    `yaml:"vertex_count"`
	FrameCount  uint16                    `yaml:"frame_count"`
	Param1      uint16                    `yaml:"param_1"` // usually contains 100
	Param2      uint16                    `yaml:"param_2"` // usually contains 0
	Scale       uint16                    `yaml:"scale"`
	Frames      []WldFragMeshAnimatedBone `yaml:"frames"`
	Size6       uint32                    `yaml:"size_6"`
}

type WldFragMeshAnimatedBone struct {
	Position model.Vector3 `yaml:"position"`
}

func (e *WldFragDmTrackDef2) FragCode() int {
	return FragCodeDmTrackDef2
}

func (e *WldFragDmTrackDef2) Write(w io.Writer, isNewWorld bool) error {
	return nil
}

func (e *WldFragDmTrackDef2) Read(r io.ReadSeeker, isNewWorld bool) error {
	return nil
}
