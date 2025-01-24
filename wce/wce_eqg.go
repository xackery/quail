package wce

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/raw"
)

// MdsDef is an entry EQSKINNEDMODELDEF
type MdsDef struct {
	folders  []string
	Tag      string
	Version  uint32
	Vertices [][3]float32
	Normals  [][3]float32
	Tints    [][4]uint8
	UVs      [][2]float32
	UV2s     [][2]float32
	Faces    []*EQFace
}

type EQFace struct {
	Index        [3]uint32
	MaterialName string
	HexOneFlag   int
}

func (e *MdsDef) Definition() string {
	return "EQSKINNEDMODELDEF"
}

func (e *MdsDef) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		err := token.SetWriter(folder)
		if err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		for _, material := range token.wce.EQMaterialDefs {
			err = material.Write(token)
			if err != nil {
				return err
			}
		}

		fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
		fmt.Fprintf(w, "\tVERSION %d\n", e.Version)
		fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
		for _, v := range e.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", v[0], v[1], v[2])
		}
		fmt.Fprintf(w, "\tNUMUVs %d\n", len(e.UVs))
		for _, u := range e.UVs {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}
		fmt.Fprintf(w, "\tNUMUV2s %d\n", len(e.UV2s))
		for _, u := range e.UV2s {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}

		fmt.Fprintf(w, "\tNUMNORMALS %d\n", len(e.Normals))
		for _, n := range e.Normals {
			fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", n[0], n[1], n[2])
		}
		fmt.Fprintf(w, "\tNUMTINTS %d\n", len(e.Tints))
		for _, t := range e.Tints {
			fmt.Fprintf(w, "\t\tRGBA %d %d %d %d\n", t[0], t[1], t[2], t[3])
		}

		fmt.Fprintf(w, "\tNUMFACES %d\n", len(e.Faces))
		for i, face := range e.Faces {
			fmt.Fprintf(w, "\t\tFACE // %d\n", i)
			fmt.Fprintf(w, "\t\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
			fmt.Fprintf(w, "\t\t\tMATERIAL \"%s\"\n", face.MaterialName)
			fmt.Fprintf(w, "\t\t\tHEXONEFLAG %d\n", face.HexOneFlag)
		}
		fmt.Fprintf(w, "\n")

		token.TagSetIsWritten(e.Tag)
	}
	return nil
}

func (e *MdsDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("VERSION", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Version, records[1])
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := 0
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert := [3]float32{}
		err = parse(&vert, records[1:]...)
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}
		e.Vertices = append(e.Vertices, vert)
	}

	records, err = token.ReadProperty("NUMTINTS", 1)
	if err != nil {
		return err
	}
	numTints := 0
	err = parse(&numTints, records[1])
	if err != nil {
		return fmt.Errorf("num tints: %w", err)
	}

	for i := 0; i < numTints; i++ {
		records, err = token.ReadProperty("RGBA", 4)
		if err != nil {
			return err
		}
		tint := [4]uint8{}
		err = parse(&tint, records[1:]...)
		if err != nil {
			return fmt.Errorf("tint %d: %w", i, err)
		}

		e.Tints = append(e.Tints, tint)
	}
	records, err = token.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}
	numFaces := 0
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		_, err = token.ReadProperty("FACE", 0)
		if err != nil {
			return err
		}
		face := &EQFace{}
		records, err = token.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		err = parse(&face.Index, records[1:]...)
		if err != nil {
			return fmt.Errorf("triangle %d: %w", i, err)
		}

		records, err = token.ReadProperty("MATERIAL", 1)
		if err != nil {
			return err
		}
		face.MaterialName = records[1]

		records, err = token.ReadProperty("HEXONEFLAG", 1)
		if err != nil {
			return err
		}
		err = parse(&face.HexOneFlag, records[1])
		if err != nil {
			return fmt.Errorf("hexoneflag %d: %w", i, err)
		}

		e.Faces = append(e.Faces, face)
	}
	return nil
}

func (e *MdsDef) ToRaw(wce *Wce, dst *raw.Mds) error {
	dst.MetaFileName = e.Tag
	dst.Version = e.Version
	for i := 0; i < len(e.Vertices); i++ {
		v := &raw.Vertex{}
		v.Position = e.Vertices[i]
		v.Normal = e.Normals[i]
		v.Uv = e.UVs[i]
		v.Uv2 = e.UV2s[i]
		v.Tint = e.Tints[i]
		dst.Vertices = append(dst.Vertices, v)
	}
	for _, face := range e.Faces {
		rawFace := raw.Face{
			MaterialName: face.MaterialName,
			Index:        face.Index,
			Flags:        uint32(face.HexOneFlag),
		}
		dst.Faces = append(dst.Faces, rawFace)
	}

	return nil
}

