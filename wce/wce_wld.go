package wce

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

// WorldDef stores data about the world itself
type WorldDef struct {
	NewWorld   int
	Zone       int
	EqgVersion NullInt8
}

// Definition returns the definition of the WorldDef
func (e *WorldDef) Definition() string {
	return "WORLDDEF"
}

// Write writes the WorldDef to the writer
func (e *WorldDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tNEWWORLD %d\n", e.NewWorld)
	fmt.Fprintf(w, "\tZONE %d\n", e.Zone)
	fmt.Fprintf(w, "\tEQGVERSION? %s\n", wcVal(e.EqgVersion))
	fmt.Fprintf(w, "\n")
	return nil
}

// Read reads the WorldDef from the reader
func (e *WorldDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("NEWWORLD", 1)
	if err != nil {
		return err
	}
	err = parse(&e.NewWorld, records[1])
	if err != nil {
		return fmt.Errorf("newworld: %w", err)
	}

	records, err = token.ReadProperty("ZONE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Zone, records[1])
	if err != nil {
		return fmt.Errorf("zone: %w", err)
	}

	records, err = token.ReadProperty("EQGVERSION?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.EqgVersion, records[1])
	if err != nil {
		return fmt.Errorf("eqgversion: %w", err)
	}

	return nil
}

// GlobalAmbientLightDef is a declaration of GLOBALAMBIENTLIGHTDEF
type GlobalAmbientLightDef struct {
	fragID int16
	Color  [4]uint8
}

func (e *GlobalAmbientLightDef) Definition() string {
	return "GLOBALAMBIENTLIGHTDEF"
}

func (e *GlobalAmbientLightDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", e.Definition())
	fmt.Fprintf(w, "\tCOLOR %d %d %d %d\n", e.Color[0], e.Color[1], e.Color[2], e.Color[3])
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *GlobalAmbientLightDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("Color", 4)
	if err != nil {
		return err
	}
	err = parse(&e.Color, records[1:]...)
	if err != nil {
		return fmt.Errorf("color: %w", err)
	}

	return nil
}

func (e *GlobalAmbientLightDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfGlobalAmbientLightDef := &rawfrag.WldFragGlobalAmbientLightDef{
		Color: e.Color,
	}

	rawWld.Fragments = append(rawWld.Fragments, wfGlobalAmbientLightDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *GlobalAmbientLightDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragGlobalAmbientLightDef) error {
	if wce.GlobalAmbientLightDef != nil {
		return fmt.Errorf("duplicate globalambientlightdef found")
	}
	e.Color = frag.Color
	wce.GlobalAmbientLightDef = e

	return nil
}

// DMSpriteDef2 is a declaration of DMSpriteDef2
type DMSpriteDef2 struct {
	fragID                int16
	Tag                   string
	DmTrackTag            string
	Params2               [3]uint32
	BoundingBoxMin        [3]float32
	BoundingBoxMax        [3]float32
	CenterOffset          [3]float32
	Vertices              [][3]float32
	UVs                   [][2]float32
	VertexNormals         [][3]float32
	VertexColors          [][4]uint8
	SkinAssignmentGroups  [][2]int16
	MaterialPaletteTag    string
	Faces                 []*Face
	MeshOps               []*MeshOp
	FaceMaterialGroups    [][2]uint16
	VertexMaterialGroups  [][2]int16
	BoundingRadius        float32
	FPScale               uint16
	PolyhedronTag         string
	HexOneFlag            uint16
	HexTwoFlag            uint16
	HexFourThousandFlag   uint16
	HexEightThousandFlag  uint16
	HexTenThousandFlag    uint16
	HexTwentyThousandFlag uint16
}

type Face struct {
	Passable int
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

func (e *DMSpriteDef2) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if token.TagIsWritten(e.Tag) {
		return nil
	}
	token.TagSetIsWritten(e.Tag)

	if e.MaterialPaletteTag != "" {
		palette := token.wce.ByTag(e.MaterialPaletteTag)
		if palette == nil {
			return fmt.Errorf("material palette %s not found", e.MaterialPaletteTag)
		}
		err = palette.Write(token)
		if err != nil {
			return fmt.Errorf("material palette %s: %w", e.MaterialPaletteTag, err)
		}
	}

	if e.DmTrackTag != "" {
		dmTrack := token.wce.ByTag(e.DmTrackTag)
		if dmTrack == nil {
			return fmt.Errorf("dmtrack %s not found", e.DmTrackTag)
		}
		switch dmTrackDef := dmTrack.(type) {
		case *DMTrackDef2:
			err = dmTrackDef.Write(token)
			if err != nil {
				return fmt.Errorf("dmtrack %s: %w", dmTrackDef.Tag, err)
			}
		default:
			return fmt.Errorf("dmtrack %s unknown type %T", e.DmTrackTag, dmTrack)
		}
	}

	if e.PolyhedronTag != "" && e.PolyhedronTag != "NEGATIVE_TWO" && e.PolyhedronTag != "SPECIAL_COLLISION" {
		poly := token.wce.ByTag(e.PolyhedronTag)
		if poly == nil {
			return fmt.Errorf("polyhedron %s not found", e.PolyhedronTag)
		}
		switch polyDef := poly.(type) {
		case *PolyhedronDefinition:
			err = polyDef.Write(token)
			if err != nil {
				return fmt.Errorf("polyhedron %s: %w", polyDef.Tag, err)
			}
		case *Sprite3DDef:
			err = polyDef.Write(token)
			if err != nil {
				return fmt.Errorf("sprite 3d %s: %w", polyDef.Tag, err)
			}
		default:
			return fmt.Errorf("polyhedron %s unknown type %T", e.PolyhedronTag, poly)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tCENTEROFFSET %0.8e %0.8e %0.8e\n", e.CenterOffset[0], e.CenterOffset[1], e.CenterOffset[2])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMUVS %d\n", len(e.UVs))
	for _, uv := range e.UVs {
		fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", uv[0], uv[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTEXNORMALS %d\n", len(e.VertexNormals))
	for _, vn := range e.VertexNormals {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vn[0], vn[1], vn[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTEXCOLORS %d\n", len(e.VertexColors))
	for _, color := range e.VertexColors {
		fmt.Fprintf(w, "\t\tRGBA %d %d %d %d\n", color[0], color[1], color[2], color[3])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tSKINASSIGNMENTGROUPS %d", len(e.SkinAssignmentGroups))
	for _, sa := range e.SkinAssignmentGroups {
		fmt.Fprintf(w, " %d %d", sa[0], sa[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tMATERIALPALETTE \"%s\"\n", e.MaterialPaletteTag)
	fmt.Fprintf(w, "\tDMTRACKINST \"%s\"\n", e.DmTrackTag)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tPOLYHEDRON\n")
	fmt.Fprintf(w, "\t\tDEFINITION \"%s\"\n", e.PolyhedronTag)
	fmt.Fprintf(w, "\tNUMFACE2S %d\n", len(e.Faces))
	for i, face := range e.Faces {
		fmt.Fprintf(w, "\t\tDMFACE2 //%d\n", i)
		fmt.Fprintf(w, "\t\t\tPASSABLE %d\n", face.Passable)
		fmt.Fprintf(w, "\t\t\tTRIANGLE %d %d %d\n", face.Triangle[0], face.Triangle[1], face.Triangle[2])
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\tNUMMESHOPS %d\n", len(e.MeshOps))
	for _, meshOp := range e.MeshOps {
		fmt.Fprintf(w, "\tMESHOP %d %d %0.8f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
	}
	fmt.Fprintf(w, "\n")
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
	fmt.Fprintf(w, "\tBOUNDINGBOXMIN %0.8e %0.8e %0.8e\n", e.BoundingBoxMin[0], e.BoundingBoxMin[1], e.BoundingBoxMin[2])
	fmt.Fprintf(w, "\tBOUNDINGBOXMAX %0.8e %0.8e %0.8e\n", e.BoundingBoxMax[0], e.BoundingBoxMax[1], e.BoundingBoxMax[2])

	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.8e\n", e.BoundingRadius)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFPSCALE %d\n", e.FPScale)
	fmt.Fprintf(w, "\tHEXONEFLAG %d\n", e.HexOneFlag)
	fmt.Fprintf(w, "\tHEXTWOFLAG %d\n", e.HexTwoFlag)
	fmt.Fprintf(w, "\tHEXFOURTHOUSANDFLAG %d\n", e.HexFourThousandFlag)
	fmt.Fprintf(w, "\tHEXEIGHTTHOUSANDFLAG %d\n", e.HexEightThousandFlag)
	fmt.Fprintf(w, "\tHEXTENTHOUSANDFLAG %d\n", e.HexTenThousandFlag)
	fmt.Fprintf(w, "\tHEXTWENTYTHOUSANDFLAG %d\n", e.HexTwentyThousandFlag)

	fmt.Fprintf(w, "\n")
	return nil
}

func (e *DMSpriteDef2) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("CENTEROFFSET", 3)
	if err != nil {
		return err
	}
	err = parse(&e.CenterOffset, records[1:]...)
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := int(0)
	err = parse(&numVertices, records[1])
	if err != nil {
		return err
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
	numUVs := int(0)
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

	records, err = token.ReadProperty("NUMVERTEXNORMALS", 1)
	if err != nil {
		return err
	}
	numNormals := int(0)
	err = parse(&numNormals, records[1])
	if err != nil {
		return fmt.Errorf("num normals: %w", err)
	}

	for i := 0; i < numNormals; i++ {
		records, err = token.ReadProperty("XYZ", 3)
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

	records, err = token.ReadProperty("NUMVERTEXCOLORS", 1)
	if err != nil {
		return err
	}
	numColors := int(0)
	err = parse(&numColors, records[1])
	if err != nil {
		return fmt.Errorf("num colors: %w", err)
	}

	for i := 0; i < numColors; i++ {
		records, err = token.ReadProperty("RGBA", 4)
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

	records, err = token.ReadProperty("SKINASSIGNMENTGROUPS", -1)
	if err != nil {
		return err
	}
	numSkinAssignments := int(0)
	err = parse(&numSkinAssignments, records[1])
	if err != nil {
		return fmt.Errorf("num skin assignments: %w", err)
	}

	for i := 0; i < numSkinAssignments*2; i++ {
		val1 := int16(0)
		err = parse(&val1, records[i+2])
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}

		val2 := int16(0)
		err = parse(&val2, records[i+3])
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}
		e.SkinAssignmentGroups = append(e.SkinAssignmentGroups, [2]int16{int16(val1), int16(val2)})
		i++
	}

	records, err = token.ReadProperty("MATERIALPALETTE", 1)
	if err != nil {
		return err
	}
	e.MaterialPaletteTag = records[1]

	records, err = token.ReadProperty("DMTRACKINST", 1)
	if err != nil {
		return err
	}
	e.DmTrackTag = records[1]

	_, err = token.ReadProperty("POLYHEDRON", 0)
	if err != nil {
		return err
	}
	records, err = token.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.PolyhedronTag = records[1]

	records, err = token.ReadProperty("NUMFACE2S", 1)
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
		_, err = token.ReadProperty("DMFACE2", 0)
		if err != nil {
			return err
		}
		records, err = token.ReadProperty("PASSABLE", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Passable, records[1])
		if err != nil {
			return fmt.Errorf("face %d hex ten flag: %w", i, err)
		}

		records, err = token.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		err = parse(&face.Triangle, records[1:]...)
		if err != nil {
			return fmt.Errorf("face %d triangle: %w", i, err)
		}

		e.Faces = append(e.Faces, face)
	}

	records, err = token.ReadProperty("NUMMESHOPS", 1)
	if err != nil {
		return err
	}
	numMeshOps := int(0)
	err = parse(&numMeshOps, records[1])
	if err != nil {
		return fmt.Errorf("num mesh ops: %w", err)
	}

	for i := 0; i < numMeshOps; i++ {
		meshOp := &MeshOp{}
		records, err = token.ReadProperty("MESHOP", 5)
		if err != nil {
			return err
		}
		err = parse(&meshOp.Index1, records[1])
		if err != nil {
			return fmt.Errorf("mesh op %d index1: %w", i, err)
		}
		err = parse(&meshOp.Index2, records[2])
		if err != nil {
			return fmt.Errorf("mesh op %d index2: %w", i, err)
		}
		err = parse(&meshOp.Offset, records[3])
		if err != nil {
			return fmt.Errorf("mesh op %d offset: %w", i, err)
		}
		err = parse(&meshOp.Param1, records[4])
		if err != nil {
			return fmt.Errorf("mesh op %d param1: %w", i, err)
		}
		err = parse(&meshOp.TypeField, records[5])
		if err != nil {
			return fmt.Errorf("mesh op %d typefield: %w", i, err)
		}
		e.MeshOps = append(e.MeshOps, meshOp)
	}

	records, err = token.ReadProperty("FACEMATERIALGROUPS", -1)
	if err != nil {
		return err
	}
	numFaceMaterialGroups := int(0)
	err = parse(&numFaceMaterialGroups, records[1])
	if err != nil {
		return fmt.Errorf("num face material groups: %w", err)
	}

	for i := 0; i < numFaceMaterialGroups*2; i++ {
		val1, err := strconv.ParseUint(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		val2, err := strconv.ParseUint(records[i+3], 10, 16)
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		e.FaceMaterialGroups = append(e.FaceMaterialGroups, [2]uint16{uint16(val1), uint16(val2)})
		i++
	}

	records, err = token.ReadProperty("VERTEXMATERIALGROUPS", -1)
	if err != nil {
		return err
	}
	numVertexMaterialGroups := int(0)
	err = parse(&numVertexMaterialGroups, records[1])
	if err != nil {
		return fmt.Errorf("num vertex material groups: %w", err)
	}

	for i := 0; i < numVertexMaterialGroups*2; i++ {
		val1, err := strconv.ParseInt(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		val2, err := strconv.ParseInt(records[i+3], 10, 16)
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		e.VertexMaterialGroups = append(e.VertexMaterialGroups, [2]int16{int16(val1), int16(val2)})
		i++
	}

	records, err = token.ReadProperty("BOUNDINGBOXMIN", 3)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingBoxMin, records[1:]...)
	if err != nil {
		return fmt.Errorf("bounding box min: %w", err)
	}

	records, err = token.ReadProperty("BOUNDINGBOXMAX", 3)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingBoxMax, records[1:]...)
	if err != nil {
		return fmt.Errorf("bounding box max: %w", err)
	}

	records, err = token.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = token.ReadProperty("FPSCALE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.FPScale, records[1])
	if err != nil {
		return fmt.Errorf("fpscale: %w", err)
	}

	records, err = token.ReadProperty("HEXONEFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexOneFlag, records[1])
	if err != nil {
		return fmt.Errorf("hexoneflag: %w", err)
	}

	records, err = token.ReadProperty("HEXTWOFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwoFlag, records[1])
	if err != nil {
		return fmt.Errorf("hextwoflag: %w", err)
	}

	records, err = token.ReadProperty("HEXFOURTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexFourThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hexfourthousandflag: %w", err)
	}

	records, err = token.ReadProperty("HEXEIGHTTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexEightThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hexeightthousandflag: %w", err)
	}

	records, err = token.ReadProperty("HEXTENTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTenThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hextenthousandflag: %w", err)
	}

	records, err = token.ReadProperty("HEXTWENTYTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwentyThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hextwentythousandflag: %w", err)
	}

	return nil
}

func (e *DMSpriteDef2) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	var err error

	if e.fragID != 0 {
		return e.fragID, nil
	}

	materialPaletteRef := int16(0)
	if e.MaterialPaletteTag != "" {
		palette := wce.ByTag(e.MaterialPaletteTag)
		if palette == nil {
			return -1, fmt.Errorf("material palette %s not found", e.MaterialPaletteTag)
		}

		materialPaletteRef, err = palette.ToRaw(wce, rawWld)
		if err != nil {
			return -1, fmt.Errorf("material palette %s to raw: %w", e.MaterialPaletteTag, err)
		}
	}

	dmSpriteDef := &rawfrag.WldFragDmSpriteDef2{
		MaterialPaletteRef:   uint32(materialPaletteRef),
		CenterOffset:         e.CenterOffset,
		Params2:              e.Params2,
		BoundingRadius:       e.BoundingRadius,
		BoundingBoxMin:       e.BoundingBoxMin,
		BoundingBoxMax:       e.BoundingBoxMax,
		Scale:                e.FPScale,
		Colors:               e.VertexColors,
		FaceMaterialGroups:   e.FaceMaterialGroups,
		VertexMaterialGroups: e.VertexMaterialGroups,
	}

	if e.DmTrackTag != "" {
		dmTrackDef := wce.ByTag(e.DmTrackTag)
		if dmTrackDef == nil {
			return -1, fmt.Errorf("dmtrackdef %s not found", e.DmTrackTag)
		}

		switch dmTrack := dmTrackDef.(type) {
		case *DMTrackDef2:
			dmTrackRef, err := dmTrack.ToRaw(wce, rawWld)
			if err != nil {
				return -1, fmt.Errorf("dmtrackdef %s to raw: %w", e.DmTrackTag, err)
			}
			dmSpriteDef.DMTrackRef = int32(dmTrackRef)
		default:
			return -1, fmt.Errorf("dmtrackdef %s unknown type %T", e.DmTrackTag, dmTrackDef)
		}
	}

	if e.PolyhedronTag != "" { //&& (!strings.HasPrefix(e.Tag, "R") || !wld.isZone)
		if strings.HasPrefix(e.Tag, "R") && wce.WorldDef.Zone == 1 {
			if e.PolyhedronTag == "NEGATIVE_TWO" {
				dmSpriteDef.Fragment4Ref = -2
			}
			if dmSpriteDef.Fragment4Ref != -2 {
				return -1, fmt.Errorf("zone region polyhedron should be NEGATIVE_TWO, not %s", e.PolyhedronTag)
			}
			/* 			for i, frag := range rawWld.Fragments {
				_, ok := frag.(*rawfrag.WldFragBMInfo)
				if !ok {
					continue
				}
				dmSpriteDef.Fragment4Ref = int32(i) + 1
				break
			} */
		} else {
			if e.PolyhedronTag == "NEGATIVE_TWO" {
				dmSpriteDef.Fragment4Ref = -2
			} else if e.PolyhedronTag == "SPECIAL_COLLISION" {
			} else {
				polyhedronFrag := wce.ByTag(e.PolyhedronTag)
				if polyhedronFrag == nil {
					return -1, fmt.Errorf("polyhedron %s not found", e.PolyhedronTag)
				}

				switch polyhedron := polyhedronFrag.(type) {
				case *PolyhedronDefinition:

					polyhedronRef, err := polyhedron.ToRaw(wce, rawWld)
					if err != nil {
						return -1, fmt.Errorf("polyhedron %s to raw: %w", e.PolyhedronTag, err)
					}

					wfPoly := &rawfrag.WldFragPolyhedron{
						FragmentRef: int32(polyhedronRef),
					}
					rawWld.Fragments = append(rawWld.Fragments, wfPoly)

					dmSpriteDef.Fragment4Ref = int32(len(rawWld.Fragments))
				default:
					return -1, fmt.Errorf("polyhedrontag %T unhandled", polyhedron)
				}
			}
		}

		if dmSpriteDef.Fragment4Ref == 0 {
			return -1, fmt.Errorf("polyhedron polygon %s not found", e.PolyhedronTag)
		}
	}

	if e.HexOneFlag != 0 {
		dmSpriteDef.Flags |= 0x1
	}
	if e.HexTwoFlag != 0 {
		dmSpriteDef.Flags |= 0x2
	}
	if e.HexFourThousandFlag != 0 {
		dmSpriteDef.Flags |= 0x4000
	}
	if e.HexEightThousandFlag != 0 {
		dmSpriteDef.Flags |= 0x8000
	}
	if e.HexTenThousandFlag != 0 {
		dmSpriteDef.Flags |= 0x10000
	}
	if e.HexTwentyThousandFlag != 0 {
		dmSpriteDef.Flags |= 0x20000
	}

	dmSpriteDef.NameRef = rawWld.NameAdd(e.Tag)

	/* for i, frag := range rawWld.Fragments {
		_, ok := frag.(*rawfrag.WldFragBMInfo)
		if !ok {
			continue
		}
		dmSpriteDef.Fragment4Ref = int32(i) + 1
	} */

	scale := float32(1 / float32(int(1)<<int(e.FPScale)))

	for _, vert := range e.Vertices {
		dmSpriteDef.Vertices = append(dmSpriteDef.Vertices, [3]int16{
			int16(vert[0] / scale),
			int16(vert[1] / scale),
			int16(vert[2] / scale),
		})
	}

	for _, uv := range e.UVs {
		if wce.WorldDef.NewWorld > 0 {
			dmSpriteDef.UVs = append(dmSpriteDef.UVs, [2]float32{
				float32(uv[0]),
				float32(uv[1]),
			})
		} else {
			dmSpriteDef.UVs = append(dmSpriteDef.UVs, [2]float32{
				float32(uv[0] * 256),
				float32(uv[1] * 256),
			})
		}
	}

	for _, normal := range e.VertexNormals {
		dmSpriteDef.VertexNormals = append(dmSpriteDef.VertexNormals, [3]int8{
			int8(normal[0] * 128),
			int8(normal[1] * 128),
			int8(normal[2] * 128),
		})
	}

	dmSpriteDef.Colors = e.VertexColors
	for _, face := range e.Faces {
		wfFace := &rawfrag.WldFragMeshFaceEntry{
			Index: face.Triangle,
		}
		if face.Passable != 0 {
			wfFace.Flags |= 0x10
		}

		dmSpriteDef.Faces = append(dmSpriteDef.Faces, *wfFace)
	}

	dmSpriteDef.FaceMaterialGroups = e.FaceMaterialGroups
	dmSpriteDef.SkinAssignmentGroups = e.SkinAssignmentGroups
	dmSpriteDef.VertexMaterialGroups = e.VertexMaterialGroups

	for _, meshOp := range e.MeshOps {
		dmSpriteDef.MeshOps = append(dmSpriteDef.MeshOps, rawfrag.WldFragMeshOpEntry{
			Index1:    meshOp.Index1,
			Index2:    meshOp.Index2,
			Offset:    meshOp.Offset,
			Param1:    meshOp.Param1,
			TypeField: meshOp.TypeField,
		})
	}

	rawWld.Fragments = append(rawWld.Fragments, dmSpriteDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *DMSpriteDef2) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragDmSpriteDef2) error {
	if frag == nil {
		return fmt.Errorf("frag is not dmspritedef2 (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	if frag.MaterialPaletteRef > 0 {
		if len(rawWld.Fragments) < int(frag.MaterialPaletteRef) {
			return fmt.Errorf("materialpalette ref %d out of bounds", frag.MaterialPaletteRef)
		}
		materialPalette, ok := rawWld.Fragments[frag.MaterialPaletteRef].(*rawfrag.WldFragMaterialPalette)
		if !ok {
			return fmt.Errorf("materialpalette ref %d not found", frag.MaterialPaletteRef)
		}
		e.MaterialPaletteTag = rawWld.Name(materialPalette.NameRef)
	}

	if frag.DMTrackRef != 0 {
		if len(rawWld.Fragments) < int(frag.DMTrackRef) {
			return fmt.Errorf("dmtrack ref %d out of bounds", frag.DMTrackRef)
		}
		dmTrack, ok := rawWld.Fragments[frag.DMTrackRef].(*rawfrag.WldFragDMTrack)
		if !ok {
			return fmt.Errorf("dmtrack ref %d not valid", frag.DMTrackRef)
		}
		if len(rawWld.Fragments) < int(dmTrack.TrackRef) {
			return fmt.Errorf("dmtrack name ref %d out of bounds", dmTrack.TrackRef)
		}
		dmTrackDef, ok := rawWld.Fragments[dmTrack.TrackRef].(*rawfrag.WldFragDmTrackDef2)
		if !ok {
			return fmt.Errorf("dmtrackdef2 name ref %d not valid", dmTrack.TrackRef)
		}
		e.DmTrackTag = rawWld.Name(dmTrackDef.NameRef)
	}

	if frag.Fragment4Ref != 0 {
		if frag.Fragment4Ref == -2 {
			e.PolyhedronTag = "NEGATIVE_TWO"
		} else {
			if len(rawWld.Fragments) < int(frag.Fragment4Ref) {
				return fmt.Errorf("fragment4 (bminfo) ref %d out of bounds", frag.Fragment4Ref)
			}
			frag4 := rawWld.Fragments[frag.Fragment4Ref]
			switch frag4Def := frag4.(type) {
			case *rawfrag.WldFragPolyhedron:
				if len(rawWld.Fragments) < int(frag4Def.FragmentRef) {
					return fmt.Errorf("fragment4 (polygon) ref %d out of bounds", frag4Def.FragmentRef)
				}

				frag4 = rawWld.Fragments[frag4Def.FragmentRef]
				switch frag4Def := frag4.(type) {
				case *rawfrag.WldFragPolyhedronDef:
					e.PolyhedronTag = rawWld.Name(frag4Def.NameRef)
				default:
					return fmt.Errorf("fragment4 wanted polyhedrondef, got unknown type %T", frag4)
				}
			default:
				return fmt.Errorf("fragment4 unknown type %T", frag4)
			}
		}

	}
	e.CenterOffset = frag.CenterOffset
	e.Params2 = frag.Params2
	e.BoundingRadius = frag.BoundingRadius
	e.BoundingBoxMin = frag.BoundingBoxMin
	e.BoundingBoxMax = frag.BoundingBoxMax
	e.FPScale = frag.Scale

	scale := 1.0 / float32(int(1<<frag.Scale))

	for _, vert := range frag.Vertices {
		e.Vertices = append(e.Vertices, [3]float32{
			float32(vert[0]) * scale,
			float32(vert[1]) * scale,
			float32(vert[2]) * scale,
		})
	}
	for _, uv := range frag.UVs {
		if rawWld.IsNewWorld {
			e.UVs = append(e.UVs, [2]float32{
				float32(uv[0]),
				float32(uv[1]),
			})
		} else {
			e.UVs = append(e.UVs, [2]float32{
				float32(uv[0]) / float32(256),
				float32(uv[1]) / float32(256),
			})
		}
	}
	for _, vn := range frag.VertexNormals {
		e.VertexNormals = append(e.VertexNormals, [3]float32{
			float32(vn[0]) / float32(128),
			float32(vn[1]) / float32(128),
			float32(vn[2]) / float32(128),
		})
	}

	e.VertexColors = frag.Colors

	for _, face := range frag.Faces {
		f := &Face{
			Triangle: face.Index,
		}
		if face.Flags&0x10 != 0 {
			f.Passable = 1
		}
		e.Faces = append(e.Faces, f)
	}

	if frag.Flags&0x1 != 0 {
		e.HexOneFlag = 1
	}
	if frag.Flags&0x2 != 0 {
		e.HexTwoFlag = 1
	}
	if frag.Flags&0x4000 != 0 {
		e.HexFourThousandFlag = 1
	}
	if frag.Flags&0x8000 != 0 {
		e.HexEightThousandFlag = 1
	}
	if frag.Flags&0x10000 != 0 {
		e.HexTenThousandFlag = 1
	}
	if frag.Flags&0x20000 != 0 {
		e.HexTwentyThousandFlag = 1
	}

	e.FaceMaterialGroups = frag.FaceMaterialGroups
	e.SkinAssignmentGroups = frag.SkinAssignmentGroups
	e.VertexMaterialGroups = frag.VertexMaterialGroups

	for _, mop := range frag.MeshOps {
		e.MeshOps = append(e.MeshOps, &MeshOp{
			Index1:    mop.Index1,
			Index2:    mop.Index2,
			Offset:    mop.Offset,
			Param1:    mop.Param1,
			TypeField: mop.TypeField,
		})
	}
	return nil
}

// DMSpriteDef is a declaration of DMSPRITEDEF
type DMSpriteDef struct {
	fragID               int16
	Tag                  string
	Fragment1            int16
	MaterialPaletteTag   string
	Fragment3            uint32
	Center               NullFloat32Slice3
	Params1              [3]float32
	Vertices             [][3]float32
	TexCoords            [][2]float32
	Normals              [][3]float32
	Colors               []int32
	Faces                []*DMSpriteDefFace
	Meshops              []*DMSpriteDefMeshOp
	SkinAssignmentGroups [][2]uint16
	Data8                []uint32 // 0x200 flag
	FaceMaterialGroups   [][2]int16
	VertexMaterialGroups [][2]int16
	Params2              NullFloat32Slice3
	Params3              NullFloat32Slice6
	HexTwoHundredFlag    int16
	HexEightHundredFlag  int16
	HexOneThousandFlag   int16
}

type DMSpriteDefFace struct {
	Flags         uint16
	Data          [4]uint16
	VertexIndexes [3]uint16
}

type DMSpriteDefMeshOp struct {
	TypeField   uint32
	VertexIndex uint32
	Offset      float32
	Param1      uint16
	Param2      uint16
}

func (e *DMSpriteDef) Definition() string {
	return "DMSPRITEDEFINITION"
}

func (e *DMSpriteDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if e.MaterialPaletteTag != "" {
		materialPalette := token.wce.ByTag(e.MaterialPaletteTag)
		if materialPalette == nil {
			return fmt.Errorf("material palette %s not found", e.MaterialPaletteTag)
		}
		err = materialPalette.Write(token)
		if err != nil {
			return fmt.Errorf("material palette %s write: %w", e.MaterialPaletteTag, err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tFRAGMENT1 %d\n", e.Fragment1)
	fmt.Fprintf(w, "\tMATERIALPALETTE \"%s\"\n", e.MaterialPaletteTag)
	fmt.Fprintf(w, "\tFRAGMENT3 %d\n", e.Fragment3)
	fmt.Fprintf(w, "\tCENTER? %s\n", wcVal(e.Center))
	fmt.Fprintf(w, "\tPARAMS1 %0.8e %0.8e %0.8e\n", e.Params1[0], e.Params1[1], e.Params1[2])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMTEXCOORDS %d\n", len(e.TexCoords))
	for _, tex := range e.TexCoords {
		fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", tex[0], tex[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMNORMALS %d\n", len(e.Normals))
	for _, normal := range e.Normals {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", normal[0], normal[1], normal[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMCOLORS %d\n", len(e.Colors))
	for _, color := range e.Colors {
		fmt.Fprintf(w, "\t\tRGBA %d %d %d %d\n", color>>24, (color>>16)&0xff, (color>>8)&0xff, color&0xff)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMFACES %d\n", len(e.Faces))
	for i, face := range e.Faces {
		fmt.Fprintf(w, "\t\tDMFACE //%d\n", i)
		fmt.Fprintf(w, "\t\t\tFLAG %d\n", face.Flags)
		fmt.Fprintf(w, "\t\t\tDATA %d %d %d %d\n", face.Data[0], face.Data[1], face.Data[2], face.Data[3])
		fmt.Fprintf(w, "\t\t\tTRIANGLE %d %d %d\n", face.VertexIndexes[0], face.VertexIndexes[1], face.VertexIndexes[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMMESHOPS %d\n", len(e.Meshops))
	for _, meshop := range e.Meshops {
		if meshop.TypeField >= 1 && meshop.TypeField <= 3 {
			// TypeField 1-3: Offset is NULL, VertexIndex is printed
			fmt.Fprintf(w, "\t\tMESHOP %d %d NULL %d %d\n", meshop.TypeField, meshop.VertexIndex, meshop.Param1, meshop.Param2)
		} else if meshop.TypeField == 4 {
			// TypeField 4: VertexIndex is NULL, Offset is printed
			fmt.Fprintf(w, "\t\tMESHOP %d NULL %0.8f %d %d\n", meshop.TypeField, meshop.Offset, meshop.Param1, meshop.Param2)
		}
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tSKINASSIGNMENTGROUPS %d", len(e.SkinAssignmentGroups))
	for _, sa := range e.SkinAssignmentGroups {
		fmt.Fprintf(w, " %d %d", sa[0], sa[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tDATA8 %d", len(e.Data8))
	for _, d8 := range e.Data8 {
		fmt.Fprintf(w, " %d", d8)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFACEMATERIALGROUPS %d", len(e.FaceMaterialGroups))
	for _, fmg := range e.FaceMaterialGroups {
		fmt.Fprintf(w, " %d %d", fmg[0], fmg[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tVERTEXMATERIALGROUPS %d", len(e.VertexMaterialGroups))
	for _, vmg := range e.VertexMaterialGroups {
		fmt.Fprintf(w, " %d %d", vmg[0], vmg[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tPARAMS2? %s\n", wcVal(e.Params2))
	fmt.Fprintf(w, "\tPARAMS3? %s\n", wcVal(e.Params3))
	fmt.Fprintf(w, "\tHEXTWOHUNDREDFLAG %d\n", e.HexTwoHundredFlag)
	fmt.Fprintf(w, "\tHEXEIGHTHUNDREDFLAG %d\n", e.HexEightHundredFlag)
	fmt.Fprintf(w, "\tHEXONETHOUSANDFLAG %d\n", e.HexOneThousandFlag)

	fmt.Fprintf(w, "\n")
	return nil
}

func (e *DMSpriteDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("FRAGMENT1", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Fragment1, records[1])
	if err != nil {
		return fmt.Errorf("fragment1: %w", err)
	}

	records, err = token.ReadProperty("MATERIALPALETTE", 1)
	if err != nil {
		return err
	}
	e.MaterialPaletteTag = records[1]

	records, err = token.ReadProperty("FRAGMENT3", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Fragment3, records[1])
	if err != nil {
		return fmt.Errorf("fragment3: %w", err)
	}

	records, err = token.ReadProperty("CENTER?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.Center, records[1:]...)
	if err != nil {
		return fmt.Errorf("center: %w", err)
	}

	records, err = token.ReadProperty("PARAMS1", 3)
	if err != nil {
		return err
	}
	err = parse(&e.Params1, records[1:]...)
	if err != nil {
		return fmt.Errorf("params1: %w", err)
	}

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := int(0)
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

	records, err = token.ReadProperty("NUMTEXCOORDS", 1)
	if err != nil {
		return err
	}

	numUVs := int(0)
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
		e.TexCoords = append(e.TexCoords, uv)
	}

	records, err = token.ReadProperty("NUMNORMALS", 1)
	if err != nil {
		return err
	}

	numNormals := int(0)
	err = parse(&numNormals, records[1])
	if err != nil {
		return fmt.Errorf("num normals: %w", err)
	}

	for i := 0; i < numNormals; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		normal := [3]float32{}
		err = parse(&normal, records[1:]...)
		if err != nil {
			return fmt.Errorf("normal %d: %w", i, err)
		}
		e.Normals = append(e.Normals, normal)
	}

	records, err = token.ReadProperty("NUMCOLORS", 1)
	if err != nil {
		return err
	}

	numColors := int(0)
	err = parse(&numColors, records[1])
	if err != nil {
		return fmt.Errorf("num colors: %w", err)
	}

	for i := 0; i < numColors; i++ {
		records, err = token.ReadProperty("RGBA", 4)
		if err != nil {
			return err
		}
		color := int32(0)
		err = parse(&color, records[1:]...)
		if err != nil {
			return fmt.Errorf("color %d: %w", i, err)
		}
		e.Colors = append(e.Colors, color)
	}

	records, err = token.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}

	numFaces := int(0)
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		face := &DMSpriteDefFace{}
		_, err = token.ReadProperty("DMFACE", 0)
		if err != nil {
			return err
		}
		records, err = token.ReadProperty("FLAG", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Flags, records[1])
		if err != nil {
			return fmt.Errorf("face %d 0x004b flag: %w", i, err)
		}

		records, err = token.ReadProperty("DATA", 4)
		if err != nil {
			return err
		}
		err = parse(&face.Data, records[1:]...)
		if err != nil {
			return fmt.Errorf("face %d data: %w", i, err)
		}

		records, err = token.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		err = parse(&face.VertexIndexes, records[1:]...)
		if err != nil {
			return fmt.Errorf("face %d triangle: %w", i, err)
		}

		e.Faces = append(e.Faces, face)
	}

	records, err = token.ReadProperty("NUMMESHOPS", 1)
	if err != nil {
		return err
	}
	numMeshOps := int(0)
	err = parse(&numMeshOps, records[1])
	if err != nil {
		return fmt.Errorf("num mesh ops: %w", err)
	}

	for i := 0; i < numMeshOps; i++ {
		meshOp := &DMSpriteDefMeshOp{}
		records, err = token.ReadProperty("MESHOP", 5)
		if err != nil {
			return err
		}
		err = parse(&meshOp.TypeField, records[1])
		if err != nil {
			return fmt.Errorf("mesh op %d typefield: %w", i, err)
		}

		// Handle conditional NULL values for VertexIndex and Offset
		if meshOp.TypeField >= 1 && meshOp.TypeField <= 3 {
			// TypeField 1-3: Offset is NULL, VertexIndex is valid
			err = parse(&meshOp.VertexIndex, records[2])
			if err != nil {
				return fmt.Errorf("mesh op %d vertex index: %w", i, err)
			}
			meshOp.Offset = 0 // Offset is NULL
		} else if meshOp.TypeField == 4 {
			// TypeField 4: VertexIndex is NULL, Offset is valid
			err = parse(&meshOp.Offset, records[3])
			if err != nil {
				return fmt.Errorf("mesh op %d offset: %w", i, err)
			}
			meshOp.VertexIndex = 0 // VertexIndex is NULL
		} else {
			return fmt.Errorf("mesh op %d invalid typefield: %d", i, meshOp.TypeField)
		}

		err = parse(&meshOp.Param1, records[4])
		if err != nil {
			return fmt.Errorf("mesh op %d param1: %w", i, err)
		}
		err = parse(&meshOp.Param2, records[5])
		if err != nil {
			return fmt.Errorf("mesh op %d param2: %w", i, err)
		}

		e.Meshops = append(e.Meshops, meshOp)
	}

	records, err = token.ReadProperty("SKINASSIGNMENTGROUPS", -1)
	if err != nil {
		return err
	}

	numSkinAssignments := int(0)
	err = parse(&numSkinAssignments, records[1])
	if err != nil {
		return fmt.Errorf("num skin assignments: %w", err)
	}

	for i := 0; i < numSkinAssignments*2; i++ {

		val1 := uint16(0)
		err = parse(&val1, records[i+2])
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}

		val2 := uint16(0)
		err = parse(&val2, records[i+3])
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}
		e.SkinAssignmentGroups = append(e.SkinAssignmentGroups, [2]uint16{val1, val2})
		i++
	}

	records, err = token.ReadProperty("DATA8", -1)
	if err != nil {
		return err
	}

	for _, record := range records[1:] {
		val := uint32(0)
		err = parse(&val, record)
		if err != nil {
			return fmt.Errorf("data8: %w", err)
		}
		e.Data8 = append(e.Data8, val)
	}

	records, err = token.ReadProperty("FACEMATERIALGROUPS", -1)
	if err != nil {
		return err
	}

	numFaceMaterialGroups := int(0)
	err = parse(&numFaceMaterialGroups, records[1])
	if err != nil {
		return fmt.Errorf("num face material groups: %w", err)
	}

	for i := 0; i < numFaceMaterialGroups*2; i++ {
		val1 := int16(0)
		err = parse(&val1, records[i+2])
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}

		val2 := int16(0)
		err = parse(&val2, records[i+3])
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		e.FaceMaterialGroups = append(e.FaceMaterialGroups, [2]int16{val1, val2})
		i++
	}

	records, err = token.ReadProperty("VERTEXMATERIALGROUPS", -1)
	if err != nil {
		return err
	}

	numVertexMaterialGroups := int(0)
	err = parse(&numVertexMaterialGroups, records[1])
	if err != nil {
		return fmt.Errorf("num vertex material groups: %w", err)
	}

	for i := 0; i < numVertexMaterialGroups*2; i++ {
		val1 := int16(0)
		err = parse(&val1, records[i+2])
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}

		val2 := int16(0)
		err = parse(&val2, records[i+3])
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		e.VertexMaterialGroups = append(e.VertexMaterialGroups, [2]int16{val1, val2})
		i++
	}

	records, err = token.ReadProperty("PARAMS2?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.Params2, records[1:]...)
	if err != nil {
		return fmt.Errorf("params2: %w", err)
	}

	records, err = token.ReadProperty("PARAMS3?", 6)
	if err != nil {
		return err
	}

	err = parse(&e.Params3, records[1:]...)
	if err != nil {
		return fmt.Errorf("params3: %w", err)
	}

	records, err = token.ReadProperty("HEXTWOHUNDREDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwoHundredFlag, records[1])
	if err != nil {
		return fmt.Errorf("hextwohundredflag: %w", err)
	}

	records, err = token.ReadProperty("HEXEIGHTHUNDREDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexEightHundredFlag, records[1])
	if err != nil {
		return fmt.Errorf("hexeighthundredflag: %w", err)
	}

	records, err = token.ReadProperty("HEXONETHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexOneThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hexonethousandflag: %w", err)
	}

	return nil
}

func (e *DMSpriteDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	var err error
	if e.fragID != 0 {
		return e.fragID, nil
	}

	materialPaletteRef := int16(0)
	if e.MaterialPaletteTag != "" {
		palette := wce.ByTag(e.MaterialPaletteTag)
		if palette == nil {
			return -1, fmt.Errorf("material palette %s not found", e.MaterialPaletteTag)
		}

		materialPaletteRef, err = palette.ToRaw(wce, rawWld)
		if err != nil {
			return -1, fmt.Errorf("material palette %s to raw: %w", e.MaterialPaletteTag, err)
		}
	}

	wfDMSpriteDef := &rawfrag.WldFragDMSpriteDef{
		MaterialPaletteRef: uint32(materialPaletteRef),
		CenterOffset:       [3]float32{e.Center.Float32Slice3[0], e.Center.Float32Slice3[1], e.Center.Float32Slice3[2]},
	}
	wfDMSpriteDef.NameRef = rawWld.NameAdd(e.Tag)
	wfDMSpriteDef.Fragment1 = e.Fragment1
	wfDMSpriteDef.Fragment3 = e.Fragment3
	wfDMSpriteDef.Params1 = e.Params1
	wfDMSpriteDef.Vertices = e.Vertices
	wfDMSpriteDef.TexCoords = e.TexCoords
	wfDMSpriteDef.Normals = e.Normals
	wfDMSpriteDef.Colors = e.Colors
	for _, face := range e.Faces {
		wfDMSpriteDef.Faces = append(wfDMSpriteDef.Faces, rawfrag.WldFragDMSpriteDefFace{
			Flags:         face.Flags,
			Data:          face.Data,
			VertexIndexes: face.VertexIndexes,
		})
	}
	for _, meshop := range e.Meshops {
		wfDMSpriteDef.Meshops = append(wfDMSpriteDef.Meshops, rawfrag.WldFragDMSpriteDefMeshOp{
			TypeField:   meshop.TypeField,
			VertexIndex: meshop.VertexIndex,
			Offset:      meshop.Offset,
			Param1:      meshop.Param1,
			Param2:      meshop.Param2,
		})
	}
	wfDMSpriteDef.SkinAssignmentGroups = e.SkinAssignmentGroups
	wfDMSpriteDef.Data8 = e.Data8
	wfDMSpriteDef.FaceMaterialGroups = e.FaceMaterialGroups
	wfDMSpriteDef.VertexMaterialGroups = e.VertexMaterialGroups
	if e.Params2.Valid {
		wfDMSpriteDef.Flags |= 0x2000
		wfDMSpriteDef.Params2 = e.Params2.Float32Slice3
	}
	if e.Params3.Valid {
		wfDMSpriteDef.Flags |= 0x4000
		wfDMSpriteDef.Params3 = e.Params3.Float32Slice6
	}

	if e.Center.Valid {
		wfDMSpriteDef.Flags |= 0x1
	}
	if e.Params2.Valid {
		wfDMSpriteDef.Flags |= 0x2
	}
	if e.HexEightHundredFlag != 0 {
		wfDMSpriteDef.Flags |= 0x800
	}
	if e.HexOneThousandFlag != 0 {
		wfDMSpriteDef.Flags |= 0x1000
	}
	if e.Params3.Valid {
		wfDMSpriteDef.Flags |= 0x2000
	}

	rawWld.Fragments = append(rawWld.Fragments, wfDMSpriteDef)
	e.fragID = int16(len(rawWld.Fragments))

	return int16(len(rawWld.Fragments)), nil
}

func (e *DMSpriteDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragDMSpriteDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not dmspritedef (wrong fragcode?)")
	}
	e.Tag = rawWld.Name(frag.NameRef)
	e.Fragment1 = frag.Fragment1
	if frag.MaterialPaletteRef > 0 {
		if len(rawWld.Fragments) < int(frag.MaterialPaletteRef) {
			return fmt.Errorf("materialpalette ref %d out of bounds", frag.MaterialPaletteRef)
		}
		materialPalette, ok := rawWld.Fragments[frag.MaterialPaletteRef].(*rawfrag.WldFragMaterialPalette)
		if !ok {
			return fmt.Errorf("materialpalette ref %d not found", frag.MaterialPaletteRef)
		}
		e.MaterialPaletteTag = rawWld.Name(materialPalette.NameRef)
	}
	e.Fragment3 = frag.Fragment3
	e.Params1 = frag.Params1
	e.Vertices = frag.Vertices
	e.TexCoords = frag.TexCoords
	e.Normals = frag.Normals
	e.Colors = frag.Colors
	for _, face := range frag.Faces {
		e.Faces = append(e.Faces, &DMSpriteDefFace{
			Flags:         face.Flags,
			Data:          face.Data,
			VertexIndexes: face.VertexIndexes,
		})
	}
	for _, meshop := range frag.Meshops {
		e.Meshops = append(e.Meshops, &DMSpriteDefMeshOp{
			TypeField:   meshop.TypeField,
			VertexIndex: meshop.VertexIndex,
			Offset:      meshop.Offset,
			Param1:      meshop.Param1,
			Param2:      meshop.Param2,
		})
	}
	e.SkinAssignmentGroups = frag.SkinAssignmentGroups
	e.Data8 = frag.Data8
	e.FaceMaterialGroups = frag.FaceMaterialGroups
	e.VertexMaterialGroups = frag.VertexMaterialGroups
	if frag.Flags&0x01 != 0 {
		e.Center.Valid = true
		e.Center.Float32Slice3 = frag.CenterOffset
	}
	if frag.Flags&0x02 != 0 {
		e.Params2.Valid = true
		e.Params2.Float32Slice3 = frag.Params2
	}

	if frag.Flags&0x200 != 0 {
		e.HexTwoHundredFlag = 1
	}

	if frag.Flags&0x2000 != 0 {
		e.Params3.Valid = true
		e.Params3.Float32Slice6 = frag.Params3
	}

	return nil

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

func (e *MaterialPalette) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	if token.TagIsWritten(e.Tag) {
		return nil
	}
	token.TagSetIsWritten(e.Tag)

	for _, materialTag := range e.Materials {
		materialDef := token.wce.ByTag(materialTag)
		if materialDef == nil {
			return fmt.Errorf("material %s not found", materialTag)
		}

		err = materialDef.Write(token)
		if err != nil {
			return fmt.Errorf("write materialdef %s: %w", materialTag, err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tNUMMATERIALS %d\n", len(e.Materials))
	for _, mat := range e.Materials {
		fmt.Fprintf(w, "\tMATERIAL \"%s\"\n", mat)
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *MaterialPalette) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("NUMMATERIALS", 1)
	if err != nil {
		return fmt.Errorf("NUMMATERIALS: %w", err)
	}
	numMaterials := int(0)
	err = parse(&numMaterials, records[1])
	if err != nil {
		return fmt.Errorf("num materials: %w", err)
	}

	for i := 0; i < numMaterials; i++ {
		records, err = token.ReadProperty("MATERIAL", 1)
		if err != nil {
			return fmt.Errorf("MATERIAL: %w", err)
		}
		e.Materials = append(e.Materials, records[1])
	}

	return nil
}

func (e *MaterialPalette) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfPalette := &rawfrag.WldFragMaterialPalette{
		Flags: e.flags,
	}
	for _, mat := range e.Materials {

		srcMat := wce.ByTag(mat)
		if srcMat == nil {
			return -1, fmt.Errorf("material %s not found", mat)
		}

		matRef, err := srcMat.ToRaw(wce, rawWld)
		if err != nil {
			return -1, fmt.Errorf("material %s to raw: %w", mat, err)
		}

		wfPalette.MaterialRefs = append(wfPalette.MaterialRefs, uint32(matRef))
	}

	wfPalette.NameRef = rawWld.NameAdd(e.Tag)
	rawWld.Fragments = append(rawWld.Fragments, wfPalette)
	e.fragID = int16(len(rawWld.Fragments))

	return int16(len(rawWld.Fragments)), nil
}

func (e *MaterialPalette) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragMaterialPalette) error {
	if frag == nil {
		return fmt.Errorf("frag is not materialpalette (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.flags = frag.Flags
	for _, materialRef := range frag.MaterialRefs {
		if len(rawWld.Fragments) < int(materialRef) {
			return fmt.Errorf("material ref %d not found", materialRef)
		}
		material, ok := rawWld.Fragments[materialRef].(*rawfrag.WldFragMaterialDef)
		if !ok {
			return fmt.Errorf("invalid materialdef fragment at offset %d", materialRef)
		}
		e.Materials = append(e.Materials, rawWld.Name(material.NameRef))
	}

	return nil
}

// MaterialDef is an entry MATERIALDEFINITION
type MaterialDef struct {
	fragID             int16
	model              string
	Tag                string
	Variation          int
	SpriteHexFiftyFlag int
	RenderMethod       string
	RGBPen             [4]uint8
	Brightness         float32
	ScaledAmbient      float32
	SimpleSpriteTag    string
	Pair1              NullUint32
	Pair2              NullFloat32
	DoubleSided        int
}

func (e *MaterialDef) Definition() string {
	return "MATERIALDEFINITION"
}

func (e *MaterialDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if token.TagIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.Variation)) {
		return nil
	}

	token.TagSetIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.Variation))

	if e.SimpleSpriteTag != "" {
		simpleSprite := token.wce.ByTagWithIndex(e.SimpleSpriteTag, e.Variation)
		if simpleSprite == nil {
			return fmt.Errorf("simple sprite %s not found", e.SimpleSpriteTag)
		}
		err = simpleSprite.Write(token)
		if err != nil {
			return fmt.Errorf("simple sprite %s: %w", e.SimpleSpriteTag, err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tVARIATION %d\n", e.Variation)
	fmt.Fprintf(w, "\tRENDERMETHOD \"%s\"\n", e.RenderMethod)
	fmt.Fprintf(w, "\tRGBPEN %d %d %d %d\n", e.RGBPen[0], e.RGBPen[1], e.RGBPen[2], e.RGBPen[3])
	fmt.Fprintf(w, "\tBRIGHTNESS %0.8e\n", e.Brightness)
	fmt.Fprintf(w, "\tSCALEDAMBIENT %0.8e\n", e.ScaledAmbient)
	fmt.Fprintf(w, "\tSIMPLESPRITEINST\n")
	fmt.Fprintf(w, "\t\tTAG \"%s\"\n", e.SimpleSpriteTag)
	fmt.Fprintf(w, "\t\tHEXFIFTYFLAG %d\n", e.SpriteHexFiftyFlag)
	fmt.Fprintf(w, "\tPAIRS? %s %s\n", wcVal(e.Pair1), wcVal(e.Pair2))
	fmt.Fprintf(w, "\tDOUBLESIDED %d\n", e.DoubleSided)
	fmt.Fprintf(w, "\n")

	token.TagSetIsWritten(e.Tag)
	return nil
}

func (e *MaterialDef) Read(token *AsciiReadToken) error {

	e.model = token.wce.lastReadModelTag

	records, err := token.ReadProperty("VARIATION", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Variation, records[1])
	if err != nil {
		return fmt.Errorf("variation: %w", err)
	}

	records, err = token.ReadProperty("RENDERMETHOD", 1)
	if err != nil {
		return err
	}
	e.RenderMethod = records[1]

	records, err = token.ReadProperty("RGBPEN", 4)
	if err != nil {
		return err
	}
	err = parse(&e.RGBPen, records[1:]...)
	if err != nil {
		return fmt.Errorf("rgbpen: %w", err)
	}

	records, err = token.ReadProperty("BRIGHTNESS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Brightness, records[1])
	if err != nil {
		return fmt.Errorf("brightness: %w", err)
	}

	records, err = token.ReadProperty("SCALEDAMBIENT", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ScaledAmbient, records[1])
	if err != nil {
		return fmt.Errorf("scaled ambient: %w", err)
	}

	_, err = token.ReadProperty("SIMPLESPRITEINST", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	e.SimpleSpriteTag = records[1]

	records, err = token.ReadProperty("HEXFIFTYFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpriteHexFiftyFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex fifty flag: %w", err)
	}

	records, err = token.ReadProperty("PAIRS?", 2)
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

	records, err = token.ReadProperty("DOUBLESIDED", 1)
	if err != nil {
		return err
	}
	err = parse(&e.DoubleSided, records[1])
	if err != nil {
		return fmt.Errorf("doublesided: %w", err)
	}

	token.wce.variationMaterialDefs[token.wce.lastReadModelTag] = append(token.wce.variationMaterialDefs[token.wce.lastReadModelTag], e)
	return nil
}

func (e *MaterialDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if !wce.isVariationMaterial && e.fragID != 0 {
		return e.fragID, nil
	}

	wfMaterialDef := &rawfrag.WldFragMaterialDef{
		RenderMethod:  model.RenderMethodInt(e.RenderMethod),
		RGBPen:        e.RGBPen,
		Brightness:    e.Brightness,
		ScaledAmbient: e.ScaledAmbient,
	}

	if e.DoubleSided > 0 {
		wfMaterialDef.Flags |= 0x01
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
		spriteDef := wce.ByTagWithIndex(e.SimpleSpriteTag, e.Variation)
		if spriteDef == nil {
			return -1, fmt.Errorf("simple sprite %s not found", e.SimpleSpriteTag)
		}

		spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
		if err != nil {
			return -1, fmt.Errorf("simple sprite %s to raw: %w", e.SimpleSpriteTag, err)
		}

		wfSprite := &rawfrag.WldFragSimpleSprite{
			//NameRef:   rawWld.NameAdd(m.SimpleSpriteTag),
			SpriteRef: uint32(spriteDefRef),
		}

		if e.SpriteHexFiftyFlag > 0 {
			wfSprite.Flags |= 0x50
		}
		rawWld.Fragments = append(rawWld.Fragments, wfSprite)

		spriteRef := int16(len(rawWld.Fragments))

		wfMaterialDef.SimpleSpriteRef = uint32(spriteRef)
	}

	wfMaterialDef.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfMaterialDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *MaterialDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragMaterialDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not materialdef (wrong fragcode?)")
	}

	e.model = wce.lastReadModelTag

	if frag.SimpleSpriteRef > 0 {
		if len(rawWld.Fragments) < int(frag.SimpleSpriteRef) {
			return fmt.Errorf("simplesprite ref %d out of bounds", frag.SimpleSpriteRef)
		}
		simpleSprite, ok := rawWld.Fragments[frag.SimpleSpriteRef].(*rawfrag.WldFragSimpleSprite)
		if !ok {
			return fmt.Errorf("simplesprite ref %d not found", frag.SimpleSpriteRef)
		}
		if len(rawWld.Fragments) < int(simpleSprite.SpriteRef) {
			return fmt.Errorf("sprite ref %d out of bounds", simpleSprite.SpriteRef)
		}
		spriteDef, ok := rawWld.Fragments[simpleSprite.SpriteRef].(*rawfrag.WldFragSimpleSpriteDef)
		if !ok {
			return fmt.Errorf("material's simple sprite ref %d not found", simpleSprite.SpriteRef)
		}
		if simpleSprite.Flags&0x50 != 0 {
			e.SpriteHexFiftyFlag = 1
		}

		e.SimpleSpriteTag = rawWld.Name(spriteDef.NameRef)
	}
	e.Tag = rawWld.Name(frag.NameRef)
	e.RenderMethod = model.RenderMethodStr(frag.RenderMethod)
	e.RGBPen = frag.RGBPen
	e.Brightness = frag.Brightness
	e.ScaledAmbient = frag.ScaledAmbient
	if wce.isVariationMaterial {
		e.Variation = 1
	}
	if frag.Flags&0x01 != 0 {
		e.DoubleSided = 1
	}
	if frag.Flags&0x02 != 0 {
		e.Pair1.Valid = true
		e.Pair1.Uint32 = frag.Pair1
		e.Pair2.Valid = true
		e.Pair2.Float32 = frag.Pair2
	}

	wce.variationMaterialDefs[wce.lastReadModelTag] = append(wce.variationMaterialDefs[wce.lastReadModelTag], e)
	return nil
}

// BlitSpriteDef is a declaration of BLITSPRITEDEF
type BlitSpriteDef struct {
	fragID      int16
	Tag         string
	SpriteTag   string
	Unknown     int32
	Transparent int16
}

func (e *BlitSpriteDef) Definition() string {
	return "BLITSPRITEDEF"
}

func (e *BlitSpriteDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if token.TagIsWritten(e.Tag) {
		return nil
	}

	token.TagSetIsWritten(e.Tag)

	if e.SpriteTag != "" {
		spriteDef := token.wce.ByTag(e.SpriteTag)
		if spriteDef == nil {
			return fmt.Errorf("sprite %s not found", e.SpriteTag)
		}
		err = spriteDef.Write(token)
		if err != nil {
			return fmt.Errorf("write sprite %s: %w", e.SpriteTag, err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tSPRITE \"%s\"\n", e.SpriteTag)
	fmt.Fprintf(w, "\tUNKNOWN %d\n", e.Unknown)
	fmt.Fprintf(w, "\tTRANSPARENT %d\n", e.Transparent)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *BlitSpriteDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("SPRITE", 1)
	if err != nil {
		return fmt.Errorf("SPRITE: %w", err)
	}
	e.SpriteTag = records[1]

	records, err = token.ReadProperty("UNKNOWN", 1)
	if err != nil {
		return fmt.Errorf("UNKNOWN: %w", err)
	}

	err = parse(&e.Unknown, records[1])
	if err != nil {
		return fmt.Errorf("unknown: %w", err)
	}

	records, err = token.ReadProperty("TRANSPARENT", 1)
	if err != nil {
		return fmt.Errorf("TRANSPARENT: %w", err)
	}

	err = parse(&e.Transparent, records[1])
	if err != nil {
		return fmt.Errorf("transparent: %w", err)
	}

	return nil
}

func (e *BlitSpriteDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	if e.SpriteTag == "" {
		return -1, fmt.Errorf("sprite tag not set")
	}

	spriteDef := wce.ByTag(e.SpriteTag)
	if spriteDef == nil {
		return -1, fmt.Errorf("sprite %s not found", e.SpriteTag)
	}

	spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
	if err != nil {
		return -1, fmt.Errorf("sprite %s to raw: %w", e.SpriteTag, err)
	}

	wfSprite := &rawfrag.WldFragSimpleSprite{
		//NameRef:   rawWld.NameAdd(m.SimpleSpriteTag),
		SpriteRef: uint32(spriteDefRef),
	}

	rawWld.Fragments = append(rawWld.Fragments, wfSprite)

	spriteRef := int16(len(rawWld.Fragments))

	wfBlitSpriteDef := &rawfrag.WldFragBlitSpriteDef{
		Unknown:           e.Unknown,
		SpriteInstanceRef: uint32(spriteRef),
	}

	if e.Transparent > 0 {
		wfBlitSpriteDef.Flags |= 0x100
	}

	wfBlitSpriteDef.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfBlitSpriteDef)
	e.fragID = int16(len(rawWld.Fragments))

	return e.fragID, nil
}

func (e *BlitSpriteDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragBlitSpriteDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not blitspritedef (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	if frag.SpriteInstanceRef > 0 {
		if len(rawWld.Fragments) < int(frag.SpriteInstanceRef) {
			return fmt.Errorf("sprite ref %d out of bounds", frag.SpriteInstanceRef)
		}

		spriteInst, ok := rawWld.Fragments[frag.SpriteInstanceRef].(*rawfrag.WldFragSimpleSprite)
		if !ok {
			return fmt.Errorf("sprite ref %d not found", frag.SpriteInstanceRef)
		}

		if len(rawWld.Fragments) < int(spriteInst.SpriteRef) {
			return fmt.Errorf("sprite ref %d out of bounds", spriteInst.SpriteRef)
		}

		spriteDef, ok := rawWld.Fragments[spriteInst.SpriteRef].(*rawfrag.WldFragSimpleSpriteDef)
		if !ok {
			return fmt.Errorf("spritedef ref %d not found", spriteInst.SpriteRef)
		}

		e.SpriteTag = rawWld.Name(spriteDef.NameRef)
	}

	e.Unknown = frag.Unknown
	e.Transparent = int16(frag.Flags & 0x100)
	return nil
}

// SimpleSpriteDef is a declaration of SIMPLESPRITEDEF
type SimpleSpriteDef struct {
	fragID             int16
	Tag                string
	Variation          int
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

func (e *SimpleSpriteDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	token.TagSetIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.Variation))

	token.TagSetIsWritten(e.Tag)
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tVARIATION %d\n", e.Variation)
	fmt.Fprintf(w, "\tSKIPFRAMES? %s\n", wcVal(e.SkipFrames))
	fmt.Fprintf(w, "\tANIMATED? %s\n", wcVal(e.Animated))
	fmt.Fprintf(w, "\tSLEEP? %s\n", wcVal(e.Sleep))
	fmt.Fprintf(w, "\tCURRENTFRAME? %s\n", wcVal(e.CurrentFrame))
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.SimpleSpriteFrames))
	for _, frame := range e.SimpleSpriteFrames {
		fmt.Fprintf(w, "\t\tFRAME \"%s\" \"%s\"\n", frame.TextureFile, frame.TextureTag)
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *SimpleSpriteDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("VARIATION", 1)
	if err != nil {
		return fmt.Errorf("VARIATION: %w", err)
	}
	err = parse(&e.Variation, records[1])
	if err != nil {
		return fmt.Errorf("variation: %w", err)
	}

	records, err = token.ReadProperty("SKIPFRAMES?", 1)
	if err != nil {
		return fmt.Errorf("SKIPFRAMES?: %w", err)
	}

	err = parse(&e.SkipFrames, records[1])
	if err != nil {
		return fmt.Errorf("skip frames: %w", err)
	}

	records, err = token.ReadProperty("ANIMATED?", 1)
	if err != nil {
		return fmt.Errorf("ANIMATED?: %w", err)
	}
	err = parse(&e.Animated, records[1])
	if err != nil {
		return fmt.Errorf("animated: %w", err)
	}

	records, err = token.ReadProperty("SLEEP?", 1)
	if err != nil {
		return fmt.Errorf("SLEEP?: %w", err)
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = token.ReadProperty("CURRENTFRAME?", 1)
	if err != nil {
		return fmt.Errorf("CURRENTFRAME?: %w", err)
	}

	err = parse(&e.CurrentFrame, records[1])
	if err != nil {
		return fmt.Errorf("current frame: %w", err)
	}

	records, err = token.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return fmt.Errorf("NUMFRAMES: %w", err)
	}
	numFrames := int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = token.ReadProperty("FRAME", 2)
		if err != nil {
			return fmt.Errorf("FRAME: %w", err)
		}
		e.SimpleSpriteFrames = append(e.SimpleSpriteFrames, SimpleSpriteFrame{
			TextureFile: records[1],
			TextureTag:  records[2],
		})
	}

	return nil
}

func (e *SimpleSpriteDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {

	/* if !wce.isVariationMaterial && e.fragID != 0 {
		return e.fragID, nil
	} */

	flags := uint32(0)
	wfSimpleSpriteDef := &rawfrag.WldFragSimpleSpriteDef{
		Sleep: e.Sleep.Uint32,
	}

	if e.SkipFrames.Valid {
		flags |= 0x02
	}
	//flags |= 0x04
	//if len(e.SimpleSpriteFrames) > 1 {
	//	flags |= 0x08
	//}

	if e.Animated.Valid {
		flags |= 0x08
	}

	if e.Sleep.Valid {
		flags |= 0x10
	}
	if e.CurrentFrame.Valid {
		flags |= 0x20
	}

	wfSimpleSpriteDef.Flags = flags

	if len(e.SimpleSpriteFrames) > 0 {
		wfBMInfo := &rawfrag.WldFragBMInfo{}
		for _, frame := range e.SimpleSpriteFrames {
			nameRef := rawWld.NameAdd(frame.TextureTag)
			if wfBMInfo.NameRef != 0 && nameRef != wfBMInfo.NameRef {
				rawWld.Fragments = append(rawWld.Fragments, wfBMInfo)
				wfSimpleSpriteDef.BitmapRefs = append(wfSimpleSpriteDef.BitmapRefs, uint32(len(rawWld.Fragments)))
				wfBMInfo = &rawfrag.WldFragBMInfo{}
			}
			wfBMInfo.NameRef = nameRef
			wfBMInfo.TextureNames = append(wfBMInfo.TextureNames, frame.TextureFile+"\x00")
		}

		rawWld.Fragments = append(rawWld.Fragments, wfBMInfo)
		wfSimpleSpriteDef.BitmapRefs = append(wfSimpleSpriteDef.BitmapRefs, uint32(len(rawWld.Fragments)))
	}

	wfSimpleSpriteDef.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfSimpleSpriteDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *SimpleSpriteDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragSimpleSpriteDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not simplespritedef (wrong fragcode?)")
	}
	e.Tag = rawWld.Name(frag.NameRef)
	if frag.Flags&0x02 == 0x02 {
		e.SkipFrames.Valid = true
	}
	if frag.Flags&0x08 == 0x08 {
		e.Animated.Valid = true
	}
	if frag.Flags&0x10 == 0x10 {
		e.Sleep.Valid = true
		e.Sleep.Uint32 = frag.Sleep
	}
	if frag.Flags&0x20 == 0x20 {
		e.CurrentFrame.Valid = true
		e.CurrentFrame.Int32 = frag.CurrentFrame
	}

	if wce.isVariationMaterial {
		e.Variation = 1
	}

	for _, bitmapRef := range frag.BitmapRefs {
		if bitmapRef == 0 {
			return nil
		}
		if len(rawWld.Fragments) < int(bitmapRef) {
			return fmt.Errorf("bitmap ref %d not found", bitmapRef)
		}
		bitmap := rawWld.Fragments[bitmapRef]
		bmInfo, ok := bitmap.(*rawfrag.WldFragBMInfo)
		if !ok {
			return fmt.Errorf("invalid bitmap ref %d", bitmapRef)
		}
		for _, name := range bmInfo.TextureNames {
			e.SimpleSpriteFrames = append(e.SimpleSpriteFrames, SimpleSpriteFrame{
				TextureTag:  rawWld.Name(bmInfo.NameRef),
				TextureFile: name,
			})
		}
	}
	return nil
}

// ActorDef is a declaration of ACTORDEF
type ActorDef struct {
	fragID         int16
	Tag            string
	Callback       string
	BoundsRef      int32
	CurrentAction  NullUint32
	Location       NullFloat32Slice6
	ActiveGeometry NullUint32
	Unk1           uint32
	Actions        []ActorAction
	Unk2           uint32
	HasEightyFlag  int
}

// ActorAction is a declaration of ACTION
type ActorAction struct {
	Unk1           uint32
	LevelOfDetails []ActorLevelOfDetail
}

// ActorLevelOfDetail is a declaration of LEVELOFDETAIL
type ActorLevelOfDetail struct {
	SpriteTag      string
	SpriteTagIndex int
	SpriteFlags    uint32
	MinDistance    float32
}

func (e *ActorDef) Definition() string {
	return "ACTORDEF"
}

func (e *ActorDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if token.TagIsWritten(e.Tag) {
		return nil
	}

	token.TagSetIsWritten(e.Tag)

	for _, action := range e.Actions {
		for lodIndex, lod := range action.LevelOfDetails {
			if lod.SpriteTag == "" {
				continue
			}

			spriteFrag := token.wce.ByTag(lod.SpriteTag)
			if spriteFrag == nil {
				return fmt.Errorf("lod %d sprite %s not found", lodIndex, lod.SpriteTag)
			}

			switch sprite := spriteFrag.(type) {
			case *SimpleSpriteDef:
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d spritedef %s: %w", lodIndex, sprite.Tag, err)
				}
			case *Sprite3DDef:
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d 3dspritedef %s: %w", lodIndex, sprite.Tag, err)
				}
			case *Sprite2DDef:
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d 2dspritedef %s: %w", lodIndex, sprite.Tag, err)
				}
			case *BlitSpriteDef: // particle effects ues this
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d blitspritedef %s: %w", lodIndex, sprite.Tag, err)
				}
			case *HierarchicalSpriteDef:
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d hsprite %s: %w", lodIndex, sprite.Tag, err)
				}

				variations := token.wce.variationMaterialDefs[token.wce.lastReadModelTag]
				sort.Slice(variations, func(i, j int) bool {
					return variations[i].Tag < variations[j].Tag
				})

				for _, variation := range variations {
					err = variation.Write(token)
					if err != nil {
						return fmt.Errorf("lod %d hsprite %s variation %s: %w", lodIndex, sprite.Tag, variation.Tag, err)
					}
				}

				if err != nil {
					return fmt.Errorf("lod %d hsprite %s material: %w", lodIndex, sprite.Tag, err)
				}
			case *DMSpriteDef:
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d dmspritedef %s: %w", lodIndex, sprite.Tag, err)
				}
			case *DMSpriteDef2:
				err = sprite.Write(token)
				if err != nil {
					return fmt.Errorf("lod %d dmspritedef %s: %w", lodIndex, sprite.Tag, err)
				}

			default:
				return fmt.Errorf("lod %d unknown sprite type %T", lodIndex, sprite)
			}
		}
	}

	for _, action := range e.Actions {
		for _, lod := range action.LevelOfDetails {
			if lod.SpriteTag == "" {
				continue
			}

			sprite := token.wce.ByTag(lod.SpriteTag)
			if sprite == nil {
				return fmt.Errorf("lod sprite %s not found", lod.SpriteTag)
			}

			err = sprite.Write(token)
			if err != nil {
				return fmt.Errorf("lod sprite %s: %w", lod.SpriteTag, err)
			}
		}
	}

	baseTag := strings.TrimSuffix(e.Tag, "_ACTORDEF")
	for _, sprite := range token.wce.DMSpriteDef2s {
		if !strings.HasPrefix(sprite.Tag, baseTag) {
			continue
		}
		err = sprite.Write(token)
		if err != nil {
			return fmt.Errorf("dmspritedef %s: %w", sprite.Tag, err)

		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tCALLBACK \"%s\"\n", e.Callback)
	fmt.Fprintf(w, "\tBOUNDSREF %d\n", e.BoundsRef)
	fmt.Fprintf(w, "\tCURRENTACTION? %s\n", wcVal(e.CurrentAction))
	fmt.Fprintf(w, "\tLOCATION? %s\n", wcVal(e.Location))
	fmt.Fprintf(w, "\tACTIVEGEOMETRY? %s\n", wcVal(e.ActiveGeometry))
	fmt.Fprintf(w, "\tNUMACTIONS %d\n", len(e.Actions))
	for _, action := range e.Actions {
		fmt.Fprintf(w, "\t\tACTION\n")
		fmt.Fprintf(w, "\t\t\tUNK1 %d\n", action.Unk1)
		fmt.Fprintf(w, "\t\t\tNUMLEVELSOFDETAIL %d\n", len(action.LevelOfDetails))
		for _, lod := range action.LevelOfDetails {
			fmt.Fprintf(w, "\t\t\t\tLEVELOFDETAIL\n")
			fmt.Fprintf(w, "\t\t\t\t\tSPRITE \"%s\"\n", lod.SpriteTag)
			fmt.Fprintf(w, "\t\t\t\t\tSPRITEINDEX %d\n", lod.SpriteTagIndex)
			fmt.Fprintf(w, "\t\t\t\t\tMINDISTANCE %0.8e\n", lod.MinDistance)
		}
	}
	fmt.Fprintf(w, "\tUNK2 %d\n", e.Unk2)
	fmt.Fprintf(w, "\tHASEIGHTYFLAG %d\n", e.HasEightyFlag)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *ActorDef) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("CALLBACK", 1)
	if err != nil {
		return err
	}
	e.Callback = records[1]

	records, err = token.ReadProperty("BOUNDSREF", 1)
	if err != nil {
		return err
	}

	err = parse(&e.BoundsRef, records[1])
	if err != nil {
		return fmt.Errorf("bounds ref: %w", err)
	}

	records, err = token.ReadProperty("CURRENTACTION?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentAction, records[1])
	if err != nil {
		return fmt.Errorf("current action: %w", err)
	}

	records, err = token.ReadProperty("LOCATION?", 6)
	if err != nil {
		return err
	}
	err = parse(&e.Location, records[1:]...)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = token.ReadProperty("ACTIVEGEOMETRY?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ActiveGeometry, records[1])
	if err != nil {
		return fmt.Errorf("active geometry: %w", err)
	}

	records, err = token.ReadProperty("NUMACTIONS", 1)
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
		_, err = token.ReadProperty("ACTION", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("UNK1", 1)
		if err != nil {
			return err
		}
		err = parse(&action.Unk1, records[1])
		if err != nil {
			return fmt.Errorf("unk1: %w", err)
		}

		records, err = token.ReadProperty("NUMLEVELSOFDETAIL", 1)
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
			_, err = token.ReadProperty("LEVELOFDETAIL", 0)
			if err != nil {
				return err
			}

			records, err = token.ReadProperty("SPRITE", 1)
			if err != nil {
				return err
			}
			lod.SpriteTag = records[1]

			records, err = token.ReadProperty("SPRITEINDEX", 1)
			if err != nil {
				return err
			}
			err = parse(&lod.SpriteTagIndex, records[1])
			if err != nil {
				return fmt.Errorf("sprite index: %w", err)
			}

			records, err = token.ReadProperty("MINDISTANCE", 1)
			if err != nil {
				return err
			}

			err = parse(&lod.MinDistance, records[1])
			if err != nil {
				return fmt.Errorf("min distance: %w", err)
			}

			action.LevelOfDetails = append(action.LevelOfDetails, lod)
		}

		e.Actions = append(e.Actions, action)

	}

	records, err = token.ReadProperty("UNK2", 1)
	if err != nil {
		return err
	}

	err = parse(&e.Unk2, records[1])
	if err != nil {
		return fmt.Errorf("unk2: %w", err)
	}

	records, err = token.ReadProperty("HASEIGHTYFLAG", 1)
	if err != nil {
		return err
	}

	err = parse(&e.HasEightyFlag, records[1])
	if err != nil {
		return fmt.Errorf("has eighty flag: %w", err)
	}

	return nil
}

func (e *ActorDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	var err error
	if e.fragID != 0 {
		return e.fragID, nil
	}

	actorDef := &rawfrag.WldFragActorDef{
		BoundsRef:     e.BoundsRef,
		CurrentAction: e.CurrentAction.Uint32,
	}

	if e.CurrentAction.Valid {
		actorDef.Flags |= 0x01
		actorDef.CurrentAction = e.CurrentAction.Uint32
	}

	if e.Location.Valid {
		actorDef.Flags |= 0x02
		actorDef.Location = e.Location.Float32Slice6
	}

	if e.ActiveGeometry.Valid {
		actorDef.Flags |= 0x40
		//actorDef.ActiveGeometry = e.ActiveGeometry.Uint32
	}

	if e.HasEightyFlag > 0 {
		actorDef.Flags |= 0x80
	}

	wce.lastReadModelTag = strings.TrimSuffix(e.Tag, "_ACTORDEF")

	for _, action := range e.Actions {
		actorAction := rawfrag.WldFragModelAction{
			Unk1: action.Unk1,
		}

		for lodIndex, lod := range action.LevelOfDetails {
			if lod.SpriteTag == "" {
				continue
			}

			var spriteRef int16
			spriteVar := wce.ByTag(lod.SpriteTag)
			if spriteVar == nil {
				return -1, fmt.Errorf("lod sprite %s not found", lod.SpriteTag)
			}
			switch spriteDef := spriteVar.(type) {
			case *DMSpriteDef:
				spriteRef, err = spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef %s to raw: %w", lod.SpriteTag, err)
				}
				sprite := &rawfrag.WldFragDMSprite{
					DMSpriteRef: int32(spriteRef),
				}

				rawWld.Fragments = append(rawWld.Fragments, sprite)
				spriteRef = int16(len(rawWld.Fragments))
			case *DMSpriteDef2:
				spriteRef, err = spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef2 %s to raw: %w", lod.SpriteTag, err)
				}
				sprite := &rawfrag.WldFragDMSprite{
					DMSpriteRef: int32(spriteRef),
				}

				rawWld.Fragments = append(rawWld.Fragments, sprite)
				spriteRef = int16(len(rawWld.Fragments))
			case *Sprite3DDef:
				spriteRef, err = spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("3dspritedef %s to raw: %w", lod.SpriteTag, err)
				}
				sprite := &rawfrag.WldFragSprite3D{
					Flags:          lod.SpriteFlags,
					Sprite3DDefRef: int32(spriteRef),
				}

				rawWld.Fragments = append(rawWld.Fragments, sprite)
				spriteRef = int16(len(rawWld.Fragments))
			case *HierarchicalSpriteDef:
				spriteRef, err = spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("hierchcicalspritedef %s to raw: %w", lod.SpriteTag, err)
				}

				variations := wce.variationMaterialDefs[wce.lastReadModelTag]
				sort.Slice(variations, func(i, j int) bool {
					return variations[i].Tag < variations[j].Tag
				})

				for _, variation := range variations {
					_, err = variation.ToRaw(wce, rawWld)
					if err != nil {
						return -1, fmt.Errorf("lod %d hsprite %s variation %s: %w", lodIndex, spriteDef.Tag, variation.Tag, err)
					}
				}

				sprite := &rawfrag.WldFragHierarchicalSprite{
					//NameRef
					HierarchicalSpriteRef: uint32(spriteRef),
					Param:                 0,
				}

				rawWld.Fragments = append(rawWld.Fragments, sprite)
				spriteRef = int16(len(rawWld.Fragments))

			case *Sprite2DDef:
				spriteRef, err = spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("2dspritedef %s to raw: %w", lod.SpriteTag, err)
				}

				sprite := &rawfrag.WldFragSprite2D{
					TwoDSpriteRef: uint32(spriteRef),
				}

				rawWld.Fragments = append(rawWld.Fragments, sprite)
				spriteRef = int16(len(rawWld.Fragments))
			case *BlitSpriteDef:
				spriteRef, err = spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("blitspritedef %s to raw: %w", lod.SpriteTag, err)
				}

				sprite := &rawfrag.WldFragBlitSprite{
					BlitSpriteRef: int32(spriteRef),
				}

				rawWld.Fragments = append(rawWld.Fragments, sprite)
				spriteRef = int16(len(rawWld.Fragments))
			default:
				return -1, fmt.Errorf("actordef %s lod %s unknown sprite type %T", e.Tag, lod.SpriteTag, spriteDef)
			}
			if err != nil {
				return -1, fmt.Errorf("sprite %s to raw: %w", lod.SpriteTag, err)
			}

			actorAction.Lods = append(actorAction.Lods, lod.MinDistance)
			actorDef.FragmentRefs = append(actorDef.FragmentRefs, uint32(spriteRef))
		}

		actorDef.Actions = append(actorDef.Actions, actorAction)
	}

	actorDef.CallbackNameRef = rawWld.NameAdd(e.Callback)
	actorDef.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, actorDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), err
}

func (e *ActorDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragActorDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not actordef (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.Callback = rawWld.Name(frag.CallbackNameRef)
	e.BoundsRef = frag.BoundsRef
	e.Unk1 = frag.Unk1

	if frag.Flags&0x01 == 0x01 {
		e.CurrentAction.Valid = true
		e.CurrentAction.Uint32 = frag.CurrentAction
	}
	if frag.Flags&0x02 == 0x02 {
		e.Location.Valid = true
		e.Location.Float32Slice6 = frag.Location
	}
	if frag.Flags&0x40 == 0x40 {
		e.ActiveGeometry.Valid = true
	}

	if frag.Flags&0x80 == 0x80 {
		e.HasEightyFlag = 1
	}

	if len(frag.Actions) != len(frag.FragmentRefs) {
		return fmt.Errorf("actordef actions and fragmentrefs mismatch")
	}

	fragRefIndex := 0
	for _, srcAction := range frag.Actions {
		lods := []ActorLevelOfDetail{}
		for _, srcLod := range srcAction.Lods {
			spriteTag := ""
			if len(frag.FragmentRefs) > fragRefIndex {
				spriteRef := frag.FragmentRefs[fragRefIndex]
				if len(rawWld.Fragments) < int(spriteRef) {
					return fmt.Errorf("actordef fragment ref %d not found", spriteRef)
				}
				switch sprite := rawWld.Fragments[spriteRef].(type) {
				case *rawfrag.WldFragSprite3D:
					if len(rawWld.Fragments) < int(sprite.Sprite3DDefRef) {
						return fmt.Errorf("sprite3ddef ref %d out of range", sprite.Sprite3DDefRef)
					}
					spriteDef, ok := rawWld.Fragments[sprite.Sprite3DDefRef].(*rawfrag.WldFragSprite3DDef)
					if !ok {
						return fmt.Errorf("sprite3ddef ref %d not found", sprite.Sprite3DDefRef)
					}
					spriteTag = rawWld.Name(spriteDef.NameRef)
				case *rawfrag.WldFragDMSprite:
					if len(rawWld.Fragments) < int(sprite.DMSpriteRef) {
						return fmt.Errorf("dmsprite ref %d out of range", sprite.DMSpriteRef)
					}
					switch spriteDef := rawWld.Fragments[sprite.DMSpriteRef].(type) {
					case *rawfrag.WldFragDMSpriteDef:
						spriteTag = rawWld.Name(spriteDef.NameRef)
					case *rawfrag.WldFragDmSpriteDef2:
						spriteTag = rawWld.Name(spriteDef.NameRef)
					default:
						return fmt.Errorf("unhandled dmsprite instance def fragment type %d (%s)", sprite.FragCode(), raw.FragName(sprite.FragCode()))
					}
				case *rawfrag.WldFragHierarchicalSprite:
					if len(rawWld.Fragments) < int(sprite.HierarchicalSpriteRef) {
						return fmt.Errorf("hierarchicalsprite def ref %d not found", sprite.HierarchicalSpriteRef)
					}
					spriteDef, ok := rawWld.Fragments[sprite.HierarchicalSpriteRef].(*rawfrag.WldFragHierarchicalSpriteDef)
					if !ok {
						return fmt.Errorf("hierarchicalsprite def ref %d not found", sprite.HierarchicalSpriteRef)
					}
					spriteTag = rawWld.Name(spriteDef.NameRef)
				case *rawfrag.WldFragBlitSprite:
					if len(rawWld.Fragments) < int(sprite.BlitSpriteRef) {
						return fmt.Errorf("blitsprite def ref %d not found", sprite.BlitSpriteRef)
					}
					spriteDef, ok := rawWld.Fragments[sprite.BlitSpriteRef].(*rawfrag.WldFragBlitSpriteDef)
					if !ok {
						return fmt.Errorf("blitsprite def ref %d not found", sprite.BlitSpriteRef)
					}
					spriteTag = rawWld.Name(spriteDef.NameRef)

				case *rawfrag.WldFragSprite2D:
					if len(rawWld.Fragments) < int(sprite.TwoDSpriteRef) {
						return fmt.Errorf("sprite2d def ref %d not found", sprite.TwoDSpriteRef)
					}
					spriteDef, ok := rawWld.Fragments[sprite.TwoDSpriteRef].(*rawfrag.WldFragSprite2DDef)
					if !ok {
						return fmt.Errorf("sprite2d def ref %d not found", sprite.TwoDSpriteRef)
					}
					spriteTag = rawWld.Name(spriteDef.NameRef)
				default:
					return fmt.Errorf("unhandled sprite instance fragment type %d (%s)", sprite.FragCode(), raw.FragName(sprite.FragCode()))
				}
			}
			lod := ActorLevelOfDetail{
				SpriteTag:   spriteTag,
				MinDistance: srcLod,
			}

			lods = append(lods, lod)
			fragRefIndex++
		}

		e.Actions = append(e.Actions, ActorAction{
			Unk1:           srcAction.Unk1,
			LevelOfDetails: lods,
		})
	}

	return nil
}

// ActorInst is a declaration of ACTORINST
type ActorInst struct {
	fragID            int16
	Tag               string
	DefinitionTag     string
	CurrentAction     NullUint32
	Location          NullFloat32Slice6
	BoundingRadius    NullFloat32
	Scale             NullFloat32
	SoundTag          NullString
	Active            NullUint32
	SpriteVolumeOnly  NullUint32
	DMRGBTrackTag     NullString
	SphereTag         string
	SphereRadius      float32
	HexTwoHundredFlag int
	UserData          string
}

func (e *ActorInst) Definition() string {
	return "ACTORINST"
}

func (e *ActorInst) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if e.DMRGBTrackTag.Valid {
		dTrack := token.wce.ByTag(e.DMRGBTrackTag.String)
		if dTrack == nil {
			return fmt.Errorf("dmrgbtrack %s not found", e.DMRGBTrackTag.String)
		}
		err = dTrack.Write(token)
		if err != nil {
			return fmt.Errorf("dmrgbtrack %s: %w", e.DMRGBTrackTag.String, err)
		}
	}

	if e.DefinitionTag == "!UNK" {
		return fmt.Errorf("actordef %s is !UNK and not found", e.DefinitionTag)
	}

	if e.DefinitionTag != "" {
		/* 	actorDef := token.wce.ByTag(e.DefinitionTag)
		if actorDef == nil {
			return fmt.Errorf("actordef %s not found", e.DefinitionTag)
		}
		err = actorDef.Write(token)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", e.DefinitionTag, err)
		} */
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
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
	fmt.Fprintf(w, "\tSPHERERADIUS %0.8e\n", e.SphereRadius)
	fmt.Fprintf(w, "\tHEXTWOHUNDREDFLAG %d\n", e.HexTwoHundredFlag)
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", e.UserData)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *ActorInst) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.DefinitionTag = records[1]

	records, err = token.ReadProperty("CURRENTACTION?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentAction, records[1])
	if err != nil {
		return fmt.Errorf("current action: %w", err)
	}

	records, err = token.ReadProperty("LOCATION?", 6)
	if err != nil {
		return err
	}
	err = parse(&e.Location, records[1:]...)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = token.ReadProperty("BOUNDINGRADIUS?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = token.ReadProperty("SCALEFACTOR?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Scale, records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	records, err = token.ReadProperty("SOUND?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SoundTag, records[1])
	if err != nil {
		return fmt.Errorf("sound: %w", err)
	}

	records, err = token.ReadProperty("ACTIVE?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Active, records[1])
	if err != nil {
		return fmt.Errorf("active: %w", err)
	}

	records, err = token.ReadProperty("SPRITEVOLUMEONLY?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpriteVolumeOnly, records[1])
	if err != nil {
		return fmt.Errorf("sprite volume only: %w", err)
	}

	records, err = token.ReadProperty("DMRGBTRACK?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.DMRGBTrackTag, records[1])
	if err != nil {
		return fmt.Errorf("dm rgb track: %w", err)
	}

	records, err = token.ReadProperty("SPHERE", 1)
	if err != nil {
		return err
	}
	e.SphereTag = records[1]

	records, err = token.ReadProperty("SPHERERADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SphereRadius, records[1])
	if err != nil {
		return fmt.Errorf("sphere radius: %w", err)
	}

	records, err = token.ReadProperty("HEXTWOHUNDREDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwoHundredFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex two hundred flag: %w", err)
	}

	records, err = token.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}
	e.UserData = records[1]

	return nil
}

func (e *ActorInst) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	var err error
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfActorInst := &rawfrag.WldFragActor{}

	if e.DefinitionTag != "" {
		actorDef := wce.ByTag(e.DefinitionTag)
		if actorDef != nil {

			actorDefRef, err := actorDef.ToRaw(wce, rawWld)
			if err != nil {
				return -1, fmt.Errorf("actor definition %s to raw: %w", e.DefinitionTag, err)
			}

			wfActorInst.ActorDefRef = int32(actorDefRef)
		} else {
			wfActorInst.ActorDefRef = rawWld.NameAdd(e.DefinitionTag)
		}
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
		wfActorInst.SoundNameRef = rawWld.NameAdd(e.SoundTag.String)
	}

	if e.Active.Valid {
		wfActorInst.Flags |= 0x20
	}

	if e.SpriteVolumeOnly.Valid {
		wfActorInst.Flags |= 0x80
	}

	if e.DMRGBTrackTag.Valid {
		wfActorInst.Flags |= 0x100
		dmRGBTrackDef := wce.ByTag(e.DMRGBTrackTag.String)
		if dmRGBTrackDef == nil {
			return -1, fmt.Errorf("dm rgb track def %s not found", e.DMRGBTrackTag.String)
		}

		dmRGBDefTrackRef, err := dmRGBTrackDef.ToRaw(wce, rawWld)
		if err != nil {
			return -1, fmt.Errorf("dm rgb track %s to raw: %w", e.DMRGBTrackTag.String, err)
		}

		wfRGBTrack := &rawfrag.WldFragDmRGBTrack{
			NameRef:  0,
			TrackRef: int32(dmRGBDefTrackRef),
			Flags:    0,
		}
		if e.DefinitionTag != "" && wfActorInst.ActorDefRef == 0 {
			// in some cases, a string ref occurs instead
			wfActorInst.ActorDefRef = rawWld.NameAdd(e.DefinitionTag)
		}
		rawWld.Fragments = append(rawWld.Fragments, wfRGBTrack)
		dmRGBTrackRef := int16(len(rawWld.Fragments))
		wfActorInst.DMRGBTrackRef = int32(dmRGBTrackRef)
	}

	if e.HexTwoHundredFlag > 0 {
		wfActorInst.Flags |= 0x200
	}

	if e.SphereRadius > 0 {
		sphere := &rawfrag.WldFragSphere{
			NameRef: rawWld.NameAdd(e.SphereTag),
			Radius:  e.SphereRadius,
		}

		rawWld.Fragments = append(rawWld.Fragments, sphere)
		wfActorInst.SphereRef = uint32(len(rawWld.Fragments))
	}

	wfActorInst.UserData = e.UserData

	rawWld.Fragments = append(rawWld.Fragments, wfActorInst)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), err
}

func (e *ActorInst) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragActor) error {
	if frag == nil {
		return fmt.Errorf("frag is not actorinst (wrong fragcode?)")
	}

	actorDefTag := ""
	if frag.ActorDefRef != 0 {
		actorDefTag = rawWld.Name(frag.ActorDefRef) // some times it's just a string ref
		if !strings.HasSuffix(actorDefTag, "_ACTORDEF") {
			if len(rawWld.Fragments) < int(frag.ActorDefRef) {
				return fmt.Errorf("actordef ref %d out of bounds", frag.ActorDefRef)
			}

			actorDef, ok := rawWld.Fragments[frag.ActorDefRef].(*rawfrag.WldFragActorDef)
			if !ok {
				return fmt.Errorf("actordef ref %d not found", frag.ActorDefRef)
			}
			actorDefTag = rawWld.Name(actorDef.NameRef)
		}

	}

	if len(rawWld.Fragments) < int(frag.SphereRef) {
		return fmt.Errorf("sphere ref %d not found", frag.SphereRef)
	}

	sphereRadius := float32(0)
	if frag.SphereRef > 0 {
		sphereDef, ok := rawWld.Fragments[frag.SphereRef].(*rawfrag.WldFragSphere)
		if !ok {
			return fmt.Errorf("sphere ref %d not found", frag.SphereRef)
		}
		sphereRadius = sphereDef.Radius
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.DefinitionTag = actorDefTag
	e.SphereRadius = sphereRadius
	e.UserData = frag.UserData

	if frag.Flags&0x01 == 0x01 {
		e.CurrentAction.Valid = true
		e.CurrentAction.Uint32 = frag.CurrentAction
	}

	if frag.Flags&0x02 == 0x02 {
		e.Location.Valid = true
		e.Location.Float32Slice6 = frag.Location
	}

	if frag.Flags&0x04 == 0x04 {
		e.BoundingRadius.Valid = true
		e.BoundingRadius.Float32 = frag.BoundingRadius
	}

	if frag.Flags&0x08 == 0x08 {
		e.Scale.Valid = true
		e.Scale.Float32 = frag.ScaleFactor
	}

	if frag.Flags&0x10 == 0x10 {
		e.SoundTag.Valid = true
		e.SoundTag.String = rawWld.Name(frag.SoundNameRef)
	}

	if frag.Flags&0x20 == 0x20 {
		e.Active.Valid = true
	}

	// 0x40 unknown
	if frag.Flags&0x80 == 0x80 {
		e.SpriteVolumeOnly.Valid = true
	}

	if frag.Flags&0x100 == 0x100 {
		e.DMRGBTrackTag.Valid = true

		trackTag := ""
		if frag.DMRGBTrackRef == 0 {
			return fmt.Errorf("dmrgbtrack flag set, but ref is 0")
		}
		if len(rawWld.Fragments) < int(frag.DMRGBTrackRef) {
			return fmt.Errorf("dmrgbtrack ref %d out of bounds", frag.DMRGBTrackRef)
		}

		track, ok := rawWld.Fragments[frag.DMRGBTrackRef].(*rawfrag.WldFragDmRGBTrack)
		if !ok {
			return fmt.Errorf("dmrgbtrack ref %d not found", frag.DMRGBTrackRef)
		}
		if len(rawWld.Fragments) < int(track.TrackRef) {
			return fmt.Errorf("dmrgbtrackdef ref %d not found", track.TrackRef)
		}

		trackDef, ok := rawWld.Fragments[track.TrackRef].(*rawfrag.WldFragDmRGBTrackDef)
		if !ok {
			return fmt.Errorf("dmrgbtrackdef ref %d not found", track.TrackRef)
		}
		if trackDef.NameRef != 0 {
			trackTag = rawWld.Name(trackDef.NameRef)
		}
		e.DMRGBTrackTag.String = trackTag
	}

	if frag.Flags&0x200 == 0x200 {
		e.HexTwoHundredFlag = 1
	}

	return nil
}

// LightDef is a declaration of LIGHTDEF
type LightDef struct {
	fragID       int16
	Tag          string
	CurrentFrame NullUint32
	Sleep        NullUint32
	SkipFrames   int
	LightLevels  []float32
	Colors       [][3]float32
}

func (e *LightDef) Definition() string {
	return "LIGHTDEFINITION"
}

func (e *LightDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s  \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tCURRENTFRAME? %s\n", wcVal(e.CurrentFrame))
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.LightLevels))
	for _, level := range e.LightLevels {
		fmt.Fprintf(w, "\t\tLIGHTLEVELS %0.8f\n", level)
	}
	fmt.Fprintf(w, "\tSLEEP? %s\n", wcVal(e.Sleep))
	fmt.Fprintf(w, "\tSKIPFRAMES %d\n", e.SkipFrames)
	fmt.Fprintf(w, "\tNUMCOLORS %d\n", len(e.Colors))
	for _, color := range e.Colors {
		fmt.Fprintf(w, "\t\tCOLOR %0.8f %0.8f %0.8f\n", color[0], color[1], color[2])
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *LightDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("CURRENTFRAME?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentFrame, records[1])
	if err != nil {
		return fmt.Errorf("current frame: %w", err)
	}

	records, err = token.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return err
	}
	numFrames := int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = token.ReadProperty("LIGHTLEVELS", 1)
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

	records, err = token.ReadProperty("SLEEP?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = token.ReadProperty("SKIPFRAMES", 1)
	if err != nil {
		return err
	}
	if records[1] == "1" {
		e.SkipFrames = 1
	}

	records, err = token.ReadProperty("NUMCOLORS", 1)
	if err != nil {
		return err
	}
	numColors := int(0)
	err = parse(&numColors, records[1])
	if err != nil {
		return fmt.Errorf("num colors: %w", err)
	}

	for i := 0; i < numColors; i++ {
		records, err = token.ReadProperty("COLOR", 3)
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

	return nil
}

func (e *LightDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	var err error

	wfLightDef := &rawfrag.WldFragLightDef{
		NameRef: rawWld.NameAdd(e.Tag),
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

	if e.SkipFrames > 0 {
		wfLightDef.Flags |= 0x08
	}

	if len(e.Colors) > 0 {
		wfLightDef.Flags |= 0x10
		wfLightDef.Colors = e.Colors
	}

	rawWld.Fragments = append(rawWld.Fragments, wfLightDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), err
}

func (e *LightDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragLightDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not lightdef (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.LightLevels = frag.LightLevels
	e.Colors = frag.Colors
	if frag.Flags&0x01 == 0x01 {
		e.CurrentFrame.Valid = true
		e.CurrentFrame.Uint32 = frag.FrameCurrentRef
	}
	if frag.Flags&0x02 == 0x02 {
		e.Sleep.Valid = true
		e.Sleep.Uint32 = frag.Sleep
	}
	if frag.Flags&0x04 == 0x04 {
		e.LightLevels = frag.LightLevels
	} else {
		if len(frag.LightLevels) > 0 {
			return fmt.Errorf("light levels found but flag 0x04 not set")
		}
	}

	if frag.Flags&0x10 != 0 {
		e.Colors = frag.Colors
	}

	return nil
}

// PointLight is a declaration of POINTLIGHT
type PointLight struct {
	fragID          int16
	Tag             string
	LightDefTag     string
	Static          int
	StaticInfluence int
	HasRegions      int
	LightFlags      uint32
	Flags           uint32
	Location        [3]float32
	Radius          float32
}

func (e *PointLight) Definition() string {
	return "POINTLIGHT"
}

func (e *PointLight) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	if e.LightDefTag != "" {
		lightDef := token.wce.ByTag(e.LightDefTag)
		if lightDef == nil {
			return fmt.Errorf("lightdef %s not found", e.LightDefTag)
		}
		err = lightDef.Write(token)
		if err != nil {
			return fmt.Errorf("lightdef %s: %w", e.LightDefTag, err)
		}
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", e.LightDefTag)
	fmt.Fprintf(w, "\tSTATIC %d\n", e.Static)
	fmt.Fprintf(w, "\tSTATICINFLUENCE %d\n", e.StaticInfluence)
	fmt.Fprintf(w, "\tHASREGIONS %d\n", e.HasRegions)
	fmt.Fprintf(w, "\tXYZ %0.8f %0.8f %0.8f\n", e.Location[0], e.Location[1], e.Location[2])
	fmt.Fprintf(w, "\tRADIUSOFINFLUENCE %0.8e\n", e.Radius)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *PointLight) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("LIGHT", 1)
	if err != nil {
		return err
	}
	e.LightDefTag = records[1]

	records, err = token.ReadProperty("STATIC", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Static, records[1])
	if err != nil {
		return fmt.Errorf("static: %w", err)
	}

	records, err = token.ReadProperty("STATICINFLUENCE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.StaticInfluence, records[1])
	if err != nil {
		return fmt.Errorf("static influence: %w", err)
	}

	records, err = token.ReadProperty("HASREGIONS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HasRegions, records[1])
	if err != nil {
		return fmt.Errorf("has regions: %w", err)
	}

	records, err = token.ReadProperty("XYZ", 3)
	if err != nil {
		return err
	}
	err = parse(&e.Location, records[1:]...)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = token.ReadProperty("RADIUSOFINFLUENCE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Radius, records[1])
	if err != nil {
		return fmt.Errorf("radius of influence: %w", err)
	}

	return nil
}

func (e *PointLight) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	if e.LightDefTag == "" {
		return -1, fmt.Errorf("lightdef tag not set")
	}

	lightDef := wce.ByTag(e.LightDefTag)
	if lightDef == nil {
		return -1, fmt.Errorf("lightdef %s not found", e.LightDefTag)
	}

	lightDefRef, err := lightDef.ToRaw(wce, rawWld)
	if err != nil {
		return -1, fmt.Errorf("lightdef %s to raw: %w", e.LightDefTag, err)
	}

	wfLightInstance := &rawfrag.WldFragLight{
		NameRef:     0,
		LightDefRef: int32(lightDefRef),
		Flags:       0,
	}

	rawWld.Fragments = append(rawWld.Fragments, wfLightInstance)

	lightInstRef := int32(len(rawWld.Fragments))

	light := &rawfrag.WldFragPointLight{
		NameRef:  rawWld.NameAdd(e.Tag),
		LightRef: int32(lightInstRef),
		Location: e.Location,
		Radius:   e.Radius,
	}

	if e.Static == 1 {
		light.Flags |= 0x20
	}

	if e.StaticInfluence == 1 {
		light.Flags |= 0x40
	}

	if e.HasRegions == 1 {
		light.Flags |= 0x80
	}

	rawWld.Fragments = append(rawWld.Fragments, light)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *PointLight) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragPointLight) error {
	if frag == nil {
		return fmt.Errorf("frag is not pointlight (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	if frag.LightRef > 0 {
		if len(rawWld.Fragments) < int(frag.LightRef) {
			return fmt.Errorf("light ref %d not found", frag.LightRef)
		}

		light, ok := rawWld.Fragments[frag.LightRef].(*rawfrag.WldFragLight)
		if !ok {
			return fmt.Errorf("light ref %d not found", frag.LightRef)
		}

		if len(rawWld.Fragments) < int(light.LightDefRef) {
			return fmt.Errorf("lightdef ref %d not found", light.LightDefRef)
		}

		lightDef, ok := rawWld.Fragments[light.LightDefRef].(*rawfrag.WldFragLightDef)
		if !ok {
			return fmt.Errorf("lightdef ref %d not found", light.LightDefRef)
		}

		e.LightDefTag = rawWld.Name(lightDef.NameRef)
	}
	e.Location = frag.Location
	e.Radius = frag.Radius

	if frag.Flags&0x20 == 0x20 {
		e.Static = 1
	}

	if frag.Flags&0x40 == 0x40 {
		e.StaticInfluence = 1
	}

	if frag.Flags&0x80 == 0x80 {
		e.HasRegions = 1
	}

	return nil
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

func (e *Sprite3DDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	if token.TagIsWritten(e.Tag) {
		return nil
	}
	token.TagSetIsWritten(e.Tag)
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tCENTEROFFSET? %s\n", wcVal(e.CenterOffset))
	fmt.Fprintf(w, "\tBOUNDINGRADIUS? %s\n", wcVal(e.BoundingRadius))
	fmt.Fprintf(w, "\tSPHERELIST \"%s\"\n", e.SphereListTag)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMBSPNODES %d\n", len(e.BSPNodes))
	for i, node := range e.BSPNodes {
		fmt.Fprintf(w, "\t\tBSPNODE //%d\n", i)
		fmt.Fprintf(w, "\t\t\tVERTEXLIST %d", len(node.Vertices))
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
		fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func (s *Sprite3DDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("CENTEROFFSET?", 3)
	if err != nil {
		return err
	}
	err = parse(&s.CenterOffset, records[1:]...)
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = token.ReadProperty("BOUNDINGRADIUS?", 1)
	if err != nil {
		return err
	}
	err = parse(&s.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = token.ReadProperty("SPHERELIST", 1)
	if err != nil {
		return err
	}
	s.SphereListTag = records[1]

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices := int(0)
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

		s.Vertices = append(s.Vertices, vert)
	}

	records, err = token.ReadProperty("NUMBSPNODES", 1)
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
		_, err = token.ReadProperty("BSPNODE", 0)
		if err != nil {
			return err
		}
		records, err = token.ReadProperty("VERTEXLIST", -1)
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

		records, err = token.ReadProperty("RENDERMETHOD", 1)
		if err != nil {
			return err
		}

		node.RenderMethod = records[1]

		_, err = token.ReadProperty("RENDERINFO", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("PEN?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.Pen, records[1])
		if err != nil {
			return fmt.Errorf("render pen: %w", err)
		}

		records, err = token.ReadProperty("BRIGHTNESS?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.Brightness, records[1])
		if err != nil {
			return fmt.Errorf("render brightness: %w", err)
		}

		records, err = token.ReadProperty("SCALEDAMBIENT?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.ScaledAmbient, records[1])
		if err != nil {
			return fmt.Errorf("render scaled ambient: %w", err)
		}

		records, err = token.ReadProperty("SPRITE?", 1)
		if err != nil {
			return err
		}
		err = parse(&node.SpriteTag, records[1])
		if err != nil {
			return fmt.Errorf("render sprite: %w", err)
		}

		records, err = token.ReadProperty("UVORIGIN?", 3)
		if err != nil {
			return err
		}
		err = parse(&node.UvOrigin, records[1:]...)
		if err != nil {
			return fmt.Errorf("render uv origin: %w", err)
		}

		records, err = token.ReadProperty("UAXIS?", 3)
		if err != nil {
			return err
		}
		err = parse(&node.UAxis, records[1:]...)
		if err != nil {
			return fmt.Errorf("render u axis: %w", err)
		}

		records, err = token.ReadProperty("VAXIS?", 3)
		if err != nil {
			return err
		}
		err = parse(&node.VAxis, records[1:]...)
		if err != nil {
			return fmt.Errorf("render v axis: %w", err)
		}

		records, err = token.ReadProperty("UVCOUNT", 1)
		if err != nil {
			return err
		}
		numUVs := int(0)
		err = parse(&numUVs, records[1])
		if err != nil {
			return fmt.Errorf("num uvs: %w", err)
		}

		for j := 0; j < numUVs; j++ {
			records, err = token.ReadProperty("UV", 2)
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

		records, err = token.ReadProperty("TWOSIDED", 1)
		if err != nil {
			return err
		}
		err = parse(&node.TwoSided, records[1])
		if err != nil {
			return fmt.Errorf("two sided: %w", err)
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

		s.BSPNodes = append(s.BSPNodes, node)
	}

	return nil
}

func (e *Sprite3DDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
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
				bnode.RenderSimpleSpriteReference = uint32(rawWld.NameAdd(node.SpriteTag.String))
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

	wfSprite3DDef.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfSprite3DDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *Sprite3DDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragSprite3DDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not sprite3ddef (wrong fragcode?)")
	}

	if len(rawWld.Fragments) < int(frag.SphereListRef) {
		return fmt.Errorf("spherelist ref %d out of bounds", frag.SphereListRef)
	}

	if frag.SphereListRef > 0 {
		sphereList, ok := rawWld.Fragments[frag.SphereListRef].(*rawfrag.WldFragSphereList)
		if !ok {
			return fmt.Errorf("spherelist ref %d not found", frag.SphereListRef)
		}
		e.SphereListTag = rawWld.Name(sphereList.NameRef)
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.Vertices = frag.Vertices

	if frag.Flags&0x01 == 0x01 {
		e.CenterOffset.Valid = true
		e.CenterOffset.Float32Slice3 = frag.CenterOffset
	}

	if frag.Flags&0x02 == 0x02 {
		e.BoundingRadius.Valid = true
		e.BoundingRadius.Float32 = frag.BoundingRadius
	}

	for _, bspNode := range frag.BspNodes {
		node := &BSPNode{
			FrontTree:    bspNode.FrontTree,
			BackTree:     bspNode.BackTree,
			Vertices:     bspNode.VertexIndexes,
			RenderMethod: model.RenderMethodStr(bspNode.RenderMethod),
		}

		if bspNode.RenderFlags&0x01 == 0x01 {
			node.Pen.Valid = true
			node.Pen.Uint32 = bspNode.RenderPen
		}

		if bspNode.RenderFlags&0x02 == 0x02 {
			node.Brightness.Valid = true
			node.Brightness.Float32 = bspNode.RenderBrightness
		}

		if bspNode.RenderFlags&0x04 == 0x04 {
			node.ScaledAmbient.Valid = true
			node.ScaledAmbient.Float32 = bspNode.RenderScaledAmbient
		}

		if bspNode.RenderFlags&0x08 == 0x08 {
			node.SpriteTag.Valid = true
			if len(rawWld.Fragments) < int(bspNode.RenderSimpleSpriteReference) {
				return fmt.Errorf("sprite ref %d not found", bspNode.RenderSimpleSpriteReference)
			}
			spriteDef := rawWld.Fragments[bspNode.RenderSimpleSpriteReference]
			switch simpleSprite := spriteDef.(type) {
			case *rawfrag.WldFragSimpleSpriteDef:
				node.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
			case *rawfrag.WldFragDMSpriteDef:
				node.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
			case *rawfrag.WldFragHierarchicalSpriteDef:
				node.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
			case *rawfrag.WldFragSprite2D:
				node.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
			default:
				return fmt.Errorf("unhandled render sprite reference fragment type %d", spriteDef.FragCode())
			}
		}

		if bspNode.RenderFlags&0x10 == 0x10 {
			// has uvinfo
			node.UvOrigin.Valid = true
			node.UAxis.Valid = true
			node.VAxis.Valid = true
			node.UvOrigin.Float32Slice3 = bspNode.RenderUVInfoOrigin
			node.UAxis.Float32Slice3 = bspNode.RenderUVInfoUAxis
			node.VAxis.Float32Slice3 = bspNode.RenderUVInfoVAxis
		}

		if bspNode.RenderFlags&0x20 == 0x20 {
			node.Uvs = bspNode.Uvs
		}

		if bspNode.RenderFlags&0x40 == 0x40 {
			node.TwoSided = 1
		}

		e.BSPNodes = append(e.BSPNodes, node)
	}

	return nil
}

type PolyhedronDefinition struct {
	fragID         int16
	Tag            string
	BoundingRadius float32
	ScaleFactor    float32
	Vertices       [][3]float32
	Faces          [][]uint32
	HexOneFlag     int
}

type PolyhedronDefinitionFace struct {
	Vertices []uint32
}

func (e *PolyhedronDefinition) Definition() string {
	return "POLYHEDRONDEFINITION"
}

func (e *PolyhedronDefinition) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.8e\n", e.BoundingRadius)
	fmt.Fprintf(w, "\tSCALEFACTOR %0.8e\n", e.ScaleFactor)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(e.Vertices))
	for _, vert := range e.Vertices {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMFACES %d\n", len(e.Faces))
	for _, faces := range e.Faces {
		fmt.Fprintf(w, "\t\tVERTEXLIST %d", len(faces))
		for _, face := range faces {
			fmt.Fprintf(w, " %d", face)
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "\tHEXONEFLAG %d\n", e.HexOneFlag)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *PolyhedronDefinition) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = token.ReadProperty("SCALEFACTOR", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ScaleFactor, records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	records, err = token.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}

	numVertices := int(0)
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

	records, err = token.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}
	numFaces := int(0)
	err = parse(&numFaces, records[1])
	if err != nil {
		return fmt.Errorf("num faces: %w", err)
	}

	for i := 0; i < numFaces; i++ {
		records, err = token.ReadProperty("VERTEXLIST", -1)
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
		faceVals := []uint32{}
		for j := 0; j < numVertices; j++ {
			val := uint32(0)
			err = parse(&val, records[j+2])
			if err != nil {
				return fmt.Errorf("vertex %d: %w", j, err)
			}
			faceVals = append(faceVals, val)
		}
		e.Faces = append(e.Faces, faceVals)
	}

	records, err = token.ReadProperty("HEXONEFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexOneFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex one flag: %w", err)
	}

	return nil
}

func (e *PolyhedronDefinition) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfPolyhedronDef := &rawfrag.WldFragPolyhedronDef{
		NameRef:        rawWld.NameAdd(e.Tag),
		BoundingRadius: e.BoundingRadius,
		ScaleFactor:    e.ScaleFactor,
		Vertices:       e.Vertices,
		Faces:          e.Faces,
	}

	if e.HexOneFlag > 0 {
		wfPolyhedronDef.Flags |= 0x01
	}

	rawWld.Fragments = append(rawWld.Fragments, wfPolyhedronDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *PolyhedronDefinition) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragPolyhedronDef) error {
	e.Tag = rawWld.Name(frag.NameRef)
	e.BoundingRadius = frag.BoundingRadius
	e.ScaleFactor = frag.ScaleFactor
	e.Vertices = frag.Vertices
	e.Faces = frag.Faces
	if frag.Flags&0x01 != 0 {
		e.HexOneFlag = 1
	}

	return nil
}

type TrackInstance struct {
	fragID             int16
	animation          string
	model              string
	Tag                string
	TagIndex           int
	SpriteTag          string
	DefinitionTag      string
	DefinitionTagIndex int
	Interpolate        int
	Reverse            int
	Sleep              NullUint32
}

func (e *TrackInstance) Definition() string {
	return "TRACKINSTANCE"
}

func (e *TrackInstance) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	if token.wce.isTrackAni(e.Tag) {
		w, err = token.UseTempWriter(e.model + "_ani")
		if err != nil {
			return err
		}
	}
	if token.TagIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.TagIndex)) {
		return nil
	}

	token.TagSetIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.TagIndex))
	if e.DefinitionTag != "" {
		trackDef := token.wce.ByTagWithIndex(e.DefinitionTag, e.DefinitionTagIndex)
		if trackDef == nil {
			return fmt.Errorf("track %s%d refers to trackdef %s%d but it does not exist", e.Tag, e.TagIndex, e.DefinitionTag, e.DefinitionTagIndex)
		}
		err = trackDef.Write(token)
		if err != nil {
			return fmt.Errorf("trackdef %s%d write: %w", e.DefinitionTag, e.DefinitionTagIndex, err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tTAGINDEX %d\n", e.TagIndex)
	fmt.Fprintf(w, "\tSPRITE \"%s_DMSPRITEDEF\"\n", e.SpriteTag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", e.DefinitionTag)
	fmt.Fprintf(w, "\tDEFINITIONINDEX %d\n", e.DefinitionTagIndex)
	fmt.Fprintf(w, "\tINTERPOLATE %d // deprecated\n", e.Interpolate)
	fmt.Fprintf(w, "\tREVERSE %d // deprecated \n", e.Reverse)
	fmt.Fprintf(w, "\tSLEEP? %s\n", wcVal(e.Sleep))
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *TrackInstance) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("TAGINDEX", 1)
	if err != nil {
		return err
	}
	err = parse(&e.TagIndex, records[1])
	if err != nil {
		return fmt.Errorf("tag index: %w", err)
	}

	records, err = token.ReadProperty("SPRITE", 1)
	if err != nil {
		return err
	}
	e.SpriteTag = strings.TrimSuffix(records[1], "_DMSPRITEDEF")

	records, err = token.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.DefinitionTag = records[1]

	records, err = token.ReadProperty("DEFINITIONINDEX", 1)
	if err != nil {
		return err
	}
	err = parse(&e.DefinitionTagIndex, records[1])
	if err != nil {
		return fmt.Errorf("definition tag index: %w", err)
	}

	records, err = token.ReadProperty("INTERPOLATE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Interpolate, records[1])
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

	records, err = token.ReadProperty("REVERSE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Reverse, records[1])
	if err != nil {
		return fmt.Errorf("reverse: %w", err)
	}

	records, err = token.ReadProperty("SLEEP?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	return nil
}

func (e *TrackInstance) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfTrack := &rawfrag.WldFragTrack{}

	if e.DefinitionTag == "" {
		return -1, fmt.Errorf("track instance %s has no definition", e.Tag)
	}

	trackDefFrag := wce.ByTagWithIndex(e.DefinitionTag, e.DefinitionTagIndex)
	if trackDefFrag == nil {
		return -1, fmt.Errorf("track instance %s refers to trackdef %s but it does not exist", e.Tag, e.DefinitionTag)
	}

	trackDef, ok := trackDefFrag.(*TrackDef)
	if !ok {
		return -1, fmt.Errorf("track instance %s refers to trackdef %s but it is not a trackdef", e.Tag, e.DefinitionTag)
	}

	trackDefRef, err := trackDef.ToRaw(wce, rawWld)
	if err != nil {
		return -1, fmt.Errorf("track instance %s refers to trackdef %s but it failed to convert: %w", e.Tag, e.DefinitionTag, err)
	}

	wfTrack.NameRef = rawWld.NameAdd(e.Tag)
	wfTrack.TrackRef = int32(trackDefRef)
	if e.Sleep.Valid {
		wfTrack.Flags |= 0x01
		wfTrack.Sleep = e.Sleep.Uint32
	}
	if e.Reverse > 0 {
		wfTrack.Flags |= 0x02
	}
	if e.Interpolate > 0 {
		wfTrack.Flags |= 0x04
	}

	rawWld.Fragments = append(rawWld.Fragments, wfTrack)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *TrackInstance) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragTrack) error {
	if frag == nil {
		return fmt.Errorf("frag is not track instance (wrong fragcode?)")
	}

	if len(rawWld.Fragments) < int(frag.TrackRef) {
		return fmt.Errorf("trackdef ref %d out of bounds", frag.TrackRef)
	}

	trackDef, ok := rawWld.Fragments[frag.TrackRef].(*rawfrag.WldFragTrackDef)
	if !ok {
		return fmt.Errorf("trackdef ref %d not found", frag.TrackRef)
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.TagIndex = wce.NextTagIndex(e.Tag)
	var seqNum int
	e.SpriteTag, _, _, seqNum = wce.trackTagAndSequence(e.Tag)
	if seqNum < 0 {
		e.SpriteTag = wce.lastReadModelTag
	}

	// Use original logic if isObj is true
	if wce.isObj {
		e.model = wce.lastReadModelTag
		e.animation = "" // Ensure animation is blank for non-animation tracks
	} else if wce.isTrackAni(e.Tag) {
		// Parse animation and model codes
		e.animation, e.model = wce.trackAnimationParse(e.Tag)
	} else {
		// Fallback to original logic
		e.model = wce.lastReadModelTag
		e.animation = "" // Ensure animation is blank for non-animation tracks
	}

	e.DefinitionTag = rawWld.Name(trackDef.NameRef)
	e.DefinitionTagIndex = wce.tagIndexes[e.DefinitionTag]

	if frag.Flags&0x01 == 0x01 {
		e.Sleep.Valid = true
		e.Sleep.Uint32 = frag.Sleep
	}
	if frag.Flags&0x02 == 0x02 {
		e.Reverse = 1
	}
	if frag.Flags&0x04 == 0x04 {
		e.Interpolate = 1
	}

	return nil
}

type TrackDef struct {
	fragID       int16
	model        string
	Tag          string
	TagIndex     int
	SpriteTag    string
	Frames       []*Frame
	LegacyFrames []*LegacyFrame
}

type Frame struct {
	XYZScale int16
	XYZ      [3]int16
	RotScale int16
	Rotation [3]int16
}

type LegacyFrame struct {
	XYZScale int16
	XYZ      [3]int16
	Rotation [4]float32
}

func (e *TrackDef) Definition() string {
	return "TRACKDEFINITION"
}

func (e *TrackDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if strings.HasSuffix(e.Tag, "_TRACKDEF") {
		trackInstVal := token.wce.ByTagWithIndex(strings.ReplaceAll(e.Tag, "_TRACKDEF", "_TRACK"), e.TagIndex)
		if trackInstVal == nil {
			return fmt.Errorf("trackdef %s%d does not have a track instance", e.Tag, e.TagIndex)
		}
		trackInst, ok := trackInstVal.(*TrackInstance)
		if !ok {
			return fmt.Errorf("trackdef %s%d does not have a track instance", e.Tag, e.TagIndex)
		}

		if token.wce.isTrackAni(trackInst.Tag) {
			w, err = token.UseTempWriter(e.model + "_ani")
			if err != nil {
				return err
			}
		}
	}

	if token.TagIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.TagIndex)) {
		return nil
	}

	token.TagSetIsWritten(fmt.Sprintf("%s_%d", e.Tag, e.TagIndex))

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tTAGINDEX %d\n", e.TagIndex)
	fmt.Fprintf(w, "\tSPRITE \"%s_DMSPRITEDEF\"\n", e.SpriteTag)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.Frames))
	for _, frame := range e.Frames {
		fmt.Fprintf(w, "\t\tFRAME %d %d %d %d ", frame.XYZScale, frame.XYZ[0], frame.XYZ[1], frame.XYZ[2])
		fmt.Fprintf(w, "%d %d %d %d\n", frame.RotScale, frame.Rotation[0], frame.Rotation[1], frame.Rotation[2])
	}
	fmt.Fprintf(w, "\tNUMLEGACYFRAMES %d\n", len(e.LegacyFrames))
	for _, frame := range e.LegacyFrames {
		fmt.Fprintf(w, "\t\tLEGACYFRAME %d %d %d %d ", frame.XYZScale, frame.XYZ[0], frame.XYZ[1], frame.XYZ[2])
		fmt.Fprintf(w, "%0.8f %0.8f %0.8f %0.8f\n", frame.Rotation[0], frame.Rotation[1], frame.Rotation[2], frame.Rotation[3])
	}
	fmt.Fprintf(w, "\n")

	return nil
}

func (e *TrackDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("TAGINDEX", 1)
	if err != nil {
		return err
	}
	err = parse(&e.TagIndex, records[1])
	if err != nil {
		return fmt.Errorf("tag index: %w", err)
	}

	records, err = token.ReadProperty("SPRITE", 1)
	if err != nil {
		return err
	}
	e.SpriteTag = strings.TrimSuffix(records[1], "_DMSPRITEDEF")

	records, err = token.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return err
	}
	numFrames := int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		frame := &Frame{}
		records, err = token.ReadProperty("FRAME", -1)
		if err != nil {
			return err
		}
		if len(records) != 9 {
			return fmt.Errorf("frame: expected 9, got %d", len(records))
		}

		err = parse(&frame.XYZScale, records[1])
		if err != nil {
			return fmt.Errorf("xyz scale: %w", err)
		}

		err = parse(&frame.XYZ, records[2:5]...)
		if err != nil {
			return fmt.Errorf("xyz: %w", err)
		}

		err = parse(&frame.RotScale, records[5])
		if err != nil {
			return fmt.Errorf("rot scale: %w", err)
		}

		err = parse(&frame.Rotation, records[6:9]...)
		if err != nil {
			return fmt.Errorf("rotabc: %w", err)
		}

		e.Frames = append(e.Frames, frame)
	}

	records, err = token.ReadProperty("NUMLEGACYFRAMES", 1)
	if err != nil {
		return err
	}
	numFrames = int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num legacy frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		frame := &LegacyFrame{}
		records, err = token.ReadProperty("LEGACYFRAME", -1)
		if err != nil {
			return err
		}

		if len(records) != 9 {
			return fmt.Errorf("legacy frame: expected 9, got %d", len(records))
		}

		err = parse(&frame.XYZScale, records[1])
		if err != nil {
			return fmt.Errorf("xyz scale: %w", err)
		}

		err = parse(&frame.XYZ, records[2:5]...)
		if err != nil {
			return fmt.Errorf("xyz: %w", err)
		}

		err = parse(&frame.Rotation, records[5:9]...)
		if err != nil {
			return fmt.Errorf("rotabc: %w", err)
		}

		e.LegacyFrames = append(e.LegacyFrames, frame)
	}

	return nil
}

func (e *TrackDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	//if e.fragID != 0 {
	//	return e.fragID, nil
	//}

	wfTrack := &rawfrag.WldFragTrackDef{}

	for _, frame := range e.Frames {
		wfFrame := rawfrag.WldFragTrackBoneTransform{
			ShiftDenominator: frame.XYZScale,
		}

		wfFrame.Shift = frame.XYZ

		wfFrame.RotateDenominator = frame.RotScale

		wfTrack.Flags |= 0x08

		wfFrame.Rotation = [4]int16{
			frame.Rotation[0],
			frame.Rotation[1],
			frame.Rotation[2],
			0,
		}

		wfTrack.FrameTransforms = append(wfTrack.FrameTransforms, wfFrame)
	}

	for _, frame := range e.LegacyFrames {
		wfFrame := rawfrag.WldFragTrackBoneTransform{
			ShiftDenominator: frame.XYZScale,
		}

		wfFrame.Shift = frame.XYZ

		wfFrame.Rotation = [4]int16{
			int16(frame.Rotation[0]),
			int16(frame.Rotation[1]),
			int16(frame.Rotation[2]),
			int16(frame.Rotation[3]),
		}

		wfTrack.FrameTransforms = append(wfTrack.FrameTransforms, wfFrame)
	}

	wfTrack.NameRef = rawWld.NameAdd(e.Tag)
	rawWld.Fragments = append(rawWld.Fragments, wfTrack)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *TrackDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragTrackDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not trackdef (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.TagIndex = wce.NextTagIndex(e.Tag)
	var seqNum int
	e.SpriteTag, _, _, seqNum = wce.trackTagAndSequence(e.Tag)
	if seqNum < 0 {
		e.SpriteTag = wce.lastReadModelTag
	}
	e.model = wce.lastReadModelTag

	for _, fragFrame := range frag.FrameTransforms {
		//scale := float32(1.0)
		/* if fragFrame.ShiftDenominator > 0 {
			scale = 1.0 / float32(int(1<<fragFrame.ShiftDenominator))
		}
		if fragFrame.ShiftDenominator < 0 {
			scale = float32(math.Pow(2, float64(-fragFrame.ShiftDenominator)))
		} */
		//frame.XYZ = fragFrame.Shift

		if frag.Flags&0x08 != 0 {
			frame := &Frame{
				XYZScale: fragFrame.ShiftDenominator,
				XYZ:      fragFrame.Shift,
			}

			frame.RotScale = fragFrame.RotateDenominator
			//if fragFrame.RotateDenominator > 0 {
			//	scale = 1.0 / float32(int(1<<fragFrame.RotateDenominator))
			//}
			frame.Rotation = [3]int16{
				fragFrame.Rotation[0],
				fragFrame.Rotation[1],
				fragFrame.Rotation[2],
			}
			e.Frames = append(e.Frames, frame)
		} else {
			frame := &LegacyFrame{
				XYZScale: fragFrame.ShiftDenominator,
				XYZ:      fragFrame.Shift,
			}
			frame.Rotation = [4]float32{
				float32(fragFrame.Rotation[0]),
				float32(fragFrame.Rotation[1]),
				float32(fragFrame.Rotation[2]),
				float32(fragFrame.Rotation[3]),
			}
			e.LegacyFrames = append(e.LegacyFrames, frame)
		}
	}

	return nil
}

type HierarchicalSpriteDef struct {
	fragID                int16
	Tag                   string
	Dags                  []Dag
	AttachedSkins         []AttachedSkin
	CenterOffset          NullFloat32Slice3 // 0x01
	BoundingRadius        NullFloat32       // 0x02
	HexTwoHundredFlag     int               // 0x200
	HexTwentyThousandFlag int               // 0x20000
	PolyhedronTag         string
}

type Dag struct {
	Tag        string
	Track      string
	TrackIndex int
	SubDags    []uint32
	SpriteTag  string
}

type AttachedSkin struct {
	DMSpriteTag               string
	LinkSkinUpdatesToDagIndex uint32
}

func (e *HierarchicalSpriteDef) Definition() string {
	return "HIERARCHICALSPRITEDEF"
}

func (e *HierarchicalSpriteDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if token.TagIsWritten(e.Tag) {
		return nil
	}

	token.TagSetIsWritten(e.Tag)

	/* for _, skin := range e.AttachedSkins {

	} */

	for _, dag := range e.Dags {
		if dag.Track != "" {

			trackDef := token.wce.ByTagWithIndex(dag.Track, dag.TrackIndex)
			if trackDef == nil {
				return fmt.Errorf("track %s_%d not found", dag.Track, dag.TrackIndex)
			}

			err = trackDef.Write(token)
			if err != nil {
				return fmt.Errorf("track %s_%d: %w", dag.Track, dag.TrackIndex, err)
			}
		}
		if dag.SpriteTag != "" {
			spriteDef := token.wce.ByTag(dag.SpriteTag)
			if spriteDef == nil {
				return fmt.Errorf("sprite %s not found", dag.SpriteTag)
			}

			if token.TagIsWritten(dag.SpriteTag) {
				continue
			}

			err = spriteDef.Write(token)
			if err != nil {
				return fmt.Errorf("sprite %s: %w", dag.SpriteTag, err)
			}
			token.TagSetIsWritten(dag.SpriteTag)
		}
	}

	for _, skin := range e.AttachedSkins {
		if skin.DMSpriteTag == "" {
			continue
		}

		dmSprite := token.wce.ByTag(skin.DMSpriteTag)
		err = dmSprite.Write(token)
		if err != nil {
			return fmt.Errorf("dmsprite %s: %w", skin.DMSpriteTag, err)
		}
	}

	if e.PolyhedronTag != "" && e.PolyhedronTag != "SPECIAL_COLLISION" {
		polyhedronDef := token.wce.ByTag(e.PolyhedronTag)
		if polyhedronDef == nil {
			return fmt.Errorf("polyhedron %s not found", e.PolyhedronTag)
		}

		err = polyhedronDef.Write(token)
		if err != nil {
			return fmt.Errorf("polyhedron %s: %w", e.PolyhedronTag, err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tNUMDAGS %d\n", len(e.Dags))
	for i, dag := range e.Dags {
		fmt.Fprintf(w, "\t\tDAG // %d\n", i)
		fmt.Fprintf(w, "\t\t\tTAG \"%s\"\n", dag.Tag)
		fmt.Fprintf(w, "\t\t\tSPRITE \"%s\"\n", dag.SpriteTag)
		fmt.Fprintf(w, "\t\t\tTRACK \"%s\"\n", dag.Track)
		fmt.Fprintf(w, "\t\t\tTRACKINDEX %d\n", dag.TrackIndex)
		fmt.Fprintf(w, "\t\t\tSUBDAGLIST %d", len(dag.SubDags))
		for _, subDag := range dag.SubDags {
			fmt.Fprintf(w, " %d", subDag)
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMATTACHEDSKINS %d\n", len(e.AttachedSkins))

	for _, skin := range e.AttachedSkins {
		fmt.Fprintf(w, "\t\tATTACHEDSKIN\n")
		fmt.Fprintf(w, "\t\t\tDMSPRITE \"%s\"\n", skin.DMSpriteTag)
		fmt.Fprintf(w, "\t\t\tLINKSKINUPDATESTODAGINDEX %d\n", skin.LinkSkinUpdatesToDagIndex)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\tPOLYHEDRON\n")
	fmt.Fprintf(w, "\t\tDEFINITION \"%s\" // refer to polyhedron tag, or SPECIAL_COLLISION = 4294967293\n", e.PolyhedronTag)

	fmt.Fprintf(w, "\tCENTEROFFSET? %s\n", wcVal(e.CenterOffset))
	fmt.Fprintf(w, "\tBOUNDINGRADIUS? %s\n", wcVal(e.BoundingRadius))
	fmt.Fprintf(w, "\tHEXTWOHUNDREDFLAG %d\n", e.HexTwoHundredFlag)
	fmt.Fprintf(w, "\tHEXTWENTYTHOUSANDFLAG %d\n", e.HexTwentyThousandFlag)

	fmt.Fprintf(w, "\n")
	return nil
}

func (e *HierarchicalSpriteDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("NUMDAGS", 1)
	if err != nil {
		return err
	}
	numDags := int(0)
	err = parse(&numDags, records[1])
	if err != nil {
		return fmt.Errorf("num dags: %w", err)
	}

	for i := 0; i < numDags; i++ {
		dag := Dag{}
		_, err = token.ReadProperty("DAG", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("TAG", 1)
		if err != nil {
			return err
		}
		dag.Tag = records[1]

		records, err = token.ReadProperty("SPRITE", 1)
		if err != nil {
			return err
		}
		dag.SpriteTag = records[1]

		records, err = token.ReadProperty("TRACK", 1)
		if err != nil {
			return err
		}
		dag.Track = records[1]

		records, err = token.ReadProperty("TRACKINDEX", 1)
		if err != nil {
			return err
		}
		err = parse(&dag.TrackIndex, records[1])
		if err != nil {
			return fmt.Errorf("track index: %w", err)
		}

		records, err = token.ReadProperty("SUBDAGLIST", -1)
		if err != nil {
			return err
		}
		numSubDags := int(0)
		err = parse(&numSubDags, records[1])
		if err != nil {
			return fmt.Errorf("num sub dags: %w", err)
		}
		if len(records) != numSubDags+2 {
			return fmt.Errorf("sub dag list: expected %d, got %d", numSubDags, len(records)-2)
		}
		for j := 0; j < numSubDags; j++ {
			val := uint32(0)
			err = parse(&val, records[j+2])
			if err != nil {
				return fmt.Errorf("sub dag %d: %w", j, err)
			}
			dag.SubDags = append(dag.SubDags, val)
		}

		e.Dags = append(e.Dags, dag)
	}

	records, err = token.ReadProperty("NUMATTACHEDSKINS", 1)
	if err != nil {
		return err
	}
	numAttachedSkins := int(0)
	err = parse(&numAttachedSkins, records[1])
	if err != nil {
		return fmt.Errorf("num attached skins: %w", err)
	}

	for i := 0; i < numAttachedSkins; i++ {
		_, err = token.ReadProperty("ATTACHEDSKIN", 0)
		if err != nil {
			return err
		}
		skin := AttachedSkin{}
		records, err = token.ReadProperty("DMSPRITE", 1)
		if err != nil {
			return err
		}
		skin.DMSpriteTag = records[1]

		records, err = token.ReadProperty("LINKSKINUPDATESTODAGINDEX", 1)
		if err != nil {
			return err
		}
		err = parse(&skin.LinkSkinUpdatesToDagIndex, records[1])
		if err != nil {
			return fmt.Errorf("link skin updates to dag index: %w", err)
		}

		e.AttachedSkins = append(e.AttachedSkins, skin)
	}

	_, err = token.ReadProperty("POLYHEDRON", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	e.PolyhedronTag = records[1]

	records, err = token.ReadProperty("CENTEROFFSET?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.CenterOffset, records[1:]...)
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = token.ReadProperty("BOUNDINGRADIUS?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = token.ReadProperty("HEXTWOHUNDREDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwoHundredFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex two hundred flag: %w", err)
	}

	records, err = token.ReadProperty("HEXTWENTYTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwentyThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex twenty thousand flag: %w", err)
	}

	return nil
}

func (e *HierarchicalSpriteDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfHierarchicalSpriteDef := &rawfrag.WldFragHierarchicalSpriteDef{}

	if e.PolyhedronTag == "SPECIAL_COLLISION" {
		wfHierarchicalSpriteDef.CollisionVolumeRef = 4294967293
	}
	if e.PolyhedronTag != "" &&
		e.PolyhedronTag != "SPECIAL_COLLISION" {
		collusionDef := wce.ByTag(e.PolyhedronTag)
		if collusionDef == nil {
			fmt.Printf("collision volume not found: %s\n", e.PolyhedronTag)
			//			return -1, fmt.Errorf("collision volume not found: %s", e.PolyhedronTag)
		}
		collisionTag := ""
		switch sprite := collusionDef.(type) {
		case *PolyhedronDefinition:
			collisionTag = sprite.Tag
			polyDefID, err := collusionDef.ToRaw(wce, rawWld)
			if err != nil {
				return -1, fmt.Errorf("collision volume to raw: %w", err)
			}

			wfPoly := &rawfrag.WldFragPolyhedron{
				FragmentRef: int32(polyDefID),
			}

			rawWld.Fragments = append(rawWld.Fragments, wfPoly)

			wfHierarchicalSpriteDef.CollisionVolumeRef = uint32(len(rawWld.Fragments))

		case *DMSpriteDef2: // chequip has this on EYE_HS_DEF
			collisionTag = sprite.Tag
			wfHierarchicalSpriteDef.CollisionVolumeRef = uint32(rawWld.NameAdd(collisionTag))
		case *DMSpriteDef:
			collisionTag = sprite.Tag
			wfHierarchicalSpriteDef.CollisionVolumeRef = uint32(rawWld.NameAdd(collisionTag))
		case nil:
		default:
			return -1, fmt.Errorf("unsupported collision volume type: %T", collusionDef)
		}

	}

	if e.CenterOffset.Valid {
		wfHierarchicalSpriteDef.Flags |= 0x01
		wfHierarchicalSpriteDef.CenterOffset = e.CenterOffset.Float32Slice3
	}

	if e.BoundingRadius.Valid {
		wfHierarchicalSpriteDef.Flags |= 0x02
		wfHierarchicalSpriteDef.BoundingRadius = e.BoundingRadius.Float32
	}

	if e.HexTwoHundredFlag > 0 {
		wfHierarchicalSpriteDef.Flags |= 0x200
	}
	if e.HexTwentyThousandFlag > 0 {
		wfHierarchicalSpriteDef.Flags |= 0x20000
	}

	dmSpriteInstances := []*rawfrag.WldFragDMSprite{}

	for _, dag := range e.Dags {
		wfDag := &rawfrag.WldFragDag{}

		if dag.SpriteTag != "" {
			spriteDefFrag := wce.ByTag(dag.SpriteTag)
			if spriteDefFrag == nil {
				return -1, fmt.Errorf("sprite instance not found: %s", dag.SpriteTag)
			}
			switch spriteDef := spriteDefFrag.(type) {
			case *SimpleSpriteDef:
				spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef to raw: %w", err)
				}

				wfSprite := &rawfrag.WldFragDMSprite{
					//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
					DMSpriteRef: int32(spriteDefRef),
					Params:      0,
				}

				rawWld.Fragments = append(rawWld.Fragments, wfSprite)
				wfDag.MeshOrSpriteOrParticleRef = uint32(len(rawWld.Fragments))
			case *DMSpriteDef:
				spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef to raw: %w", err)
				}

				wfSprite := &rawfrag.WldFragDMSprite{
					//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
					DMSpriteRef: int32(spriteDefRef),
					Params:      0,
				}

				rawWld.Fragments = append(rawWld.Fragments, wfSprite)
				wfDag.MeshOrSpriteOrParticleRef = uint32(len(rawWld.Fragments))
			case *HierarchicalSpriteDef:
				spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef to raw: %w", err)
				}

				wfSprite := &rawfrag.WldFragDMSprite{
					//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
					DMSpriteRef: int32(spriteDefRef),
					Params:      0,
				}

				rawWld.Fragments = append(rawWld.Fragments, wfSprite)
				wfDag.MeshOrSpriteOrParticleRef = uint32(len(rawWld.Fragments))
			case *Sprite3DDef:

				spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef to raw: %w", err)
				}

				wfSprite := &rawfrag.WldFragDMSprite{
					//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
					DMSpriteRef: int32(spriteDefRef),
					Params:      0,
				}

				rawWld.Fragments = append(rawWld.Fragments, wfSprite)
				wfDag.MeshOrSpriteOrParticleRef = uint32(len(rawWld.Fragments))
			case *DMSpriteDef2:
				spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
				if err != nil {
					return -1, fmt.Errorf("dmspritedef to raw: %w", err)
				}

				wfSprite := &rawfrag.WldFragDMSprite{
					//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
					DMSpriteRef: int32(spriteDefRef),
					Params:      0,
				}

				rawWld.Fragments = append(rawWld.Fragments, wfSprite)
				wfDag.MeshOrSpriteOrParticleRef = uint32(len(rawWld.Fragments))
			case *ParticleCloudDef:
				/*
					spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
					if err != nil {
						return -1, fmt.Errorf("particle to raw: %w", err)
					}
						wfSprite := &rawfrag.WldFragDMSprite{
							//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
							DMSpriteRef: int32(spriteDefRef),
							Params:      0,
						}

						rawWld.Fragments = append(rawWld.Fragments, wfSprite) */

				wfDag.MeshOrSpriteOrParticleRef = uint32(spriteDef.fragID)
			default:
				return -1, fmt.Errorf("unsupported toraw dag spritetag instance type: %T", spriteDefFrag)
			}
		}

		wfHierarchicalSpriteDef.Dags = append(wfHierarchicalSpriteDef.Dags, wfDag)
	}

	for _, skin := range e.AttachedSkins {
		wfHierarchicalSpriteDef.LinkSkinUpdatesToDagIndexes = append(wfHierarchicalSpriteDef.LinkSkinUpdatesToDagIndexes, skin.LinkSkinUpdatesToDagIndex)
		if skin.DMSpriteTag == "" {
			wfHierarchicalSpriteDef.DMSprites = append(wfHierarchicalSpriteDef.DMSprites, 0)
			continue
		}

		spriteDefFrag := wce.ByTag(skin.DMSpriteTag)
		if spriteDefFrag == nil {
			return -1, fmt.Errorf("skin sprite def not found: %s", skin.DMSpriteTag)
		}

		err := spriteVariationToRaw(wce, rawWld, spriteDefFrag)
		if err != nil {
			return -1, fmt.Errorf("sprite variation toraw: %w", err)
		}
		switch spriteDef := spriteDefFrag.(type) {
		case *DMSpriteDef2:
			spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
			if err != nil {
				return -1, fmt.Errorf("dmspritedef2 to raw: %w", err)
			}

			wfDMSprite := &rawfrag.WldFragDMSprite{
				//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
				DMSpriteRef: int32(spriteDefRef),
				Params:      0,
			}

			dmSpriteInstances = append(dmSpriteInstances, wfDMSprite)
		case *DMSpriteDef:
			spriteDefRef, err := spriteDef.ToRaw(wce, rawWld)
			if err != nil {
				return -1, fmt.Errorf("dmspritedef to raw: %w", err)
			}

			wfDMSprite := &rawfrag.WldFragDMSprite{
				//NameRef:     rawWld.NameAdd(skin.DMSpriteTag),
				DMSpriteRef: int32(spriteDefRef),
				Params:      0,
			}

			dmSpriteInstances = append(dmSpriteInstances, wfDMSprite)

		default:
			return -1, fmt.Errorf("unsupported toraw attachedskin sprite instance type: %T", spriteDefFrag)
		}

		//wfHierarchicalSpriteDef.DMSprites = append(wfHierarchicalSpriteDef.DMSprites, uint32(spriteRef))
	}

	for i, dag := range e.Dags {
		wfDag := wfHierarchicalSpriteDef.Dags[i]

		trackFrag := wce.ByTagWithIndex(dag.Track, dag.TrackIndex)
		if trackFrag == nil {
			return -1, fmt.Errorf("track not found: %s index %d", dag.Track, dag.TrackIndex)
		}

		track, ok := trackFrag.(*TrackInstance)
		if !ok {
			return -1, fmt.Errorf("invalid track type: %T", trackFrag)
		}

		trackRef, err := track.ToRaw(wce, rawWld)
		if err != nil {
			return -1, fmt.Errorf("track to raw: %w", err)
		}

		//wfDag.NameRef = rawWld.NameAdd(dag.Tag)

		wfDag.TrackRef = uint32(trackRef)
		wfDag.SubDags = dag.SubDags

	}

	for i, dag := range e.Dags {
		wfDag := wfHierarchicalSpriteDef.Dags[i]
		wfDag.NameRef = rawWld.NameAdd(dag.Tag)
	}
	wfHierarchicalSpriteDef.NameRef = rawWld.NameAdd(e.Tag)

	for _, wfDMSprite := range dmSpriteInstances {
		rawWld.Fragments = append(rawWld.Fragments, wfDMSprite)
		wfHierarchicalSpriteDef.DMSprites = append(wfHierarchicalSpriteDef.DMSprites, uint32(len(rawWld.Fragments)))
	}

	rawWld.Fragments = append(rawWld.Fragments, wfHierarchicalSpriteDef)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *HierarchicalSpriteDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragHierarchicalSpriteDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not hierarchical sprite def (wrong fragcode?)")
	}

	if frag.CollisionVolumeRef != 0 && frag.CollisionVolumeRef != 4294967293 {
		if len(rawWld.Fragments) < int(frag.CollisionVolumeRef) {
			return fmt.Errorf("collision volume ref %d out of bounds", frag.CollisionVolumeRef)
		}

		switch collision := rawWld.Fragments[frag.CollisionVolumeRef].(type) {
		case *rawfrag.WldFragPolyhedron:
			if len(rawWld.Fragments) < int(collision.FragmentRef) {
				return fmt.Errorf("collision def ref %d not found", collision.FragmentRef)
			}
			collisionFragDef := rawWld.Fragments[collision.FragmentRef]
			if collisionFragDef == nil {
				return fmt.Errorf("collision def ref %d not found", collision.FragmentRef)
			}

			collisionDef, ok := collisionFragDef.(*rawfrag.WldFragPolyhedronDef)
			if !ok {
				return fmt.Errorf("collision def ref type incorrect: %T", collisionFragDef)
			}
			e.PolyhedronTag = rawWld.Name(collisionDef.NameRef)
		default:
			return fmt.Errorf("unknown collision volume ref %d (%s)", frag.CollisionVolumeRef, raw.FragName(collision.FragCode()))
		}
	}
	if frag.CollisionVolumeRef == 4294967293 {
		e.PolyhedronTag = "SPECIAL_COLLISION"
	}
	e.Tag = rawWld.Name(frag.NameRef)
	if frag.Flags&0x01 != 0 {
		e.CenterOffset.Valid = true
		e.CenterOffset.Float32Slice3 = frag.CenterOffset
	}
	if frag.Flags&0x02 != 0 {
		e.BoundingRadius.Valid = true
		e.BoundingRadius.Float32 = frag.BoundingRadius
	}
	if frag.Flags&0x200 != 0 {
		e.HexTwoHundredFlag = 1
	}
	if frag.Flags&0x20000 != 0 {
		e.HexTwentyThousandFlag = 1
	}

	for i, dag := range frag.Dags {
		if len(rawWld.Fragments) < int(dag.TrackRef) {
			return fmt.Errorf("dag %d track ref %d not found", i, dag.TrackRef)
		}
		srcTrack, ok := rawWld.Fragments[dag.TrackRef].(*rawfrag.WldFragTrack)
		if !ok {
			return fmt.Errorf("dag %d track ref %d not found", i, dag.TrackRef)
		}

		spriteTag := ""
		if dag.MeshOrSpriteOrParticleRef > 0 {
			if len(rawWld.Fragments) < int(dag.MeshOrSpriteOrParticleRef) {
				return fmt.Errorf("dag %d mesh or sprite or particle ref %d not found", i, dag.MeshOrSpriteOrParticleRef)
			}

			spriteFrag := rawWld.Fragments[dag.MeshOrSpriteOrParticleRef]
			if spriteFrag == nil {
				return fmt.Errorf("dag %d mesh or sprite or particle ref %d not found", i, dag.MeshOrSpriteOrParticleRef)
			}

			spriteRef := int32(0)

			switch sprite := spriteFrag.(type) {
			case *rawfrag.WldFragDMSprite:
				spriteRef = sprite.DMSpriteRef
			case *rawfrag.WldFragParticleCloudDef:
				spriteTag = rawWld.Name(sprite.NameRef)
				//				spriteRef = int16(sprite.BlitSpriteRef)
			default:
				return fmt.Errorf("dag %d unhandled sprite instance or particle reference fragment type %d (%s)", i, spriteFrag.FragCode(), raw.FragName(spriteFrag.FragCode()))
			}
			if spriteTag == "" {
				if spriteRef < 0 {
					spriteRef = -spriteRef
				}

				if len(rawWld.Fragments) < int(spriteRef) {
					return fmt.Errorf("dag %d sprite instance/particle ref %d out of range", i, spriteRef)
				}

				spriteDef := rawWld.Fragments[spriteRef]
				switch simpleSprite := spriteDef.(type) {
				case *rawfrag.WldFragSimpleSpriteDef:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragDMSpriteDef:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragHierarchicalSpriteDef:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragSprite2D:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragDmSpriteDef2:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragBlitSpriteDef:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				case *rawfrag.WldFragBMInfo:
					spriteTag = rawWld.Name(simpleSprite.NameRef)
				default:
					return fmt.Errorf("dag %d unhandled mesh or sprite or particle reference fragment type %d (%s)", i, spriteDef.FragCode(), raw.FragName(spriteDef.FragCode()))
				}
			}
		}
		/* if spriteTag != "" && e.PolyhedronTag == "" {
			e.PolyhedronTag = spriteTag
		} */

		dag := Dag{
			Tag:        rawWld.Name(dag.NameRef),
			Track:      rawWld.Name(srcTrack.NameRef),
			TrackIndex: wce.tagIndexes[rawWld.Name(srcTrack.NameRef)],
			SubDags:    dag.SubDags,
			SpriteTag:  spriteTag,
		}

		e.Dags = append(e.Dags, dag)
	}

	// based on frag.Flags&0x100 == 0x100 {
	for i := 0; i < len(frag.DMSprites); i++ {
		dmSpriteTag := ""
		if len(rawWld.Fragments) < int(frag.DMSprites[i]) {
			return fmt.Errorf("dmsprite ref %d not found", frag.DMSprites[i])
		}
		dmSprite, ok := rawWld.Fragments[frag.DMSprites[i]].(*rawfrag.WldFragDMSprite)
		if !ok {
			return fmt.Errorf("dmsprite ref %d not found", frag.DMSprites[i])
		}
		if len(rawWld.Fragments) < int(dmSprite.DMSpriteRef) {
			return fmt.Errorf("dmsprite ref %d not found", dmSprite.DMSpriteRef)
		}
		switch spriteDef := rawWld.Fragments[dmSprite.DMSpriteRef].(type) {
		case *rawfrag.WldFragSimpleSpriteDef:
			dmSpriteTag = rawWld.Name(spriteDef.NameRef)
		case *rawfrag.WldFragDMSpriteDef:
			dmSpriteTag = rawWld.Name(spriteDef.NameRef)
		case *rawfrag.WldFragHierarchicalSpriteDef:
			dmSpriteTag = rawWld.Name(spriteDef.NameRef)
		case *rawfrag.WldFragSprite2D:
			dmSpriteTag = rawWld.Name(spriteDef.NameRef)
		case *rawfrag.WldFragDmSpriteDef2:
			dmSpriteTag = rawWld.Name(spriteDef.NameRef)
		default:
			return fmt.Errorf("unhandled dmsprite reference fragment type %d (%s) at offset %d", spriteDef.FragCode(), raw.FragName(spriteDef.FragCode()), i)
		}

		skin := AttachedSkin{
			DMSpriteTag:               dmSpriteTag,
			LinkSkinUpdatesToDagIndex: frag.LinkSkinUpdatesToDagIndexes[i],
		}

		e.AttachedSkins = append(e.AttachedSkins, skin)
	}

	return nil
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

func (e *WorldTree) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tNUMWORLDNODES %d\n", len(e.WorldNodes))
	for i, node := range e.WorldNodes {
		fmt.Fprintf(w, "\t\tWORLDNODE // %d\n", i)
		fmt.Fprintf(w, "\t\t\tNORMALABCD %0.8e %0.8e %0.8e %0.8e\n", node.Normals[0], node.Normals[1], node.Normals[2], node.Normals[3])
		fmt.Fprintf(w, "\t\t\tWORLDREGIONTAG \"%s\"\n", node.WorldRegionTag)
		fmt.Fprintf(w, "\t\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\t\tBACKTREE %d\n", node.BackTree)
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *WorldTree) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("NUMWORLDNODES", 1)
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
		_, err = token.ReadProperty("WORLDNODE", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("NORMALABCD", 4)
		if err != nil {
			return err
		}
		err = parse(&node.Normals, records[1:]...)
		if err != nil {
			return fmt.Errorf("normals: %w", err)
		}

		records, err = token.ReadProperty("WORLDREGIONTAG", 1)
		if err != nil {
			return err
		}
		node.WorldRegionTag = records[1]

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

		e.WorldNodes = append(e.WorldNodes, node)

	}

	return nil
}

func (e *WorldTree) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfWorldTree := &rawfrag.WldFragWorldTree{}

	for _, node := range e.WorldNodes {

		regionRef := int32(0)
		if node.WorldRegionTag != "" {
			regionFrag := wce.ByTag(node.WorldRegionTag)
			if regionFrag == nil {
				return -1, fmt.Errorf("region not found: %s", node.WorldRegionTag)
			}
			region, ok := regionFrag.(*Region)
			if !ok {
				return -1, fmt.Errorf("invalid region type: %T", regionFrag)
			}
			if !strings.HasPrefix(region.Tag, "R") {
				return -1, fmt.Errorf("invalid region tag (needs R Prefix): %s", region.Tag)
			}
			regionVal, err := strconv.Atoi(region.Tag[1:])
			if err != nil {
				return -1, fmt.Errorf("invalid region tag (Should be R########): %s", region.Tag)
			}
			regionRef = int32(regionVal)
		}
		wfNode := rawfrag.WorldTreeNode{
			Normal:    node.Normals,
			RegionRef: regionRef,
			FrontRef:  int32(node.FrontTree),
			BackRef:   int32(node.BackTree),
		}

		wfWorldTree.Nodes = append(wfWorldTree.Nodes, wfNode)
	}

	wfWorldTree.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfWorldTree)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *WorldTree) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragWorldTree) error {
	if frag == nil {
		return fmt.Errorf("frag is not world tree (wrong fragcode?)")
	}

	for _, srcNode := range frag.Nodes {
		regionTag := ""
		if srcNode.RegionRef > 0 {
			regionTag = fmt.Sprintf("R%06d", srcNode.RegionRef)
		}
		node := &WorldNode{
			Normals:        srcNode.Normal,
			WorldRegionTag: regionTag,
			FrontTree:      uint32(srcNode.FrontRef),
			BackTree:       uint32(srcNode.BackRef),
		}
		e.WorldNodes = append(e.WorldNodes, node)
	}
	return nil
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
	Ranges []byte
}

func (e *Region) Definition() string {
	return "REGION"
}

func (e *Region) Write(token *AsciiWriteToken) error {
	var err error

	if e.SpriteTag != "" {
		sprite := token.wce.ByTag(e.SpriteTag)
		if sprite == nil {
			return fmt.Errorf("sprite not found: %s", e.SpriteTag)
		}
		err = sprite.Write(token)
		if err != nil {
			return fmt.Errorf("sprite write: %w", err)
		}
	}
	err = token.SetWriter("region")
	if err != nil {
		return fmt.Errorf("set writer region: %w", err)
	}
	w, err := token.Writer()
	if err != nil {
		return fmt.Errorf("get writer region: %w", err)
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tREVERBVOLUME %0.8e\n", e.ReverbVolume)
	fmt.Fprintf(w, "\tREVERBOFFSET %d\n", e.ReverbOffset)
	fmt.Fprintf(w, "\tREGIONFOG %d\n", e.RegionFog)
	fmt.Fprintf(w, "\tGOURAND2 %d\n", e.Gouraud2)
	fmt.Fprintf(w, "\tENCODEDVISIBILITY %d\n", e.EncodedVisibility)
	fmt.Fprintf(w, "\tVISLISTBYTES %d\n", e.VisListBytes)
	fmt.Fprintf(w, "\tNUMREGIONVERTEX %d\n", len(e.RegionVertices))
	for _, vert := range e.RegionVertices {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMRENDERVERTICES %d\n", len(e.RenderVertices))
	for _, vert := range e.RenderVertices {
		fmt.Fprintf(w, "\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMWALLS %d\n", len(e.Walls))
	for i, wall := range e.Walls {
		fmt.Fprintf(w, "\t\tWALL // %d\n", i)
		fmt.Fprintf(w, "\t\t\tNORMALABCD %0.8e %0.8e %0.8e %0.8e\n", wall.Normal[0], wall.Normal[1], wall.Normal[2], wall.Normal[3])
		fmt.Fprintf(w, "\t\t\tNUMVERTICES %d\n", len(wall.Vertices))
		for _, vert := range wall.Vertices {
			fmt.Fprintf(w, "\t\t\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
		}
	}
	fmt.Fprintf(w, "\tNUMOBSTACLES %d\n", len(e.Obstacles))
	for i, obs := range e.Obstacles {
		fmt.Fprintf(w, "\t\tOBSTACLE // %d\n", i)
		fmt.Fprintf(w, "\t\t\tNORMALABCD %0.8e %0.8e %0.8e %0.8e\n", obs.Normal[0], obs.Normal[1], obs.Normal[2], obs.Normal[3])
		fmt.Fprintf(w, "\t\t\tNUMVERTICES %d\n", len(obs.Vertices))
		for _, vert := range obs.Vertices {
			fmt.Fprintf(w, "\t\t\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
		}
	}
	fmt.Fprintf(w, "\tNUMCUTTINGOBSTACLES %d\n", len(e.CuttingObstacles))
	for i, obs := range e.CuttingObstacles {
		fmt.Fprintf(w, "\t\tCUTTINGOBSTACLE // %d\n", i)
		fmt.Fprintf(w, "\t\t\tNORMALABCD %0.8e %0.8e %0.8e %0.8e\n", obs.Normal[0], obs.Normal[1], obs.Normal[2], obs.Normal[3])
		fmt.Fprintf(w, "\t\t\tNUMVERTICES %d\n", len(obs.Vertices))
		for _, vert := range obs.Vertices {
			fmt.Fprintf(w, "\t\t\t\tXYZ %0.8e %0.8e %0.8e\n", vert[0], vert[1], vert[2])
		}
	}
	fmt.Fprintf(w, "\tVISTREE\n")
	fmt.Fprintf(w, "\t\tNUMVISNODE %d\n", len(e.VisTree.VisNodes))
	for i, node := range e.VisTree.VisNodes {
		fmt.Fprintf(w, "\t\t\tVISNODE // %d\n", i)
		fmt.Fprintf(w, "\t\t\t\tNORMALABCD %0.8e %0.8e %0.8e %0.8e\n", node.Normal[0], node.Normal[1], node.Normal[2], node.Normal[3])
		fmt.Fprintf(w, "\t\t\t\tVISLISTINDEX %d\n", node.VisListIndex)
		fmt.Fprintf(w, "\t\t\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\t\t\tBACKTREE %d\n", node.BackTree)
	}
	fmt.Fprintf(w, "\t\tNUMVISIBLELIST %d\n", len(e.VisTree.VisLists))
	for i, list := range e.VisTree.VisLists {
		// Determine if the 0x80 flag is set using e.VisListBytes
		if e.VisListBytes == 1 {
			// Calculate visible regions from range values using RLE
			regions := []int{}
			currentRegion := 1

			for i < len(list.Ranges) {
				byteVal := list.Ranges[i]

				switch {
				case byteVal >= 0x00 && byteVal <= 0x3E:
					// Skip forward by this many region IDs
					currentRegion += int(byteVal)
				case byteVal == 0x3F:
					// Skip forward by the amount given in the following 16-bit WORD
					nextByte1 := list.Ranges[i+1]
					nextByte2 := list.Ranges[i+2]
					skipAmount := int(nextByte2)<<8 | int(nextByte1)
					currentRegion += skipAmount
					i += 2
				case byteVal >= 0x40 && byteVal <= 0x7F:
					// Skip forward based on bits 3..5, then include the number of IDs based on bits 0..2
					skipAmount := int((byteVal & 0b00111000) >> 3)
					includeCount := int(byteVal & 0b00000111)
					currentRegion += skipAmount
					for j := 0; j < includeCount; j++ {
						regions = append(regions, currentRegion)
						currentRegion++
					}
				case byteVal >= 0x80 && byteVal <= 0xBF:
					// Include the number of IDs based on bits 3..5, then skip forward based on bits 0..2
					includeCount := int((byteVal & 0b00111000) >> 3)
					skipAmount := int(byteVal & 0b00000111)
					for j := 0; j < includeCount; j++ {
						regions = append(regions, currentRegion)
						currentRegion++
					}
					currentRegion += skipAmount
				case byteVal >= 0xC0 && byteVal <= 0xFE:
					// Subtracting 0xC0, this many region IDs are nearby
					includeCount := int(byteVal - 0xC0)
					for j := 0; j < includeCount; j++ {
						regions = append(regions, currentRegion)
						currentRegion++
					}
				case byteVal == 0xFF:
					// Include regions by the amount given in the following 16-bit WORD
					nextByte1 := list.Ranges[i+1]
					nextByte2 := list.Ranges[i+2]
					includeCount := int(nextByte2)<<8 | int(nextByte1)
					for j := 0; j < includeCount; j++ {
						regions = append(regions, currentRegion)
						currentRegion++
					}
					i += 2
				}

				i++
			}

			// Print the REGIONS data for the 0x80 flag set case
			fmt.Fprintf(w, "\t\t\tVISLIST // %d\n", i)
			fmt.Fprintf(w, "\t\t\t\tREGIONS %d", len(regions))
			for _, region := range regions {
				fmt.Fprintf(w, " %d", region)
			}
			fmt.Fprintf(w, "\n")

		} else {
			// 0x80 flag is not set, handle as pairs of uint16 WORDs
			wordCount := len(list.Ranges) / 2
			regions := []int{}

			for j := 0; j < wordCount; j++ {
				lowByte := list.Ranges[j*2]
				highByte := list.Ranges[j*2+1]
				regionIndex := int(highByte)<<8 | int(lowByte)
				regionIndex += 1 // Convert from 0-based to 1-based indexing
				regions = append(regions, regionIndex)
			}

			// Print the REGIONS data for the non-0x80 flag case
			fmt.Fprintf(w, "\t\t\tVISLIST // %d\n", i)
			fmt.Fprintf(w, "\t\t\t\tREGIONS %d", len(regions))
			for _, region := range regions {
				fmt.Fprintf(w, " %d", region)
			}
			fmt.Fprintf(w, "\n")
		}
	}
	fmt.Fprintf(w, "\tSPHERE %0.8e %0.8e %0.8e %0.8e\n", e.Sphere[0], e.Sphere[1], e.Sphere[2], e.Sphere[3])
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", e.UserData)
	fmt.Fprintf(w, "\tSPRITE \"%s\"\n", e.SpriteTag)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *Region) Read(token *AsciiReadToken) error {
	e.VisTree = &VisTree{}
	records, err := token.ReadProperty("REVERBVOLUME", 1)
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

		// Determine if the 0x80 flag is set using e.VisListBytes
		if e.VisListBytes == 1 {
			// Handle the case where 0x80 is set
			records, err = token.ReadProperty("REGIONS", -1)
			if err != nil {
				return err
			}

			numRegions := int(0)
			err = parse(&numRegions, records[1])
			if err != nil {
				return fmt.Errorf("num regions: %w", err)
			}

			regions := []int{}
			if numRegions > 0 {
				regions = make([]int, numRegions)
				for j := 0; j < numRegions; j++ {
					err = parse(&regions[j], records[j+2])
					if err != nil {
						return fmt.Errorf("region %d: %w", j, err)
					}
				}
			}

			// Calculate groups of visible and not-visible regions
			groups := []struct {
				visible bool
				count   int
			}{}

			if len(regions) > 0 {
				currentRegion := 1
				groupStart := 1
				visible := regions[0] == currentRegion

				for currentRegion <= regions[len(regions)-1] {
					isVisible := false
					for _, region := range regions {
						if region == currentRegion {
							isVisible = true
							break
						}
					}

					if isVisible != visible {
						// Save the current group
						groups = append(groups, struct {
							visible bool
							count   int
						}{
							visible: visible,
							count:   currentRegion - groupStart,
						})
						// Update visibility and start of new group
						visible = isVisible
						groupStart = currentRegion
					}

					currentRegion++
				}

				// Save the final group
				groups = append(groups, struct {
					visible bool
					count   int
				}{
					visible: visible,
					count:   currentRegion - groupStart,
				})
			}

			// Write out the encoded bytes using RLE logic
			if len(regions) == 0 {
				// If there are no regions, still add an empty VisList with zero ranges
				list.Ranges = []byte{}
			} else {
				for g := 0; g < len(groups); g++ {
					group := groups[g]

					if group.visible {
						// Handle visible groups
						if g+1 < len(groups) && group.count <= 7 && !groups[g+1].visible && groups[g+1].count <= 7 {
							// If the current visible group is 1-7 and the next not-visible group is also 1-7
							list.Ranges = append(list.Ranges, byte(0x80|(group.count<<3)|groups[g+1].count))
							g++ // Skip the next group since it's combined with this one
						} else if group.count <= 62 {
							// If the group is visible and 1-62 in size
							list.Ranges = append(list.Ranges, byte(0xC0+group.count))
						} else {
							// If the group is visible and 63 or greater in size
							list.Ranges = append(list.Ranges, 0xFF)
							list.Ranges = append(list.Ranges, byte(group.count&0xFF))      // Lower byte
							list.Ranges = append(list.Ranges, byte((group.count>>8)&0xFF)) // Upper byte
						}
					} else {
						// Handle not-visible groups
						if g+1 < len(groups) && group.count <= 7 && groups[g+1].visible && groups[g+1].count <= 7 {
							// If the current not-visible group is 1-7 and the next visible group is also 1-7
							list.Ranges = append(list.Ranges, byte(0x40|(group.count<<3)|groups[g+1].count))
							g++ // Skip the next group since it's combined with this one
						} else if group.count <= 62 {
							// If the group is not-visible and 1-62 in size
							list.Ranges = append(list.Ranges, byte(group.count))
						} else {
							// If the group is not-visible and 63 or greater in size
							list.Ranges = append(list.Ranges, 0x3F)
							list.Ranges = append(list.Ranges, byte(group.count&0xFF))      // Lower byte
							list.Ranges = append(list.Ranges, byte((group.count>>8)&0xFF)) // Upper byte
						}
					}
				}
			}
		} else {
			// Handle the case where 0x80 is not set - read as pairs of WORDs
			records, err = token.ReadProperty("REGIONS", -1)
			if err != nil {
				return err
			}

			numRegions := int(0)
			err = parse(&numRegions, records[1])
			if err != nil {
				return fmt.Errorf("num regions: %w", err)
			}

			list.Ranges = make([]byte, numRegions*2)
			for j := 0; j < numRegions; j++ {
				regionIndex := 0
				err = parse(&regionIndex, records[j+2])
				if err != nil {
					return fmt.Errorf("region %d: %w", j, err)
				}
				// Convert 1-based index to 0-based index for writing to bytes
				regionIndex -= 1
				list.Ranges[j*2] = byte(regionIndex & 0xFF)          // Lower byte
				list.Ranges[j*2+1] = byte((regionIndex >> 8) & 0xFF) // Upper byte
			}
		}

		e.VisTree.VisLists = append(e.VisTree.VisLists, list)
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

	return nil
}

func (e *Region) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
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
		aLightDef := wce.ByTag(e.AmbientLightTag)
		if aLightDef == nil {
			return 0, fmt.Errorf("ambient light def not found: %s", e.AmbientLightTag)
		}

		aLightRef, err := aLightDef.ToRaw(wce, rawWld)
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
		spriteDef := wce.ByTag(e.SpriteTag)
		if spriteDef == nil {
			return 0, fmt.Errorf("region sprite def not found: %s", e.SpriteTag)
		}

		spriteRef, err := spriteDef.ToRaw(wce, rawWld)
		if err != nil {
			return 0, fmt.Errorf("sprite def to raw: %w", err)
		}
		wfRegion.MeshReference = int32(spriteRef)
	}
	wfRegion.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfRegion)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *Region) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragRegion) error {
	if frag == nil {
		return fmt.Errorf("frag is not region (wrong fragcode?)")
	}

	e.VisTree = &VisTree{}
	e.Tag = rawWld.Name(frag.NameRef)
	e.RegionVertices = frag.RegionVertices
	e.Sphere = frag.Sphere
	e.ReverbVolume = frag.ReverbVolume
	e.ReverbOffset = frag.ReverbOffset
	// 0x01 is sphere, we just copy
	// 0x02 has reverb volume, we just copy
	// 0x04 has reverb offset, we just copy
	if frag.Flags&0x08 == 0x08 {
		e.RegionFog = 1
	}
	if frag.Flags&0x10 == 0x10 {
		e.Gouraud2 = 1
	}
	if frag.Flags&0x20 == 0x20 {
		e.EncodedVisibility = 1
	}
	// 0x40 unknown
	if frag.Flags&0x80 == 0x80 {
		e.VisListBytes = 1
	}

	if frag.MeshReference > 0 && frag.Flags&0x100 != 0x100 {
		fmt.Printf("region mesh ref %d but flag 0x100 not set\n", frag.MeshReference)
	}

	if frag.AmbientLightRef > 0 {
		if len(rawWld.Fragments) < int(frag.AmbientLightRef) {
			return fmt.Errorf("ambient light ref %d not found", frag.AmbientLightRef)
		}

		ambientLight, ok := rawWld.Fragments[frag.AmbientLightRef].(*rawfrag.WldFragAmbientLight)
		if !ok {
			return fmt.Errorf("ambient light ref %d not found", frag.AmbientLightRef)
		}

		e.AmbientLightTag = rawWld.Name(ambientLight.NameRef)
	}

	for _, node := range frag.VisNodes {

		visNode := &VisNode{
			Normal:       node.NormalABCD,
			VisListIndex: node.VisListIndex,
			FrontTree:    node.FrontTree,
			BackTree:     node.BackTree,
		}

		e.VisTree.VisNodes = append(e.VisTree.VisNodes, visNode)
	}

	for _, visList := range frag.VisLists {
		visListData := &VisList{}
		for _, rangeVal := range visList.Ranges {
			visListData.Ranges = append(visListData.Ranges, byte(rangeVal))
		}

		e.VisTree.VisLists = append(e.VisTree.VisLists, visListData)
	}

	if frag.MeshReference > 0 {
		if len(rawWld.Fragments) < int(frag.MeshReference) {
			return fmt.Errorf("mesh ref %d not found", frag.MeshReference)
		}

		rawMesh := rawWld.Fragments[frag.MeshReference]
		switch mesh := rawMesh.(type) {
		case *rawfrag.WldFragDmSpriteDef2:
			e.SpriteTag = rawWld.Name(mesh.NameRef)
		default:
			return fmt.Errorf("unhandled mesh reference fragment type %d (%s)", rawMesh.FragCode(), raw.FragName(rawMesh.FragCode()))
		}
	}

	return nil
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

func (e *AmbientLight) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", e.LightTag)
	fmt.Fprintf(w, "\t// LIGHTFLAGS %d\n", e.LightFlags)
	fmt.Fprintf(w, "\tREGIONLIST %d", len(e.Regions))
	for _, region := range e.Regions {
		fmt.Fprintf(w, " %d", region)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *AmbientLight) Read(token *AsciiReadToken) error {

	records, err := token.ReadProperty("LIGHT", 1)
	if err != nil {
		return err
	}

	e.LightTag = records[1]

	records, err = token.ReadProperty("REGIONLIST", -1)
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

	return nil
}

func (e *AmbientLight) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfAmbientLight := &rawfrag.WldFragAmbientLight{
		Regions: e.Regions,
	}

	if len(e.LightTag) > 0 {
		lightDef := wce.ByTag(e.LightTag)
		if lightDef == nil {
			return 0, fmt.Errorf("light def not found: %s", e.LightTag)
		}

		lightDefRef, err := lightDef.ToRaw(wce, rawWld)
		if err != nil {
			return 0, fmt.Errorf("light def to raw: %w", err)
		}
		wfAmbientLight.LightRef = int32(lightDefRef)

		wfLight := &rawfrag.WldFragLight{
			//NameRef: ,
			LightDefRef: int32(lightDefRef),
			Flags:       e.LightFlags,
		}

		rawWld.Fragments = append(rawWld.Fragments, wfLight)

		wfAmbientLight.LightRef = int32(len(rawWld.Fragments))
	}

	wfAmbientLight.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfAmbientLight)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *AmbientLight) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragAmbientLight) error {
	if frag == nil {
		return fmt.Errorf("frag is not ambient light (wrong fragcode?)")
	}

	lightTag := ""
	lightFlags := uint32(0)
	if frag.LightRef > 0 {
		if len(rawWld.Fragments) < int(frag.LightRef) {
			return fmt.Errorf("lightdef ref %d out of bounds", frag.LightRef)
		}

		light, ok := rawWld.Fragments[frag.LightRef].(*rawfrag.WldFragLight)
		if !ok {
			return fmt.Errorf("lightdef ref %d not found", frag.LightRef)
		}

		lightFlags = light.Flags

		lightDef, ok := rawWld.Fragments[light.LightDefRef].(*rawfrag.WldFragLightDef)
		if !ok {
			return fmt.Errorf("lightdef ref %d not found", light.LightDefRef)
		}

		lightTag = rawWld.Name(lightDef.NameRef)
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.LightTag = lightTag
	e.LightFlags = lightFlags
	e.Regions = frag.Regions

	return nil

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

func (e *Zone) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tREGIONLIST %d", len(e.Regions))
	for _, region := range e.Regions {
		fmt.Fprintf(w, " %d", region)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", e.UserData)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *Zone) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("REGIONLIST", -1)
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

	records, err = token.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}

	e.UserData = records[1]

	return nil
}

func (e *Zone) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}

	wfZone := &rawfrag.WldFragZone{
		NameRef:  rawWld.NameAdd(e.Tag),
		Flags:    0,
		Regions:  e.Regions,
		UserData: e.UserData,
	}

	rawWld.Fragments = append(rawWld.Fragments, wfZone)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *Zone) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragZone) error {
	if frag == nil {
		return fmt.Errorf("frag is not zone (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.Regions = frag.Regions
	e.UserData = frag.UserData
	return nil
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

func (e *RGBTrackDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tDATA1 %d\n", e.Data1)
	fmt.Fprintf(w, "\tDATA2 %d\n", e.Data2)
	fmt.Fprintf(w, "\tSLEEP %d\n", e.Sleep)
	fmt.Fprintf(w, "\tDATA4 %d\n", e.Data4)
	fmt.Fprintf(w, "\tRGBDEFORMATIONFRAME\n")
	fmt.Fprintf(w, "\t\tNUMRGBAS %d\n", len(e.RGBAs))
	for _, rgba := range e.RGBAs {
		fmt.Fprintf(w, "\t\tRGBA %d %d %d %d\n", rgba[0], rgba[1], rgba[2], rgba[3])
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *RGBTrackDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("DATA1", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Data1, records[1])
	if err != nil {
		return fmt.Errorf("data1: %w", err)
	}

	records, err = token.ReadProperty("DATA2", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Data2, records[1])
	if err != nil {
		return fmt.Errorf("data2: %w", err)
	}

	records, err = token.ReadProperty("SLEEP", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = token.ReadProperty("DATA4", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Data4, records[1])
	if err != nil {
		return fmt.Errorf("data4: %w", err)
	}

	_, err = token.ReadProperty("RGBDEFORMATIONFRAME", 0)
	if err != nil {
		return err
	}

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

		err = parse(&rgba, records[1:]...)
		if err != nil {
			return fmt.Errorf("rgba: %w", err)
		}
		e.RGBAs = append(e.RGBAs, rgba)
	}

	return nil
}

func (e *RGBTrackDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfRGBTrack := &rawfrag.WldFragDmRGBTrackDef{
		RGBAs: e.RGBAs,
	}

	wfRGBTrack.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfRGBTrack)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *RGBTrackDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragDmRGBTrackDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not rgb track def (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.Data1 = frag.Data1
	e.Data2 = frag.Data2
	e.Sleep = frag.Sleep
	e.Data4 = frag.Data4
	e.RGBAs = frag.RGBAs
	return nil
}

type ParticleCloudDef struct {
	fragID                int16
	Tag                   string
	BlitSpriteDefTag      string
	SettingOne            uint32
	SettingTwo            uint32
	Movement              string
	HighOpacity           int
	FollowItem            int
	SimultaneousParticles uint32
	UnkSix                uint32
	UnkSeven              uint32
	UnkEight              uint32
	UnkNine               uint32
	UnkTen                uint32
	SpawnRadius           float32
	SpawnAngle            float32
	SpawnLifespan         uint32
	SpawnVelocity         float32
	SpawnNormal           [3]float32
	SpawnRate             uint32
	SpawnScale            float32
	Color                 [4]uint8
	HexEightyFlag         int
	HexOneHundredFlag     int
	HexFourHundredFlag    int
	HexFourThousandFlag   int
	HexEightThousandFlag  int
	HexTenThousandFlag    int
	HexTwentyThousandFlag int
}

func (e *ParticleCloudDef) Definition() string {
	return "PARTICLECLOUDDEF"
}

func (e *ParticleCloudDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}

	if e.BlitSpriteDefTag != "" {
		sDef := token.wce.ByTag(e.BlitSpriteDefTag)
		if sDef == nil {
			return fmt.Errorf("blit sprite def not found: %s", e.BlitSpriteDefTag)
		}
		err = sDef.Write(token)
		if err != nil {
			return fmt.Errorf("blit sprite def to raw: %w", err)
		}
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tBLITTAG \"%s\"\n", e.BlitSpriteDefTag)
	fmt.Fprintf(w, "\tSETTINGONE %d\n", e.SettingOne)
	fmt.Fprintf(w, "\tSETTINGTWO %d\n", e.SettingTwo)
	fmt.Fprintf(w, "\tMOVEMENT \"%s\" // SPHERE, PLANE, STREAM, NONE\n", e.Movement)
	fmt.Fprintf(w, "\tHIGHOPACITY %d\n", e.HighOpacity)
	fmt.Fprintf(w, "\tFOLLOWITEM %d\n", e.FollowItem)
	fmt.Fprintf(w, "\tSIMULTANEOUSPARTICLES %d\n", e.SimultaneousParticles)
	fmt.Fprintf(w, "\tUNKSIX %d\n", e.UnkSix)
	fmt.Fprintf(w, "\tUNKSEVEN %d\n", e.UnkSeven)
	fmt.Fprintf(w, "\tUNKEIGHT %d\n", e.UnkEight)
	fmt.Fprintf(w, "\tUNKNINE %d\n", e.UnkNine)
	fmt.Fprintf(w, "\tUNKTEN %d\n", e.UnkTen)
	fmt.Fprintf(w, "\tSPAWN\n")
	fmt.Fprintf(w, "\t\tRADIUS %0.8e\n", e.SpawnRadius)
	fmt.Fprintf(w, "\t\tANGLE %0.8e\n", e.SpawnAngle)
	fmt.Fprintf(w, "\t\tLIFESPAN %d\n", e.SpawnLifespan)
	fmt.Fprintf(w, "\t\tVELOCITY %0.8e\n", e.SpawnVelocity)
	fmt.Fprintf(w, "\t\tNORMALXYZ %0.8e %0.8e %0.8e\n", e.SpawnNormal[0], e.SpawnNormal[1], e.SpawnNormal[2])
	fmt.Fprintf(w, "\t\tRATE %d\n", e.SpawnRate)
	fmt.Fprintf(w, "\t\tSCALE %0.8e\n", e.SpawnScale)
	fmt.Fprintf(w, "\tCOLOR %d %d %d %d\n", e.Color[0], e.Color[1], e.Color[2], e.Color[3])
	fmt.Fprintf(w, "\tHEXEIGHTYFLAG %d\n", e.HexEightyFlag)
	fmt.Fprintf(w, "\tHEXONEHUNDREDFLAG %d\n", e.HexOneHundredFlag)
	fmt.Fprintf(w, "\tHEXFOURHUNDREDFLAG %d\n", e.HexFourHundredFlag)
	fmt.Fprintf(w, "\tHEXFOURTHOUSANDFLAG %d\n", e.HexFourThousandFlag)
	fmt.Fprintf(w, "\tHEXEIGHTTHOUSANDFLAG %d\n", e.HexEightThousandFlag)
	fmt.Fprintf(w, "\tHEXTENTHOUSANDFLAG %d\n", e.HexTenThousandFlag)
	fmt.Fprintf(w, "\tHEXTWENTYTHOUSANDFLAG %d\n", e.HexTwentyThousandFlag)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *ParticleCloudDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("BLITTAG", 1)
	if err != nil {
		return err
	}
	e.BlitSpriteDefTag = records[1]

	records, err = token.ReadProperty("SETTINGONE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SettingOne, records[1])
	if err != nil {
		return fmt.Errorf("setting one: %w", err)
	}

	records, err = token.ReadProperty("SETTINGTWO", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SettingTwo, records[1])
	if err != nil {
		return fmt.Errorf("setting two: %w", err)
	}

	records, err = token.ReadProperty("MOVEMENT", 1)
	if err != nil {
		return err
	}
	e.Movement = records[1]

	records, err = token.ReadProperty("HIGHOPACITY", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HighOpacity, records[1])
	if err != nil {
		return fmt.Errorf("high opacity: %w", err)
	}

	records, err = token.ReadProperty("FOLLOWITEM", 1)
	if err != nil {
		return err
	}
	err = parse(&e.FollowItem, records[1])
	if err != nil {
		return fmt.Errorf("follow item: %w", err)
	}

	records, err = token.ReadProperty("SIMULTANEOUSPARTICLES", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SimultaneousParticles, records[1])
	if err != nil {
		return fmt.Errorf("simultaneous particles: %w", err)
	}

	records, err = token.ReadProperty("UNKSIX", 1)
	if err != nil {
		return err
	}
	err = parse(&e.UnkSix, records[1])
	if err != nil {
		return fmt.Errorf("unk six: %w", err)
	}

	records, err = token.ReadProperty("UNKSEVEN", 1)
	if err != nil {
		return err
	}
	err = parse(&e.UnkSeven, records[1])
	if err != nil {
		return fmt.Errorf("unk seven: %w", err)
	}

	records, err = token.ReadProperty("UNKEIGHT", 1)
	if err != nil {
		return err
	}
	err = parse(&e.UnkEight, records[1])
	if err != nil {
		return fmt.Errorf("unk eight: %w", err)
	}

	records, err = token.ReadProperty("UNKNINE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.UnkNine, records[1])
	if err != nil {
		return fmt.Errorf("unk nine: %w", err)
	}

	records, err = token.ReadProperty("UNKTEN", 1)
	if err != nil {
		return err
	}
	err = parse(&e.UnkTen, records[1])
	if err != nil {
		return fmt.Errorf("unk ten: %w", err)
	}

	_, err = token.ReadProperty("SPAWN", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("RADIUS", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpawnRadius, records[1])
	if err != nil {
		return fmt.Errorf("spawn radius: %w", err)
	}

	records, err = token.ReadProperty("ANGLE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpawnAngle, records[1])
	if err != nil {
		return fmt.Errorf("spawn angle: %w", err)
	}

	records, err = token.ReadProperty("LIFESPAN", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpawnLifespan, records[1])
	if err != nil {
		return fmt.Errorf("spawn lifespan: %w", err)
	}

	records, err = token.ReadProperty("VELOCITY", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpawnVelocity, records[1])
	if err != nil {
		return fmt.Errorf("spawn velocity: %w", err)
	}

	records, err = token.ReadProperty("NORMALXYZ", 3)
	if err != nil {
		return err
	}
	err = parse(&e.SpawnNormal, records[1:]...)
	if err != nil {
		return fmt.Errorf("spawn normal: %w", err)
	}

	records, err = token.ReadProperty("RATE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpawnRate, records[1])
	if err != nil {
		return fmt.Errorf("spawn rate: %w", err)
	}

	records, err = token.ReadProperty("SCALE", 1)
	if err != nil {
		return err
	}

	err = parse(&e.SpawnScale, records[1])
	if err != nil {
		return fmt.Errorf("spawn scale: %w", err)
	}

	records, err = token.ReadProperty("COLOR", 4)
	if err != nil {
		return err
	}
	err = parse(&e.Color, records[1:]...)
	if err != nil {
		return fmt.Errorf("color: %w", err)
	}

	records, err = token.ReadProperty("HEXEIGHTYFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexEightyFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex eighty flag: %w", err)
	}

	records, err = token.ReadProperty("HEXONEHUNDREDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexOneHundredFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex one hundred flag: %w", err)
	}

	records, err = token.ReadProperty("HEXFOURHUNDREDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexFourHundredFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex four hundred flag: %w", err)
	}

	records, err = token.ReadProperty("HEXFOURTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexFourThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex four thousand flag: %w", err)
	}

	records, err = token.ReadProperty("HEXEIGHTTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexEightThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex eight thousand flag: %w", err)
	}

	records, err = token.ReadProperty("HEXTENTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTenThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex ten thousand flag: %w", err)
	}

	records, err = token.ReadProperty("HEXTWENTYTHOUSANDFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTwentyThousandFlag, records[1])
	if err != nil {
		return fmt.Errorf("hex twenty thousand flag: %w", err)
	}

	return nil
}

func (e *ParticleCloudDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfParticleCloud := &rawfrag.WldFragParticleCloudDef{
		SettingOne:            e.SettingOne,
		SettingTwo:            e.SettingTwo,
		SimultaneousParticles: e.SimultaneousParticles,
		Unk6:                  e.UnkSix,
		Unk7:                  e.UnkSeven,
		Unk8:                  e.UnkEight,
		Unk9:                  e.UnkNine,
		Unk10:                 e.UnkTen,
		SpawnRadius:           e.SpawnRadius,
		SpawnAngle:            e.SpawnAngle,
		SpawnLifespan:         e.SpawnLifespan,
		SpawnVelocity:         e.SpawnVelocity,
		SpawnNormal:           e.SpawnNormal,
		SpawnRate:             e.SpawnRate,
		SpawnScale:            e.SpawnScale,
		Color:                 e.Color,
	}

	if e.HighOpacity != 0 {
		wfParticleCloud.Flags |= 0x01
	}
	if e.FollowItem != 0 {
		wfParticleCloud.Flags |= 0x02
	}

	if e.HexEightyFlag != 0 {
		wfParticleCloud.Flags |= 0x80
	}

	if e.HexOneHundredFlag != 0 {
		wfParticleCloud.Flags |= 0x100
	}

	if e.HexFourHundredFlag != 0 {
		wfParticleCloud.Flags |= 0x400
	}

	if e.HexFourThousandFlag != 0 {
		wfParticleCloud.Flags |= 0x4000
	}

	if e.HexEightThousandFlag != 0 {
		wfParticleCloud.Flags |= 0x8000
	}

	if e.HexTenThousandFlag != 0 {
		wfParticleCloud.Flags |= 0x10000
	}

	if e.HexTwentyThousandFlag != 0 {
		wfParticleCloud.Flags |= 0x20000
	}

	switch e.Movement {
	case "SPHERE":
		wfParticleCloud.ParticleMovement = 1
	case "PLANE":
		wfParticleCloud.ParticleMovement = 2
	case "STREAM":
		wfParticleCloud.ParticleMovement = 3
	case "NONE":
		wfParticleCloud.ParticleMovement = 4
	default:
		return 0, fmt.Errorf("unknown movement type %s", e.Movement)
	}

	if e.BlitSpriteDefTag == "" {
		return 0, fmt.Errorf("blit sprite def tag not set")
	}

	blit := wce.ByTag(e.BlitSpriteDefTag)
	if blit == nil {
		return 0, fmt.Errorf("blit sprite def not found: %s", e.BlitSpriteDefTag)
	}

	blitFragID, err := blit.ToRaw(wce, rawWld)
	if err != nil {
		return 0, fmt.Errorf("blit sprite def to raw: %w", err)
	}

	blitSpriteDefRef := uint32(blitFragID)

	wfParticleCloud.BlitSpriteRef = blitSpriteDefRef

	wfParticleCloud.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfParticleCloud)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *ParticleCloudDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragParticleCloudDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not particle cloud def (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	if len(rawWld.Fragments) < int(frag.BlitSpriteRef) {
		return fmt.Errorf("blit sprite def ref %d out of bounds", frag.BlitSpriteRef)
	}
	bSprite, ok := rawWld.Fragments[frag.BlitSpriteRef].(*rawfrag.WldFragBlitSpriteDef)
	if !ok {
		return fmt.Errorf("blit sprite def ref %d not found", frag.BlitSpriteRef)
	}

	e.BlitSpriteDefTag = rawWld.Name(bSprite.NameRef)

	if len(rawWld.Fragments) < int(bSprite.SpriteInstanceRef) {
		return fmt.Errorf("sprite instance ref %d out of bounds", bSprite.SpriteInstanceRef)
	}
	sSprite, ok := rawWld.Fragments[bSprite.SpriteInstanceRef].(*rawfrag.WldFragSimpleSprite)
	if !ok {
		return fmt.Errorf("sprite instance ref %d not found", bSprite.SpriteInstanceRef)
	}
	if len(rawWld.Fragments) < int(sSprite.SpriteRef) {
		return fmt.Errorf("sprite def ref %d out of bounds", sSprite.SpriteRef)
	}
	e.SettingOne = frag.SettingOne
	e.SettingTwo = frag.SettingTwo
	switch frag.ParticleMovement {
	case 1:
		e.Movement = "SPHERE"
	case 2:
		e.Movement = "PLANE"
	case 3:
		e.Movement = "STREAM"
	case 4:
		e.Movement = "NONE"
	default:
		return fmt.Errorf("unknown movement type %d", frag.ParticleMovement)
	}
	if frag.Flags&0x01 != 0 {
		e.HighOpacity = 1
	}
	if frag.Flags&0x02 != 0 {
		e.FollowItem = 1
	}
	if frag.Flags&0x80 != 0 {
		e.HexEightyFlag = 1
	}
	if frag.Flags&0x100 != 0 {
		e.HexOneHundredFlag = 1
	}
	if frag.Flags&0x400 != 0 {
		e.HexFourHundredFlag = 1
	}
	if frag.Flags&0x4000 != 0 {
		e.HexFourThousandFlag = 1
	}
	if frag.Flags&0x8000 != 0 {
		e.HexEightThousandFlag = 1
	}
	if frag.Flags&0x10000 != 0 {
		e.HexTenThousandFlag = 1
	}
	if frag.Flags&0x20000 != 0 {
		e.HexTwentyThousandFlag = 1
	}
	e.SimultaneousParticles = frag.SimultaneousParticles
	e.UnkSix = frag.Unk6
	e.UnkSeven = frag.Unk7
	e.UnkEight = frag.Unk8
	e.UnkNine = frag.Unk9
	e.UnkTen = frag.Unk10
	e.SpawnRadius = frag.SpawnRadius
	e.SpawnAngle = frag.SpawnAngle
	e.SpawnLifespan = frag.SpawnLifespan
	e.SpawnVelocity = frag.SpawnVelocity
	e.SpawnNormal = frag.SpawnNormal
	e.SpawnRate = frag.SpawnRate
	e.SpawnScale = frag.SpawnScale
	e.Color = frag.Color
	return nil
}

type Sprite2DDef struct {
	fragID          int16
	Tag             string
	Scale           [2]float32
	SphereListTag   string
	DepthScale      NullFloat32
	CenterOffset    NullFloat32Slice3
	BoundingRadius  NullFloat32
	CurrentFrameRef NullInt32
	Sleep           NullUint32
	Pitches         []*Pitch
	RenderMethod    string
	Pen             NullUint32
	Brightness      NullFloat32
	ScaledAmbient   NullFloat32
	SpriteTag       NullString
	UvOrigin        NullFloat32Slice3
	UAxis           NullFloat32Slice3
	VAxis           NullFloat32Slice3
	Uvs             [][2]float32
	TwoSided        int
	HexTenFlag      uint16
}

type Pitch struct {
	PitchCap        int32
	TopOrBottomView uint32 // Only 0 or 1
	Headings        []*Heading
}

type Heading struct {
	HeadingCap     int32
	Sprite2DFrames []*Sprite2DFrame
}

type Sprite2DFrame struct {
	TextureFile string
	TextureTag  string
}

func (e *Sprite2DDef) Definition() string {
	return "SPRITE2DDEF"
}

func (e *Sprite2DDef) Write(token *AsciiWriteToken) error {
	w, err := token.Writer()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tSCALE %0.8e %0.8e\n", e.Scale[0], e.Scale[1])
	fmt.Fprintf(w, "\tSPHERELISTTAG \"%s\"\n", e.SphereListTag)
	fmt.Fprintf(w, "\tDEPTHSCALE? %s\n", wcVal(e.DepthScale))
	fmt.Fprintf(w, "\tCENTEROFFSET? %s\n", wcVal(e.CenterOffset))
	fmt.Fprintf(w, "\tBOUNDINGRADIUS? %s\n", wcVal(e.BoundingRadius))
	fmt.Fprintf(w, "\tCURRENTFRAMEREF? %s\n", wcVal(e.CurrentFrameRef))
	fmt.Fprintf(w, "\tSLEEP? %s\n", wcVal(e.Sleep))
	fmt.Fprintf(w, "\tNUMPITCHES %d\n", len(e.Pitches))
	for i, pitch := range e.Pitches {
		fmt.Fprintf(w, "\t\tPITCH // %d\n", i)
		fmt.Fprintf(w, "\t\t\tPITCHCAP %d\n", pitch.PitchCap)
		fmt.Fprintf(w, "\t\t\tTOPORBOTTOMVIEW %d\n", pitch.TopOrBottomView)
		fmt.Fprintf(w, "\t\t\tNUMHEADINGS %d\n", len(pitch.Headings))
		for i, heading := range pitch.Headings {
			fmt.Fprintf(w, "\t\t\t\tHEADING // %d\n", i)
			fmt.Fprintf(w, "\t\t\t\t\tHEADINGCAP %d\n", heading.HeadingCap)
			fmt.Fprintf(w, "\t\t\t\t\tNUMFRAMES %d\n", len(heading.Sprite2DFrames))
			for _, frame := range heading.Sprite2DFrames {
				fmt.Fprintf(w, "\t\t\t\t\t\tFRAME \"%s\" \"%s\"\n", frame.TextureFile, frame.TextureTag)
			}
		}
	}
	fmt.Fprintf(w, "\t\tRENDERMETHOD \"%s\"\n", e.RenderMethod)
	fmt.Fprintf(w, "\t\tRENDERINFO\n")
	fmt.Fprintf(w, "\t\t\tPEN? %s\n", wcVal(e.Pen))
	fmt.Fprintf(w, "\t\t\tBRIGHTNESS? %s\n", wcVal(e.Brightness))
	fmt.Fprintf(w, "\t\t\tSCALEDAMBIENT? %s\n", wcVal(e.ScaledAmbient))
	fmt.Fprintf(w, "\t\t\tSPRITE? \"%s\"\n", wcVal(e.SpriteTag))
	fmt.Fprintf(w, "\t\t\tUVORIGIN? %s\n", wcVal(e.UvOrigin))
	fmt.Fprintf(w, "\t\t\tUAXIS? %s\n", wcVal(e.UAxis))
	fmt.Fprintf(w, "\t\t\tVAXIS? %s\n", wcVal(e.VAxis))
	fmt.Fprintf(w, "\t\t\tUVCOUNT %d\n", len(e.Uvs))
	for _, uv := range e.Uvs {
		fmt.Fprintf(w, "\t\t\tUV %s\n", wcVal(uv))
	}
	fmt.Fprintf(w, "\t\t\tTWOSIDED %d\n", e.TwoSided)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tHEXTENFLAG %d\n", e.HexTenFlag)
	fmt.Fprintf(w, "\n")
	return nil
}

func (e *Sprite2DDef) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("SCALE", 2)
	if err != nil {
		return err
	}
	err = parse(&e.Scale, records[1:]...)
	if err != nil {
		return fmt.Errorf("scale: %w", err)
	}

	records, err = token.ReadProperty("SPHERELISTTAG", 1)
	if err != nil {
		return err
	}
	e.SphereListTag = records[1]

	records, err = token.ReadProperty("DEPTHSCALE?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.DepthScale, records[1])
	if err != nil {
		return fmt.Errorf("depth scale: %w", err)
	}

	records, err = token.ReadProperty("CENTEROFFSET?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.CenterOffset, records[1:]...)
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = token.ReadProperty("BOUNDINGRADIUS?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.BoundingRadius, records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = token.ReadProperty("CURRENTFRAMEREF?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.CurrentFrameRef, records[1])
	if err != nil {
		return fmt.Errorf("current frame ref: %w", err)
	}

	records, err = token.ReadProperty("SLEEP?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = token.ReadProperty("NUMPITCHES", 1)
	if err != nil {
		return err
	}
	numPitches := 0
	err = parse(&numPitches, records[1])
	if err != nil {
		return fmt.Errorf("num pitches: %w", err)
	}

	e.Pitches = []*Pitch{}
	for i := 0; i < numPitches; i++ {
		pitch := &Pitch{}
		_, err = token.ReadProperty("PITCH", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("PITCHCAP", 1)
		if err != nil {
			return err
		}
		err = parse(&pitch.PitchCap, records[1])
		if err != nil {
			return fmt.Errorf("pitch cap: %w", err)
		}

		records, err = token.ReadProperty("TOPORBOTTOMVIEW", 1)
		if err != nil {
			return err
		}
		err = parse(&pitch.TopOrBottomView, records[1])
		if err != nil {
			return fmt.Errorf("top or bottom view: %w", err)
		}

		records, err = token.ReadProperty("NUMHEADINGS", 1)
		if err != nil {
			return err
		}
		numHeadings := 0
		err = parse(&numHeadings, records[1])
		if err != nil {
			return fmt.Errorf("num headings: %w", err)
		}

		pitch.Headings = []*Heading{}
		for j := 0; j < numHeadings; j++ {
			heading := &Heading{}
			_, err = token.ReadProperty("HEADING", 0)
			if err != nil {
				return err
			}

			records, err = token.ReadProperty("HEADINGCAP", 1)
			if err != nil {
				return err
			}
			err = parse(&heading.HeadingCap, records[1])
			if err != nil {
				return fmt.Errorf("heading cap: %w", err)
			}

			records, err = token.ReadProperty("NUMFRAMES", 1)
			if err != nil {
				return err
			}
			numFrames := 0
			err = parse(&numFrames, records[1])
			if err != nil {
				return fmt.Errorf("num frames: %w", err)
			}

			heading.Sprite2DFrames = []*Sprite2DFrame{}
			for k := 0; k < numFrames; k++ {
				frame := &Sprite2DFrame{}
				records, err = token.ReadProperty("FRAME", 2)
				if err != nil {
					return fmt.Errorf("FRAME: %w", err)
				}
				frame.TextureFile = records[1]
				frame.TextureTag = records[2]
				heading.Sprite2DFrames = append(heading.Sprite2DFrames, frame)
			}
			pitch.Headings = append(pitch.Headings, heading)
		}
		e.Pitches = append(e.Pitches, pitch)
	}

	records, err = token.ReadProperty("RENDERMETHOD", 1)
	if err != nil {
		return err
	}

	e.RenderMethod = records[1]

	_, err = token.ReadProperty("RENDERINFO", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("PEN?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Pen, records[1])
	if err != nil {
		return fmt.Errorf("render pen: %w", err)
	}

	records, err = token.ReadProperty("BRIGHTNESS?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Brightness, records[1])
	if err != nil {
		return fmt.Errorf("render brightness: %w", err)
	}

	records, err = token.ReadProperty("SCALEDAMBIENT?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.ScaledAmbient, records[1])
	if err != nil {
		return fmt.Errorf("render scaled ambient: %w", err)
	}

	records, err = token.ReadProperty("SPRITE?", 1)
	if err != nil {
		return err
	}
	err = parse(&e.SpriteTag, records[1])
	if err != nil {
		return fmt.Errorf("render sprite: %w", err)
	}

	records, err = token.ReadProperty("UVORIGIN?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.UvOrigin, records[1:]...)
	if err != nil {
		return fmt.Errorf("render uv origin: %w", err)
	}

	records, err = token.ReadProperty("UAXIS?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.UAxis, records[1:]...)
	if err != nil {
		return fmt.Errorf("render u axis: %w", err)
	}

	records, err = token.ReadProperty("VAXIS?", 3)
	if err != nil {
		return err
	}
	err = parse(&e.VAxis, records[1:]...)
	if err != nil {
		return fmt.Errorf("render v axis: %w", err)
	}

	records, err = token.ReadProperty("UVCOUNT", 1)
	if err != nil {
		return err
	}
	numUVs := int(0)
	err = parse(&numUVs, records[1])
	if err != nil {
		return fmt.Errorf("num uvs: %w", err)
	}

	for j := 0; j < numUVs; j++ {
		records, err = token.ReadProperty("UV", 2)
		if err != nil {
			return err
		}
		uv := [2]float32{}
		err = parse(&uv, records[1:]...)
		if err != nil {
			return fmt.Errorf("uv %d: %w", j, err)
		}
		e.Uvs = append(e.Uvs, uv)
	}

	records, err = token.ReadProperty("TWOSIDED", 1)
	if err != nil {
		return err
	}
	err = parse(&e.TwoSided, records[1])
	if err != nil {
		return fmt.Errorf("two sided: %w", err)
	}

	records, err = token.ReadProperty("HEXTENFLAG", 1)
	if err != nil {
		return err
	}
	err = parse(&e.HexTenFlag, records[1])
	if err != nil {
		return fmt.Errorf("hextenflag: %w", err)
	}

	return nil
}

func (e *Sprite2DDef) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	if e.fragID != 0 {
		return e.fragID, nil
	}
	wfSprite2D := &rawfrag.WldFragSprite2DDef{
		Scale:        e.Scale,
		RenderMethod: model.RenderMethodInt(e.RenderMethod),
	}

	if e.DepthScale.Valid {
		wfSprite2D.Flags |= 0x80
		wfSprite2D.DepthScale = e.DepthScale.Float32
	}

	if e.CenterOffset.Valid {
		wfSprite2D.Flags |= 0x01
		wfSprite2D.CenterOffset = e.CenterOffset.Float32Slice3
	}

	if e.BoundingRadius.Valid {
		wfSprite2D.Flags |= 0x02
		wfSprite2D.BoundingRadius = e.BoundingRadius.Float32
	}

	if e.CurrentFrameRef.Valid {
		wfSprite2D.Flags |= 0x04
		wfSprite2D.CurrentFrameRef = e.CurrentFrameRef.Int32
	}

	if e.Sleep.Valid {
		wfSprite2D.Flags |= 0x08
		wfSprite2D.Sleep = e.Sleep.Uint32
	}

	wfSprite2D.Pitches = []*rawfrag.WldFragSprite2DPitch{}
	for _, pitch := range e.Pitches {
		rawPitch := &rawfrag.WldFragSprite2DPitch{
			PitchCap:        pitch.PitchCap,
			TopOrBottomView: pitch.TopOrBottomView,
		}

		rawPitch.Headings = []*rawfrag.WldFragSprite2DHeading{}
		for _, heading := range pitch.Headings {
			rawHeading := &rawfrag.WldFragSprite2DHeading{
				HeadingCap: heading.HeadingCap,
			}

			if len(heading.Sprite2DFrames) > 0 {
				wfBMInfo := &rawfrag.WldFragBMInfo{}
				for _, frame := range heading.Sprite2DFrames {
					frameRef := rawWld.NameAdd(frame.TextureTag)
					if wfBMInfo.NameRef != 0 && frameRef != wfBMInfo.NameRef {
						// Add the current BMInfo to fragments and reset it
						rawWld.Fragments = append(rawWld.Fragments, wfBMInfo)
						rawHeading.FrameRefs = append(rawHeading.FrameRefs, int32(len(rawWld.Fragments)))
						wfBMInfo = &rawfrag.WldFragBMInfo{}
					}
					wfBMInfo.NameRef = frameRef
					wfBMInfo.TextureNames = append(wfBMInfo.TextureNames, frame.TextureFile+"\x00")
				}

				rawWld.Fragments = append(rawWld.Fragments, wfBMInfo)
				rawHeading.FrameRefs = append(rawHeading.FrameRefs, int32(len(rawWld.Fragments)))
			}
			rawPitch.Headings = append(rawPitch.Headings, rawHeading)
		}
		wfSprite2D.Pitches = append(wfSprite2D.Pitches, rawPitch)
	}

	if e.Pen.Valid {
		wfSprite2D.RenderFlags |= 0x01
		wfSprite2D.RenderPen = e.Pen.Uint32
	}

	if e.Brightness.Valid {
		wfSprite2D.RenderFlags |= 0x02
		wfSprite2D.RenderBrightness = e.Brightness.Float32
	}

	if e.ScaledAmbient.Valid {
		wfSprite2D.RenderFlags |= 0x04
		wfSprite2D.RenderScaledAmbient = e.ScaledAmbient.Float32
	}

	if e.SpriteTag.Valid {
		wfSprite2D.RenderFlags |= 0x08
		wfSprite2D.RenderSimpleSpriteReference = uint32(rawWld.NameAdd(e.SpriteTag.String))
	}

	if e.UvOrigin.Valid {
		wfSprite2D.RenderFlags |= 0x10
		wfSprite2D.RenderUVInfoOrigin = e.UvOrigin.Float32Slice3
		wfSprite2D.RenderUVInfoUAxis = e.UAxis.Float32Slice3
		wfSprite2D.RenderUVInfoVAxis = e.VAxis.Float32Slice3
	}

	if len(e.Uvs) > 0 {
		wfSprite2D.RenderFlags |= 0x20
		wfSprite2D.Uvs = e.Uvs
	}

	if e.SphereListTag != "" {
		sphereList := wce.ByTag(e.SphereListTag)
		if sphereList == nil {
			return 0, fmt.Errorf("sphere list tag not found: %s", e.SphereListTag)
		}

		sphereListRef, err := sphereList.ToRaw(wce, rawWld)
		if err != nil {
			return 0, fmt.Errorf("sphere list to raw: %w", err)
		}
		wfSprite2D.SphereListRef = uint32(sphereListRef)
	}

	if e.HexTenFlag != 0 {
		wfSprite2D.Flags |= 0x10
	}

	wfSprite2D.NameRef = rawWld.NameAdd(e.Tag)

	rawWld.Fragments = append(rawWld.Fragments, wfSprite2D)
	e.fragID = int16(len(rawWld.Fragments))
	return int16(len(rawWld.Fragments)), nil
}

func (e *Sprite2DDef) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragSprite2DDef) error {
	if frag == nil {
		return fmt.Errorf("frag is not sprite 2d def (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)

	if frag.SphereListRef > 0 {
		if len(rawWld.Fragments) < int(frag.SphereListRef) {
			return fmt.Errorf("sphere list ref %d out of bounds", frag.SphereListRef)
		}
		sphereListRef := rawWld.Fragments[frag.SphereListRef]

		sphereList, ok := sphereListRef.(*rawfrag.WldFragSphereList)
		if !ok {
			return fmt.Errorf("sphere list ref %d not found", frag.SphereListRef)
		}

		e.SphereListTag = rawWld.Name(sphereList.NameRef)
	}
	e.Scale = frag.Scale

	if frag.Flags&0x80 == 0x80 {
		e.DepthScale.Valid = true
		e.DepthScale.Float32 = frag.DepthScale
	}

	if frag.Flags&0x01 == 0x01 {
		e.CenterOffset.Valid = true
		e.CenterOffset.Float32Slice3 = frag.CenterOffset
	}

	if frag.Flags&0x02 == 0x02 {
		e.BoundingRadius.Valid = true
		e.BoundingRadius.Float32 = frag.BoundingRadius
	}

	if frag.Flags&0x04 == 0x04 {
		e.CurrentFrameRef.Valid = true
		e.CurrentFrameRef.Int32 = frag.CurrentFrameRef
	}

	if frag.Flags&0x08 == 0x08 {
		e.Sleep.Valid = true
		e.Sleep.Uint32 = frag.Sleep
	}

	e.Pitches = []*Pitch{}
	for _, rawPitch := range frag.Pitches {
		pitch := &Pitch{
			PitchCap:        rawPitch.PitchCap,
			TopOrBottomView: rawPitch.TopOrBottomView,
			Headings:        []*Heading{},
		}

		for _, rawHeading := range rawPitch.Headings {
			heading := &Heading{
				HeadingCap: rawHeading.HeadingCap,
			}
			for _, frameRef := range rawHeading.FrameRefs {
				if frameRef == 0 {
					return nil
				}
				if len(rawWld.Fragments) <= int(frameRef) {
					return fmt.Errorf("frame reference %d out of range", frameRef)
				}
				frame := rawWld.Fragments[frameRef]
				bmInfo, ok := frame.(*rawfrag.WldFragBMInfo)
				if !ok {
					return fmt.Errorf("invalid frame ref %d", frameRef)
				}
				for _, textureName := range bmInfo.TextureNames {
					heading.Sprite2DFrames = append(heading.Sprite2DFrames, &Sprite2DFrame{
						TextureTag:  rawWld.Name(bmInfo.NameRef),
						TextureFile: textureName,
					})
				}
			}
			pitch.Headings = append(pitch.Headings, heading)
		}
		e.Pitches = append(e.Pitches, pitch)
	}

	e.RenderMethod = model.RenderMethodStr(frag.RenderMethod)
	if frag.RenderFlags&0x01 == 0x01 {
		e.Pen.Valid = true
		e.Pen.Uint32 = frag.RenderPen
	}

	if frag.RenderFlags&0x02 == 0x02 {
		e.Brightness.Valid = true
		e.Brightness.Float32 = frag.RenderBrightness
	}

	if frag.RenderFlags&0x04 == 0x04 {
		e.ScaledAmbient.Valid = true
		e.ScaledAmbient.Float32 = frag.RenderScaledAmbient
	}

	if frag.RenderFlags&0x08 == 0x08 {
		e.SpriteTag.Valid = true
		if len(rawWld.Fragments) < int(frag.RenderSimpleSpriteReference) {
			return fmt.Errorf("sprite2d's simple sprite ref %d not found", frag.RenderSimpleSpriteReference)
		}
		spriteDef := rawWld.Fragments[frag.RenderSimpleSpriteReference]
		switch simpleSprite := spriteDef.(type) {
		case *rawfrag.WldFragSimpleSpriteDef:
			e.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
		case *rawfrag.WldFragDMSpriteDef:
			e.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
		case *rawfrag.WldFragHierarchicalSpriteDef:
			e.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
		case *rawfrag.WldFragSprite2D:
			e.SpriteTag.String = rawWld.Name(simpleSprite.NameRef)
		default:
			return fmt.Errorf("unhandled render sprite reference fragment type %d", spriteDef.FragCode())
		}
	}

	if frag.RenderFlags&0x10 == 0x10 {
		// has uvinfo
		e.UvOrigin.Valid = true
		e.UAxis.Valid = true
		e.VAxis.Valid = true
		e.UvOrigin.Float32Slice3 = frag.RenderUVInfoOrigin
		e.UAxis.Float32Slice3 = frag.RenderUVInfoUAxis
		e.VAxis.Float32Slice3 = frag.RenderUVInfoVAxis
	}

	if frag.RenderFlags&0x20 == 0x20 {
		e.Uvs = frag.Uvs
	}

	if frag.RenderFlags&0x40 == 0x40 {
		e.TwoSided = 1
	}

	if frag.Flags&0x10 != 0 {
		e.HexTenFlag = 1
	}

	return nil
}

func spriteVariationToRaw(wce *Wce, rawWld *raw.Wld, e WldDefinitioner) error {
	var err error
	tag := ""
	switch spriteDef := e.(type) {
	case *DMSpriteDef2:
		tag = spriteDef.Tag
	case *DMSpriteDef:
		tag = spriteDef.Tag
	default:
		return fmt.Errorf("unknown type %T", e)
	}
	tag = strings.TrimSuffix(tag, "_DMSPRITEDEF")
	var index int
	if len(tag) >= 5 {
		index, err = strconv.Atoi(tag[len(tag)-2:])
		if err != nil {
			return nil
			//return fmt.Errorf("tag index: %w", err)
		}
		tag = tag[:len(tag)-2]
	} else {
		index = 0
	}

	tagLong := strings.TrimSuffix(tag, "_DMSPRITEDEF")
	if tagLong == tag {
		tagLong = ""
	}

	// check for variations
	for i := 0; i < 10; i++ {
		if i <= index {
			continue
		}
		variationTag := fmt.Sprintf("%s%02d_DMSPRITEDEF", tag, i)
		def := wce.ByTag(variationTag)
		if def == nil {
			return nil
		}
		_, err = def.ToRaw(wce, rawWld)
		if err != nil {
			return fmt.Errorf("%s to raw: %w", variationTag, err)
		}

		if tagLong != "" {
			variationTag = fmt.Sprintf("%s%02d_DMSPRITEDEF", tagLong, i)
			def := wce.ByTag(variationTag)
			if def == nil {
				return nil
			}
			_, err = def.ToRaw(wce, rawWld)
			if err != nil {
				return fmt.Errorf("%s to raw: %w", variationTag, err)
			}
		}
	}

	return nil
}

type DMTrackDef2 struct {
	fragID  int16
	model   string
	Tag     string
	Sleep   uint16
	Param2  uint16
	FPScale uint16
	Frames  [][][3]float32
	Size6   uint16
}

func (e *DMTrackDef2) Definition() string {
	return "DMTRACKDEF2"
}

func (e *DMTrackDef2) Write(token *AsciiWriteToken) error {

	w, err := token.UseTempWriter(e.model + "_ani")
	if err != nil {
		return err
	}

	if token.TagIsWritten(e.Tag) {
		return nil
	}

	token.TagSetIsWritten(e.Tag)

	if e.Sleep == 0 {
		return fmt.Errorf("sleep is 0 for dmtrackdef2 %s on %s, this isn't handled report to Xackery", e.Tag, e.model)
	}

	fmt.Fprintf(w, "%s \"%s\"\n", e.Definition(), e.Tag)
	fmt.Fprintf(w, "\tSLEEP %d\n", e.Sleep)
	fmt.Fprintf(w, "\tPARAM2 %d\n", e.Param2)
	fmt.Fprintf(w, "\tFPSCALE %d\n", e.FPScale)
	fmt.Fprintf(w, "\tSIZE6 %d\n", e.Size6)

	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(e.Frames))
	for _, vertFrames := range e.Frames {
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(vertFrames))
		for _, frame := range vertFrames {
			fmt.Fprintf(w, "\t\t\tXYZ %0.8e %0.8e %0.8e\n", frame[0], frame[1], frame[2])
		}
	}
	fmt.Fprintf(w, "\n")

	return nil
}

func (e *DMTrackDef2) Read(token *AsciiReadToken) error {
	records, err := token.ReadProperty("SLEEP", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Sleep, records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = token.ReadProperty("PARAM2", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Param2, records[1])
	if err != nil {
		return fmt.Errorf("param2: %w", err)
	}

	records, err = token.ReadProperty("FPSCALE", 1)
	if err != nil {
		return err
	}
	err = parse(&e.FPScale, records[1])
	if err != nil {
		return fmt.Errorf("fpscale: %w", err)
	}

	records, err = token.ReadProperty("SIZE6", 1)
	if err != nil {
		return err
	}
	err = parse(&e.Size6, records[1])
	if err != nil {
		return fmt.Errorf("size6: %w", err)
	}

	records, err = token.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return err
	}
	numFrames := int(0)
	err = parse(&numFrames, records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	originalVerts := 0
	for i := 0; i < numFrames; i++ {
		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}
		numVertices := int(0)
		err = parse(&numVertices, records[1])
		if err != nil {
			return fmt.Errorf("frame %d num vertices: %w", i, err)
		}

		if i == 0 {
			originalVerts = numVertices
		}
		if originalVerts != numVertices {
			return fmt.Errorf("frame %d has different number of vertices than original frame", i)
		}
		var vertFrames [][3]float32
		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}
			frame := [3]float32{}
			err = parse(&frame, records[1:]...)
			if err != nil {
				return fmt.Errorf("frame %d vertex %d: %w", i, j, err)
			}
			vertFrames = append(vertFrames, frame)
		}
		e.Frames = append(e.Frames, vertFrames)
	}

	return nil
}

func (e *DMTrackDef2) ToRaw(wce *Wce, rawWld *raw.Wld) (int16, error) {
	//if e.fragID != 0 {
	//	return e.fragID, nil
	//}

	wfTrack2 := &rawfrag.WldFragDmTrackDef2{
		Sleep:  e.Sleep,
		Param2: e.Param2,
		Scale:  e.FPScale,
		Size6:  e.Size6,
	}

	scale := float32(1 / float32(int(1)<<int(e.FPScale)))

	for _, frame := range e.Frames {
		frames := make([][3]int16, 0)
		for _, vert := range frame {
			frames = append(frames, [3]int16{
				int16(vert[0] / scale),
				int16(vert[1] / scale),
				int16(vert[2] / scale),
			})
		}
		wfTrack2.Frames = append(wfTrack2.Frames, frames)
	}

	wfTrack2.NameRef = rawWld.NameAdd(e.Tag)
	// flags?
	rawWld.Fragments = append(rawWld.Fragments, wfTrack2)
	e.fragID = int16(len(rawWld.Fragments))

	wfTrack := &rawfrag.WldFragDMTrack{
		TrackRef: int32(e.fragID),
		Flags:    0,
	}
	rawWld.Fragments = append(rawWld.Fragments, wfTrack)
	return int16(e.fragID), nil
}

func (e *DMTrackDef2) FromRaw(wce *Wce, rawWld *raw.Wld, frag *rawfrag.WldFragDmTrackDef2) error {
	if frag == nil {
		return fmt.Errorf("frag is not trackdef (wrong fragcode?)")
	}

	e.Tag = rawWld.Name(frag.NameRef)
	e.Sleep = frag.Sleep
	e.Param2 = frag.Param2
	e.FPScale = frag.Scale
	e.Size6 = frag.Size6

	scale := 1.0 / float32(int(1<<frag.Scale))

	for _, frame := range frag.Frames {
		frames := make([][3]float32, 0)
		for _, vert := range frame {
			frames = append(frames, [3]float32{
				float32(vert[0]) * scale,
				float32(vert[1]) * scale,
				float32(vert[2]) * scale,
			})
		}
		e.Frames = append(e.Frames, frames)
	}

	e.model = wce.lastReadModelTag
	if frag.Flags != 0 {
		return fmt.Errorf("unknown flags %d", frag.Flags)
	}

	return nil
}
