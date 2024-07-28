package wld

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xackery/quail/helper"
)

// DMSpriteDef2 is a declaration of DMSpriteDef2
type DMSpriteDef2 struct {
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
	SkinAssignmentGroups [][2]uint16
	MaterialPaletteTag   string
	Colors               [][4]uint8
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

func (d *DMSpriteDef2) Definition() string {
	return "DMSPRITEDEF2"
}

func (d *DMSpriteDef2) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", d.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", d.Tag)
	fmt.Fprintf(w, "\t// FLAGS \"%d\" // need to assess\n", d.Flags)
	fmt.Fprintf(w, "\tCENTEROFFSET %0.7e %0.7e %0.7e\n", d.CenterOffset[0], d.CenterOffset[1], d.CenterOffset[2])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(d.Vertices))
	for _, vert := range d.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMUVS %d\n", len(d.UVs))
	for _, uv := range d.UVs {
		fmt.Fprintf(w, "\tUV %0.7e %0.7e\n", uv[0], uv[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTEXNORMALS %d\n", len(d.VertexNormals))
	for _, vn := range d.VertexNormals {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vn[0], vn[1], vn[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tSKINASSIGNMENTGROUPS %d", len(d.SkinAssignmentGroups))
	for _, sa := range d.SkinAssignmentGroups {
		fmt.Fprintf(w, " %d %d", sa[0], sa[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tMATERIALPALETTE \"%s\"\n", d.MaterialPaletteTag)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tPOLYHEDRON\n")
	fmt.Fprintf(w, "\t\tDEFINITION \"%s\"\n", d.PolyhedronTag)
	fmt.Fprintf(w, "\tENDPOLYHEDRON\n\n")
	fmt.Fprintf(w, "\tNUMFACE2S %d\n", len(d.Faces))
	fmt.Fprintf(w, "\n")
	for i, face := range d.Faces {
		fmt.Fprintf(w, "\tDMFACE2 //%d\n", i+1)
		fmt.Fprintf(w, "\t\tFLAGS %d\n", face.Flags)
		fmt.Fprintf(w, "\t\tTRIANGLE   %d %d %d\n", face.Triangle[0], face.Triangle[1], face.Triangle[2])
		fmt.Fprintf(w, "\tENDDMFACE2 //%d\n\n", i+1)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\t// meshops are not supported\n")
	fmt.Fprintf(w, "\t// NUMMESHOPS %d\n", len(d.MeshOps))
	for _, meshOp := range d.MeshOps {
		fmt.Fprintf(w, "\t// TODO: MESHOP %d %d %0.7f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFACEMATERIALGROUPS %d", len(d.FaceMaterialGroups))
	for _, group := range d.FaceMaterialGroups {
		fmt.Fprintf(w, " %d %d", group[0], group[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tVERTEXMATERIALGROUPS %d", len(d.VertexMaterialGroups))
	for _, group := range d.VertexMaterialGroups {
		fmt.Fprintf(w, " %d %d", group[0], group[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7e\n", d.BoundingRadius)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFPSCALE %d\n", d.FPScale)
	fmt.Fprintf(w, "ENDDMSPRITEDEF2\n\n")
	return nil
}

func (d *DMSpriteDef2) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	d.Tag = records[1]

	records, err = r.ReadProperty("CENTEROFFSET", 3)
	if err != nil {
		return err
	}
	d.CenterOffset, err = helper.ParseFloat32Slice3(records[1:])
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = r.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices, err := helper.ParseInt(records[1])
	if err != nil {
		return err
	}
	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}
		d.Vertices = append(d.Vertices, vert)
	}

	records, err = r.ReadProperty("NUMUVS", 1)
	if err != nil {
		return err
	}
	numUVs, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num uvs: %w", err)
	}

	for i := 0; i < numUVs; i++ {
		records, err = r.ReadProperty("UV", 2)
		if err != nil {
			return err
		}
		uv, err := helper.ParseFloat32Slice2(records[1:])
		if err != nil {
			return fmt.Errorf("uv %d: %w", i, err)
		}
		d.UVs = append(d.UVs, uv)
	}

	records, err = r.ReadProperty("NUMVERTEXNORMALS", 1)
	if err != nil {
		return err
	}
	numNormals, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num normals: %w", err)
	}

	for i := 0; i < numNormals; i++ {
		records, err = r.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		norm, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("normal %d: %w", i, err)
		}
		d.VertexNormals = append(d.VertexNormals, norm)
	}

	records, err = r.ReadProperty("SKINASSIGNMENTGROUPS", 1)
	if err != nil {
		return err
	}
	numSkinAssignments, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num skin assignments: %w", err)
	}

	for i := 0; i < numSkinAssignments; i++ {
		val1, err := strconv.ParseUint(records[i+1], 10, 16)
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}
		val2, err := strconv.ParseUint(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("skin assignment %d: %w", i, err)
		}
		d.SkinAssignmentGroups = append(d.SkinAssignmentGroups, [2]uint16{uint16(val1), uint16(val2)})
	}

	records, err = r.ReadProperty("MATERIALPALETTE", 1)
	if err != nil {
		return err
	}
	d.MaterialPaletteTag = records[1]

	_, err = r.ReadProperty("POLYHEDRON", 0)
	if err != nil {
		return err
	}
	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	d.PolyhedronTag = records[1]
	_, err = r.ReadProperty("ENDPOLYHEDRON", 0)
	if err != nil {
		return err
	}

	records, err = r.ReadProperty("NUMFACE2S", 1)
	if err != nil {
		return err
	}

	numFaces, err := helper.ParseInt(records[1])
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
		face.Flags, err = helper.ParseUint16(records[1])
		if err != nil {
			return fmt.Errorf("face %d flags: %w", i, err)
		}

		records, err = r.ReadProperty("TRIANGLE", 3)
		if err != nil {
			return err
		}
		face.Triangle, err = helper.ParseUint16Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("face %d triangle: %w", i, err)
		}

		_, err = r.ReadProperty("ENDDMFACE2", 0)
		if err != nil {
			return err
		}

		d.Faces = append(d.Faces, face)
	}

	records, err = r.ReadProperty("FACEMATERIALGROUPS", 1)
	if err != nil {
		return err
	}
	numFaceMaterialGroups, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num face material groups: %w", err)
	}

	for i := 0; i < numFaceMaterialGroups; i++ {
		val1, err := strconv.ParseUint(records[i+1], 10, 16)
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		val2, err := strconv.ParseUint(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("face material group %d: %w", i, err)
		}
		d.FaceMaterialGroups = append(d.FaceMaterialGroups, [2]uint16{uint16(val1), uint16(val2)})
	}

	records, err = r.ReadProperty("VERTEXMATERIALGROUPS", 1)
	if err != nil {
		return err
	}
	numVertexMaterialGroups, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vertex material groups: %w", err)
	}

	for i := 0; i < numVertexMaterialGroups; i++ {
		val1, err := strconv.ParseInt(records[i+1], 10, 16)
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		val2, err := strconv.ParseInt(records[i+2], 10, 16)
		if err != nil {
			return fmt.Errorf("vertex material group %d: %w", i, err)
		}
		d.VertexMaterialGroups = append(d.VertexMaterialGroups, [2]int16{int16(val1), int16(val2)})
	}

	records, err = r.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	d.BoundingRadius, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("FPSCALE", 1)
	if err != nil {
		return err
	}
	d.FPScale, err = helper.ParseUint16(records[1])
	if err != nil {
		return fmt.Errorf("fpscale: %w", err)
	}

	_, err = r.ReadProperty("ENDDMSPRITEDEF2", 0)
	if err != nil {
		return err
	}

	return nil
}

// DMSpriteDef is a declaration of DMSPRITEDEF
type DMSpriteDef struct {
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

// DMSprite is a declaration of DMSPRITEINSTANCE
type DMSprite struct {
	Tag           string
	DefinitionTag string
	Param         uint32
}

func (d *DMSprite) Definition() string {
	return "DMSPRITEINSTANCE"
}

func (d *DMSprite) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", d.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", d.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", d.DefinitionTag)
	fmt.Fprintf(w, "\tPARAM %d\n", d.Param)
	fmt.Fprintf(w, "ENDDMSPRITEINSTANCE\n\n")
	return nil
}

func (d *DMSprite) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	d.Tag = records[1]

	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}

	d.DefinitionTag = records[1]

	records, err = r.ReadProperty("PARAM", 1)
	if err != nil {
		return err
	}

	d.Param, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("param: %w", err)
	}

	_, err = r.ReadProperty("ENDDMSPRITEINSTANCE", 0)
	if err != nil {
		return err
	}

	return nil
}