func (e *MdsDef) FromRaw(wce *Wce, src *raw.Mds) error {
	folder := strings.TrimSuffix(strings.ToLower(wce.FileName), ".eqg")
	e.folders = append(e.folders, folder)
	e.Tag = string(src.FileName())
	e.Version = src.Version
	for _, v := range src.Vertices {

		e.Vertices = append(e.Vertices, v.Position)
		e.Normals = append(e.Normals, v.Normal)
		e.UVs = append(e.UVs, v.Uv)
		e.UV2s = append(e.UV2s, v.Uv2)
		e.Tints = append(e.Tints, v.Tint)
	}
	for _, face := range src.Faces {
		eqFace := &EQFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flags&0x01 != 0 {
			eqFace.HexOneFlag = 1
		}
		e.Faces = append(e.Faces, eqFace)
	}

	return nil
}

// ModDef is an entry EQMODELDEF
type ModDef struct {
	folders  []string
	Tag      string
	Version  uint32
	Vertices [][3]float32
	Normals  [][3]float32
	Tints    [][4]uint8
	UVs      [][2]float32
	UV2s     [][2]float32
	Faces    []*EQFace
}

func (e *ModDef) Definition() string {
	return "EQMODELDEF"
}

func (e *ModDef) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		err := token.SetWriter(folder)
		if err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		if token.TagIsWritten(e.Tag) {
			return nil
		}

		token.TagSetIsWritten(e.Tag)

		for _, material := range token.wce.EQMaterialDefs {
			err = material.Write(token)
			if err != nil {
				return err
			}
		}

		fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
		fmt.Fprintf(w, "\tVERSION %d\n", e.Version)
		fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
		for _, v := range e.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", v[0], v[1], v[2])
		}
		fmt.Fprintf(w, "\tNUMUVs %d\n", len(e.UVs))
		for _, u := range e.UVs {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}
		fmt.Fprintf(w, "\tNUMUV2s %d\n", len(e.UV2s))
		for _, u := range e.UV2s {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}

		fmt.Fprintf(w, "\tNUMNORMALS %d\n", len(e.Normals))
		for _, n := range e.Normals {
			fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", n[0], n[1], n[2])
		}
		fmt.Fprintf(w, "\tNUMTINTS %d\n", len(e.Tints))
		for _, t := range e.Tints {
			fmt.Fprintf(w, "\t\tRGBA %d %d %d %d\n", t[0], t[1], t[2], t[3])
		}

		fmt.Fprintf(w, "\tNUMFACES %d\n", len(e.Faces))
		for i, face := range e.Faces {
			fmt.Fprintf(w, "\t\tFACE // %d\n", i)
			fmt.Fprintf(w, "\t\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
			fmt.Fprintf(w, "\t\t\tMATERIAL \"%s\"\n", face.MaterialName)
			fmt.Fprintf(w, "\t\t\tHEXONEFLAG %d\n", face.HexOneFlag)
		}

		fmt.Fprintf(w, "\n")

		token.TagSetIsWritten(e.Tag)
	}
	return nil
}

func (e *ModDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("VERSION", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Version, records[1])
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := 0
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert := [3]float32{}
		err = parse(&vert, records[1:]...)
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}
		e.Vertices = append(e.Vertices, vert)
	}

	records, err = token.ReadProperty("NUMUVS", 1)
	if err != nil {
		return err
	}
	numUVs := 0
	err = parse(&numUVs, records[1])
	if err != nil {
		return fmt.Errorf("num uvs: %w", err)
	}

	for i := 0; i < numUVs; i++ {
		records, err = token.ReadProperty("UV", 2)
		if err != nil {
			return err
		}
		uv := [2]float32{}
		err = parse(&uv, records[1:]...)
		if err != nil {
			return fmt.Errorf("uv %d: %w", i, err)
		}

		e.UVs = append(e.UVs, uv)
	}
	records, err = token.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}
	numFaces := 0
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		_, err = token.ReadProperty("FACE", 0)
		if err != nil {
			return err
		}
		face := &EQFace{}
		records, err = token.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		err = parse(&face.Index, records[1:]...)
		if err != nil {
			return fmt.Errorf("triangle %d: %w", i, err)
		}

		records, err = token.ReadProperty("MATERIAL", 1)
		if err != nil {
			return err
		}
		face.MaterialName = records[1]

		records, err = token.ReadProperty("HEXONEFLAG", 1)
		if err != nil {
			return err
		}
		err = parse(&face.HexOneFlag, records[1])
		if err != nil {
			return fmt.Errorf("hexoneflag %d: %w", i, err)
		}

		e.Faces = append(e.Faces, face)
	}
	return nil
}

