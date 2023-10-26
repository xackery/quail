package quail

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
)

func (e *Quail) WldMarshal(w io.Writer) error {

	world := common.NewWld(e.Header.Name)

	fragIndex := 1
	textureRefs := make(map[string]int)
	materials := make(map[string]int)

	// as per gfaydark_obj
	// every material has:
	//texturelist bminfo [x]
	//texture simplespritedef [ ]
	//textureref simplesprite
	//material materialdef

	for _, model := range e.Models {
		mesh := &wld.Mesh{}

		materialList := &wld.MaterialList{}
		for _, srcMat := range model.Materials {
			matRef, ok := materials[srcMat.Name]
			if ok {
				materialList.MaterialRefs = append(materialList.MaterialRefs, uint32(matRef))
				continue
			}

			dstMat := &wld.Material{}
			for _, srcProp := range srcMat.Properties {
				if !strings.Contains(srcProp.Name, "texture") {
					continue
				}
				textureRef, ok := textureRefs[srcProp.Value]
				if !ok {
					ext := filepath.Ext(srcProp.Value)
					baseName := strings.TrimSuffix(srcProp.Value, ext)

					dstTextureList := &wld.TextureList{ // aka BmInfo
						NameRef:      world.NameAdd(baseName),
						TextureNames: []string{srcProp.Value},
					}
					world.Fragments[fragIndex] = dstTextureList
					fragIndex++

					texture := &wld.Texture{ // aka SimpleSpriteDef
						NameRef:        world.NameAdd(srcProp.Value),
						Flags:          0x00000000,
						TextureCurrent: 0,
						Sleep:          0,
						TextureRefs:    []uint32{uint32(fragIndex - 1)},
					}

					world.Fragments[fragIndex] = texture
					fragIndex++

					textureRefInst := &wld.TextureRef{ // aka SimpleSprite
						NameRef:    world.NameAdd(srcProp.Value),
						TextureRef: int16(fragIndex - 1),
						Flags:      0x00000000,
					}

					world.Fragments[fragIndex] = textureRefInst
					fragIndex++

					textureRefs[srcProp.Value] = fragIndex - 1
					textureRef = fragIndex - 1
				}

				dstMat.TextureRef = uint32(textureRef)
				dstMat.NameRef = world.NameAdd(srcMat.Name)
				dstMat.Flags = 2
				dstMat.RenderMethod = 0x00000001
				dstMat.RGBPen = 0x00FFFFFF
				dstMat.Brightness = 0.0
				dstMat.ScaledAmbient = 0.75

				world.Fragments[fragIndex] = dstMat
				fragIndex++
				materialList.MaterialRefs = append(materialList.MaterialRefs, uint32(fragIndex-1))
			}

			materialList.NameRef = world.NameAdd(srcMat.Name)
			materialList.Flags = 0x00014003
			if model.FileType == "ter" {
				materialList.Flags = 0x00018003
			}
			world.Fragments[fragIndex] = materialList
			fragIndex++
			mesh.MaterialListRef = uint32(fragIndex - 1)
		}
		mesh.NameRef = world.NameAdd(model.Header.Name)
		mesh.Flags = 0x00014003

		mesh.AnimationRef = 0 // for anims later
		mesh.Fragment3Ref = 0
		mesh.Fragment4Ref = 0
		//mesh.Center
		//mesh.Params2

		//mesh.MaxDistance
		//mesh.Min
		//mesh.Max
		mesh.RawScale = 13
		scale := float32(1 / float32(int(1)<<int(mesh.RawScale)))

		for _, srcVert := range model.Vertices {
			mesh.Vertices = append(mesh.Vertices, [3]int16{int16(srcVert.Position.X / scale), int16(srcVert.Position.Y / scale), int16(srcVert.Position.Z / scale)})
			mesh.Normals = append(mesh.Normals, [3]int8{int8(srcVert.Normal.X), int8(srcVert.Normal.Y), int8(srcVert.Normal.Z)})
			mesh.Colors = append(mesh.Colors, srcVert.Tint)
			mesh.UVs = append(mesh.UVs, [2]int16{int16(srcVert.Uv.X * 256), int16(srcVert.Uv.Y * 256)})
		}

		for _, srcTriangle := range model.Triangles {
			entry := wld.MeshTriangleEntry{
				Flags: uint16(srcTriangle.Flag),
				Index: [3]uint16{uint16(srcTriangle.Index.X), uint16(srcTriangle.Index.Y), uint16(srcTriangle.Index.Z)},
			}

			mesh.Triangles = append(mesh.Triangles, entry)
		}

		world.Fragments[fragIndex] = mesh
		fragIndex++

		// materialist materialpalette is used to group materials for dmspritedef2 (aka mesh)
		// dmspriteref  instance of def2 is dmsprite
		// model lastly an actordef
	}

	//bminfo
	//simplesprite

	/*
		fragIndex := 1

		// TODO: material fragments

			for _, model := range e.Models {
				mesh := &wld.Mesh{}
				name := model.Header.Name

				mesh.NameRef = int32(pos)
				mesh.Flags = 0x00014003
				mesh.MaterialListRef = uint32(0) // TODO: add proper refs
				mesh.AnimationRef = int32(0)     // TODO: add proper refs
				mesh.Vertices = model.Vertices
				world.Fragments[fragIndex] = mesh
				fragIndex++
			}
	*/
	return wld.Encode(world, w)
}