// MaterialPalette is a declaration of MATERIALPALETTE
type MaterialPalette struct {
	Tag       string
	flags     uint32
	Materials []string
}

func (m *MaterialPalette) Definition() string {
	return "MATERIALPALETTE"
}

func (m *MaterialPalette) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", m.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", m.Tag)
	fmt.Fprintf(w, "\tNUMMATERIALS %d\n", len(m.Materials))
	for _, mat := range m.Materials {
		fmt.Fprintf(w, "\tMATERIAL \"%s\"\n", mat)
	}
	fmt.Fprintf(w, "ENDMATERIALPALETTE\n\n")
	return nil
}

func (m *MaterialPalette) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	m.Tag = records[1]

	records, err = r.ReadProperty("NUMMATERIALS", 1)
	if err != nil {
		return fmt.Errorf("NUMMATERIALS: %w", err)
	}
	numMaterials, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num materials: %w", err)
	}

	for i := 0; i < numMaterials; i++ {
		records, err = r.ReadProperty("MATERIAL", 1)
		if err != nil {
			return fmt.Errorf("MATERIAL: %w", err)
		}
		m.Materials = append(m.Materials, records[1])
	}

	_, err = r.ReadProperty("ENDMATERIALPALETTE", 0)
	if err != nil {
		return fmt.Errorf("ENDMATERIALPALETTE: %w", err)
	}

	return nil
}

// MaterialDef is an entry MATERIALDEFINITION
type MaterialDef struct {
	Tag                  string   // TAG %s
	Flags                uint32   // FLAGS %d
	RenderMethod         string   // RENDERMETHOD %s
	RGBPen               [4]uint8 // RGBPEN %d %d %d
	Brightness           float32  // BRIGHTNESS %0.7f
	ScaledAmbient        float32  // SCALEDAMBIENT %0.7f
	SimpleSpriteInstTag  string   // SIMPLESPRITEINST
	SimpleSpriteInstFlag uint32   // FLAGS %d
	Pair1                uint32
	Pair2                float32
}

func (m *MaterialDef) Definition() string {
	return "MATERIALDEFINITION"
}

func (m *MaterialDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", m.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", m.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", m.Flags)
	fmt.Fprintf(w, "\tRENDERMETHOD %s\n", m.RenderMethod)
	fmt.Fprintf(w, "\tRGBPEN %d %d %d\n", m.RGBPen[0], m.RGBPen[1], m.RGBPen[2])
	fmt.Fprintf(w, "\tBRIGHTNESS %0.7f\n", m.Brightness)
	fmt.Fprintf(w, "\tSCALEDAMBIENT %0.7f\n", m.ScaledAmbient)
	fmt.Fprintf(w, "\tSIMPLESPRITEINST\n")
	fmt.Fprintf(w, "\t\tTAG \"%s\"\n", m.SimpleSpriteInstTag)
	fmt.Fprintf(w, "\t\t// FLAGS %d\n", m.SimpleSpriteInstFlag)
	fmt.Fprintf(w, "\tENDSIMPLESPRITEINST\n")
	fmt.Fprintf(w, "\t// PAIR1 %d\n", m.Pair1)
	fmt.Fprintf(w, "\t// PAIR2 %0.7f\n", m.Pair2)
	fmt.Fprintf(w, "ENDMATERIALDEFINITION\n\n")
	return nil
}

