package raw

import (
	"fmt"
	"io"
)

// FragmentReader is used to read a fragment in wld format
type FragmentReader interface {
	Encode(w io.Writer) error
	FragCode() int
}

// FragmentWriter is used to write a fragment in wld format
type FragmentWriter interface {
}

var (
	fragNames = map[int]string{
		0:  "Default",
		1:  "PaletteFile",
		2:  "UserData",
		3:  "TextureList",
		4:  "Texture",
		5:  "TextureRef",
		6:  "TwoDSprite",
		7:  "TwoDSpriteRef",
		8:  "ThreeDSprite",
		9:  "ThreeDSpriteRef",
		10: "FourDSprite",
		11: "FourDSpriteRef",
		12: "ParticleSprite",
		13: "ParticleSpriteRef",
		14: "CompositeSprite",
		15: "CompositeSpriteRef",
		16: "SkeletonTrack",
		17: "SkeletonTrackRef",
		18: "Track",
		19: "TrackRef",
		20: "Model",
		21: "ModelRef",
		22: "Sphere",
		23: "Polyhedron",
		24: "PolyhedronRef",
		25: "SphereList",
		26: "SphereListRef",
		27: "Light",
		28: "LightRef",
		29: "PointLightOld",
		30: "PointLightOldRef",
		31: "Sound",
		32: "SoundRef",
		33: "WorldTree",
		34: "Region",
		35: "ActiveGeoRegion",
		36: "SkyRegion",
		37: "DirectionalLightOld",
		38: "BlitSprite",
		39: "BlitSpriteRef",
		40: "PointLight",
		41: "Zone",
		42: "AmbientLight",
		43: "DirectionalLight",
		44: "DMSprite",
		45: "DMSpriteRef",
		46: "DMTrack",
		47: "DMTrackRef",
		48: "Material",
		49: "MaterialList",
		50: "DMRGBTrack",
		51: "DMRGBTrackRef",
		52: "ParticleCloud",
		53: "First",
		54: "Mesh",
		55: "MeshAnimated",
	}
)

// FragName returns the name of a fragment
func FragName(fragCode int) string {
	name, ok := fragNames[fragCode]
	if ok {
		return name
	}
	return fmt.Sprintf("unknownFrag%d", fragCode)
}

// FragIndex returns the index of a fragment
func FragIndex(name string) int {
	for k, v := range fragNames {
		if v == name {
			return k
		}
	}
	return -1
}
