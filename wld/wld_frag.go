package wld

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

// GlobalAmbientLightDef is a declaration of GLOBALAMBIENTLIGHTDEF
type GlobalAmbientLightDef struct {
	fragID int16
	Tag    string
}

func (e *GlobalAmbientLightDef) Definition() string {
	return "GLOBALAMBIENTLIGHTDEF"
}

func (e *GlobalAmbientLightDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "ENDGLOBALAMBIENTLIGHTDEF\n\n")
	return nil
}

func (e *GlobalAmbientLightDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	_, err = r.ReadProperty("ENDGLOBALAMBIENTLIGHTDEF", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *GlobalAmbientLightDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfGlobalAmbientLightDef := &rawfrag.WldFragGlobalAmbientLightDef{
		NameRef: -16777216,
	}
	if e.Tag != "DEFAULT_AMBIENTLIGHT" {
		wfGlobalAmbientLightDef.NameRef = raw.NameAdd(e.Tag)
	}

	dst.Fragments = append(dst.Fragments, wfGlobalAmbientLightDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

// DMSpriteDef2 is a declaration of DMSpriteDef2
type DMSpriteDef2 struct {
	fragID               int16
	Tag                  string
	Flags                uint32
	DmTrackTag           string
	Fragment3Ref         int32
	Fragment4Ref         int32
	Params2              [3]uint32
	MaxDistance          float32
	Min                  [3]float32
	Max                  [3]float32
	CenterOffset         [3]float32
	Vertices             [][3]float32
	UVs                  [][2]float32
	VertexNormals        [][3]float32
	VertexColors         [][4]uint8
	SkinAssignmentGroups [][2]uint16
	MaterialPaletteTag   string
	Faces                []*Face
	MeshOps              []*MeshOp
	FaceMaterialGroups   [][2]uint16
	VertexMaterialGroups [][2]int16
	BoundingRadius       float32
	FPScale              uint16
	PolyhedronTag        string
}

type Face struct {
	Flags    uint16
	Triangle [3]uint16
}

type MeshOp struct {
	Index1    uint16
	Index2    uint16
	Offset    float32
	Param1    uint8
	TypeField uint8
}

func (e *DMSpriteDef2) Definition() string {
	return "DMSPRITEDEF2"
}

func (e *DMSpriteDef2) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\t// FLAGS \"%d\" // need to assess\n", e.Flags)
	fmt.Fprintf(w, "\tCENTEROFFSET %0.7e %0.7e %0.7e\n", e.CenterOffset[0], e.CenterOffset[1], e.CenterOffset[2])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMUVS %d\n", len(e.UVs))
	for _, uv := range e.UVs {
		fmt.Fprintf(w, "\tUV %0.7e %0.7e\n", uv[0], uv[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTEXNORMALS %d\n", len(e.VertexNormals))
	for _, vn := range e.VertexNormals {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vn[0], vn[1], vn[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTEXCOLORS %d\n", len(e.VertexColors))
	for _, color := range e.VertexColors {
		fmt.Fprintf(w, "\tRGBA %d %d %d %d\n", color[0], color[1], color[2], color[3])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tSKINASSIGNMENTGROUPS %d", len(e.SkinAssignmentGroups))
	for _, sa := range e.SkinAssignmentGroups {
		fmt.Fprintf(w, " %d %d", sa[0], sa[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tMATERIALPALETTE \"%s\"\n", e.MaterialPaletteTag)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tPOLYHEDRON\n")
	fmt.Fprintf(w, "\t\tDEFINITION \"%s\"\n", e.PolyhedronTag)
	fmt.Fprintf(w, "\tENDPOLYHEDRON\n\n")
	fmt.Fprintf(w, "\tNUMFACE2S %d\n", len(e.Faces))
	fmt.Fprintf(w, "\n")
	for i, face := range e.Faces {
		fmt.Fprintf(w, "\tDMFACE2 //%d\n", i+1)
		fmt.Fprintf(w, "\t\tFLAGS %d\n", face.Flags)
		fmt.Fprintf(w, "\t\tTRIANGLE %d %d %d\n", face.Triangle[0], face.Triangle[1], face.Triangle[2])
		fmt.Fprintf(w, "\tENDDMFACE2 //%d\n\n", i+1)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\t// meshops are not supported\n")
	fmt.Fprintf(w, "\t// NUMMESHOPS %d\n", len(e.MeshOps))
	for _, meshOp := range e.MeshOps {
		fmt.Fprintf(w, "\t// TODO: MESHOP %d %d %0.7f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tMAXDISTANCE %0.7e\n", e.MaxDistance)
	fmt.Fprintf(w, "\tFACEMATERIALGROUPS %d", len(e.FaceMaterialGroups))
	for _, group := range e.FaceMaterialGroups {
		fmt.Fprintf(w, " %d %d", group[0], group[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tVERTEXMATERIALGROUPS %d", len(e.VertexMaterialGroups))
	for _, group := range e.VertexMaterialGroups {
		fmt.Fprintf(w, " %d %d", group[0], group[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7e\n", e.BoundingRadius)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFPSCALE %d\n", e.FPScale)
	fmt.Fprintf(w, "ENDDMSPRITEDEF2\n\n")
	return nil
}

func (e *DMSpriteDef2) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("CENTEROFFSET", 3)
	if err != nil {
		return err
	}
	err = parse(&e.CenterOffset, records[1:]...)
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = r.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := int(0)
	err = parse(&numVertices, records[1])
	if err != nil {
		return err
	}
	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ", 3)
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

	records, err = r.ReadProperty("NUMUVS", 1)
	if err != nil {
		return err
	}
	numUVs := int(0)
	err = parse(&numUVs, records[1])
	if err != nil {
		return fmt.Errorf("num uvs: %w", err)
	}

	for i := 0; i < numUVs; i++ {
		records, err = r.ReadProperty("UV", 2)
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

	records, err = r.ReadProperty("NUMVERTEXNORMALS", 1)
	if err != nil {
		return err
	}
	numNormals := int(0)
	err = parse(&numNormals, records[1])
	if err != nil {
		return fmt.Errorf("num normals: %w", err)
	}

	for i := 0; i < numNormals; i++ {
		records, err = r.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		normal := [3]float32{}
		err = parse(&normal, records[1:]...)
		if err != nil {
			return fmt.Errorf("normal %d: %w", i, err)
		}
		e.VertexNormals = append(e.VertexNormals, normal)
	}

	records, err = r.ReadProperty("NUMVERTEXCOLORS", 1)
	if err != nil {
		return err
	}
	numColors := int(0)
	err = parse(&numColors, records[1])
	if err != nil {
		return fmt.Errorf("num colors: %w", err)
	}

	for i := 0; i < numColors; i++ {
		records, err = r.ReadProperty("RGBA", 4)
		if err != nil {
			return err
		}
		color := [4]uint8{}
		err = parse(&color, records[1:]...)
		if err != nil {
			return fmt.Errorf("color %d: %w", i, err)
		}
		e.VertexColors = append(e.VertexColors, color)
	}

	records, err = r.ReadProperty("SKINASSIGNMENTGROUPS", 1)
	if err != nil {
		return err
	}
	numSkinAssignments := int(0)
	err = parse(&numSkinAssignments, records[1])
	if err != nil {
		return fmt.Errorf("num skin assignments: %w", err)
	}

	for i := 0; i < numSkinAssignments; i++ {
		val1 := uint8(0)
		err = parse(&val1, records[i+1])
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}

		val2 := uint8(0)
		err = parse(&val2, records[i+2])
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}
		e.SkinAssignmentGroups = append(e.SkinAssignmentGroups, [2]uint16{uint16(val1), uint16(val2)})
	}

	records, err = r.ReadProperty("MATERIALPALETTE", 1)
	if err != nil {
		return err
	}
	e.MaterialPaletteTag = records[1]

	_, err = r.ReadProperty("POLYHEDRON", 0)
	if err != nil {
		return err
	}
	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.PolyhedronTag = records[1]
	_, err = r.ReadProperty("ENDPOLYHEDRON", 0)
	if err != nil {
		return err
	}

	records, err = r.ReadProperty("NUMFACE2S", 1)
	if err != nil {
		return err
	}

	numFaces := int(0)
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		face := &Face{}
		_, err = r.ReadProperty("DMFACE2", 0)
		if err != nil {
			return err
		}
		records, err = r.ReadProperty("FLAGS", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Flags, records[1])
		if err != nil {
			return fmt.Errorf("face %d flags: %w", i, err)
		}

		records, err = r.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		err = parse(&face.Triangle, records[1:]...)
		if err != nil {
			return fmt.Errorf("face %d triangle: %w", i, err)
		}

		_, err = r.ReadProperty("ENDDMFACE2", 0)
		if err != nil {
			return err
		}

		e.Faces = append(e.Faces, face)
	}

	records, err = r.ReadProperty("MAXDISTANCE", 1)
	if err != nil {
		return err
	}

	err = parse(&e.MaxDistance, records[1])
	if err != nil {
		return fmt.Errorf("max distance: %w", err)
	}

	records, err = r.ReadProperty("FACEMATERIALGROUPS", -1)
	if err != nil {
		return err
	}
	numFaceMaterialGroups := int(0)
	err = parse(&numFaceMaterialGroups, records[1])
	if err != nil {
		return fmt.Errorf("num face material groups: %w", err)
	}

	for i := 0; i < numFaceMaterialGroups; i++ {
		val1, err := strconv.ParseUint(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		val2, err := strconv.ParseUint(records[i+3], 10, 16)
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		e.FaceMaterialGroups = append(e.FaceMaterialGroups, [2]uint16{uint16(val1), uint16(val2)})
	}

	records, err = r.ReadProperty("VERTEXMATERIALGROUPS", -1)
	if err != nil {
		return err
	}
	numVertexMaterialGroups := int(0)
	err = parse(&numVertexMaterialGroups, records[1])
	if err != nil {
		return fmt.Errorf("num vertex material groups: %w", err)
	}

	for i := 0; i < numVertexMaterialGroups; i++ {
		val1, err := strconv.ParseInt(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		val2, err := strconv.ParseInt(records[i+3], 10, 16)
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		e.VertexMaterialGroups = append(e.VertexMaterialGroups, [2]int16{int16(val1), int16(val2)})
	}

	records, err = r.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("FPSCALE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.FPScale, records[1])
	if err != nil {
		return fmt.Errorf("fpscale: %w", err)
	}

	_, err = r.ReadProperty("ENDDMSPRITEDEF2", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *DMSpriteDef2) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	var err error

	if e.fragID != 0 {
		return e.fragID, nil
	}

	materialPaletteRef := int16(0)
	if e.MaterialPaletteTag != "" {
		palette := srcWld.ByTag(e.MaterialPaletteTag)
		if palette == nil {
			return -1, fmt.Errorf("material palette %s not found", e.MaterialPaletteTag)
		}

		materialPaletteRef, err = palette.ToRaw(srcWld, dst)
		if err != nil {
			return -1, fmt.Errorf("material palette %s to raw: %w", e.MaterialPaletteTag, err)
		}
	}
	dmSpriteDef := &rawfrag.WldFragDmSpriteDef2{
		NameRef:              raw.NameAdd(e.Tag),
		MaterialPaletteRef:   uint32(materialPaletteRef),
		CenterOffset:         e.CenterOffset,
		Params2:              e.Params2,
		MaxDistance:          e.MaxDistance,
		Min:                  e.Min,
		Max:                  e.Max,
		Scale:                e.FPScale,
		Colors:               e.VertexColors,
		FaceMaterialGroups:   e.FaceMaterialGroups,
		VertexMaterialGroups: e.VertexMaterialGroups,
	}
	if srcWld.isZone {
		dmSpriteDef.Flags = 0x00018003
	} else {
		dmSpriteDef.Flags = 0x00014003
	}

	for i, frag := range dst.Fragments {
		_, ok := frag.(*rawfrag.WldFragBMInfo)
		if !ok {
			continue
		}
		dmSpriteDef.Fragment4Ref = int32(i) + 1
	}

	scale := float32(1 / float32(int(1)<<int(e.FPScale)))

	for _, face := range e.Faces {
		dmSpriteDef.Faces = append(dmSpriteDef.Faces, rawfrag.WldFragMeshFaceEntry{
			Flags: face.Flags,
			Index: face.Triangle,
		})
	}

	for _, vert := range e.Vertices {
		dmSpriteDef.Vertices = append(dmSpriteDef.Vertices, [3]int16{
			int16(vert[0] / scale),
			int16(vert[1] / scale),
			int16(vert[2] / scale),
		})
	}

	for _, uv := range e.UVs {
		dmSpriteDef.UVs = append(dmSpriteDef.UVs, [2]int16{
			int16(uv[0] / scale),
			int16(uv[1] / scale),
		})
	}

	for _, normal := range e.VertexNormals {
		dmSpriteDef.VertexNormals = append(dmSpriteDef.VertexNormals, [3]int8{
			int8(normal[0] / scale),
			int8(normal[1] / scale),
			int8(normal[2] / scale),
		})
	}

	dst.Fragments = append(dst.Fragments, dmSpriteDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

// DMSpriteDef is a declaration of DMSPRITEDEF
type DMSpriteDef struct {
	fragID         int16
	Tag            string
	Flags          uint32
	Fragment1Maybe int16
	Material       string
	Fragment3      uint32
	CenterPosition [3]float32
	Params2        uint32
	Something2     uint32
	Something3     uint32
	Verticies      [][3]float32
	TexCoords      [][3]float32
	Normals        [][3]float32
	Colors         []int32
	Polygons       []*DMSpriteDefSpritePolygon
	VertexPieces   []*DMSpriteDefVertexPiece
	PostVertexFlag uint32
	RenderGroups   []*DMSpriteDefRenderGroup
	VertexTex      [][2]float32
	Size6Pieces    []*DMSpriteDefSize6Entry
}

type DMSpriteDefSpritePolygon struct {
	Flag int16
	Unk1 int16
	Unk2 int16
	Unk3 int16
	Unk4 int16
	I1   int16
	I2   int16
	I3   int16
}

type DMSpriteDefVertexPiece struct {
	Count  int16
	Offset int16
}

type DMSpriteDefRenderGroup struct {
	PolygonCount int16
	MaterialId   int16
}

type DMSpriteDefSize6Entry struct {
	Unk1 uint32
	Unk2 uint32
	Unk3 uint32
	Unk4 uint32
	Unk5 uint32
}

func (e *DMSpriteDef) Definition() string {
	return "DMSPRITEDEF"
}

func (e *DMSpriteDef) Write(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func (e *DMSpriteDef) Read(r *AsciiReadToken) error {
	return fmt.Errorf("not implemented")
}

func (e *DMSpriteDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	return -1, fmt.Errorf("not implemented")
}

// DMSprite is a declaration of DMSPRITEINSTANCE
type DMSprite struct {
	fragID        int16
	Tag           string
	DefinitionTag string
	Param         uint32
}

func (e *DMSprite) Definition() string {
	return "DMSPRITEINSTANCE"
}

func (e *DMSprite) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", e.DefinitionTag)
	fmt.Fprintf(w, "\tPARAM %d\n", e.Param)
	fmt.Fprintf(w, "ENDDMSPRITEINSTANCE\n\n")
	return nil
}

func (e *DMSprite) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}

	e.DefinitionTag = records[1]

	records, err = r.ReadProperty("PARAM", 1)
	if err != nil {
		return err
	}

	err = parse(&e.Param, records[1])
	if err != nil {
		return fmt.Errorf("param: %w", err)
	}

	_, err = r.ReadProperty("ENDDMSPRITEINSTANCE", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *DMSprite) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	return -1, fmt.Errorf("not implemented")
}

// MaterialPalette is a declaration of MATERIALPALETTE
type MaterialPalette struct {
	fragID    int16
	Tag       string
	flags     uint32
	Materials []string
}

func (e *MaterialPalette) Definition() string {
	return "MATERIALPALETTE"
}

func (e *MaterialPalette) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tNUMMATERIALS %d\n", len(e.Materials))
	for _, mat := range e.Materials {
		fmt.Fprintf(w, "\tMATERIAL \"%s\"\n", mat)
	}
	fmt.Fprintf(w, "ENDMATERIALPALETTE\n\n")
	return nil
}

func (e *MaterialPalette) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("NUMMATERIALS", 1)
	if err != nil {
		return fmt.Errorf("NUMMATERIALS: %w", err)
	}
	numMaterials := int(0)
	err = parse(&numMaterials, records[1])
	if err != nil {
		return fmt.Errorf("num materials: %w", err)
	}

	for i := 0; i < numMaterials; i++ {
		records, err = r.ReadProperty("MATERIAL", 1)
		if err != nil {
			return fmt.Errorf("MATERIAL: %w", err)
		}
		e.Materials = append(e.Materials, records[1])
	}

	_, err = r.ReadProperty("ENDMATERIALPALETTE", 0)
	if err != nil {
		return fmt.Errorf("ENDMATERIALPALETTE: %w", err)
	}

	return nil
}

func (e *MaterialPalette) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfPalette := &rawfrag.WldFragMaterialPalette{
		Flags: e.flags,
	}
	for _, mat := range e.Materials {

		srcMat := srcWld.ByTag(mat)
		if srcMat == nil {
			return -1, fmt.Errorf("material %s not found", mat)
		}

		matRef, err := srcMat.ToRaw(srcWld, dst)
		if err != nil {
			return -1, fmt.Errorf("material %s to raw: %w", mat, err)
		}

		wfPalette.MaterialRefs = append(wfPalette.MaterialRefs, uint32(matRef))
	}

	wfPalette.NameRef = raw.NameAdd(e.Tag)
	dst.Fragments = append(dst.Fragments, wfPalette)
	e.fragID = int16(len(dst.Fragments))

	return int16(len(dst.Fragments)), nil
}

// MaterialDef is an entry MATERIALDEFINITION
type MaterialDef struct {
	fragID          int16
	Tag             string   // TAG %s
	RenderMethod    string   // RENDERMETHOD %s
	RGBPen          [4]uint8 // RGBPEN %d %d %d
	Brightness      float32  // BRIGHTNESS %0.7f
	ScaledAmbient   float32  // SCALEDAMBIENT %0.7f
	SimpleSpriteTag string   // SIMPLESPRITEINST
	Pair1           NullUint32
	Pair2           NullFloat32
}

func (e *MaterialDef) Definition() string {
	return "MATERIALDEFINITION"
}

func (e *MaterialDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tRENDERMETHOD \"%s\"\n", e.RenderMethod)
	fmt.Fprintf(w, "\tRGBPEN %d %d %d %d\n", e.RGBPen[0], e.RGBPen[1], e.RGBPen[2], e.RGBPen[3])
	fmt.Fprintf(w, "\tBRIGHTNESS %0.7f\n", e.Brightness)
	fmt.Fprintf(w, "\tSCALEDAMBIENT %0.7f\n", e.ScaledAmbient)
	fmt.Fprintf(w, "\tSIMPLESPRITEINST\n")
	fmt.Fprintf(w, "\t\tTAG \"%s\"\n", e.SimpleSpriteTag)
	fmt.Fprintf(w, "\tENDSIMPLESPRITEINST\n")
	fmt.Fprintf(w, "\tPAIRS? %s %s\n", wcVal(e.Pair1), wcVal(e.Pair2))
	fmt.Fprintf(w, "ENDMATERIALDEFINITION\n\n")
	return nil
}

func (e *MaterialDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("RENDERMETHOD", 1)
	if err != nil {
		return err
	}
	e.RenderMethod = records[1]

	records, err = r.ReadProperty("RGBPEN", 4)
	if err != nil {
		return err
	}
	err = parse(&e.RGBPen, records[1:]...)
	if err != nil {
		return fmt.Errorf("rgbpen: %w", err)
	}

	records, err = r.ReadProperty("BRIGHTNESS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Brightness, records[1])
	if err != nil {
		return fmt.Errorf("brightness: %w", err)
	}

	records, err = r.ReadProperty("SCALEDAMBIENT", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ScaledAmbient, records[1])
	if err != nil {
		return fmt.Errorf("scaled ambient: %w", err)
	}

	_, err = r.ReadProperty("SIMPLESPRITEINST", 0)
	if err != nil {
		return err
	}

	records, err = r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.SimpleSpriteTag = records[1]

	_, err = r.ReadProperty("ENDSIMPLESPRITEINST", 0)
	if err != nil {
		return err
	}

	records, err = r.ReadProperty("PAIRS?", 2)
	if err != nil {
		return err
	}

	err = parse(&e.Pair1, records[1])
	if err != nil {
		return fmt.Errorf("has pairs: %w", err)
	}

	err = parse(&e.Pair2, records[2])
	if err != nil {
		return fmt.Errorf("pair1: %w", err)
	}

	_, err = r.ReadProperty("ENDMATERIALDEFINITION", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *MaterialDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfMaterialDef := &rawfrag.WldFragMaterialDef{
		RenderMethod:  model.RenderMethodInt(e.RenderMethod),
		RGBPen:        e.RGBPen,
		Brightness:    e.Brightness,
		ScaledAmbient: e.ScaledAmbient,
	}

	if e.Pair1.Valid && e.Pair2.Valid {
		wfMaterialDef.Pair1 = e.Pair1.Uint32
		wfMaterialDef.Pair2 = e.Pair2.Float32
		wfMaterialDef.Flags |= 0x02
	} else {
		wfMaterialDef.Pair1 = 0
		wfMaterialDef.Pair2 = 0
	}

	if e.SimpleSpriteTag != "" {
		spriteDef := srcWld.ByTag(e.SimpleSpriteTag)
		if spriteDef == nil {
			return -1, fmt.Errorf("simple sprite %s not found", e.SimpleSpriteTag)
		}

		spriteDefRef, err := spriteDef.ToRaw(srcWld, dst)
		if err != nil {
			return -1, fmt.Errorf("simple sprite %s to raw: %w", e.SimpleSpriteTag, err)
		}

		wfSprite := &rawfrag.WldFragSimpleSprite{
			//NameRef:   raw.NameAdd(m.SimpleSpriteTag),
			SpriteRef: uint32(spriteDefRef),
		}

		wfSprite.Flags = 0x50 // just a static generic thing
		dst.Fragments = append(dst.Fragments, wfSprite)

		spriteRef := int16(len(dst.Fragments))

		wfMaterialDef.SimpleSpriteRef = uint32(spriteRef)
	}

	wfMaterialDef.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfMaterialDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

// SimpleSpriteDef is a declaration of SIMPLESPRITEDEF
type SimpleSpriteDef struct {
	fragID             int16
	Tag                string
	SkipFrames         NullUint32
	Sleep              NullUint32
	CurrentFrame       NullInt32
	Animated           NullUint32
	SimpleSpriteFrames []SimpleSpriteFrame
}

type SimpleSpriteFrame struct {
	TextureFile string
	TextureTag  string
}

func (e *SimpleSpriteDef) Definition() string {
	return "SIMPLESPRITEDEF"
}

func (e *SimpleSpriteDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tSIMPLESPRITETAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tSKIPFRAMES? %s\n", wcVal(e.SkipFrames))
	fmt.Fprintf(w, "\tANIMATED? %s\n", wcVal(e.Animated))
	fmt.Fprintf(w, "\tSLEEP? %s\n", wcVal(e.Sleep))
	fmt.Fprintf(w, "\tCURRENTFRAME? %s\n", wcVal(e.CurrentFrame))
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.SimpleSpriteFrames))
	for _, frame := range e.SimpleSpriteFrames {
		fmt.Fprintf(w, "\tFRAME \"%s\" \"%s\"\n", frame.TextureFile, frame.TextureTag)
	}
	fmt.Fprintf(w, "ENDSIMPLESPRITEDEF\n\n")
	return nil
}

func (e *SimpleSpriteDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("SIMPLESPRITETAG", 0)
	if err != nil {
		return fmt.Errorf("SIMPLESPRITETAG: %w", err)
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("SKIPFRAMES?", 1)
	if err != nil {
		return fmt.Errorf("SKIPFRAMES?: %w", err)
	}

	err = parse(&e.SkipFrames, records[1])
	if err != nil {
		return fmt.Errorf("skip frames: %w", err)
	}

	records, err = r.ReadProperty("ANIMATED?", 1)
	if err != nil {
		return fmt.Errorf("ANIMATED?: %w", err)
	}
	err = parse(&e.Animated, records[1])
	if err != nil {
		return fmt.Errorf("animated: %w", err)
	}

	records, err = r.ReadProperty("SLEEP?", 1)
	if err != nil {
		return fmt.Errorf("SLEEP?: %w", err)
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = r.ReadProperty("CURRENTFRAME?", 1)
	if err != nil {
		return fmt.Errorf("CURRENTFRAME?: %w", err)
	}

	err = parse(&e.CurrentFrame, records[1])
	if err != nil {
		return fmt.Errorf("current frame: %w", err)
	}

	records, err = r.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return fmt.Errorf("NUMFRAMES: %w", err)
	}
	numFrames := int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = r.ReadProperty("FRAME", 2)
		if err != nil {
			return fmt.Errorf("FRAME: %w", err)
		}
		e.SimpleSpriteFrames = append(e.SimpleSpriteFrames, SimpleSpriteFrame{
			TextureFile: records[1],
			TextureTag:  records[2],
		})
	}

	_, err = r.ReadProperty("ENDSIMPLESPRITEDEF", 0)
	if err != nil {
		return fmt.Errorf("ENDSIMPLESPRITEDEF: %w", err)
	}
	return nil
}

func (e *SimpleSpriteDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	//if e.fragID != 0 {
	//	return e.fragID, nil
	//}
	flags := uint32(0)
	wfSimpleSpriteDef := &rawfrag.WldFragSimpleSpriteDef{
		Sleep: e.Sleep.Uint32,
	}

	if e.SkipFrames.Valid {
		flags |= 0x02
	}
	//flags |= 0x04
	if len(e.SimpleSpriteFrames) > 1 {
		flags |= 0x08
	}

	if e.Animated.Valid {
		flags |= 0x08
	} else {
		if len(e.SimpleSpriteFrames) > 1 {
			return -1, fmt.Errorf("simple sprite %s is not animated but has more than one frame", e.Tag)
		}
	}

	if e.Sleep.Valid {
		flags |= 0x10
	}
	if e.CurrentFrame.Valid {
		flags |= 0x20
	}

	wfSimpleSpriteDef.Flags = flags

	for _, frame := range e.SimpleSpriteFrames {
		wfBMInfo := &rawfrag.WldFragBMInfo{
			NameRef:      raw.NameAdd(frame.TextureTag),
			TextureNames: []string{frame.TextureFile + "\x00"},
		}

		dst.Fragments = append(dst.Fragments, wfBMInfo)
		wfSimpleSpriteDef.BitmapRefs = append(wfSimpleSpriteDef.BitmapRefs, uint32(len(dst.Fragments)))
	}

	wfSimpleSpriteDef.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfSimpleSpriteDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

// ActorDef is a declaration of ACTORDEF
type ActorDef struct {
	fragID        int16
	Tag           string
	Callback      string
	BoundsRef     int32
	CurrentAction NullUint32
	Location      [6]float32
	Unk1          uint32
	Actions       []ActorAction
	Unk2          uint32
}

func (e *ActorDef) Definition() string {
	return "ACTORDEF"
}

func (e *ActorDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tACTORTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tCALLBACK \"%s\"\n", e.Callback)
	fmt.Fprintf(w, "\t// BOUNDSREF %d\n", e.BoundsRef)
	fmt.Fprintf(w, "\tCURRENTACTION? %s\n", wcVal(e.CurrentAction))
	fmt.Fprintf(w, "\tLOCATION %0.7f %0.7f %0.7f %0.7f %0.7f %0.7f\n", e.Location[0], e.Location[1], e.Location[2], e.Location[3], e.Location[4], e.Location[5])
	fmt.Fprintf(w, "\tNUMACTIONS %d\n", len(e.Actions))
	for _, action := range e.Actions {
		fmt.Fprintf(w, "\tACTION\n")
		fmt.Fprintf(w, "\t\t// UNK1 %d\n", action.Unk1)
		fmt.Fprintf(w, "\t\tNUMLEVELSOFDETAIL %d\n", len(action.LevelOfDetails))
		for _, lod := range action.LevelOfDetails {
			fmt.Fprintf(w, "\t\tLEVELOFDETAIL\n")
			fmt.Fprintf(w, "\t\t\tSPRITE \"%s\"\n", lod.SpriteTag)
			fmt.Fprintf(w, "\t\t\tMINDISTANCE %0.7f\n", lod.MinDistance)
			fmt.Fprintf(w, "\t\tENDLEVELOFDETAIL\n")
		}
		fmt.Fprintf(w, "\tENDACTION\n")
	}
	fmt.Fprintf(w, "\t// UNK2 %d\n", e.Unk2)
	fmt.Fprintf(w, "ENDACTORDEF\n\n")
	return nil
}

func (e *ActorDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("ACTORTAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("CALLBACK", 1)
	if err != nil {
		return err
	}
	e.Callback = records[1]

	records, err = r.ReadProperty("CURRENTACTION?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentAction, records[1])
	if err != nil {
		return fmt.Errorf("current action: %w", err)
	}

	records, err = r.ReadProperty("LOCATION", 6)
	if err != nil {
		return err
	}
	err = parse(&e.Location, records[1:]...)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = r.ReadProperty("NUMACTIONS", 1)
	if err != nil {
		return err
	}
	numActions := int(0)
	err = parse(&numActions, records[1])
	if err != nil {
		return fmt.Errorf("num actions: %w", err)
	}

	for i := 0; i < numActions; i++ {
		action := ActorAction{}
		_, err = r.ReadProperty("ACTION", 0)
		if err != nil {
			return err
		}

		records, err = r.ReadProperty("NUMLEVELSOFDETAIL", 1)
		if err != nil {
			return err
		}

		numLod := int(0)
		err = parse(&numLod, records[1])
		if err != nil {
			return fmt.Errorf("num lod: %w", err)
		}

		for j := 0; j < numLod; j++ {
			lod := ActorLevelOfDetail{}
			_, err = r.ReadProperty("LEVELOFDETAIL", 0)
			if err != nil {
				return err
			}

			records, err = r.ReadProperty("SPRITE", 1)
			if err != nil {
				return err
			}
			lod.SpriteTag = records[1]

			records, err = r.ReadProperty("MINDISTANCE", 1)
			if err != nil {
				return err
			}

			err = parse(&lod.MinDistance, records[1])
			if err != nil {
				return fmt.Errorf("min distance: %w", err)
			}

			_, err = r.ReadProperty("ENDLEVELOFDETAIL", 0)
			if err != nil {
				return err
			}

			action.LevelOfDetails = append(action.LevelOfDetails, lod)
		}

		_, err = r.ReadProperty("ENDACTION", 0)
		if err != nil {
			return err
		}

		e.Actions = append(e.Actions, action)

	}

	_, err = r.ReadProperty("ENDACTORDEF", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *ActorDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	var err error
	if e.fragID != 0 {
		return e.fragID, nil
	}

	actorDef := &rawfrag.WldFragActorDef{
		BoundsRef:     e.BoundsRef,
		CurrentAction: e.CurrentAction.Uint32,
		Location:      e.Location,
	}

	for _, action := range e.Actions {
		actorAction := rawfrag.WldFragModelAction{
			Unk1: action.Unk1,
		}

		for _, lod := range action.LevelOfDetails {
			if lod.SpriteTag == "" {
				continue
			}

			var spriteRef int16
			spriteVar := srcWld.ByTag(lod.SpriteTag)
			if spriteVar == nil {
				return -1, fmt.Errorf("sprite %s not found", lod.SpriteTag)
			}
			switch spriteDef := spriteVar.(type) {
			case *DMSpriteDef2:
				spriteRef, err = spriteDef.ToRaw(srcWld, dst)
			case *Sprite3DDef:
				spriteRef, err = spriteDef.ToRaw(srcWld, dst)
				if err != nil {
					return -1, fmt.Errorf("sprite %s to raw: %w", lod.SpriteTag, err)
				}
				sprite := &rawfrag.WldFragSprite3D{
					Flags:          lod.SpriteFlags,
					Sprite3DDefRef: int32(spriteRef),
				}

				dst.Fragments = append(dst.Fragments, sprite)
				spriteRef = int16(len(dst.Fragments))
			default:
				return -1, fmt.Errorf("unknown sprite type %T", spriteDef)
			}
			if err != nil {
				return -1, fmt.Errorf("sprite %s to raw: %w", lod.SpriteTag, err)
			}

			actorAction.Lods = append(actorAction.Lods, lod.MinDistance)
			actorDef.FragmentRefs = append(actorDef.FragmentRefs, uint32(spriteRef))
		}

		actorDef.Actions = append(actorDef.Actions, actorAction)
	}

	actorDef.NameRef = raw.NameAdd(e.Tag)
	actorDef.CallbackNameRef = raw.NameAdd(e.Callback)

	dst.Fragments = append(dst.Fragments, actorDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), err
}

// ActorAction is a declaration of ACTION
type ActorAction struct {
	Unk1           uint32
	LevelOfDetails []ActorLevelOfDetail
}

// ActorLevelOfDetail is a declaration of LEVELOFDETAIL
type ActorLevelOfDetail struct {
	SpriteTag   string
	SpriteFlags uint32
	MinDistance float32
}

// ActorInst is a declaration of ACTORINST
type ActorInst struct {
	fragID           int16
	Tag              string
	DefinitionTag    string
	CurrentAction    NullUint32
	Location         NullFloat32Slice6
	BoundingRadius   NullFloat32
	Scale            NullFloat32
	SoundTag         NullString
	Active           NullUint32
	SpriteVolumeOnly NullUint32
	DMRGBTrackTag    NullString
	SphereTag        string
	SphereRadius     float32
	UserData         string
}

func (e *ActorInst) Definition() string {
	return "ACTORINST"
}

func (e *ActorInst) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", e.DefinitionTag)
	fmt.Fprintf(w, "\tCURRENTACTION? %s\n", wcVal(e.CurrentAction))
	fmt.Fprintf(w, "\tLOCATION? %s\n", wcVal(e.Location))
	fmt.Fprintf(w, "\tBOUNDINGRADIUS? %s\n", wcVal(e.BoundingRadius))
	fmt.Fprintf(w, "\tSCALEFACTOR? %s\n", wcVal(e.Scale))
	fmt.Fprintf(w, "\tSOUND? \"%s\"\n", wcVal(e.SoundTag))
	fmt.Fprintf(w, "\tACTIVE? %s\n", wcVal(e.Active))
	fmt.Fprintf(w, "\tSPRITEVOLUMEONLY? %s\n", wcVal(e.SpriteVolumeOnly))
	fmt.Fprintf(w, "\tDMRGBTRACK? \"%s\"\n", wcVal(e.DMRGBTrackTag))
	fmt.Fprintf(w, "\tSPHERE \"%s\"\n", e.SphereTag)
	fmt.Fprintf(w, "\tSPHERERADIUS %0.7f\n", e.SphereRadius)
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", e.UserData)
	fmt.Fprintf(w, "ENDACTORINST\n\n")
	return nil
}

func (e *ActorInst) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.DefinitionTag = records[1]

	records, err = r.ReadProperty("CURRENTACTION?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentAction, records[1])
	if err != nil {
		return fmt.Errorf("current action: %w", err)
	}

	records, err = r.ReadProperty("LOCATION?", 6)
	if err != nil {
		return err
	}
	err = parse(&e.Location, records[1:]...)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = r.ReadProperty("BOUNDINGRADIUS?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("SCALEFACTOR?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Scale, records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	records, err = r.ReadProperty("SOUND?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SoundTag, records[1])
	if err != nil {
		return fmt.Errorf("sound: %w", err)
	}

	records, err = r.ReadProperty("ACTIVE?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Active, records[1])
	if err != nil {
		return fmt.Errorf("active: %w", err)
	}

	records, err = r.ReadProperty("SPRITEVOLUMEONLY?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpriteVolumeOnly, records[1])
	if err != nil {
		return fmt.Errorf("sprite volume only: %w", err)
	}

	records, err = r.ReadProperty("DMRGBTRACK?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.DMRGBTrackTag, records[1])
	if err != nil {
		return fmt.Errorf("dm rgb track: %w", err)
	}

	records, err = r.ReadProperty("SPHERE", 1)
	if err != nil {
		return err
	}
	e.SphereTag = records[1]

	records, err = r.ReadProperty("SPHERERADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SphereRadius, records[1])
	if err != nil {
		return fmt.Errorf("sphere radius: %w", err)
	}

	records, err = r.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}
	e.UserData = records[1]

	_, err = r.ReadProperty("ENDACTORINST", 0)
	if err != nil {
		return err
	}
	return nil
}

func (e *ActorInst) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	var err error
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfActorInst := &rawfrag.WldFragActor{}

	if e.DefinitionTag != "" {
		actorDef := srcWld.ByTag(e.DefinitionTag)
		if actorDef == nil {
			return -1, fmt.Errorf("actor definition %s not found", e.DefinitionTag)
		}

		actorDefRef, err := actorDef.ToRaw(srcWld, dst)
		if err != nil {
			return -1, fmt.Errorf("actor definition %s to raw: %w", e.DefinitionTag, err)
		}

		wfActorInst.ActorDefRef = int32(actorDefRef)
	}

	if e.CurrentAction.Valid {
		wfActorInst.Flags |= 0x01
		wfActorInst.CurrentAction = e.CurrentAction.Uint32
	}

	if e.Location.Valid {
		wfActorInst.Flags |= 0x02
		wfActorInst.Location = e.Location.Float32Slice6
	}

	if e.BoundingRadius.Valid {
		wfActorInst.Flags |= 0x04
		wfActorInst.BoundingRadius = e.BoundingRadius.Float32
	}

	if e.Scale.Valid {
		wfActorInst.Flags |= 0x08
		wfActorInst.ScaleFactor = e.Scale.Float32
	}

	if e.SoundTag.Valid {
		wfActorInst.Flags |= 0x10
		wfActorInst.SoundNameRef = raw.NameAdd(e.SoundTag.String)
	}

	if e.Active.Valid {
		wfActorInst.Flags |= 0x20
	}

	if e.SpriteVolumeOnly.Valid {
		wfActorInst.Flags |= 0x40
	}

	if e.DMRGBTrackTag.Valid {
		dmRGBTrackDef := srcWld.ByTag(e.DMRGBTrackTag.String)
		if dmRGBTrackDef == nil {
			return -1, fmt.Errorf("dm rgb track def %s not found", e.DMRGBTrackTag.String)
		}

		dmRGBTrackRef, err := dmRGBTrackDef.ToRaw(srcWld, dst)
		if err != nil {
			return -1, fmt.Errorf("dm rgb track %s to raw: %w", e.DMRGBTrackTag.String, err)
		}

		wfActorInst.DMRGBTrackRef = int32(dmRGBTrackRef)
	}

	if e.SphereRadius > 0 {
		sphere := &rawfrag.WldFragSphere{
			NameRef: raw.NameAdd(e.SphereTag),
			Radius:  e.SphereRadius,
		}

		dst.Fragments = append(dst.Fragments, sphere)
		wfActorInst.SphereRef = uint32(len(dst.Fragments))
	}

	wfActorInst.UserData = e.UserData

	dst.Fragments = append(dst.Fragments, wfActorInst)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), err
}

// LightDef is a declaration of LIGHTDEF
type LightDef struct {
	fragID       int16
	Tag          string
	Flags        uint32
	CurrentFrame NullUint32
	Sleep        NullUint32
	LightLevels  []float32
	Colors       [][3]float32
}

func (e *LightDef) Definition() string {
	return "LIGHTDEFINITION"
}

func (e *LightDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tCURRENTFRAME? %s\n", wcVal(e.CurrentFrame))
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.LightLevels))
	for _, level := range e.LightLevels {
		fmt.Fprintf(w, "\tLIGHTLEVELS %0.6f\n", level)
	}
	fmt.Fprintf(w, "\tSLEEP? %s\n", wcVal(e.Sleep))
	isSkipFrames := 0
	if e.Flags&0x08 == 0x08 {
		isSkipFrames = 1
	}
	fmt.Fprintf(w, "\tSKIPFRAMES %d\n", isSkipFrames)
	fmt.Fprintf(w, "\tNUMCOLORS %d\n", len(e.Colors))
	for _, color := range e.Colors {
		fmt.Fprintf(w, "\tCOLOR %0.6f %0.6f %0.6f\n", color[0], color[1], color[2])
	}
	fmt.Fprintf(w, "ENDLIGHTDEFINITION\n\n")
	return nil
}

func (e *LightDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("CURRENTFRAME?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentFrame, records[1])
	if err != nil {
		return fmt.Errorf("current frame: %w", err)
	}

	records, err = r.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return err
	}
	numFrames := int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = r.ReadProperty("LIGHTLEVELS", 1)
		if err != nil {
			return err
		}
		level := float32(0)
		err = parse(&level, records[1])
		if err != nil {
			return fmt.Errorf("light level: %w", err)
		}
		e.LightLevels = append(e.LightLevels, level)
	}

	records, err = r.ReadProperty("SLEEP?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = r.ReadProperty("SKIPFRAMES", 1)
	if err != nil {
		return err
	}
	if records[1] == "1" {
		e.Flags |= 0x08
	}

	records, err = r.ReadProperty("NUMCOLORS", 1)
	if err != nil {
		return err
	}
	numColors := int(0)
	err = parse(&numColors, records[1])
	if err != nil {
		return fmt.Errorf("num colors: %w", err)
	}

	for i := 0; i < numColors; i++ {
		records, err = r.ReadProperty("COLOR", 3)
		if err != nil {
			return err
		}
		color := [3]float32{}
		err = parse(&color, records[1:]...)
		if err != nil {
			return fmt.Errorf("color: %w", err)
		}

		e.Colors = append(e.Colors, color)
	}

	_, err = r.ReadProperty("ENDLIGHTDEFINITION", 0)
	if err != nil {
		return err
	}
	return nil
}

func (e *LightDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	var err error

	wfLightDef := &rawfrag.WldFragLightDef{
		NameRef: raw.NameAdd(e.Tag),
		Flags:   e.Flags,
	}

	if e.CurrentFrame.Valid {
		wfLightDef.Flags |= 0x01
		wfLightDef.FrameCurrentRef = e.CurrentFrame.Uint32
	}

	if e.Sleep.Valid {
		wfLightDef.Flags |= 0x02
		wfLightDef.Sleep = e.Sleep.Uint32
	}

	if len(e.LightLevels) > 0 {
		wfLightDef.Flags |= 0x04
		wfLightDef.LightLevels = e.LightLevels
	}
	// TODO: skip_frames 0x08

	if len(e.Colors) > 1 {
		wfLightDef.Flags |= 0x10
		wfLightDef.Colors = e.Colors
	}

	dst.Fragments = append(dst.Fragments, wfLightDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), err
}

// PointLight is a declaration of POINTLIGHT
type PointLight struct {
	fragID     int16
	Tag        string
	LightFlags uint32
	LightTag   string
	Flags      uint32
	Location   [3]float32
	Radius     float32
}

func (e *PointLight) Definition() string {
	return "POINTLIGHT"
}

func (e *PointLight) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", e.Flags)
	fmt.Fprintf(w, "\tXYZ %0.6f %0.6f %0.6f\n", e.Location[0], e.Location[1], e.Location[2])
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", e.LightTag)
	fmt.Fprintf(w, "\t// LIGHTFLAGS %d\n", e.LightFlags)
	fmt.Fprintf(w, "\tRADIUSOFINFLUENCE %0.7f\n", e.Radius)
	fmt.Fprintf(w, "ENDPOINTLIGHT\n\n")
	return nil
}

func (e *PointLight) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("LIGHT", 1)
	if err != nil {
		return err
	}
	e.LightTag = records[1]

	records, err = r.ReadProperty("XYZ", 3)
	if err != nil {
		return err
	}
	err = parse(&e.Location, records[1:]...)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = r.ReadProperty("RADIUSOFINFLUENCE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Radius, records[1])
	if err != nil {
		return fmt.Errorf("radius of influence: %w", err)
	}

	return nil
}

func (e *PointLight) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	return -1, fmt.Errorf("not implemented")
}

// Sprite3DDef is a declaration of SPRITE3DDEF
type Sprite3DDef struct {
	fragID         int16
	Tag            string
	CenterOffset   NullFloat32Slice3
	BoundingRadius NullFloat32
	SphereListTag  string
	Vertices       [][3]float32
	BSPNodes       []*BSPNode
}

// BSPNode is a declaration of BSPNODE
type BSPNode struct {
	Vertices      []uint32
	RenderMethod  string
	Pen           NullUint32
	Brightness    NullFloat32
	ScaledAmbient NullFloat32
	SpriteTag     NullString
	UvOrigin      NullFloat32Slice3
	UAxis         NullFloat32Slice3
	VAxis         NullFloat32Slice3
	Uvs           [][2]float32
	TwoSided      int
	FrontTree     uint32
	BackTree      uint32
}

func (e *Sprite3DDef) Definition() string {
	return "3DSPRITEDEF"
}

func (e *Sprite3DDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tCENTEROFFSET? %s\n", wcVal(e.CenterOffset))
	fmt.Fprintf(w, "\tBOUNDINGRADIUS? %s\n", wcVal(e.BoundingRadius))
	fmt.Fprintf(w, "\tSPHERELIST \"%s\"\n", e.SphereListTag)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMBSPNODES %d\n", len(e.BSPNodes))
	for i, node := range e.BSPNodes {
		fmt.Fprintf(w, "\tBSPNODE //%d\n", i+1)
		fmt.Fprintf(w, "\t\tVERTEXLIST %d", len(node.Vertices))
		for _, vert := range node.Vertices {
			fmt.Fprintf(w, " %d", vert)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\t\tRENDERMETHOD \"%s\"\n", node.RenderMethod)
		fmt.Fprintf(w, "\t\tRENDERINFO\n")
		fmt.Fprintf(w, "\t\t\tPEN? %s\n", wcVal(node.Pen))
		fmt.Fprintf(w, "\t\t\tBRIGHTNESS? %s\n", wcVal(node.Brightness))
		fmt.Fprintf(w, "\t\t\tSCALEDAMBIENT? %s\n", wcVal(node.ScaledAmbient))
		fmt.Fprintf(w, "\t\t\tSPRITE? \"%s\"\n", wcVal(node.SpriteTag))
		fmt.Fprintf(w, "\t\t\tUVORIGIN? %s\n", wcVal(node.UvOrigin))
		fmt.Fprintf(w, "\t\t\tUAXIS? %s\n", wcVal(node.UAxis))
		fmt.Fprintf(w, "\t\t\tVAXIS? %s\n", wcVal(node.VAxis))
		fmt.Fprintf(w, "\t\t\tUVCOUNT %d\n", len(node.Uvs))
		for _, uv := range node.Uvs {
			fmt.Fprintf(w, "\t\t\tUV %s\n", wcVal(uv))
		}
		fmt.Fprintf(w, "\t\t\tTWOSIDED %d\n", node.TwoSided)
		fmt.Fprintf(w, "\t\tENDRENDERINFO\n")
		fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\tENDBSPNODE // %d\n", i+1)
	}
	fmt.Fprintf(w, "END3DSPRITEDEF\n\n")
	return nil
}

func (s *Sprite3DDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	s.Tag = records[1]

	records, err = r.ReadProperty("CENTEROFFSET?", 3)
	if err != nil {
		return err
	}
	err = parse(&s.CenterOffset, records[1:]...)
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = r.ReadProperty("BOUNDINGRADIUS?", 1)
	if err != nil {
		return err
	}
	err = parse(&s.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("SPHERELIST", 1)
	if err != nil {
		return err
	}
	s.SphereListTag = records[1]

	records, err = r.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := int(0)
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert := [3]float32{}
		err = parse(&vert, records[1:]...)
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}

		s.Vertices = append(s.Vertices, vert)
	}

	records, err = r.ReadProperty("NUMBSPNODES", 1)
	if err != nil {
		return err
	}
	numBSPNodes := int(0)
	err = parse(&numBSPNodes, records[1])
	if err != nil {
		return fmt.Errorf("num bsp nodes: %w", err)
	}

	for i := 0; i < numBSPNodes; i++ {
		node := &BSPNode{}
		_, err = r.ReadProperty("BSPNODE", 0)
		if err != nil {
			return err
		}
		records, err = r.ReadProperty("VERTEXLIST", -1)
		if err != nil {
			return err
		}
		numVertices := int(0)
		err = parse(&numVertices, records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}
		if len(records) != numVertices+2 {
			return fmt.Errorf("vertex list: expected %d, got %d", numVertices, len(records)-2)
		}
		for j := 0; j < numVertices; j++ {
			val := uint32(0)
			err = parse(&val, records[j+2])
			if err != nil {
				return fmt.Errorf("vertex %d: %w", j, err)
			}
			node.Vertices = append(node.Vertices, val)
		}

		records, err = r.ReadProperty("RENDERMETHOD", 1)
		if err != nil {
			return err
		}

		node.RenderMethod = records[1]

		_, err = r.ReadProperty("RENDERINFO", 0)
		if err != nil {
			return err
		}

		records, err = r.ReadProperty("PEN?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.Pen, records[1])
		if err != nil {
			return fmt.Errorf("render pen: %w", err)
		}

		records, err = r.ReadProperty("BRIGHTNESS?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.Brightness, records[1])
		if err != nil {
			return fmt.Errorf("render brightness: %w", err)
		}

		records, err = r.ReadProperty("SCALEDAMBIENT?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.ScaledAmbient, records[1])
		if err != nil {
			return fmt.Errorf("render scaled ambient: %w", err)
		}

		records, err = r.ReadProperty("SPRITE?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.SpriteTag, records[1])
		if err != nil {
			return fmt.Errorf("render sprite: %w", err)
		}

		records, err = r.ReadProperty("UVORIGIN?", 3)
		if err != nil {
			return err
		}
		err = parse(&node.UvOrigin, records[1:]...)
		if err != nil {
			return fmt.Errorf("render uv origin: %w", err)
		}

		records, err = r.ReadProperty("UAXIS?", 3)
		if err != nil {
			return err
		}
		err = parse(&node.UAxis, records[1:]...)
		if err != nil {
			return fmt.Errorf("render u axis: %w", err)
		}

		records, err = r.ReadProperty("VAXIS?", 3)
		if err != nil {
			return err
		}
		err = parse(&node.VAxis, records[1:]...)
		if err != nil {
			return fmt.Errorf("render v axis: %w", err)
		}

		records, err = r.ReadProperty("UVCOUNT", 1)
		if err != nil {
			return err
		}
		numUVs := int(0)
		err = parse(&numUVs, records[1])
		if err != nil {
			return fmt.Errorf("num uvs: %w", err)
		}

		for j := 0; j < numUVs; j++ {
			records, err = r.ReadProperty("UV", 2)
			if err != nil {
				return err
			}
			uv := [2]float32{}
			err = parse(&uv, records[1:]...)
			if err != nil {
				return fmt.Errorf("uv %d: %w", j, err)
			}
			node.Uvs = append(node.Uvs, uv)
		}

		records, err = r.ReadProperty("TWOSIDED", 1)
		if err != nil {
			return err
		}
		err = parse(&node.TwoSided, records[1])
		if err != nil {
			return fmt.Errorf("two sided: %w", err)
		}

		_, err = r.ReadProperty("ENDRENDERINFO", 0)
		if err != nil {
			return err
		}

		records, err = r.ReadProperty("FRONTTREE", 1)
		if err != nil {
			return err
		}

		err = parse(&node.FrontTree, records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = r.ReadProperty("BACKTREE", 1)
		if err != nil {
			return err
		}

		err = parse(&node.BackTree, records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = r.ReadProperty("ENDBSPNODE", 0)
		if err != nil {
			return err
		}

		s.BSPNodes = append(s.BSPNodes, node)
	}

	_, err = r.ReadProperty("END3DSPRITEDEF", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *Sprite3DDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfSprite3DDef := &rawfrag.WldFragSprite3DDef{
		Vertices: e.Vertices,
	}

	if e.CenterOffset.Valid {
		wfSprite3DDef.Flags |= 0x01
		wfSprite3DDef.CenterOffset = e.CenterOffset.Float32Slice3
	}

	if e.BoundingRadius.Valid {
		wfSprite3DDef.Flags |= 0x02
		wfSprite3DDef.BoundingRadius = e.BoundingRadius.Float32
	}

	if len(e.BSPNodes) > 0 {

		for _, node := range e.BSPNodes {
			bnode := rawfrag.WldFragThreeDSpriteBspNode{
				FrontTree:     node.FrontTree,
				BackTree:      node.BackTree,
				VertexIndexes: node.Vertices,

				RenderMethod: model.RenderMethodInt(node.RenderMethod),
			}

			if node.Pen.Valid {
				bnode.RenderFlags |= 0x01
				bnode.RenderPen = node.Pen.Uint32
			}

			if node.Brightness.Valid {
				bnode.RenderFlags |= 0x02
				bnode.RenderBrightness = node.Brightness.Float32
			}

			if node.ScaledAmbient.Valid {
				bnode.RenderFlags |= 0x04
				bnode.RenderScaledAmbient = node.ScaledAmbient.Float32
			}

			if node.SpriteTag.Valid {
				bnode.RenderFlags |= 0x08
				bnode.RenderSimpleSpriteReference = uint32(raw.NameAdd(node.SpriteTag.String))
			}

			if node.UvOrigin.Valid {
				bnode.RenderFlags |= 0x10
				bnode.RenderUVInfoOrigin = node.UvOrigin.Float32Slice3
				bnode.RenderUVInfoUAxis = node.UAxis.Float32Slice3
				bnode.RenderUVInfoVAxis = node.VAxis.Float32Slice3
			}

			if len(node.Uvs) > 0 {
				bnode.RenderFlags |= 0x20
				bnode.Uvs = node.Uvs
			}

			wfSprite3DDef.BspNodes = append(wfSprite3DDef.BspNodes, bnode)
		}
	}

	wfSprite3DDef.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfSprite3DDef)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

type PolyhedronDefinition struct {
	fragID         int16
	Tag            string
	Flags          uint32
	BoundingRadius float32
	ScaleFactor    float32
	Vertices       [][3]float32
	Faces          []*PolyhedronDefinitionFace
}

type PolyhedronDefinitionFace struct {
	Vertices []uint32
}

func (e *PolyhedronDefinition) Definition() string {
	return "POLYHEDRONDEFINITION"
}

func (e *PolyhedronDefinition) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", e.Flags)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", e.BoundingRadius)
	fmt.Fprintf(w, "\tSCALEFACTOR %0.7f\n", e.ScaleFactor)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMFACES %d\n", len(e.Faces))
	for i, face := range e.Faces {
		fmt.Fprintf(w, "\tFACE %d\n", i+1)
		fmt.Fprintf(w, "\t\tVERTEXLIST %d", len(face.Vertices))
		for _, vert := range face.Vertices {
			fmt.Fprintf(w, " %d", vert)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\tENDFACE %d\n", i+1)
	}
	fmt.Fprintf(w, "ENDPOLYHEDRONDEFINITION\n\n")
	return nil
}

func (e *PolyhedronDefinition) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("SCALEFACTOR", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ScaleFactor, records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	records, err = r.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}

	numVertices := int(0)
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ", 3)
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

	records, err = r.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}
	numFaces := int(0)
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		face := &PolyhedronDefinitionFace{}
		_, err = r.ReadProperty("FACE", 0)
		if err != nil {
			return err
		}

		records, err = r.ReadProperty("VERTEXLIST", 1)
		if err != nil {
			return err
		}
		numVertices := int(0)
		err = parse(&numVertices, records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}
		if len(records) != numVertices+2 {
			return fmt.Errorf("vertex list: expected %d, got %d", numVertices, len(records)-2)
		}
		for j := 0; j < numVertices; j++ {
			val := uint32(0)
			err = parse(&val, records[j+2])
			if err != nil {
				return fmt.Errorf("vertex %d: %w", j, err)
			}
			face.Vertices = append(face.Vertices, val)
		}

		_, err = r.ReadProperty("ENDFACE", 0)
		if err != nil {
			return err
		}

		e.Faces = append(e.Faces, face)
	}

	_, err = r.ReadProperty("ENDPOLYHEDRONDEFINITION", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *PolyhedronDefinition) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	return -1, fmt.Errorf("TODO: PolyhedronDefinition.ToRaw")
}

type TrackInstance struct {
	fragID        int16
	Tag           string
	DefinitionTag string
	Interpolate   int
	Reverse       int
	Sleep         uint32
}

func (e *TrackInstance) Definition() string {
	return "TRACKINSTANCE"
}

func (e *TrackInstance) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", e.DefinitionTag)
	fmt.Fprintf(w, "\tINTERPOLATE %d\n", e.Interpolate)
	fmt.Fprintf(w, "\tREVERSE %d\n", e.Reverse)
	fmt.Fprintf(w, "\tSLEEP %d\n", e.Sleep)
	fmt.Fprintf(w, "ENDTRACKINSTANCE\n\n")
	return nil
}

func (e *TrackInstance) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.DefinitionTag = records[1]

	records, err = r.ReadProperty("INTERPOLATE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Interpolate, records[1])
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

	records, err = r.ReadProperty("REVERSE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Reverse, records[1])
	if err != nil {
		return fmt.Errorf("reverse: %w", err)
	}

	records, err = r.ReadProperty("SLEEP", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	_, err = r.ReadProperty("ENDTRACKINSTANCE", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *TrackInstance) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	return -1, fmt.Errorf("TODO: TrackInstance.ToRaw")
}

type TrackDef struct {
	fragID          int16
	Tag             string
	Flags           uint32
	FrameTransforms []TrackFrameTransform
}

type TrackFrameTransform struct {
	PositionDenom float32
	Rotation      [3]int16
	Position      [3]float32
}

func (e *TrackDef) Definition() string {
	return "TRACKDEFINITION"
}

func (e *TrackDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", e.Flags)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.FrameTransforms))
	for _, frame := range e.FrameTransforms {
		fmt.Fprintf(w, "\tFRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f\n", frame.PositionDenom, frame.Rotation[0], frame.Rotation[1], frame.Rotation[2], frame.Position[0], frame.Position[1], frame.Position[2])
	}
	fmt.Fprintf(w, "ENDTRACKDEFINITION\n\n")
	return nil
}

func (e *TrackDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	return fmt.Errorf("TODO: TrackDef.Read")
}

func (e *TrackDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	return -1, fmt.Errorf("TODO: TrackDef.ToRaw")
}

type HierarchicalSpriteDef struct {
	fragID         int16
	Tag            string
	Dags           []Dag
	AttachedSkins  []AttachedSkin
	CenterOffset   [3]float32
	BoundingRadius float32
	DMSpriteTag    string
	DagCollision   int
}

type Dag struct {
	Tag       string
	Track     string
	SubDags   []uint32
	SpriteTag string
}

type AttachedSkin struct {
	DMSpriteTag               string
	LinkSkinUpdatesToDagIndex uint32
}

func (e *HierarchicalSpriteDef) Definition() string {
	return "HIERARCHICALSPRITEDEF"
}

func (e *HierarchicalSpriteDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tNUMDAGS %d\n", len(e.Dags))
	for i, dag := range e.Dags {
		fmt.Fprintf(w, "\tDAG // %d\n", i+1)
		fmt.Fprintf(w, "\t\tTAG \"%s\"\n", dag.Tag)
		fmt.Fprintf(w, "\t\tSPRITE \"%s\"\n", dag.SpriteTag)
		fmt.Fprintf(w, "\t\tTRACK \"%s\"\n", dag.Track)
		fmt.Fprintf(w, "\t\tSUBDAGLIST %d", len(dag.SubDags))
		for _, subDag := range dag.SubDags {
			fmt.Fprintf(w, " %d", subDag)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\tENDDAG // %d\n", i+1)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMATTACHEDSKINS %d\n", len(e.AttachedSkins))
	for _, skin := range e.AttachedSkins {
		fmt.Fprintf(w, "\tDMSPRITE \"%s\"\n", skin.DMSpriteTag)
		fmt.Fprintf(w, "\tLINKSKINUPDATESTODAGINDEX %d\n", skin.LinkSkinUpdatesToDagIndex)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\tCENTEROFFSET %0.1f %0.1f %0.1f\n", e.CenterOffset[0], e.CenterOffset[1], e.CenterOffset[2])

	fmt.Fprintf(w, "\tDAGCOLLISION %d\n", e.DagCollision)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7e\n", e.BoundingRadius)

	fmt.Fprintf(w, "ENDHIERARCHICALSPRITEDEF\n\n")
	return nil
}

func (e *HierarchicalSpriteDef) Read(r *AsciiReadToken) error {
	return fmt.Errorf("TODO: HierarchicalSpriteDef.Read")
}

func (e *HierarchicalSpriteDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	return -1, fmt.Errorf("TODO: HierarchicalSpriteDef.ToRaw")
}

type WorldTree struct {
	fragID     int16
	Tag        string
	WorldNodes []*WorldNode
}

type WorldNode struct {
	Normals        [4]float32
	WorldRegionTag string
	FrontTree      uint32
	BackTree       uint32
	Distance       float32
}

func (e *WorldTree) Definition() string {
	return "WORLDTREE"
}

func (e *WorldTree) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tNUMWORLDNODES %d\n", len(e.WorldNodes))
	for i, node := range e.WorldNodes {
		fmt.Fprintf(w, "\tWORLDNODE // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", node.Normals[0], node.Normals[1], node.Normals[2], node.Normals[3])
		fmt.Fprintf(w, "\t\tWORLDREGIONTAG \"%s\"\n", node.WorldRegionTag)
		fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\tENDWORLDNODE // %d\n", i+1)
	}
	fmt.Fprintf(w, "ENDWORLDTREE\n\n")
	return nil
}

func (e *WorldTree) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = r.ReadProperty("NUMWORLDNODES", 1)
	if err != nil {
		return err
	}

	numNodes := int(0)
	err = parse(&numNodes, records[1])
	if err != nil {
		return fmt.Errorf("num world nodes: %w", err)
	}

	for i := 0; i < numNodes; i++ {
		node := &WorldNode{}
		_, err = r.ReadProperty("WORLDNODE", 0)
		if err != nil {
			return err
		}

		records, err = r.ReadProperty("NORMALABCD", 4)
		if err != nil {
			return err
		}
		err = parse(&node.Normals, records[1:]...)
		if err != nil {
			return fmt.Errorf("normals: %w", err)
		}

		records, err = r.ReadProperty("WORLDREGIONTAG", 1)
		if err != nil {
			return err
		}
		node.WorldRegionTag = records[1]

		records, err = r.ReadProperty("FRONTTREE", 1)
		if err != nil {
			return err
		}
		err = parse(&node.FrontTree, records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = r.ReadProperty("BACKTREE", 1)
		if err != nil {
			return err
		}
		err = parse(&node.BackTree, records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = r.ReadProperty("ENDWORLDNODE", 0)
		if err != nil {
			return err
		}

		e.WorldNodes = append(e.WorldNodes, node)

	}

	_, err = r.ReadProperty("ENDWORLDTREE", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *WorldTree) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfWorldTree := &rawfrag.WldFragWorldTree{}

	for _, node := range e.WorldNodes {
		wfNode := rawfrag.WorldTreeNode{
			Normal:    node.Normals,
			RegionRef: -raw.NameAdd(node.WorldRegionTag),
			FrontRef:  int32(node.FrontTree),
			BackRef:   int32(node.BackTree),
		}

		wfWorldTree.Nodes = append(wfWorldTree.Nodes, wfNode)
	}

	wfWorldTree.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfWorldTree)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

type Region struct {
	fragID            int16
	Tag               string
	RegionFog         int
	Gouraud2          int
	EncodedVisibility int
	VisListBytes      int
	AmbientLightTag   string
	RegionVertices    [][3]float32
	RenderVertices    [][3]float32
	Walls             []*Wall
	Obstacles         []*Obstacle
	CuttingObstacles  []*Obstacle
	VisTree           *VisTree
	Sphere            [4]float32
	ReverbVolume      float32
	ReverbOffset      int32
	UserData          string
	SpriteTag         string
}

type Wall struct {
	Normal   [4]float32
	Vertices [][3]float32
}

type Obstacle struct {
	Normal   [4]float32
	Vertices [][3]float32
}

type VisTree struct {
	VisNodes []*VisNode
	VisLists []*VisList
}

type VisNode struct {
	Normal       [4]float32
	VisListIndex uint32
	FrontTree    uint32
	BackTree     uint32
}

type VisList struct {
	Ranges []int8
}

func (e *Region) Definition() string {
	return "REGION"
}

func (e *Region) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tREGIONTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tREVERBVOLUME %0.7f\n", e.ReverbVolume)
	fmt.Fprintf(w, "\tREVERBOFFSET %d\n", e.ReverbOffset)
	fmt.Fprintf(w, "\tREGIONFOG %d\n", e.RegionFog)
	fmt.Fprintf(w, "\tGOURAND2 %d\n", e.Gouraud2)
	fmt.Fprintf(w, "\tENCODEDVISIBILITY %d\n", e.EncodedVisibility)
	fmt.Fprintf(w, "\tVISLISTBYTES %d\n", e.VisListBytes)
	fmt.Fprintf(w, "\tNUMREGIONVERTEX %d\n", len(e.RegionVertices))
	for _, vert := range e.RegionVertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMRENDERVERTICES %d\n", len(e.RenderVertices))
	for _, vert := range e.RenderVertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMWALLS %d\n", len(e.Walls))
	for i, wall := range e.Walls {
		fmt.Fprintf(w, "\tWALL // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", wall.Normal[0], wall.Normal[1], wall.Normal[2], wall.Normal[3])
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(wall.Vertices))
		for _, vert := range wall.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\tENDWALL // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tNUMOBSTACLES %d\n", len(e.Obstacles))
	for i, obs := range e.Obstacles {
		fmt.Fprintf(w, "\tOBSTACLE // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", obs.Normal[0], obs.Normal[1], obs.Normal[2], obs.Normal[3])
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(obs.Vertices))
		for _, vert := range obs.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\tENDOBSTACLE // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tNUMCUTTINGOBSTACLES %d\n", len(e.CuttingObstacles))
	for i, obs := range e.CuttingObstacles {
		fmt.Fprintf(w, "\tCUTTINGOBSTACLE // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", obs.Normal[0], obs.Normal[1], obs.Normal[2], obs.Normal[3])
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(obs.Vertices))
		for _, vert := range obs.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\tENDCUTTINGOBSTACLE // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tVISTREE\n")
	fmt.Fprintf(w, "\t\tNUMVISNODE %d\n", len(e.VisTree.VisNodes))
	for i, node := range e.VisTree.VisNodes {
		fmt.Fprintf(w, "\t\tVISNODE // %d\n", i+1)
		fmt.Fprintf(w, "\t\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", node.Normal[0], node.Normal[1], node.Normal[2], node.Normal[3])
		fmt.Fprintf(w, "\t\t\tVISLISTINDEX %d\n", node.VisListIndex)
		fmt.Fprintf(w, "\t\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\t\tENDVISNODE // %d\n", i+1)
	}
	fmt.Fprintf(w, "\t\tNUMVISIBLELIST %d\n", len(e.VisTree.VisLists))
	for i, list := range e.VisTree.VisLists {
		fmt.Fprintf(w, "\t\tVISLIST // %d\n", i+1)
		fmt.Fprintf(w, "\t\t\tRANGE %d", len(list.Ranges))
		for _, val := range list.Ranges {
			fmt.Fprintf(w, " %d", val)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\t\tENDVISIBLELIST // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tENDVISTREE\n")
	fmt.Fprintf(w, "\tSPHERE %0.7f %0.7f %0.7f %0.7f\n", e.Sphere[0], e.Sphere[1], e.Sphere[2], e.Sphere[3])
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", e.UserData)
	fmt.Fprintf(w, "\tSPRITE \"%s\"\n", e.SpriteTag)
	fmt.Fprintf(w, "ENDREGION\n\n")
	return nil
}

func (e *Region) Read(token *AsciiReadToken) error {
	e.VisTree = &VisTree{}
	records, err := token.ReadProperty("REGIONTAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = token.ReadProperty("REVERBVOLUME", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ReverbVolume, records[1])
	if err != nil {
		return fmt.Errorf("reverb volume: %w", err)
	}

	records, err = token.ReadProperty("REVERBOFFSET", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ReverbOffset, records[1])
	if err != nil {
		return fmt.Errorf("reverb offset: %w", err)
	}

	records, err = token.ReadProperty("REGIONFOG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.RegionFog, records[1])
	if err != nil {
		return fmt.Errorf("region fog: %w", err)
	}

	records, err = token.ReadProperty("GOURAND2", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Gouraud2, records[1])
	if err != nil {
		return fmt.Errorf("gourand2: %w", err)
	}

	records, err = token.ReadProperty("ENCODEDVISIBILITY", 1)
	if err != nil {
		return err
	}
	err = parse(&e.EncodedVisibility, records[1])
	if err != nil {
		return fmt.Errorf("encoded visibility: %w", err)
	}

	records, err = token.ReadProperty("VISLISTBYTES", 1)
	if err != nil {
		return err
	}
	err = parse(&e.VisListBytes, records[1])
	if err != nil {
		return fmt.Errorf("vis list bytes: %w", err)
	}
	if e.VisListBytes != 0 && e.VisListBytes != 1 {
		return fmt.Errorf("vis list bytes: expected 0 or 1, got %d", e.VisListBytes)
	}

	records, err = token.ReadProperty("NUMREGIONVERTEX", 1)
	if err != nil {
		return err
	}
	numVertices := int(0)
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num region vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert := [3]float32{}
		err = parse(&vert, records[1:]...)
		if err != nil {
			return fmt.Errorf("region vertex %d: %w", i, err)
		}
		e.RegionVertices = append(e.RegionVertices, vert)
	}

	records, err = token.ReadProperty("NUMRENDERVERTICES", 1)
	if err != nil {
		return err
	}
	err = parse(&numVertices, records[1])
	if err != nil {
		return fmt.Errorf("num render vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert := [3]float32{}
		err = parse(&vert, records[1:]...)
		if err != nil {
			return fmt.Errorf("render vertex %d: %w", i, err)
		}
		e.RenderVertices = append(e.RenderVertices, vert)
	}

	records, err = token.ReadProperty("NUMWALLS", 1)
	if err != nil {
		return err
	}
	numWalls := int(0)
	err = parse(&numWalls, records[1])
	if err != nil {
		return fmt.Errorf("num walls: %w", err)
	}

	for i := 0; i < numWalls; i++ {
		wall := &Wall{}
		_, err = token.ReadProperty("WALL", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("NORMALABCD", 4)
		if err != nil {
			return err
		}
		err = parse(&wall.Normal, records[1:]...)
		if err != nil {
			return fmt.Errorf("wall normal: %w", err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}
		err = parse(numVertices, records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}
			vert := [3]float32{}
			err = parse(&vert, records[1:]...)
			if err != nil {
				return fmt.Errorf("wall vertex %d: %w", j, err)
			}

			wall.Vertices = append(wall.Vertices, vert)
		}

		_, err = token.ReadProperty("ENDWALL", 0)
		if err != nil {
			return err
		}

		e.Walls = append(e.Walls, wall)
	}

	records, err = token.ReadProperty("NUMOBSTACLES", 1)
	if err != nil {
		return err
	}
	numObstacles := int(0)
	err = parse(&numObstacles, records[1])
	if err != nil {
		return fmt.Errorf("num obstacles: %w", err)
	}

	for i := 0; i < numObstacles; i++ {
		obs := &Obstacle{}
		_, err = token.ReadProperty("OBSTACLE", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("NORMALABCD", 4)
		if err != nil {
			return err
		}
		err = parse(&obs.Normal, records[1:]...)
		if err != nil {
			return fmt.Errorf("obstacle normal: %w", err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}
		err = parse(&numVertices, records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}
			vert := [3]float32{}
			err = parse(&vert, records[1:]...)
			if err != nil {
				return fmt.Errorf("obstacle vertex %d: %w", j, err)
			}

			obs.Vertices = append(obs.Vertices, vert)
		}

		_, err = token.ReadProperty("ENDOBSTACLE", 0)
		if err != nil {
			return err
		}

		e.Obstacles = append(e.Obstacles, obs)
	}

	records, err = token.ReadProperty("NUMCUTTINGOBSTACLES", 1)
	if err != nil {
		return err
	}
	err = parse(&numObstacles, records[1])
	if err != nil {
		return fmt.Errorf("num cutting obstacles: %w", err)
	}

	for i := 0; i < numObstacles; i++ {
		obs := &Obstacle{}
		_, err = token.ReadProperty("CUTTINGOBSTACLE", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("NORMALABCD", 4)
		if err != nil {
			return err
		}

		err = parse(&obs.Normal, records[1:]...)
		if err != nil {
			return fmt.Errorf("cutting obstacle normal: %w", err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}

		err = parse(&numVertices, records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}

			vert := [3]float32{}
			err = parse(&vert, records[1:]...)
			if err != nil {
				return fmt.Errorf("cutting obstacle vertex %d: %w", j, err)
			}

			obs.Vertices = append(obs.Vertices, vert)
		}

		_, err = token.ReadProperty("ENDCUTTINGOBSTACLE", 0)
		if err != nil {
			return err
		}

		e.CuttingObstacles = append(e.CuttingObstacles, obs)
	}

	_, err = token.ReadProperty("VISTREE", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("NUMVISNODE", 1)
	if err != nil {
		return err
	}

	numNodes := int(0)
	err = parse(&numNodes, records[1])
	if err != nil {
		return fmt.Errorf("num vis nodes: %w", err)
	}

	for i := 0; i < numNodes; i++ {
		node := &VisNode{}
		_, err = token.ReadProperty("VISNODE", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("NORMALABCD", 4)
		if err != nil {
			return err
		}

		err = parse(&node.Normal, records[1:]...)
		if err != nil {
			return fmt.Errorf("vis node normal: %w", err)
		}

		records, err = token.ReadProperty("VISLISTINDEX", 1)
		if err != nil {
			return err
		}

		err = parse(&node.VisListIndex, records[1])
		if err != nil {
			return fmt.Errorf("vis list index: %w", err)
		}

		records, err = token.ReadProperty("FRONTTREE", 1)
		if err != nil {
			return err
		}

		err = parse(&node.FrontTree, records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = token.ReadProperty("BACKTREE", 1)
		if err != nil {
			return err
		}

		err = parse(&node.BackTree, records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = token.ReadProperty("ENDVISNODE", 0)
		if err != nil {
			return err
		}

		e.VisTree.VisNodes = append(e.VisTree.VisNodes, node)

	}

	records, err = token.ReadProperty("NUMVISIBLELIST", 1)
	if err != nil {
		return err
	}

	numLists := int(0)
	err = parse(&numLists, records[1])
	if err != nil {
		return fmt.Errorf("num visible lists: %w", err)
	}

	for i := 0; i < numLists; i++ {
		list := &VisList{}
		_, err = token.ReadProperty("VISLIST", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("RANGE", -1)
		if err != nil {
			return err
		}

		numRanges := int(0)
		err = parse(&numRanges, records[1])
		if err != nil {
			return fmt.Errorf("num ranges: %w", err)
		}

		for j := 0; j < numRanges; j++ {
			val := int8(0)
			err = parse(&val, records[j+2])
			if err != nil {
				return fmt.Errorf("range %d: %w", j, err)
			}

			list.Ranges = append(list.Ranges, val)
		}

		_, err = token.ReadProperty("ENDVISIBLELIST", 0)
		if err != nil {
			return err
		}

		e.VisTree.VisLists = append(e.VisTree.VisLists, list)
	}

	_, err = token.ReadProperty("ENDVISTREE", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("SPHERE", 4)
	if err != nil {
		return err
	}

	err = parse(&e.Sphere, records[1:]...)
	if err != nil {
		return fmt.Errorf("sphere: %w", err)
	}

	records, err = token.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}

	e.UserData = records[1]

	records, err = token.ReadProperty("SPRITE", 1)
	if err != nil {
		return err
	}
	e.SpriteTag = records[1]

	_, err = token.ReadProperty("ENDREGION", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *Region) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfRegion := &rawfrag.WldFragRegion{
		RegionVertices: e.RegionVertices,
		Sphere:         e.Sphere,
		ReverbVolume:   e.ReverbVolume,
		ReverbOffset:   e.ReverbOffset,
	}

	if wfRegion.Sphere != [4]float32{0, 0, 0, 0} {
		wfRegion.Flags |= 0x01
	}

	if e.ReverbVolume != 0 {
		wfRegion.Flags |= 0x02
	}

	if e.ReverbOffset != 0 {
		wfRegion.Flags |= 0x04
	}

	if e.RegionFog != 0 {
		wfRegion.Flags |= 0x08
	}

	if e.Gouraud2 != 0 {
		wfRegion.Flags |= 0x10
	}

	if e.EncodedVisibility != 0 {
		wfRegion.Flags |= 0x20
	}

	//0x40

	if e.VisListBytes != 0 {
		wfRegion.Flags |= 0x80
	}

	if len(e.AmbientLightTag) > 0 {
		aLightDef := srcWld.ByTag(e.AmbientLightTag)
		if aLightDef == nil {
			return 0, fmt.Errorf("ambient light def not found: %s", e.AmbientLightTag)
		}

		aLightRef, err := aLightDef.ToRaw(srcWld, dst)
		if err != nil {
			return 0, fmt.Errorf("ambient light def to raw: %w", err)
		}
		wfRegion.AmbientLightRef = int32(aLightRef)
	}

	for _, node := range e.VisTree.VisNodes {
		visNode := rawfrag.VisNode{
			NormalABCD:   node.Normal,
			VisListIndex: node.VisListIndex,
			FrontTree:    node.FrontTree,
			BackTree:     node.BackTree,
		}
		wfRegion.VisNodes = append(wfRegion.VisNodes, visNode)
	}

	for _, list := range e.VisTree.VisLists {
		visList := rawfrag.VisList{}

		for _, rng := range list.Ranges {
			visList.Ranges = append(visList.Ranges, byte(rng))
		}

		wfRegion.VisLists = append(wfRegion.VisLists, visList)
	}

	if e.SpriteTag != "" {
		wfRegion.Flags |= 0x100
		spriteDef := srcWld.ByTag(e.SpriteTag)
		if spriteDef == nil {
			return 0, fmt.Errorf("sprite def not found: %s", e.SpriteTag)
		}

		spriteRef, err := spriteDef.ToRaw(srcWld, dst)
		if err != nil {
			return 0, fmt.Errorf("sprite def to raw: %w", err)
		}
		wfRegion.MeshReference = int32(spriteRef)
	}
	wfRegion.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfRegion)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

type AmbientLight struct {
	fragID     int16
	Tag        string
	LightTag   string
	LightFlags uint32
	Regions    []uint32
}

func (e *AmbientLight) Definition() string {
	return "AMBIENTLIGHT"
}

func (e *AmbientLight) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", e.LightTag)
	fmt.Fprintf(w, "\t// LIGHTFLAGS %d\n", e.LightFlags)
	fmt.Fprintf(w, "\tREGIONLIST %d", len(e.Regions))
	for _, region := range e.Regions {
		fmt.Fprintf(w, " %d", region)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "ENDAMBIENTLIGHT\n\n")
	return nil
}

func (e *AmbientLight) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}

	e.Tag = records[1]

	records, err = r.ReadProperty("LIGHT", 1)
	if err != nil {
		return err
	}

	e.LightTag = records[1]

	records, err = r.ReadProperty("REGIONLIST", -1)
	if err != nil {
		return err
	}

	numRegions := int(0)
	err = parse(&numRegions, records[1])
	if err != nil {
		return fmt.Errorf("num regions: %w", err)
	}

	for i := 0; i < numRegions; i++ {
		val := uint32(0)
		err = parse(&val, records[i+2])
		if err != nil {
			return fmt.Errorf("region %d: %w", i, err)
		}

		e.Regions = append(e.Regions, val)
	}

	_, err = r.ReadProperty("ENDAMBIENTLIGHT", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *AmbientLight) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfAmbientLight := &rawfrag.WldFragAmbientLight{
		Regions: e.Regions,
	}

	if len(e.LightTag) > 0 {
		lightDef := srcWld.ByTag(e.LightTag)
		if lightDef == nil {
			return 0, fmt.Errorf("light def not found: %s", e.LightTag)
		}

		lightDefRef, err := lightDef.ToRaw(srcWld, dst)
		if err != nil {
			return 0, fmt.Errorf("light def to raw: %w", err)
		}
		wfAmbientLight.LightRef = int32(lightDefRef)

		wfLight := &rawfrag.WldFragLight{
			//NameRef: ,
			LightDefRef: int32(lightDefRef),
			Flags:       e.LightFlags,
		}

		dst.Fragments = append(dst.Fragments, wfLight)

		wfAmbientLight.LightRef = int32(len(dst.Fragments))
	}

	wfAmbientLight.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfAmbientLight)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

type Zone struct {
	fragID   int16
	Tag      string
	Regions  []uint32
	UserData string
}

func (e *Zone) Definition() string {
	return "ZONE"
}

func (e *Zone) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tREGIONLIST %d", len(e.Regions))
	for _, region := range e.Regions {
		fmt.Fprintf(w, " %d", region)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", e.UserData)
	fmt.Fprintf(w, "ENDZONE\n\n")
	return nil
}

func (e *Zone) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}

	e.Tag = records[1]

	records, err = r.ReadProperty("REGIONLIST", -1)
	if err != nil {
		return err
	}

	numRegions := int(0)
	err = parse(&numRegions, records[1])
	if err != nil {
		return fmt.Errorf("num regions: %w", err)
	}

	for i := 0; i < numRegions; i++ {
		val := uint32(0)
		err = parse(&val, records[i+2])
		if err != nil {
			return fmt.Errorf("region %d: %w", i, err)
		}

		e.Regions = append(e.Regions, val)
	}

	records, err = r.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}

	e.UserData = records[1]

	_, err = r.ReadProperty("ENDZONE", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *Zone) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfZone := &rawfrag.WldFragZone{
		NameRef:  raw.NameAdd(e.Tag),
		Flags:    0,
		Regions:  e.Regions,
		UserData: e.UserData,
	}

	dst.Fragments = append(dst.Fragments, wfZone)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

type RGBTrackDef struct {
	fragID int16
	Tag    string
	Data1  uint32
	Data2  uint32
	Data4  uint32
	Sleep  uint32
	RGBAs  [][4]uint8
}

func (e *RGBTrackDef) Definition() string {
	return "RGBDEFORMATIONTRACKDEF"
}

func (e *RGBTrackDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\t// NUMFRAMES %d // if this isn't 1, let xack know\n", e.Data1)
	fmt.Fprintf(w, "\t// DATA2 %d // if this isn't 1, let xack know\n", e.Data2)
	fmt.Fprintf(w, "\t// NUMVERTICES %d // // if this isn't 0, let xack know\n", e.Data4)
	fmt.Fprintf(w, "\tSLEEP %d\n", e.Sleep)
	fmt.Fprintf(w, "\tRGBDEFORMATIONFRAME")
	fmt.Fprintf(w, "\t\tNUMRGBAS %d\n", len(e.RGBAs))
	for i, rgba := range e.RGBAs {
		fmt.Fprintf(w, "\t\tRGBA %d %d %d %d %d\n", i+1, rgba[0], rgba[1], rgba[2], rgba[3])
	}
	fmt.Fprintf(w, "\tENDRGBDEFORMATIONFRAME\n")
	fmt.Fprintf(w, "ENDRGBDEFORMATIONTRACKDEF\n\n")
	return nil
}

func (e *RGBTrackDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	/* records, err = token.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return err
	}
	r.Data1, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	} */

	/* 	records, err = token.ReadProperty("DATA2", 1)
	   	if err != nil {
	   		return err
	   	}
	   	r.Data2, err = helper.ParseUint32(records[1])
	   	if err != nil {
	   		return fmt.Errorf("data2: %w", err)
	   	} */

	/* records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	r.Data4, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	} */

	records, err = token.ReadProperty("SLEEP", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	_, err = token.ReadProperty("RGBDEFORMATIONFRAME", 0)
	if err != nil {
		return err
	}

	for {

		records, err = token.ReadProperty("NUMRGBAS", 1)
		if err != nil {
			return err
		}

		numRGBAs := int(0)
		err = parse(&numRGBAs, records[1])
		if err != nil {
			return fmt.Errorf("num rgbas: %w", err)
		}

		for i := 0; i < numRGBAs; i++ {
			records, err = token.ReadProperty("RGBA", 4)
			if err != nil {
				return err
			}

			rgba := [4]uint8{}

			err = parse(&rgba[0], records[1:]...)
			if err != nil {
				return fmt.Errorf("rgba: %w", err)
			}
			e.RGBAs = append(e.RGBAs, rgba)
		}

		if records[1] == "ENDRGBDEFORMATIONFRAME" {
			break
		}
	}

	_, err = token.ReadProperty("ENDRGBDEFORMATIONTRACKDEF", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *RGBTrackDef) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfRGBTrack := &rawfrag.WldFragDmRGBTrackDef{
		RGBAs: e.RGBAs,
	}

	wfRGBTrack.NameRef = raw.NameAdd(e.Tag)

	dst.Fragments = append(dst.Fragments, wfRGBTrack)
	e.fragID = int16(len(dst.Fragments))
	return int16(len(dst.Fragments)), nil
}

// RGBTrack is a track instance for RGB deformation tracks
type RGBTrack struct {
	fragID        int16
	Tag           string
	DefinitionTag string
	Flags         uint32
}

func (e *RGBTrack) Definition() string {
	return "RGBDEFORMATIONTRACKINSTANCE"
}

func (e *RGBTrack) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", e.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", e.DefinitionTag)
	fmt.Fprintf(w, "\tFLAGS %d\n", e.Flags)
	fmt.Fprintf(w, "ENDRGBDEFORMATIONTRACKINSTANCE\n\n")
	return nil
}

func (e *RGBTrack) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.Tag = records[1]

	records, err = token.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.DefinitionTag = records[1]

	records, err = token.ReadProperty("FLAGS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Flags, records[1])
	if err != nil {
		return fmt.Errorf("flags: %w", err)
	}

	_, err = token.ReadProperty("ENDRGBDEFORMATIONTRACKINSTANCE", 0)
	if err != nil {
		return err
	}

	return nil
}

func (e *RGBTrack) ToRaw(srcWld *Wld, dst *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	return -1, fmt.Errorf("rgb track not implemented")
}