func (m *MaterialDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	m.Tag = records[1]

	records, err = r.ReadProperty("RENDERMETHOD", 1)
	if err != nil {
		return err
	}
	m.RenderMethod = records[1]

	records, err = r.ReadProperty("RGBPEN", 3)
	if err != nil {
		return err
	}
	m.RGBPen, err = helper.ParseUint8Slice4(records[1:])
	if err != nil {
		return fmt.Errorf("rgbpen: %w", err)
	}

	records, err = r.ReadProperty("BRIGHTNESS", 1)
	if err != nil {
		return err
	}
	m.Brightness, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("brightness: %w", err)
	}

	records, err = r.ReadProperty("SCALEDAMBIENT", 1)
	if err != nil {
		return err
	}
	m.ScaledAmbient, err = helper.ParseFloat32(records[1])
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
	m.SimpleSpriteInstTag = records[1]

	_, err = r.ReadProperty("ENDSIMPLESPRITEINST", 0)
	if err != nil {
		return err
	}

	_, err = r.ReadProperty("ENDMATERIALDEFINITION", 0)
	if err != nil {
		return err
	}

	return nil
}

// SimpleSpriteDef is a declaration of SIMPLESPRITEDEF
type SimpleSpriteDef struct {
	Tag                string
	SimpleSpriteFrames []SimpleSpriteFrame
}

type SimpleSpriteFrame struct {
	TextureFile string
	TextureTag  string
}

func (s *SimpleSpriteDef) Definition() string {
	return "SIMPLESPRITEDEF"
}

func (s *SimpleSpriteDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", s.Definition())
	fmt.Fprintf(w, "\tSIMPLESPRITETAG \"%s\"\n", s.Tag)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(s.SimpleSpriteFrames))
	for _, frame := range s.SimpleSpriteFrames {
		fmt.Fprintf(w, "\tFRAME \"%s\" \"%s\"\n", frame.TextureFile, frame.TextureTag)
	}
	fmt.Fprintf(w, "ENDSIMPLESPRITEDEF\n\n")
	return nil
}

func (s *SimpleSpriteDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("SIMPLESPRITETAG", 0)
	if err != nil {
		return fmt.Errorf("SIMPLESPRITETAG: %w", err)
	}
	s.Tag = records[1]

	records, err = r.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return fmt.Errorf("NUMFRAMES: %w", err)
	}
	numFrames, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = r.ReadProperty("FRAME", 2)
		if err != nil {
			return fmt.Errorf("FRAME: %w", err)
		}
		s.SimpleSpriteFrames = append(s.SimpleSpriteFrames, SimpleSpriteFrame{
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

// ActorDef is a declaration of ACTORDEF
type ActorDef struct {
	Tag           string
	Callback      string
	BoundsRef     int32
	CurrentAction uint32
	Location      [6]float32
	Unk1          uint32
	Actions       []ActorAction
	Unk2          uint32
}

func (a *ActorDef) Definition() string {
	return "ACTORDEF"
}

func (a *ActorDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", a.Definition())
	fmt.Fprintf(w, "\tACTORTAG \"%s\"\n", a.Tag)
	fmt.Fprintf(w, "\tCALLBACK \"%s\"\n", a.Callback)
	fmt.Fprintf(w, "\t// BOUNDSREF %d\n", a.BoundsRef)
	fmt.Fprintf(w, "\tCURRENTACTION %d\n", a.CurrentAction)
	fmt.Fprintf(w, "\tLOCATION %0.7f %0.7f %0.7f %0.7f %0.7f %0.7f\n", a.Location[0], a.Location[1], a.Location[2], a.Location[3], a.Location[4], a.Location[5])
	fmt.Fprintf(w, "\tNUMACTIONS %d\n", len(a.Actions))
	for _, action := range a.Actions {
		fmt.Fprintf(w, "\tACTION\n")
		fmt.Fprintf(w, "\t\t// UNK1 %d\n", action.Unk1)
		fmt.Fprintf(w, "\t\tNUMLEVELSOFDETAIL %d\n", len(action.LevelOfDetails))
		for _, lod := range action.LevelOfDetails {
			fmt.Fprintf(w, "\t\tLEVELOFDETAIL\n")
			fmt.Fprintf(w, "\t\t\t2DSPRITE \"%s\"\n", lod.Sprite2DTag)
			fmt.Fprintf(w, "\t\t\t3DSPRITE \"%s\"\n", lod.Sprite3DTag)
			fmt.Fprintf(w, "\t\t\tDMSPRITE \"%s\"\n", lod.DMSpriteTag)
			fmt.Fprintf(w, "\t\t\tHIERARCHICALSPRITE \"%s\"\n", lod.HierarchicalSpriteTag)
			fmt.Fprintf(w, "\t\t\tMINDISTANCE %0.7f\n", lod.MinDistance)
			fmt.Fprintf(w, "\t\tENDLEVELOFDETAIL\n")
		}
		fmt.Fprintf(w, "\tENDACTION\n")
	}
	fmt.Fprintf(w, "\t// UNK2 %d\n", a.Unk2)
	fmt.Fprintf(w, "ENDACTORDEF\n\n")
	return nil
}

func (a *ActorDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("ACTORTAG", 1)
	if err != nil {
		return err
	}
	a.Tag = records[1]

	records, err = r.ReadProperty("CALLBACK", 1)
	if err != nil {
		return err
	}
	a.Callback = records[1]

	records, err = r.ReadProperty("CURRENTACTION", 1)
	if err != nil {
		return err
	}
	a.CurrentAction, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("current action: %w", err)
	}

	records, err = r.ReadProperty("LOCATION", 3)
	if err != nil {
		return err
	}
	a.Location, err = helper.ParseFloat32Slice6(records[1:])
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = r.ReadProperty("NUMACTIONS", 1)
	if err != nil {
		return err
	}
	numActions, err := helper.ParseInt(records[1])
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

		numLod, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num lod: %w", err)
		}

		for j := 0; j < numLod; j++ {
			lod := ActorLevelOfDetail{}
			_, err = r.ReadProperty("LEVELOFDETAIL", 0)
			if err != nil {
				return err
			}

			records, err = r.ReadProperty("2DSPRITE", 1)
			if err != nil {
				return err
			}
			lod.Sprite2DTag = records[1]

			records, err = r.ReadProperty("3DSPRITE", 1)
			if err != nil {
				return err
			}
			lod.Sprite3DTag = records[1]

			records, err = r.ReadProperty("DMSPRITE", 1)
			if err != nil {
				return err
			}
			lod.DMSpriteTag = records[1]

			records, err = r.ReadProperty("HIERARCHICALSPRITE", 1)
			if err != nil {
				return err
			}
			lod.HierarchicalSpriteTag = records[1]

			records, err = r.ReadProperty("MINDISTANCE", 1)
			if err != nil {
				return err
			}

			lod.MinDistance, err = helper.ParseFloat32(records[1])
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

		a.Actions = append(a.Actions, action)

	}

	return nil
}

