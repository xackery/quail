package ascii

import "io"

type Wad struct {
	bitmaps   []*BMInfo       // essentially materials
	sprites   []*Sprite       // essentially meshes
	dmsprites []*DMSpriteInfo // essentially meshes
}

// BMInfo is basically BITMAPINFO
type BMInfo struct {
	Tag      string   // Tag is also known as Name
	Textures []string //
}

// Sprite is basically 2DSPRITEDEFINITION
type Sprite struct {
	Tag          string
	Flags        uint32
	CurrentFrame int32
	Sleep        uint32
	Frames       []*BMInfo
	Instances    []*SpriteInstance
}

// SpriteInstance is basically 2DSPRITEINSTANCE
type SpriteInstance struct {
	Tag       string
	Flags     uint32
	Materials []*SpriteInstanceMaterial
}

type SpriteInstanceMaterial struct {
	Tag           string
	Flags         uint32
	RenderMethod  uint32
	RGBPen        uint32
	Brightness    float32
	ScaledAmbient float32
	Pairs         [2]uint32
}

type DMSpriteInfo struct {
	Tag             string
	Flags           uint32
	MaterialPalette *DMSpriteInfoMaterialPalette
	BitmapInfoTag   string
	AnimationTag    string
	Center          [3]float32
	Params2         [3]uint32
	MaxDistance     float32
	Min             [3]float32
	Max             [3]float32
	RawScale        uint16
	MeshopCount     uint16
	Scale           float32
	Vertices        [][3]int16
	UVs             [][2]int16
	Normals         [][3]int8
	Colors          [4]int
	Triangles       []*Triangle
	//VertexPieces    []WldFragMeshVertexPiece
	//VertexMaterials []WldFragMeshVertexPiece
	//MeshOps         []WldFragMeshOpEntry
	Tracks []*DMSpriteInfoTrack
}

type Triangle struct {
	Flag     int
	Index    [3]uint16
	Material string
}

type DMSpriteInfoMaterialPalette struct {
	Tag          string
	Flags        int
	Size1        int
	MaterialTags []string
}

type DMSpriteInfoTrack struct {
	Tag    string
	Flags  uint32
	Scale  uint16
	Frames [][3]int16
	Size   uint16
}

type WldFragmentReadWriter interface {
	WldFragmentReader
	WldFragmentWriter
}

type WldFragmentReader interface {
	Read(w io.ReadSeeker) error
	FragCode() int
}

type WldFragmentWriter interface {
	Write(w io.Writer) error
}