func (e *ModDef) ToRaw(wce *Wce, dst *raw.Mod) error {

	return nil
}

func (e *ModDef) FromRaw(wce *Wce, src *raw.Mod) error {
	e.Tag = string(src.FileName())
	folder := strings.TrimSuffix(strings.ToLower(wce.FileName), ".eqg")
	if wce.WorldDef.Zone == 1 {
		folder = "obj/" + e.Tag
	}
	e.folders = append(e.folders, folder)

	for _, material := range src.Materials {
		eqMaterialDef := &EQMaterialDef{}
		err := eqMaterialDef.FromRaw(wce, material)
		if err != nil {
			return err
		}
		wce.EQMaterialDefs = append(wce.EQMaterialDefs, eqMaterialDef)
	}

	e.Version = src.Version
	for _, v := range src.Vertices {
		e.Vertices = append(e.Vertices, v.Position)
		e.Normals = append(e.Normals, v.Normal)
		e.UVs = append(e.UVs, v.Uv)
		e.UV2s = append(e.UV2s, v.Uv2)
		e.Tints = append(e.Tints, v.Tint)
	}

	for _, face := range src.Faces {
		eqFace := &EQFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flags&0x01 != 0 {
			eqFace.HexOneFlag = 1
		}
		e.Faces = append(e.Faces, eqFace)
	}

	return nil
}

// TerDef is an entry EQTERRAINDEF
type TerDef struct {
	folders  []string
	Tag      string
	Version  uint32
	Vertices [][3]float32
	Normals  [][3]float32
	Tints    [][4]uint8
	UVs      [][2]float32
	UV2s     [][2]float32
	Faces    []*EQFace
}

func (e *TerDef) Definition() string {
	return "EQTERDEF"
}

func (e *TerDef) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		err := token.SetWriter(folder)
		if err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		if token.TagIsWritten(e.Tag) {
			return nil
		}

		token.TagSetIsWritten(e.Tag)

		fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
		fmt.Fprintf(w, "\tVERSION %d\n", e.Version)
		fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
		for _, v := range e.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", v[0], v[1], v[2])
		}
		fmt.Fprintf(w, "\tNUMUVs %d\n", len(e.UVs))
		for _, u := range e.UVs {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}
		fmt.Fprintf(w, "\tNUMUV2s %d\n", len(e.UV2s))
		for _, u := range e.UV2s {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}

		fmt.Fprintf(w, "\tNUMNORMALS %d\n", len(e.Normals))
		for _, n := range e.Normals {
			fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", n[0], n[1], n[2])
		}
		fmt.Fprintf(w, "\tNUMTINTS %d\n", len(e.Tints))
		for _, t := range e.Tints {
			fmt.Fprintf(w, "\t\tRGBA %d %d %d %d\n", t[0], t[1], t[2], t[3])
		}

		fmt.Fprintf(w, "\tNUMFACES %d\n", len(e.Faces))
		for i, face := range e.Faces {
			fmt.Fprintf(w, "\t\tFACE // %d\n", i)
			fmt.Fprintf(w, "\t\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
			fmt.Fprintf(w, "\t\t\tMATERIAL \"%s\"\n", face.MaterialName)
			fmt.Fprintf(w, "\t\t\tHEXONEFLAG %d\n", face.HexOneFlag)
		}

		fmt.Fprintf(w, "\n")

		token.TagSetIsWritten(e.Tag)
	}
	return nil
}