// ActorAction is a declaration of ACTION
type ActorAction struct {
	Unk1           uint32
	LevelOfDetails []ActorLevelOfDetail
}

// ActorLevelOfDetail is a declaration of LEVELOFDETAIL
type ActorLevelOfDetail struct {
	Sprite3DTag           string
	HierarchicalSpriteTag string
	DMSpriteTag           string
	Sprite2DTag           string
	MinDistance           float32
}

// ActorInst is a declaration of ACTORINST
type ActorInst struct {
	Tag              string
	Active           int
	SpriteVolumeOnly int
	SphereTag        string
	SoundTag         string
	DefinitionTag    string
	CurrentAction    uint32
	Location         [6]float32
	Unk1             uint32
	BoundingRadius   float32
	Scale            float32
	DMRGBTrackTag    string
	UserData         string
}

func (a *ActorInst) Definition() string {
	return "ACTORINST"
}

func (a *ActorInst) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", a.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", a.Tag)
	fmt.Fprintf(w, "\tACTIVE %d\n", a.Active)
	fmt.Fprintf(w, "\tSPRITEVOLUMEONLY %d\n", a.SpriteVolumeOnly)
	fmt.Fprintf(w, "\tLOCATION %0.7f %0.7f %0.7f %0.7f %0.7f %0.7f\n", a.Location[0], a.Location[1], a.Location[2], a.Location[3], a.Location[4], a.Location[5])
	fmt.Fprintf(w, "\t// UNK1 %d\n", a.Unk1)
	fmt.Fprintf(w, "\tCURRENTACTION %d\n", a.CurrentAction)
	fmt.Fprintf(w, "\tSPHERE \"%s\"\n", a.SphereTag)
	fmt.Fprintf(w, "\tSOUND \"%s\"\n", a.SoundTag)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", a.BoundingRadius)
	fmt.Fprintf(w, "\tSCALEFACTOR %0.7f\n", a.Scale)
	fmt.Fprintf(w, "\tDMRGBTRACK \"%s\"\n", a.DMRGBTrackTag)
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", a.UserData)
	fmt.Fprintf(w, "ENDACTORINST\n\n")
	return nil
}

func (a *ActorInst) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	a.Tag = records[1]

	records, err = r.ReadProperty("ACTIVE", 1)
	if err != nil {
		return err
	}
	a.Active, err = helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("active: %w", err)
	}

	records, err = r.ReadProperty("SPRITEVOLUMEONLY", 1)
	if err != nil {
		return err
	}
	a.SpriteVolumeOnly, err = helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("sprite volume only: %w", err)
	}

	records, err = r.ReadProperty("LOCATION", 6)
	if err != nil {
		return err
	}
	a.Location, err = helper.ParseFloat32Slice6(records[1:])
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = r.ReadProperty("CURRENTACTION", 1)
	if err != nil {
		return err
	}
	a.CurrentAction, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("current action: %w", err)
	}

	records, err = r.ReadProperty("SPHERE", 1)
	if err != nil {
		return err
	}
	a.SphereTag = records[1]

	records, err = r.ReadProperty("SOUND", 1)
	if err != nil {
		return err
	}
	a.SoundTag = records[1]

	records, err = r.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	a.BoundingRadius, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("SCALEFACTOR", 1)
	if err != nil {
		return err
	}
	a.Scale, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	records, err = r.ReadProperty("DMRGBTRACK", 1)
	if err != nil {
		return err
	}
	a.DMRGBTrackTag = records[1]

	records, err = r.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}
	a.UserData = records[1]

	_, err = r.ReadProperty("ENDACTORINST", 0)
	if err != nil {
		return err
	}
	return nil
}

// LightDef is a declaration of LIGHTDEF
type LightDef struct {
	Tag             string
	Flags           uint32
	FrameCurrentRef uint32
	Sleep           uint32
	LightLevels     []float32
	Colors          [][3]float32
}

func (l *LightDef) Definition() string {
	return "LIGHTDEFINITION"
}

func (l *LightDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", l.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", l.Tag)

	fmt.Fprintf(w, "\tCURRENTFRAME %d\n", l.FrameCurrentRef)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(l.LightLevels))
	for _, level := range l.LightLevels {
		fmt.Fprintf(w, "\tLIGHTLEVELS %0.6f\n", level)
	}
	fmt.Fprintf(w, "\tSLEEP %d\n", l.Sleep)
	isSkipFrames := 0
	if l.Flags&0x08 == 0x08 {
		isSkipFrames = 1
	}
	fmt.Fprintf(w, "\tSKIPFRAMES %d\n", isSkipFrames)
	fmt.Fprintf(w, "\tNUMCOLORS %d\n", len(l.Colors))
	for _, color := range l.Colors {
		fmt.Fprintf(w, "\tCOLOR %0.6f %0.6f %0.6f\n", color[0], color[1], color[2])
	}
	fmt.Fprintf(w, "ENDLIGHTDEFINITION\n\n")
	return nil
}

func (l *LightDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	l.Tag = records[1]

	records, err = r.ReadProperty("CURRENTFRAME", 1)
	if err != nil {
		return err
	}
	l.FrameCurrentRef, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("current frame: %w", err)
	}

	records, err = r.ReadProperty("NUMFRAMES", 1)
	if err != nil {
		return err
	}
	numFrames, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = r.ReadProperty("LIGHTLEVELS", 1)
		if err != nil {
			return err
		}
		level, err := helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("light level: %w", err)
		}
		l.LightLevels = append(l.LightLevels, level)
	}

	records, err = r.ReadProperty("SLEEP", 1)
	if err != nil {
		return err
	}
	l.Sleep, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	records, err = r.ReadProperty("SKIPFRAMES", 1)
	if err != nil {
		return err
	}
	if records[1] == "1" {
		l.Flags |= 0x08
	}

	records, err = r.ReadProperty("NUMCOLORS", 1)
	if err != nil {
		return err
	}
	numColors, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num colors: %w", err)
	}

	for i := 0; i < numColors; i++ {
		records, err = r.ReadProperty("COLOR", 3)
		if err != nil {
			return err
		}
		color, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("color: %w", err)
		}

		l.Colors = append(l.Colors, color)
	}

	_, err = r.ReadProperty("ENDLIGHTDEFINITION", 0)
	if err != nil {
		return err
	}
	return nil
}

