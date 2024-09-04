package wce

import (
	"fmt"

	"github.com/xackery/quail/raw"
)

// MdsDef is an entry EQSKINNEDMODELDEF
type MdsDef struct {
	model    string
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
	for _, face := range e.Faces {
		fmt.Fprintf(w, "\t\tFACE\n")
		fmt.Fprintf(w, "\t\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
		fmt.Fprintf(w, "\t\t\tMATERIAL \"%s\"\n", face.MaterialName)
		fmt.Fprintf(w, "\t\t\tHEXONEFLAG %d\n", face.HexOneFlag)
	}
	fmt.Fprintf(w, "\n")

	token.TagSetIsWritten(e.Tag)
	return nil
}

func (e *MdsDef) Read(token *AsciiReadToken) error {

	e.model = token.wce.lastReadModelTag

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

	e.Vertices = make([][3]float32, numVertices)
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

	e.Tints = make([][4]uint8, numTints)
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

	return nil
}

func (e *MdsDef) ToRaw(wce *Wce, dst *raw.Mds) error {

	return nil
}

func (e *MdsDef) FromRaw(wce *Wce, src *raw.Mds) error {
	e.Tag = string(src.FileName())
	e.Version = src.Version
	e.Vertices = make([][3]float32, len(src.Vertices))
	e.Normals = make([][3]float32, len(src.Vertices))
	e.UVs = make([][2]float32, len(src.Vertices))
	e.UV2s = make([][2]float32, len(src.Vertices))
	e.Tints = make([][4]uint8, len(src.Vertices))
	for i, v := range src.Vertices {
		e.Vertices[i] = v.Position
		e.Normals[i] = v.Normal
		e.UVs[i] = v.Uv
		e.UV2s[i] = v.Uv2
		e.Tints[i] = v.Tint
	}
	for _, face := range src.Triangles {
		eqFace := &EQFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flag&0x01 != 0 {
			eqFace.HexOneFlag = 1
		}
		e.Faces = append(e.Faces, eqFace)
	}

	return nil
}

// ModDef is an entry EQMODELDEF
type ModDef struct {
	model    string
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
		fmt.Fprintf(w, "\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
		fmt.Fprintf(w, "\t\tMATERIAL \"%s\"\n", face.MaterialName)
		fmt.Fprintf(w, "\t\tHEXONEFLAG %d\n", face.HexOneFlag)
	}

	fmt.Fprintf(w, "\n")

	token.TagSetIsWritten(e.Tag)
	return nil
}

func (e *ModDef) Read(token *AsciiReadToken) error {

	e.model = token.wce.lastReadModelTag

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

	e.Vertices = make([][3]float32, numVertices)
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

	e.UVs = make([][2]float32, numUVs)
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

	return nil
}

func (e *ModDef) ToRaw(wce *Wce, dst *raw.Mod) error {

	return nil
}

func (e *ModDef) FromRaw(wce *Wce, src *raw.Mod) error {
	e.Tag = string(src.FileName())

	for _, material := range src.Materials {
		eqMaterialDef := &EQMaterialDef{}
		err := eqMaterialDef.FromRaw(wce, material)
		if err != nil {
			return err
		}
		wce.EQMaterialDefs = append(wce.EQMaterialDefs, eqMaterialDef)
	}

	e.Version = src.Version
	e.Vertices = make([][3]float32, len(src.Vertices))
	e.Normals = make([][3]float32, len(src.Vertices))
	e.UVs = make([][2]float32, len(src.Vertices))
	e.UV2s = make([][2]float32, len(src.Vertices))
	e.Tints = make([][4]uint8, len(src.Vertices))
	for i, v := range src.Vertices {
		e.Vertices[i] = v.Position
		e.Normals[i] = v.Normal
		e.UVs[i] = v.Uv
		e.UV2s[i] = v.Uv2
		e.Tints[i] = v.Tint
	}

	for _, face := range src.Triangles {
		eqFace := &EQFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flag&0x01 != 0 {
			eqFace.HexOneFlag = 1
		}
		e.Faces = append(e.Faces, eqFace)
	}

	return nil
}

// TerDef is an entry EQTERRAINDEF
type TerDef struct {
	model    string
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
	return "EQMODELDEF"
}

func (e *TerDef) Write(token *AsciiWriteToken) error {
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
		fmt.Fprintf(w, "\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
		fmt.Fprintf(w, "\t\tMATERIAL \"%s\"\n", face.MaterialName)
		fmt.Fprintf(w, "\t\tHEXONEFLAG %d\n", face.HexOneFlag)
	}

	fmt.Fprintf(w, "\n")

	token.TagSetIsWritten(e.Tag)
	return nil
}

func (e *TerDef) Read(token *AsciiReadToken) error {

	e.model = token.wce.lastReadModelTag

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

	e.Vertices = make([][3]float32, numVertices)
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

	e.UVs = make([][2]float32, numUVs)
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

	e.Faces = make([]*EQFace, numFaces)
	for i := 0; i < numFaces; i++ {
		records, err = token.ReadProperty("FACE", 0)
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
	e.Tag = string(src.FileName())
	e.Version = src.Version
	e.Vertices = make([][3]float32, len(src.Vertices))
	e.Normals = make([][3]float32, len(src.Vertices))
	e.UVs = make([][2]float32, len(src.Vertices))
	e.UV2s = make([][2]float32, len(src.Vertices))
	e.Tints = make([][4]uint8, len(src.Vertices))
	for i, v := range src.Vertices {
		e.Vertices[i] = v.Position
		e.Normals[i] = v.Normal
		e.UVs[i] = v.Uv
		e.UV2s[i] = v.Uv2
		e.Tints[i] = v.Tint
	}

	for _, face := range src.Triangles {
		eqFace := &EQFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flag&0x01 != 0 {
			eqFace.HexOneFlag = 1
		}
		e.Faces = append(e.Faces, eqFace)
	}

	return nil
}

// EQMaterialDef is an entry EQMATERIALDEF
type EQMaterialDef struct {
	model             string
	Tag               string
	ShaderTag         string
	HexOneFlag        int
	Properties        []*MaterialProperty
	AnimationSleep    uint32
	AnimationTextures []string
}

type MaterialProperty struct {
	Name     string
	Category uint32
	Value    string
}

func (e *EQMaterialDef) Definition() string {
	return "EQMATERIALDEF"
}

func (e *EQMaterialDef) Write(token *AsciiWriteToken) error {
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
		fmt.Fprintf(w, "\t\tPROPERTY \"%s\" %d \"%s\"\n", prop.Name, prop.Category, prop.Value)
	}

	fmt.Fprintf(w, "\tANIMSLEEP %d\n", e.AnimationSleep)
	fmt.Fprintf(w, "\tANIMTEXTURES %d\n", len(e.AnimationTextures))
	for _, anim := range e.AnimationTextures {
		fmt.Fprintf(w, " \"%s\"", anim)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\n")

	token.TagSetIsWritten(e.Tag)
	return nil
}

func (e *EQMaterialDef) Read(token *AsciiReadToken) error {

	e.model = token.wce.lastReadModelTag

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

	e.Properties = make([]*MaterialProperty, numProperties)
	for i := 0; i < numProperties; i++ {
		records, err = token.ReadProperty("PROPERTY", 3)
		if err != nil {
			return err
		}
		prop := &MaterialProperty{}
		prop.Name = records[1]
		err = parse(&prop.Category, records[2])
		if err != nil {

			return fmt.Errorf("property category: %w", err)
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

	return nil
}

func (e *EQMaterialDef) FromRaw(wce *Wce, src *raw.Material) error {
	e.Tag = src.Name
	e.ShaderTag = src.ShaderName
	if src.Flag&0x01 != 0 {
		e.HexOneFlag = 1
	}
	e.Properties = make([]*MaterialProperty, len(src.Properties))
	for i, prop := range src.Properties {
		e.Properties[i] = &MaterialProperty{
			Name:     prop.Name,
			Category: prop.Category,
			Value:    prop.Value,
		}
	}
	e.AnimationSleep = src.Animation.Sleep
	e.AnimationTextures = src.Animation.Textures
	return nil
}