func (e *TerDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("VERSION", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Version, records[1])
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := 0
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert := [3]float32{}
		err = parse(&vert, records[1:]...)
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}
		e.Vertices = append(e.Vertices, vert)
	}

	records, err = token.ReadProperty("NUMUVS", 1)
	if err != nil {
		return err
	}
	numUVs := 0
	err = parse(&numUVs, records[1])
	if err != nil {
		return fmt.Errorf("num uvs: %w", err)
	}

	for i := 0; i < numUVs; i++ {
		records, err = token.ReadProperty("UV", 2)
		if err != nil {
			return err
		}
		uv := [2]float32{}
		err = parse(&uv, records[1:]...)
		if err != nil {
			return fmt.Errorf("uv %d: %w", i, err)
		}

		e.UVs = append(e.UVs, uv)
	}

	records, err = token.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}
	numFaces := 0
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		_, err = token.ReadProperty("FACE", 0)
		if err != nil {
			return err
		}
		face := &EQFace{}
		records, err = token.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		err = parse(&face.Index, records[1:]...)
		if err != nil {
			return fmt.Errorf("triangle %d: %w", i, err)
		}

		records, err = token.ReadProperty("MATERIAL", 1)
		if err != nil {
			return err
		}
		face.MaterialName = records[1]

		records, err = token.ReadProperty("HEXONEFLAG", 1)
		if err != nil {
			return err
		}
		err = parse(&face.HexOneFlag, records[1])
		if err != nil {
			return fmt.Errorf("hexoneflag %d: %w", i, err)
		}

		e.Faces = append(e.Faces, face)
	}

	return nil
}

func (e *TerDef) ToRaw(wce *Wce, dst *raw.Ter) error {

	return nil
}

func (e *TerDef) FromRaw(wce *Wce, src *raw.Ter) error {
	e.folders = append(e.folders, "world")
	e.Tag = string(src.FileName())
	e.Version = src.Version
	for _, v := range src.Vertices {
		e.Vertices = append(e.Vertices, v.Position)
		e.Normals = append(e.Normals, v.Normal)
		e.UVs = append(e.UVs, v.Uv)
		e.UV2s = append(e.UV2s, v.Uv2)
		e.Tints = append(e.Tints, v.Tint)
	}

	for _, face := range src.Triangles {
		eqFace := &EQFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flags&0x01 != 0 {
			eqFace.HexOneFlag = 1
		}
		e.Faces = append(e.Faces, eqFace)
	}

	return nil
}

// EQMaterialDef is an entry EQMATERIALDEF
type EQMaterialDef struct {
	folders           []string
	Tag               string
	ShaderTag         string
	HexOneFlag        int
	Properties        []*MaterialProperty
	AnimationSleep    uint32
	AnimationTextures []string
}

type MaterialProperty struct {
	Name  string
	Type  raw.MaterialParamType
	Value string
}

func (e *EQMaterialDef) Definition() string {
	return "EQMATERIALDEF"
}

func (e *EQMaterialDef) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		if token.TagIsWritten(e.Tag) {
			continue
		}
		err := token.SetWriter(folder)
		if err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}

		if token.TagIsWritten(e.Tag) {
			return nil
		}

		token.TagSetIsWritten(e.Tag)

		fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
		fmt.Fprintf(w, "\tSHADERTAG \"%s\"\n", e.ShaderTag)
		fmt.Fprintf(w, "\tHEXONEFLAG %d\n", e.HexOneFlag)

		fmt.Fprintf(w, "\tNUMPROPERTIES %d\n", len(e.Properties))
		for _, prop := range e.Properties {
			fmt.Fprintf(w, "\t\tPROPERTY \"%s\" %d \"%s\"\n", prop.Name, prop.Type, prop.Value)
		}

		fmt.Fprintf(w, "\tANIMSLEEP %d\n", e.AnimationSleep)
		fmt.Fprintf(w, "\tANIMTEXTURES %d\n", len(e.AnimationTextures))
		for _, anim := range e.AnimationTextures {
			fmt.Fprintf(w, " \"%s\"", anim)
		}
		fmt.Fprintf(w, "\n")

		token.TagSetIsWritten(e.Tag)
	}
	return nil
}