// PointLight is a declaration of POINTLIGHT
type PointLight struct {
	Tag         string // TAG "%s"
	LightDefTag string // LIGHT "%s"
	Flags       uint32 // FLAGS %d
	Location    [3]float32
	Radius      float32 // RADIUSOFINFLUENCE %0.7f
}

func (p *PointLight) Definition() string {
	return "POINTLIGHT"
}

func (p *PointLight) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", p.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", p.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", p.Flags)
	fmt.Fprintf(w, "\tXYZ %0.6f %0.6f %0.6f\n", p.Location[0], p.Location[1], p.Location[2])
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", p.LightDefTag)
	fmt.Fprintf(w, "\tRADIUSOFINFLUENCE %0.7f\n", p.Radius)
	fmt.Fprintf(w, "ENDPOINTLIGHT\n\n")
	return nil
}

func (p *PointLight) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	p.Tag = records[1]

	records, err = r.ReadProperty("LIGHT", 1)
	if err != nil {
		return err
	}
	p.LightDefTag = records[1]

	records, err = r.ReadProperty("XYZ", 3)
	if err != nil {
		return err
	}

	p.Location, err = helper.ParseFloat32Slice3(records[1:])
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}

	records, err = r.ReadProperty("RADIUSOFINFLUENCE", 1)
	if err != nil {
		return err
	}
	p.Radius, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("radius of influence: %w", err)
	}

	return nil
}

// Sprite3DDef is a declaration of SPRITE3DDEF
type Sprite3DDef struct {
	Tag            string
	CenterOffset   [3]float32
	BoundingRadius float32
	SphereListTag  string
	Vertices       [][3]float32
	BSPNodes       []*BSPNode
}

// BSPNode is a declaration of BSPNODE
type BSPNode struct {
	Vertices           []uint32
	RenderMethod       string
	Flags              uint8
	Pen                uint32
	Brightness         float32
	ScaledAmbient      float32
	SpriteTag          string
	Origin             [3]float32
	UAxis              [3]float32
	VAxis              [3]float32
	RenderUVMapEntries []BspNodeUVInfo
	FrontTree          uint32
	BackTree           uint32
}

// BspNodeUVInfo is a declaration of UV
type BspNodeUVInfo struct {
	UvOrigin [3]float32 // UV %0.7f %0.7f %0.7f
	UAxis    [3]float32 // UAXIS %0.7f %0.7f %0.7f
	VAxis    [3]float32 // VAXIS %0.7f %0.7f %0.7f
}

func (s *Sprite3DDef) Definition() string {
	return "3DSPRITEDEF"
}

func (s *Sprite3DDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", s.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", s.Tag)
	fmt.Fprintf(w, "\tCENTEROFFSET %0.7f %0.7f %0.7f\n", s.CenterOffset[0], s.CenterOffset[1], s.CenterOffset[2])
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", s.BoundingRadius)
	fmt.Fprintf(w, "\tSPHERELIST \"%s\"\n", s.SphereListTag)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(s.Vertices))
	for _, vert := range s.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMBSPNODES %d\n", len(s.BSPNodes))
	for i, node := range s.BSPNodes {
		fmt.Fprintf(w, "\tBSPNODE //%d\n", i+1)
		fmt.Fprintf(w, "\t\tVERTEXLIST %d", len(node.Vertices))
		for _, vert := range node.Vertices {
			fmt.Fprintf(w, " %d", vert)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\t\tRENDERMETHOD %s\n", node.RenderMethod)
		fmt.Fprintf(w, "\t\tRENDERINFO\n")
		fmt.Fprintf(w, "\t\t\tPEN %d\n", node.Pen)
		fmt.Fprintf(w, "\t\tENDRENDERINFO\n")
		fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\t\tSPRITE \"%s\"\n", node.SpriteTag)
		fmt.Fprintf(w, "\tENDBSPNODE\n")
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

	records, err = r.ReadProperty("CENTEROFFSET", 3)
	if err != nil {
		return err
	}
	s.CenterOffset, err = helper.ParseFloat32Slice3(records[1:])
	if err != nil {
		return fmt.Errorf("center offset: %w", err)
	}

	records, err = r.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	s.BoundingRadius, err = helper.ParseFloat32(records[1])
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
	numVertices, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}
		s.Vertices = append(s.Vertices, vert)
	}

	records, err = r.ReadProperty("NUMBSPNODES", 1)
	if err != nil {
		return err
	}
	numBSPNodes, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num bsp nodes: %w", err)
	}

	for i := 0; i < numBSPNodes; i++ {
		node := &BSPNode{}
		_, err = r.ReadProperty("BSPNODE", 0)
		if err != nil {
			return err
		}
		records, err = r.ReadProperty("VERTEXLIST", 1)
		if err != nil {
			return err
		}
		numVertices, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}
		if len(records) != numVertices+2 {
			return fmt.Errorf("vertex list: expected %d, got %d", numVertices, len(records)-2)
		}
		for j := 0; j < numVertices; j++ {
			val, err := helper.ParseUint32(records[j+2])
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

		records, err = r.ReadProperty("PEN", 1)
		if err != nil {
			return err
		}
		node.Pen, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("render pen: %w", err)
		}

		_, err = r.ReadProperty("ENDRENDERINFO", 0)
		if err != nil {
			return err
		}

		records, err = r.ReadProperty("FRONTTREE", 1)
		if err != nil {
			return err
		}

		node.FrontTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = r.ReadProperty("BACKTREE", 1)
		if err != nil {
			return err
		}

		node.BackTree, err = helper.ParseUint32(records[1])
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

