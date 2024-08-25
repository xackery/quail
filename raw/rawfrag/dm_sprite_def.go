package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDMSpriteDef is DmSpriteDef in libeq, Alternate Mesh in openzone, DMSPRITEDEF in wld, LegacyMesh in lantern
type WldFragDMSpriteDef struct {
	NameRef              int32
	Flags                uint32
	Fragment1            int16
	MaterialPaletteRef   uint32
	Fragment3            uint32
	CenterOffset         [3]float32
	Params1              [3]float32
	Vertices             [][3]float32
	TexCoords            [][2]float32
	Normals              [][3]float32
	Colors               []int32
	Faces                []WldFragDMSpriteDefFace
	Meshops              []WldFragDMSpriteDefMeshOp
	SkinAssignmentGroups [][2]uint16
	Data8                []uint32
	FaceMaterialGroups   [][2]int16
	VertexMaterialGroups [][2]int16
	Params2              [3]float32
	Params3              [6]float32
}

type WldFragDMSpriteDefFace struct {
	Flags         uint16
	Data          [4]uint16
	VertexIndexes [3]uint16
}

type WldFragDMSpriteDefMeshOp struct {
	TypeField   uint32
	VertexIndex uint32
	Offset      float32
	Param1      uint16
	Param2      uint16
}

func (e *WldFragDMSpriteDef) FragCode() int {
	return FragCodeDMSpriteDef
}

