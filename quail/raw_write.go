package quail

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func (e *Quail) RawWrite(out raw.Writer) error {
	if e == nil {
		return fmt.Errorf("quail is nil")
	}
	switch val := out.(type) {
	case *raw.Lay:
		return e.layWrite(val)
	case *raw.Wld:
		return e.wldWrite(val)
	default:
		return fmt.Errorf("unknown type %T", val)
	}
}

func RawWrite(out raw.Writer, e *Quail) error {
	if e == nil {
		return fmt.Errorf("quail is nil")
	}
	return e.RawWrite(out)
}

func (e *Quail) layWrite(lay *raw.Lay) error {
	if e == nil {
		return fmt.Errorf("quail is nil")
	}
	if lay == nil {
		return fmt.Errorf("layer is nil")
	}

	if e.Header == nil {
		e.Header = &common.Header{}
	}
	lay.Version = uint32(e.Header.Version)

	for _, model := range e.Models {
		for _, entry := range model.Layers {
			entry := &raw.LayEntry{
				Material: entry.Material,
				Diffuse:  entry.Diffuse,
				Normal:   entry.Normal,
			}

			lay.Entries = append(lay.Entries, entry)
		}
	}

	return nil
}

func (e *Quail) wldWrite(wld *raw.Wld) error {
	if wld == nil {
		return fmt.Errorf("wld is nil")
	}
	if wld.Fragments == nil {
		wld.Fragments = []model.FragmentReadWriter{&rawfrag.WldFragDefault{}}
	}

	fragIndex := 1
	textureRefs := make(map[string]int)
	materials := make(map[string]int)

	// as per gfaydark_obj
	// every material has:
	//texturelist bminfo [x]
	//texture simplespritedef [ ]
	//textureref simplesprite
	//material materialdef

	for _, mod := range e.Models {
		mesh := &rawfrag.WldFragDmSpriteDef2{}

		materialList := &rawfrag.WldFragMaterialPalette{}
		for _, srcMat := range mod.Materials {
			matRef, ok := materials[srcMat.Name]
			if ok {
				materialList.MaterialRefs = append(materialList.MaterialRefs, uint32(matRef))
				continue
			}

			dstMat := &rawfrag.WldFragMaterialDef{}
			for _, srcProp := range srcMat.Properties {
				if !strings.Contains(srcProp.Name, "texture") {
					continue
				}
				textureRef, ok := textureRefs[srcProp.Value]
				if !ok {
					ext := filepath.Ext(srcProp.Value)
					baseName := strings.TrimSuffix(srcProp.Value, ext)

					dstTextureList := &rawfrag.WldFragBMInfo{ // aka BmInfo
						NameRef:      raw.NameAdd(baseName),
						TextureNames: []string{srcProp.Value},
					}
					wld.Fragments[fragIndex] = dstTextureList
					fragIndex++

					texture := &rawfrag.WldFragSimpleSpriteDef{ // aka SimpleSpriteDef
						NameRef:      raw.NameAdd(srcProp.Value),
						Flags:        0x00000000,
						CurrentFrame: 0,
						Sleep:        0,
						BitmapRefs:   []uint32{uint32(fragIndex - 1)},
					}

					wld.Fragments[fragIndex] = texture
					fragIndex++

					textureRefInst := &rawfrag.WldFragSimpleSprite{ // aka SimpleSprite
						NameRef:   raw.NameAdd(srcProp.Value),
						SpriteRef: uint32(fragIndex - 1),
						Flags:     0x00000000,
					}

					wld.Fragments[fragIndex] = textureRefInst
					fragIndex++

					textureRefs[srcProp.Value] = fragIndex - 1
					textureRef = fragIndex - 1
				}

				dstMat.SimpleSpriteRef = uint32(textureRef)
				dstMat.NameRef = raw.NameAdd(srcMat.Name)
				dstMat.Flags = 2
				dstMat.RenderMethod = 0x00000001
				// //0x00FFFFFF
				dstMat.RGBPen = [4]uint8{0, 0xFF, 0xFF, 0xFF}
				dstMat.Brightness = 0.0
				dstMat.ScaledAmbient = 0.75

				wld.Fragments[fragIndex] = dstMat
				fragIndex++
				materialList.MaterialRefs = append(materialList.MaterialRefs, uint32(fragIndex-1))
			}

			materialList.NameRef = raw.NameAdd(srcMat.Name)
			materialList.Flags = 0x00014003
			if mod.FileType == "ter" {
				materialList.Flags = 0x00018003
			}
			wld.Fragments[fragIndex] = materialList
			fragIndex++
			mesh.MaterialPaletteRef = uint32(fragIndex - 1)
		}
		mesh.NameRef = raw.NameAdd(mod.Header.Name)
		mesh.Flags = 0x00014003

		mesh.DMTrackRef = 0 // for anims later
		mesh.Fragment3Ref = 0
		mesh.Fragment4Ref = 0
		//mesh.Center
		//mesh.Params2

		//mesh.MaxDistance
		//mesh.Min
		//mesh.Max
		mesh.Scale = 13
		scale := float32(1 / float32(int(1)<<int(mesh.Scale)))

		for _, srcVert := range mod.Vertices {
			mesh.Vertices = append(mesh.Vertices, [3]int16{int16(srcVert.Position.X / scale), int16(srcVert.Position.Y / scale), int16(srcVert.Position.Z / scale)})
			mesh.VertexNormals = append(mesh.VertexNormals, [3]int8{int8(srcVert.Normal.X), int8(srcVert.Normal.Y), int8(srcVert.Normal.Z)})
			mesh.Colors = append(mesh.Colors, [4]uint8{srcVert.Tint.R, srcVert.Tint.G, srcVert.Tint.B, srcVert.Tint.A})
			mesh.UVs = append(mesh.UVs, [2]int16{int16(srcVert.Uv.X * 256), int16(srcVert.Uv.Y * 256)})
		}

		for _, srcTriangle := range mod.Triangles {
			entry := rawfrag.WldFragMeshFaceEntry{
				Flags: uint16(srcTriangle.Flag),
				Index: [3]uint16{uint16(srcTriangle.Index.X), uint16(srcTriangle.Index.Y), uint16(srcTriangle.Index.Z)},
			}

			mesh.Faces = append(mesh.Faces, entry)
		}

		wld.Fragments[fragIndex] = mesh
		fragIndex++

		// materialist materialpalette is used to group materials for dmspritedef2 (aka mesh)
		// dmspriteref  instance of def2 is dmsprite
		// model lastly an actordef
	}

	return nil
}