type PolyhedronDefinition struct {
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

func (p *PolyhedronDefinition) Definition() string {
	return "POLYHEDRONDEFINITION"
}

func (p *PolyhedronDefinition) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", p.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", p.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", p.Flags)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", p.BoundingRadius)
	fmt.Fprintf(w, "\tSCALEFACTOR %0.7f\n", p.ScaleFactor)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(p.Vertices))
	for _, vert := range p.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMFACES %d\n", len(p.Faces))
	for i, face := range p.Faces {
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

func (p *PolyhedronDefinition) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	p.Tag = records[1]

	records, err = r.ReadProperty("BOUNDINGRADIUS", 1)
	if err != nil {
		return err
	}
	p.BoundingRadius, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("SCALEFACTOR", 1)
	if err != nil {
		return err
	}
	p.ScaleFactor, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	records, err = r.ReadProperty("NUMVERTICES", 1)
	if err != nil {
		return err
	}

	numVertices, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("vertex %d: %w", i, err)
		}
		p.Vertices = append(p.Vertices, vert)
	}

	records, err = r.ReadProperty("NUMFACES", 1)
	if err != nil {
		return err
	}
	numFaces, err := helper.ParseInt(records[1])
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
		numVertices, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}
		if len(records) != numVertices+2 {
			return fmt.Errorf("vertex list: expected %d, got %d", numVertices, len(records)-2)
		}
		for j := 0; j < numVertices; j++ {
			val, err := helper.ParseUint32(records[j+2])
			if err != nil {
				return fmt.Errorf("vertex %d: %w", j, err)
			}
			face.Vertices = append(face.Vertices, val)
		}

		_, err = r.ReadProperty("ENDFACE", 0)
		if err != nil {
			return err
		}

		p.Faces = append(p.Faces, face)
	}

	_, err = r.ReadProperty("ENDPOLYHEDRONDEFINITION", 0)
	if err != nil {
		return err
	}

	return nil
}

type TrackInstance struct {
	Tag           string
	DefinitionTag string
	Interpolate   int
	Reverse       int
	Sleep         uint32
}

func (t *TrackInstance) Definition() string {
	return "TRACKINSTANCE"
}

func (t *TrackInstance) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", t.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", t.DefinitionTag)
	fmt.Fprintf(w, "\tINTERPOLATE %d\n", t.Interpolate)
	fmt.Fprintf(w, "\tREVERSE %d\n", t.Reverse)
	fmt.Fprintf(w, "\tSLEEP %d\n", t.Sleep)
	fmt.Fprintf(w, "ENDTRACKINSTANCE\n\n")
	return nil
}

func (t *TrackInstance) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	t.Tag = records[1]

	records, err = r.ReadProperty("DEFINITION", 1)
	if err != nil {
		return err
	}
	t.DefinitionTag = records[1]

	records, err = r.ReadProperty("INTERPOLATE", 1)
	if err != nil {
		return err
	}
	t.Interpolate, err = helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

	records, err = r.ReadProperty("REVERSE", 1)
	if err != nil {
		return err
	}
	t.Reverse, err = helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("reverse: %w", err)
	}

	records, err = r.ReadProperty("SLEEP", 1)
	if err != nil {
		return err
	}
	t.Sleep, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	_, err = r.ReadProperty("ENDTRACKINSTANCE", 0)
	if err != nil {
		return err
	}

	return nil
}

type TrackDef struct {
	Tag             string
	Flags           uint32
	FrameTransforms []TrackFrameTransform
}

type TrackFrameTransform struct {
	PositionDenom float32
	Rotation      [3]int16
	Position      [3]float32
}

func (t *TrackDef) Definition() string {
	return "TRACKDEFINITION"
}

func (t *TrackDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", t.Tag)
	fmt.Fprintf(w, "\t// FLAGS %d\n", t.Flags)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(t.FrameTransforms))
	for _, frame := range t.FrameTransforms {
		fmt.Fprintf(w, "\tFRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f\n", frame.PositionDenom, frame.Rotation[0], frame.Rotation[1], frame.Rotation[2], frame.Position[0], frame.Position[1], frame.Position[2])
	}
	fmt.Fprintf(w, "ENDTRACKDEFINITION\n\n")
	return nil
}

func (t *TrackDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	t.Tag = records[1]

	return nil
}

type HierarchicalSpriteDef struct {
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

func (h *HierarchicalSpriteDef) Definition() string {
	return "HIERARCHICALSPRITEDEF"
}

func (h *HierarchicalSpriteDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", h.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", h.Tag)
	fmt.Fprintf(w, "\tNUMDAGS %d\n", len(h.Dags))
	for i, dag := range h.Dags {
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
	fmt.Fprintf(w, "\tNUMATTACHEDSKINS %d\n", len(h.AttachedSkins))
	for _, skin := range h.AttachedSkins {
		fmt.Fprintf(w, "\tDMSPRITE \"%s\"\n", skin.DMSpriteTag)
		fmt.Fprintf(w, "\tLINKSKINUPDATESTODAGINDEX %d\n", skin.LinkSkinUpdatesToDagIndex)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\tCENTEROFFSET %0.1f %0.1f %0.1f\n", h.CenterOffset[0], h.CenterOffset[1], h.CenterOffset[2])

	fmt.Fprintf(w, "\tDAGCOLLISION %d\n", h.DagCollision)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7e\n", h.BoundingRadius)

	fmt.Fprintf(w, "ENDHIERARCHICALSPRITEDEF\n\n")
	return nil
}

func (h *HierarchicalSpriteDef) Read(r *AsciiReadToken) error {

	return nil
}

type WorldTree struct {
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

func (t *WorldTree) Definition() string {
	return "WORLDTREE"
}

func (t *WorldTree) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", t.Tag)
	fmt.Fprintf(w, "\tNUMWORLDNODES %d\n", len(t.WorldNodes))
	for i, node := range t.WorldNodes {
		fmt.Fprintf(w, "\tWORLDNODE %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", node.Normals[0], node.Normals[1], node.Normals[2], node.Normals[3])
		fmt.Fprintf(w, "\t\tWORLDREGIONTAG \"%s\"\n", node.WorldRegionTag)
		fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\tENDWORLDNODE %d\n", i+1)
	}
	fmt.Fprintf(w, "ENDWORLDTREE\n\n")
	return nil
}

func (t *WorldTree) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}
	t.Tag = records[1]

	records, err = r.ReadProperty("NUMWORLDNODES", 1)
	if err != nil {
		return err
	}

	numNodes, err := helper.ParseInt(records[1])
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
		node.Normals, err = helper.ParseFloat32Slice4(records[1:])
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
		node.FrontTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = r.ReadProperty("BACKTREE", 1)
		if err != nil {
			return err
		}
		node.BackTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = r.ReadProperty("ENDWORLDNODE", 0)
		if err != nil {
			return err
		}

		t.WorldNodes = append(t.WorldNodes, node)

	}

	_, err = r.ReadProperty("ENDWORLDTREE", 0)
	if err != nil {
		return err
	}

	return nil
}

