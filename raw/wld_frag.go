package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// FragmentReadWriter is used to read a fragment in wld format
type FragmentReadWriter interface {
	FragmentReader
	FragmentWriter
}

type FragmentReader interface {
	Read(w io.ReadSeeker) error
	FragCode() int
}

// FragmentWriter2 is used to write a fragment in wld format
type FragmentWriter interface {
	Write(w io.Writer) error
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

// NewFrag takes a reader, analyzes the first 4 bytes, and returns a new fragment struct based on it
func NewFrag(r io.ReadSeeker) FragmentReadWriter {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	fragCode := dec.Int32()
	err := dec.Error()
	if err != nil {
		return nil
	}
	//r.Seek(0, io.SeekStart)
	switch fragCode {
	case 0x00:
		return &WldFragDefault{}
	case 0x01:
		return &WldFragPaletteFile{}
	case 0x02:
		return &WldFragUserData{}
	case 0x03:
		return &WldFragTextureList{}
	case 0x04:
		return &WldFragTexture{}
	case 0x05:
		return &WldFragTextureRef{}
	case 0x06:
		return &WldFragTwoDSprite{}
	case 0x07:
		return &WldFragTwoDSpriteRef{}
	case 0x08:
		return &WldFragThreeDSprite{}
	case 0x09:
		return &WldFragThreeDSpriteRef{}
	case 0x0A:
		return &WldFragFourDSprite{}
	case 0x0B:
		return &WldFragFourDSpriteRef{}
	case 0x0C:
		return &WldFragParticleSprite{}
	case 0x0D:
		return &WldFragParticleSpriteRef{}
	case 0x0E:
		return &WldFragCompositeSprite{}
	case 0x0F:
		return &WldFragCompositeSpriteRef{}
	case 0x10:
		return &WldFragSkeletonTrack{}
	case 0x11:
		return &WldFragSkeletonTrackRef{}
	case 0x12:
		return &WldFragTrack{}
	case 0x13:
		return &WldFragTrackRef{}
	case 0x14:
		return &WldFragModel{}
	case 0x15:
		return &WldFragModelRef{}
	case 0x16:
		return &WldFragSphere{}
	case 0x17:
		return &WldFragPolyhedron{}
	case 0x18:
		return &WldFragPolyhedronRef{}
	case 0x19:
		return &WldFragSphereList{}
	case 0x1A:
		return &WldFragSphereListRef{}
	case 0x1B:
		return &WldFragLight{}
	case 0x1C:
		return &WldFragLightRef{}
	case 0x1D:
		return &WldFragPointLightOld{}
	case 0x1E:
		return &WldFragPointLightOldRef{}
	case 0x1F:
		return &WldFragSound{}
	case 0x20:
		return &WldFragSoundRef{}
	case 0x21:
		return &WldFragWorldTree{}
	case 0x22:
		return &WldFragRegion{}
	case 0x23:
		return &WldFragActiveGeoRegion{}
	case 0x24:
		return &WldFragSkyRegion{}
	case 0x25:
		return &WldFragDirectionalLightOld{}
	case 0x26:
		return &WldFragBlitSprite{}
	case 0x27:
		return &WldFragBlitSpriteRef{}
	case 0x28:
		return &WldFragPointLight{}
	case 0x29:
		return &WldFragZone{}
	case 0x2A:
		return &WldFragAmbientLight{}
	case 0x2B:
		return &WldFragDirectionalLight{}
	case 0x2C:
		return &WldFragDMSprite{}
	case 0x2D:
		return &WldFragDMSpriteRef{}
	case 0x2E:
		return &WldFragDMTrack{}
	case 0x2F:
		return &WldFragDMTrackRef{}
	case 0x30:
		return &WldFragMaterial{}
	case 0x31:
		return &WldFragMaterialList{}
	case 0x32:
		return &WldFragDMRGBTrack{}
	case 0x33:
		return &WldFragDMRGBTrackRef{}
	case 0x34:
		return &WldFragParticleCloud{}
	case 0x35:
		return &WldFragFirst{}
	case 0x36:
		return &WldFragMesh{}
	case 0x37:
		return &WldFragMeshAnimated{}
	}
	return nil
}