func (e *EQMaterialDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("SHADERTAG", 1)
	if err != nil {
		return err
	}
	e.ShaderTag = records[1]

	records, err = token.ReadProperty("HEXONEFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexOneFlag, records[1])
	if err != nil {
		return fmt.Errorf("hexoneflag: %w", err)
	}

	records, err = token.ReadProperty("NUMPROPERTIES", 1)
	if err != nil {
		return err
	}
	numProperties := 0
	err = parse(&numProperties, records[1])
	if err != nil {
		return fmt.Errorf("num properties: %w", err)
	}

	for i := 0; i < numProperties; i++ {
		records, err = token.ReadProperty("PROPERTY", 3)
		if err != nil {
			return err
		}
		prop := &MaterialProperty{}
		prop.Name = records[1]
		err = parse(&prop.Type, records[2])
		if err != nil {
			return fmt.Errorf("property param: %w", err)
		}

		prop.Value = records[3]

		e.Properties = append(e.Properties, prop)
	}

	records, err = token.ReadProperty("ANIMSLEEP", 1)

	if err != nil {
		return err
	}
	err = parse(&e.AnimationSleep, records[1])
	if err != nil {
		return fmt.Errorf("animsleep: %w", err)
	}

	records, err = token.ReadProperty("ANIMTEXTURES", -1)
	if err != nil {
		return err
	}
	numAnimTextures := 0
	err = parse(&numAnimTextures, records[1])
	if err != nil {
		return fmt.Errorf("num animtextures: %w", err)
	}

	for i := 0; i < numAnimTextures; i++ {
		e.AnimationTextures = append(e.AnimationTextures, records[i+2])
	}

	return nil
}

func (e *EQMaterialDef) ToRaw(wce *Wce, dst *raw.Material) error {
	dst.Name = e.Tag
	dst.EffectName = e.ShaderTag
	if e.HexOneFlag == 1 {
		dst.Flag = 0x01
	}
	for _, prop := range e.Properties {
		mp := &raw.MaterialParam{
			Name:  prop.Name,
			Type:  prop.Type,
			Value: prop.Value,
		}
		dst.Properties = append(dst.Properties, mp)
	}
	dst.Animation.Sleep = e.AnimationSleep
	dst.Animation.Textures = e.AnimationTextures

	return nil
}

func (e *EQMaterialDef) FromRaw(wce *Wce, src *raw.Material) error {
	folder := strings.TrimSuffix(strings.ToLower(wce.FileName), ".eqg")
	if wce.WorldDef.Zone == 1 {
		folder = "world"
	}
	e.folders = append(e.folders, folder)

	e.Tag = src.Name
	e.ShaderTag = src.EffectName
	if src.Flag&0x01 != 0 {
		e.HexOneFlag = 1
	}
	for _, prop := range src.Properties {
		mp := &MaterialProperty{
			Name:  prop.Name,
			Type:  prop.Type,
			Value: prop.Value,
		}
		e.Properties = append(e.Properties, mp)
	}
	e.AnimationSleep = src.Animation.Sleep
	e.AnimationTextures = src.Animation.Textures
	return nil
}

// AniDef represents an eqg .ani file
type AniDef struct {
	folders []string
	Tag     string
	Version uint32
	Bones   []*AniBone
	Strict  int
}

type AniBone struct {
	Name   string
	Frames []*AniBoneFrame
}

type AniBoneFrame struct {
	Milliseconds uint32
	Translation  [3]float32
	Rotation     [4]float32
	Scale        [3]float32
}

func (e *AniDef) Definition() string {
	return "EQANIDEF"
}

func (e *AniDef) Write(token *AsciiWriteToken) error {
	for _, folder := range e.folders {
		err := token.SetWriter(folder)
		if err != nil {
			return err
		}
		w, err := token.Writer()
		if err != nil {
			return err
		}
		if token.TagIsWritten(e.Tag) {
			return nil
		}

		token.TagSetIsWritten(e.Tag)

		fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
		fmt.Fprintf(w, "\tVERSION %d\n", e.Version)
		fmt.Fprintf(w, "\tSTRICT %d\n", e.Strict)
		fmt.Fprintf(w, "\tNUMBONES %d\n", len(e.Bones))
		for _, bone := range e.Bones {
			fmt.Fprintf(w, "\t\tBONE \"%s\"\n", bone.Name)
			fmt.Fprintf(w, "\t\tNUMFRAMES %d\n", len(bone.Frames))
			for i, frame := range bone.Frames {
				fmt.Fprintf(w, "\t\t\tFRAME // %d\n", i)
				fmt.Fprintf(w, "\t\t\t\tMILLISECONDS %d\n", frame.Milliseconds)
				fmt.Fprintf(w, "\t\t\t\tTRANSLATION %0.8e %0.8e %0.8e\n", frame.Translation[0], frame.Translation[1], frame.Translation[2])
				fmt.Fprintf(w, "\t\t\t\tROTATION %0.8e %0.8e %0.8e %0.8e\n", frame.Rotation[0], frame.Rotation[1], frame.Rotation[2], frame.Rotation[3])
				fmt.Fprintf(w, "\t\t\t\tSCALE %0.8e %0.8e %0.8e\n", frame.Scale[0], frame.Scale[1], frame.Scale[2])
			}
		}
		fmt.Fprintf(w, "\n")

		token.TagSetIsWritten(e.Tag)
	}
	return nil

}

func (e *AniDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("VERSION", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Version, records[1])
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}

	records, err = token.ReadProperty("STRICT", 1)
	if err != nil {
		return err
	}

	err = parse(&e.Strict, records[1])
	if err != nil {
		return fmt.Errorf("strict: %w", err)
	}

	records, err = token.ReadProperty("NUMBONES", 1)
	if err != nil {
		return err
	}
	numBones := 0
	err = parse(&numBones, records[1])
	if err != nil {
		return fmt.Errorf("num bones: %w", err)
	}

	for i := 0; i < numBones; i++ {
		_, err = token.ReadProperty("BONE", 1)
		if err != nil {
			return fmt.Errorf("bone %d: %w", i, err)
		}
		bone := &AniBone{}
		records, err = token.ReadProperty("NUMFRAMES", 1)
		if err != nil {
			return err
		}
		numFrames := 0
		err = parse(&numFrames, records[1])
		if err != nil {
			return fmt.Errorf("num frames: %w", err)
		}

		for j := 0; j < numFrames; j++ {
			frame := &AniBoneFrame{}

			_, err = token.ReadProperty("FRAME", 0)
			if err != nil {
				return err
			}

			_, err = token.ReadProperty("MILLISECONDS", 1)
			if err != nil {
				return err
			}
			err = parse(&frame.Milliseconds, records[1])
			if err != nil {
				return fmt.Errorf("milliseconds %d: %w", j, err)
			}
			records, err = token.ReadProperty("TRANSLATION", 3)
			if err != nil {
				return err
			}
			err = parse(&frame.Translation, records[1:]...)
			if err != nil {
				return fmt.Errorf("translation %d: %w", j, err)
			}

			records, err = token.ReadProperty("ROTATION", 4)
			if err != nil {
				return err
			}
			err = parse(&frame.Rotation, records[1:]...)
			if err != nil {
				return fmt.Errorf("rotation %d: %w", j, err)
			}

			records, err = token.ReadProperty("SCALE", 3)
			if err != nil {
				return err
			}

			err = parse(&frame.Scale, records[1:]...)
			if err != nil {
				return fmt.Errorf("scale %d: %w", j, err)
			}

			bone.Frames = append(bone.Frames, frame)
		}
		e.Bones = append(e.Bones, bone)
	}

	return nil
}