type Region struct {
	RegionTag        string
	AmbientLightTag  string
	RegionVertices   [][3]float32
	RenderVertices   [][3]float32
	Walls            []*Wall
	Obstacles        []*Obstacle
	CuttingObstacles []*Obstacle
	VisTree          *VisTree
	Sphere           [4]float32
	UserData         string
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
	Range []uint32
}

func (r *Region) Definition() string {
	return "REGION"
}

func (r *Region) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", r.Definition())
	fmt.Fprintf(w, "\tREGIONTAG \"%s\"\n", r.RegionTag)
	fmt.Fprintf(w, "\tNUMREGIONVERTEX %d\n", len(r.RegionVertices))
	for _, vert := range r.RegionVertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMRENDERVERTICES %d\n", len(r.RenderVertices))
	for _, vert := range r.RenderVertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMWALLS %d\n", len(r.Walls))
	for i, wall := range r.Walls {
		fmt.Fprintf(w, "\tWALL // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", wall.Normal[0], wall.Normal[1], wall.Normal[2], wall.Normal[3])
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(wall.Vertices))
		for _, vert := range wall.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\tENDWALL // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tNUMOBSTACLES %d\n", len(r.Obstacles))
	for i, obs := range r.Obstacles {
		fmt.Fprintf(w, "\tOBSTACLE // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", obs.Normal[0], obs.Normal[1], obs.Normal[2], obs.Normal[3])
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(obs.Vertices))
		for _, vert := range obs.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\tENDOBSTACLE // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tNUMCUTTINGOBSTACLES %d\n", len(r.CuttingObstacles))
	for i, obs := range r.CuttingObstacles {
		fmt.Fprintf(w, "\tCUTTINGOBSTACLE // %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", obs.Normal[0], obs.Normal[1], obs.Normal[2], obs.Normal[3])
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(obs.Vertices))
		for _, vert := range obs.Vertices {
			fmt.Fprintf(w, "\t\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\tENDCUTTINGOBSTACLE // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tVISTREE\n")
	fmt.Fprintf(w, "\t\tNUMVISNODE %d\n", len(r.VisTree.VisNodes))
	for i, node := range r.VisTree.VisNodes {
		fmt.Fprintf(w, "\t\tVISNODE // %d\n", i+1)
		fmt.Fprintf(w, "\t\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", node.Normal[0], node.Normal[1], node.Normal[2], node.Normal[3])
		fmt.Fprintf(w, "\t\t\tVISLISTINDEX %d\n", node.VisListIndex)
		fmt.Fprintf(w, "\t\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\t\tENDVISNODE // %d\n", i+1)
	}
	fmt.Fprintf(w, "\t\tNUMVISIBLELIST %d\n", len(r.VisTree.VisLists))
	for i, list := range r.VisTree.VisLists {
		fmt.Fprintf(w, "\t\tVISLIST // %d\n", i+1)
		fmt.Fprintf(w, "\t\t\tRANGE %d", len(list.Range))
		for _, val := range list.Range {
			fmt.Fprintf(w, " %d", val)
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\t\tENDVISIBLELIST // %d\n", i+1)
	}
	fmt.Fprintf(w, "\tENDVISTREE\n")
	fmt.Fprintf(w, "\tSPHERE %0.7f %0.7f %0.7f %0.7f\n", r.Sphere[0], r.Sphere[1], r.Sphere[2], r.Sphere[3])
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", r.UserData)
	fmt.Fprintf(w, "ENDREGION\n\n")
	return nil
}

func (r *Region) Read(token *AsciiReadToken) error {
	r.VisTree = &VisTree{}
	records, err := token.ReadProperty("REGIONTAG", 1)
	if err != nil {
		return err
	}
	r.RegionTag = records[1]

	records, err = token.ReadProperty("NUMREGIONVERTEX", 1)
	if err != nil {
		return err
	}
	numVertices, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num region vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("region vertex %d: %w", i, err)
		}
		r.RegionVertices = append(r.RegionVertices, vert)
	}

	records, err = token.ReadProperty("NUMRENDERVERTICES", 1)
	if err != nil {
		return err
	}
	numVertices, err = helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num render vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ", 3)
		if err != nil {
			return err
		}
		vert, err := helper.ParseFloat32Slice3(records[1:])
		if err != nil {
			return fmt.Errorf("render vertex %d: %w", i, err)
		}
		r.RenderVertices = append(r.RenderVertices, vert)
	}

	records, err = token.ReadProperty("NUMWALLS", 1)
	if err != nil {
		return err
	}
	numWalls, err := helper.ParseInt(records[1])
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
		wall.Normal, err = helper.ParseFloat32Slice4(records[1:])
		if err != nil {
			return fmt.Errorf("wall normal: %w", err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}
		numVertices, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}
			vert, err := helper.ParseFloat32Slice3(records[1:])
			if err != nil {
				return fmt.Errorf("wall vertex %d: %w", j, err)
			}

			wall.Vertices = append(wall.Vertices, vert)
		}

		_, err = token.ReadProperty("ENDWALL", 0)
		if err != nil {
			return err
		}

		r.Walls = append(r.Walls, wall)
	}

	records, err = token.ReadProperty("NUMOBSTACLES", 1)
	if err != nil {
		return err
	}
	numObstacles, err := helper.ParseInt(records[1])
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
		obs.Normal, err = helper.ParseFloat32Slice4(records[1:])
		if err != nil {
			return fmt.Errorf("obstacle normal: %w", err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}
		numVertices, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}
			vert, err := helper.ParseFloat32Slice3(records[1:])
			if err != nil {
				return fmt.Errorf("obstacle vertex %d: %w", j, err)
			}

			obs.Vertices = append(obs.Vertices, vert)
		}

		_, err = token.ReadProperty("ENDOBSTACLE", 0)
		if err != nil {
			return err
		}

		r.Obstacles = append(r.Obstacles, obs)
	}

	records, err = token.ReadProperty("NUMCUTTINGOBSTACLES", 1)
	if err != nil {
		return err
	}
	numObstacles, err = helper.ParseInt(records[1])
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
		obs.Normal, err = helper.ParseFloat32Slice4(records[1:])
		if err != nil {
			return fmt.Errorf("cutting obstacle normal: %w", err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return err
		}

		numVertices, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return err
			}

			vert, err := helper.ParseFloat32Slice3(records[1:])
			if err != nil {
				return fmt.Errorf("cutting obstacle vertex %d: %w", j, err)
			}

			obs.Vertices = append(obs.Vertices, vert)
		}

		_, err = token.ReadProperty("ENDCUTTINGOBSTACLE", 0)
		if err != nil {
			return err
		}

		r.CuttingObstacles = append(r.CuttingObstacles, obs)
	}

	_, err = token.ReadProperty("VISTREE", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("NUMVISNODE", 1)
	if err != nil {
		return err
	}

	numNodes, err := helper.ParseInt(records[1])
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

		node.Normal, err = helper.ParseFloat32Slice4(records[1:])
		if err != nil {
			return fmt.Errorf("vis node normal: %w", err)
		}

		records, err = token.ReadProperty("VISLISTINDEX", 1)
		if err != nil {
			return err
		}

		node.VisListIndex, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("vis list index: %w", err)
		}

		records, err = token.ReadProperty("FRONTTREE", 1)
		if err != nil {
			return err
		}

		node.FrontTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = token.ReadProperty("BACKTREE", 1)
		if err != nil {
			return err
		}

		node.BackTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = token.ReadProperty("ENDVISNODE", 0)
		if err != nil {
			return err
		}

		r.VisTree.VisNodes = append(r.VisTree.VisNodes, node)

	}

	records, err = token.ReadProperty("NUMVISIBLELIST", 1)
	if err != nil {
		return err
	}

	numLists, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num visible lists: %w", err)
	}

	for i := 0; i < numLists; i++ {
		list := &VisList{}
		_, err = token.ReadProperty("VISLIST", 0)
		if err != nil {
			return err
		}

		records, err = token.ReadProperty("RANGE", 1)
		if err != nil {
			return err
		}

		numRanges, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num ranges: %w", err)
		}

		for j := 0; j < numRanges; j++ {
			records, err = token.ReadProperty("RANGE", 1)
			if err != nil {
				return err
			}

			val, err := helper.ParseUint32(records[1])
			if err != nil {
				return fmt.Errorf("range %d: %w", j, err)
			}

			list.Range = append(list.Range, val)
		}

		_, err = token.ReadProperty("ENDVISIBLELIST", 0)
		if err != nil {
			return err
		}

		r.VisTree.VisLists = append(r.VisTree.VisLists, list)
	}

	_, err = token.ReadProperty("ENDVISTREE", 0)
	if err != nil {
		return err
	}

	records, err = token.ReadProperty("SPHERE", 4)
	if err != nil {
		return err
	}

	r.Sphere, err = helper.ParseFloat32Slice4(records[1:])
	if err != nil {
		return fmt.Errorf("sphere: %w", err)
	}

	records, err = token.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}

	r.UserData = records[1]

	_, err = token.ReadProperty("ENDREGION", 0)
	if err != nil {
		return err
	}

	return nil
}

