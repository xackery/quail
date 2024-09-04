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

	return nil
}