func (e *AniDef) ToRaw(wce *Wce, dst *raw.Ani) error {
	dst.MetaFileName = e.Tag
	dst.Version = e.Version
	for _, bone := range e.Bones {
		rawBone := &raw.AniBone{
			Name: bone.Name,
		}
		for _, frame := range bone.Frames {
			rawBoneFrame := &raw.AniBoneFrame{
				Milliseconds: frame.Milliseconds,
				Translation:  frame.Translation,
				Rotation:     frame.Rotation,
				Scale:        frame.Scale,
			}
			rawBone.Frames = append(rawBone.Frames, rawBoneFrame)
		}
		dst.Bones = append(dst.Bones, rawBone)
	}
	dst.IsStrict = e.Strict == 1
	return nil
}

func (e *AniDef) FromRaw(wce *Wce, src *raw.Ani) error {
	folder := strings.TrimSuffix(strings.ToLower(wce.FileName), ".eqg")
	e.folders = append(e.folders, folder+"_ani")
	e.Tag = src.MetaFileName
	e.Version = src.Version
	for _, bone := range src.Bones {
		aniBone := &AniBone{
			Name: bone.Name,
		}
		for _, frame := range bone.Frames {
			aniBoneFrame := &AniBoneFrame{
				Milliseconds: frame.Milliseconds,
				Translation:  frame.Translation,
				Rotation:     frame.Rotation,
				Scale:        frame.Scale,
			}

			aniBone.Frames = append(aniBone.Frames, aniBoneFrame)
		}
		e.Bones = append(e.Bones, aniBone)
	}
	if src.IsStrict {
		e.Strict = 1
	}

	return nil
}