type AmbientLight struct {
	Tag      string
	LightTag string
	Regions  []uint32
}

func (a *AmbientLight) Definition() string {
	return "AMBIENTLIGHT"
}

func (a *AmbientLight) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", a.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", a.Tag)
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", a.LightTag)
	fmt.Fprintf(w, "\tREGIONLIST %d", len(a.Regions))
	for _, region := range a.Regions {
		fmt.Fprintf(w, " %d", region)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "ENDAMBIENTLIGHT\n\n")
	return nil
}

func (a *AmbientLight) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}

	a.Tag = records[1]

	records, err = r.ReadProperty("LIGHT", 1)
	if err != nil {
		return err
	}

	a.LightTag = records[1]

	records, err = r.ReadProperty("REGIONLIST", 1)
	if err != nil {
		return err
	}

	numRegions, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num regions: %w", err)
	}

	for i := 0; i < numRegions; i++ {
		val, err := helper.ParseUint32(records[i+2])
		if err != nil {
			return fmt.Errorf("region %d: %w", i, err)
		}

		a.Regions = append(a.Regions, val)
	}

	_, err = r.ReadProperty("ENDAMBIENTLIGHT", 0)
	if err != nil {
		return err
	}

	return nil
}

type Zone struct {
	Tag      string
	Regions  []uint32
	UserData string
}

func (z *Zone) Definition() string {
	return "ZONE"
}

func (z *Zone) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", z.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", z.Tag)
	fmt.Fprintf(w, "\tREGIONLIST %d", len(z.Regions))
	for _, region := range z.Regions {
		fmt.Fprintf(w, " %d", region)
	}
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", z.UserData)
	fmt.Fprintf(w, "ENDZONE\n\n")
	return nil
}

func (z *Zone) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG", 1)
	if err != nil {
		return err
	}

	z.Tag = records[1]

	records, err = r.ReadProperty("REGIONLIST", 1)
	if err != nil {
		return err
	}

	numRegions, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num regions: %w", err)
	}

	for i := 0; i < numRegions; i++ {
		val, err := helper.ParseUint32(records[i+2])
		if err != nil {
			return fmt.Errorf("region %d: %w", i, err)
		}

		z.Regions = append(z.Regions, val)
	}

	records, err = r.ReadProperty("USERDATA", 1)
	if err != nil {
		return err
	}

	z.UserData = records[1]

	_, err = r.ReadProperty("ENDZONE", 0)
	if err != nil {
		return err
	}

	return nil
}