func (e *WldFragDMSpriteDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Vertices)))
	enc.Uint32(uint32(len(e.TexCoords)))
	enc.Uint32(uint32(len(e.Normals)))
	enc.Uint32(uint32(len(e.Colors)))
	enc.Uint32(uint32(len(e.Faces)))
	enc.Uint16(uint16(len(e.Meshops)))
	enc.Int16(e.Fragment1)
	enc.Uint32(uint32(len(e.SkinAssignmentGroups)))
	enc.Uint32(e.MaterialPaletteRef)
	enc.Uint32(e.Fragment3)
	enc.Float32(e.CenterOffset[0])
	enc.Float32(e.CenterOffset[1])
	enc.Float32(e.CenterOffset[2])
	enc.Float32(e.Params1[0])
	enc.Float32(e.Params1[1])
	enc.Float32(e.Params1[2])

	for _, vertex := range e.Vertices {
		enc.Float32(vertex[0])
		enc.Float32(vertex[1])
		enc.Float32(vertex[2])
	}

	for _, uv := range e.TexCoords {
		enc.Float32(uv[0])
		enc.Float32(uv[1])
	}

	for _, normal := range e.Normals {
		enc.Float32(normal[0])
		enc.Float32(normal[1])
		enc.Float32(normal[2])
	}

	for _, color := range e.Colors {
		enc.Int32(color)
	}

	for _, face := range e.Faces {
		enc.Uint16(face.Flags)
		enc.Uint16(face.Data[0])
		enc.Uint16(face.Data[1])
		enc.Uint16(face.Data[2])
		enc.Uint16(face.Data[3])
		enc.Uint16(face.VertexIndexes[0])
		enc.Uint16(face.VertexIndexes[1])
		enc.Uint16(face.VertexIndexes[2])
	}

	for _, meshop := range e.Meshops {
		enc.Uint32(meshop.TypeField)
		enc.Uint32(meshop.VertexIndex)
		enc.Float32(meshop.Offset)
		enc.Uint16(meshop.Param1)
		enc.Uint16(meshop.Param2)
	}

	for _, skinAssignmentGroup := range e.SkinAssignmentGroups {
		enc.Uint16(skinAssignmentGroup[0])
		enc.Uint16(skinAssignmentGroup[1])
	}

	if e.Flags&0x200 != 0 {
		enc.Uint32(uint32(len(e.Data8)))
		for _, data8 := range e.Data8 {
			enc.Uint32(data8)
		}
	}

	if e.Flags&0x800 != 0 {
		enc.Uint32(uint32(len(e.FaceMaterialGroups)))
		for _, faceMaterialGroup := range e.FaceMaterialGroups {
			enc.Int16(faceMaterialGroup[0])
			enc.Int16(faceMaterialGroup[1])
		}
	}

	if e.Flags&0x1000 != 0 {
		enc.Uint32(uint32(len(e.VertexMaterialGroups)))
		for _, vertexMaterialGroup := range e.VertexMaterialGroups {
			enc.Int16(vertexMaterialGroup[0])
			enc.Int16(vertexMaterialGroup[1])
		}
	}

	if e.Flags&0x2000 != 0 {
		enc.Float32(e.Params2[0])
		enc.Float32(e.Params2[1])
		enc.Float32(e.Params2[2])
	}

	if e.Flags&0x4000 != 0 {
		enc.Float32(e.Params3[0])
		enc.Float32(e.Params3[1])
		enc.Float32(e.Params3[2])
		enc.Float32(e.Params3[3])
		enc.Float32(e.Params3[4])
		enc.Float32(e.Params3[5])
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDMSpriteDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Uint32()
	texCoordCount := dec.Uint32()
	normalCount := dec.Uint32()
	colorCount := dec.Uint32()
	faceCount := dec.Uint32()
	meshopCount := dec.Uint16()
	e.Fragment1 = dec.Int16()
	skinAssignmentGroupCount := dec.Uint32()
	e.MaterialPaletteRef = dec.Uint32()
	e.Fragment3 = dec.Uint32()
	e.CenterOffset[0] = dec.Float32()
	e.CenterOffset[1] = dec.Float32()
	e.CenterOffset[2] = dec.Float32()
	e.Params1[0] = dec.Float32()
	e.Params1[1] = dec.Float32()
	e.Params1[2] = dec.Float32()

	for i := uint32(0); i < vertexCount; i++ {
		e.Vertices = append(e.Vertices, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}

	for i := uint32(0); i < texCoordCount; i++ {
		e.TexCoords = append(e.TexCoords, [2]float32{dec.Float32(), dec.Float32()})
	}

	for i := uint32(0); i < normalCount; i++ {
		e.Normals = append(e.Normals, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}

	for i := uint32(0); i < colorCount; i++ {
		e.Colors = append(e.Colors, dec.Int32())
	}

	for i := uint32(0); i < faceCount; i++ {
		p := WldFragDMSpriteDefFace{
			Flags:         dec.Uint16(),
			Data:          [4]uint16{dec.Uint16(), dec.Uint16(), dec.Uint16(), dec.Uint16()},
			VertexIndexes: [3]uint16{dec.Uint16(), dec.Uint16(), dec.Uint16()},
		}
		e.Faces = append(e.Faces, p)
	}

	for i := uint16(0); i < meshopCount; i++ {
		s := WldFragDMSpriteDefMeshOp{
			TypeField:   dec.Uint32(),
			VertexIndex: dec.Uint32(),
			Offset:      dec.Float32(),
			Param1:      dec.Uint16(),
			Param2:      dec.Uint16(),
		}
		e.Meshops = append(e.Meshops, s)
	}

	for i := uint32(0); i < skinAssignmentGroupCount; i++ {
		e.SkinAssignmentGroups = append(e.SkinAssignmentGroups, [2]uint16{dec.Uint16(), dec.Uint16()})
	}

	if e.Flags&0x200 != 0 {
		numData8 := dec.Uint32()
		for i := uint32(0); i < numData8; i++ {
			e.Data8 = append(e.Data8, dec.Uint32())
		}
	}

	if e.Flags&0x800 != 0 {
		faceMaterialCount := dec.Uint32()
		for i := uint32(0); i < faceMaterialCount; i++ {
			e.FaceMaterialGroups = append(e.FaceMaterialGroups, [2]int16{dec.Int16(), dec.Int16()})
		}
	}

	if e.Flags&0x1000 != 0 {
		vertexMaterialCount := dec.Uint32()
		for i := uint32(0); i < vertexMaterialCount; i++ {
			e.VertexMaterialGroups = append(e.VertexMaterialGroups, [2]int16{dec.Int16(), dec.Int16()})
		}
	}

	if e.Flags&0x2000 != 0 {
		e.Params2 = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	}

	if e.Flags&0x4000 != 0 {
		e.Params3 = [6]float32{dec.Float32(), dec.Float32(), dec.Float32(), dec.Float32(), dec.Float32(), dec.Float32()}
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil

}
