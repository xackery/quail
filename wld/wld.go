// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

import (
	"fmt"
	"sync"
)

// Wld is a struct representing a Wld file
type Wld struct {
	FileName           string
	GlobalAmbientLight string
	Version            uint32
	SimpleSpriteDefs   []*SimpleSpriteDef
	MaterialDefs       []*MaterialDef
	MaterialPalettes   []*MaterialPalette
	DMSpriteDef2s      []*DMSpriteDef2
	ActorDefs          []*ActorDef
	ActorInsts         []*ActorInst

	//writing temporary files
	mu                sync.RWMutex
	writtenPalettes   map[string]bool
	writtenMaterials  map[string]bool
	writtenSpriteDefs map[string]bool
	writtenActorDefs  map[string]bool
	writtenActorInsts map[string]bool
}

// DMSpriteDef2 is a declaration of DMSpriteDef2
type DMSpriteDef2 struct {
	Tag          string     // TAG "%s"
	Flags        uint32     // FLAGS %d
	DmTrackTag   string     // DMTRACK "%s"
	Fragment3Ref int32      // ?? FRAGMENT3REF %d
	Fragment4Ref int32      // ?? FRAGMENT4REF %d
	Params2      [3]uint32  // ?? PARAMS2 %d %d %d
	MaxDistance  float32    // ?? MAXDISTANCE %0.7f
	Min          [3]float32 // ?? MIN %0.7f %0.7f %0.7f
	Max          [3]float32 // ?? MAX %0.7f %0.7f %0.7f

	CenterOffset [3]float32 // CENTEROFFSET %0.7f %0.7f %0.7f
	// NUMVERTICES %d
	Vertices [][3]float32 // XYZ %0.7f %0.7f %0.7f
	// NUMUVS %d
	UVs [][2]float32 // UV %0.7f %0.7f
	// NUMVERTEXNORMALS %d
	VertexNormals        [][3]float32 // XYZ %0.7f %0.7f %0.7f
	SkinAssignmentGroups [][2]uint16  // SKINASSIGNMENTGROUPS %d %d
	MaterialPaletteTag   string       // MATERIALPALETTE "%s"
	// NUMFACE2S %d
	Faces []*Face // DMFACE2S
	// NUMMESHOPS %d
	MeshOps              []*MeshOp   // MESHOP
	FaceMaterialGroups   [][2]uint16 // FACEMATERIALGROUPS %d %d
	VertexMaterialGroups [][2]int16  // VERTEXMATERIALGROUPS %d %d
	BoundingRadius       float32     // BOUNDINGRADIUS %0.7f
	FPScale              uint16      // FPScale %d
}

// Ascii returns the ascii representation of a DMSpriteDef2
func (d *DMSpriteDef2) Ascii() string {
	out := "DMSPRITEDEF2\n"
	out += fmt.Sprintf("\tTAG \"%s\"\n", d.Tag)
	if d.Flags != 0 {
		out += fmt.Sprintf("\tFLAGS %d\n", d.Flags)
	}
	out += fmt.Sprintf("\tCENTEROFFSET %0.7f %0.7f %0.7f\n", d.CenterOffset[0], d.CenterOffset[1], d.CenterOffset[2])
	out += "\n"
	out += fmt.Sprintf("\tNUMVERTICES %d\n", len(d.Vertices))
	for _, vert := range d.Vertices {
		out += fmt.Sprintf("\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	out += "\n"
	out += fmt.Sprintf("\tNUMUVS %d\n", len(d.UVs))
	for _, uv := range d.UVs {
		out += fmt.Sprintf("\tUV %0.7f %0.7f\n", uv[0], uv[1])
	}
	out += "\n"
	out += fmt.Sprintf("\tNUMVERTEXNORMALS %d\n", len(d.VertexNormals))
	for _, vn := range d.VertexNormals {
		out += fmt.Sprintf("\tXYZ %0.7f %0.7f %0.7f\n", vn[0], vn[1], vn[2])
	}
	assigments := ""
	for _, sa := range d.SkinAssignmentGroups {
		assigments += fmt.Sprintf("%d %d, ", sa[0], sa[1])
	}
	if len(assigments) > 0 {
		assigments = assigments[:len(assigments)-2]
	}
	out += "\n"
	out += fmt.Sprintf("\tSKINASSIGNMENTGROUPS %s\n", assigments)
	out += fmt.Sprintf("\tMATERIALPALETTE \"%s\"\n", d.MaterialPaletteTag)
	out += "\n"
	out += fmt.Sprintf("\tNUMFACE2S %d\n", len(d.Faces))
	out += "\n"
	for i, face := range d.Faces {
		out += fmt.Sprintf("\tDMFACE2S //%d\n", i+1)
		if face.Flags != 0 {
			out += fmt.Sprintf("\t\tFLAGS %d\n", face.Flags)
		}
		out += fmt.Sprintf("\t\tTRIANGLE   %d, %d, %d\n", face.Triangle[0], face.Triangle[1], face.Triangle[2])
		out += fmt.Sprintf("\tENDFACE //%d\n\n", i+1)
	}
	out += "\n"
	out += "\tNUMMESHOPS 0\n"
	//out += fmt.Sprintf("\tNUMMESHOPS %d\n", len(d.MeshOps))
	/* for _, meshOp := range d.MeshOps {
		out += fmt.Sprintf("\tMESHOP %d %d %0.7f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
		// TODO: figure out MESHOPS
		// MESHOP_VA %d
	} */
	out += "\n"
	groups := ""
	for _, group := range d.FaceMaterialGroups {
		groups += fmt.Sprintf("%d %d, ", group[0], group[1])
	}
	if len(groups) > 0 {
		groups = groups[:len(groups)-2]
	}
	out += fmt.Sprintf("\tFACEMATERIALGROUPS %s\n", groups)
	groups = ""
	for _, group := range d.VertexMaterialGroups {
		groups += fmt.Sprintf("%d %d, ", group[0], group[1])
	}
	if len(groups) > 0 {
		groups = groups[:len(groups)-2]
	}
	out += fmt.Sprintf("\tVERTEXMATERIALGROUPS %s\n", groups)
	out += fmt.Sprintf("\tBOUNDINGRADIUS %0.7f\n", d.BoundingRadius)
	out += "\n"
	out += fmt.Sprintf("\tFPSCALE %d\n", d.FPScale)
	out += "ENDDMSPRITEDEF2\n"
	out += "\n"
	return out
}

type Face struct {
	Flags    uint16    // FLAGS %d
	Triangle [3]uint16 // TRIANGLE %d %d %d
}

type MeshOp struct {
	Index1    uint16  `yaml:"index_1"`
	Index2    uint16  `yaml:"index_2"`
	Offset    float32 `yaml:"offset"`
	Param1    uint8   `yaml:"param_1"`
	TypeField uint8   `yaml:"type_field"`
}

// MaterialPalette is a declaration of MATERIALPALETTE
type MaterialPalette struct {
	Tag   string // TAG "%s"
	Flags uint32 // ?? FLAGS %d
	// NUMMATERIALS %d
	Materials []string // MATERIAL "%s"
}

// Ascii returns the ascii representation of a MaterialPalette
func (m *MaterialPalette) Ascii() string {
	out := "MATERIALPALETTE\n"
	out += fmt.Sprintf("\tTAG \"%s\"\n", m.Tag)
	out += fmt.Sprintf("\tNUMMATERIALS %d\n", len(m.Materials))
	for _, mat := range m.Materials {
		out += fmt.Sprintf("\tMATERIAL \"%s\"\n", mat)
	}
	out += "ENDMATERIALPALETTE\n\n"
	return out
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
}

// Ascii returns the ascii representation of a MaterialDef
func (m *MaterialDef) Ascii() string {
	out := "MATERIALDEFINITION\n"
	out += fmt.Sprintf("\tTAG \"%s\"\n", m.Tag)
	out += fmt.Sprintf("\tFLAGS %d\n", m.Flags)
	out += fmt.Sprintf("\tRENDERMETHOD %s\n", m.RenderMethod)
	out += fmt.Sprintf("\tRGBPEN %d %d %d\n", m.RGBPen[0], m.RGBPen[1], m.RGBPen[2])
	out += fmt.Sprintf("\tBRIGHTNESS %0.7f\n", m.Brightness)
	out += fmt.Sprintf("\tSCALEDAMBIENT %0.7f\n", m.ScaledAmbient)
	if m.SimpleSpriteInstTag != "" {
		out += "\tSIMPLESPRITEINST\n"
		out += fmt.Sprintf("\t\tTAG \"%s\"\n", m.SimpleSpriteInstTag)
		if m.SimpleSpriteInstFlag != 0 {
			out += fmt.Sprintf("\t\tFLAGS %d\n", m.SimpleSpriteInstFlag)
		}
		out += "\tENDSIMPLESPRITEINST\n"
	}

	out += "ENDMATERIALDEFINITION\n\n"
	return out
}

// SimpleSpriteDef is a declaration of SIMPLESPRITEDEF
type SimpleSpriteDef struct {
	Tag string // SIMPLESPRITETAG "%s"
	// NUMFRAMES %d
	BMInfos [][2]string // BMINFO "%s" "%s"
}

// Ascii returns the ascii representation of a SimpleSpriteDef
func (s *SimpleSpriteDef) Ascii() string {
	out := "SIMPLESPRITEDEF\n"
	out += fmt.Sprintf("\tSIMPLESPRITETAG \"%s\"\n", s.Tag)
	out += fmt.Sprintf("\tNUMFRAMES %d\n", len(s.BMInfos))
	for _, bm := range s.BMInfos {
		out += fmt.Sprintf("\tBMINFO \"%s\" \"%s\"\n", bm[0], bm[1])
	}
	out += "ENDSIMPLESPRITEDEF\n\n"
	return out
}

// ActorDef is a declaration of ACTORDEF
type ActorDef struct {
	Tag           string     // ACTORTAG "%s"
	Callback      string     // CALLBACK "%s"
	BoundsRef     int32      // ?? BOUNDSTAG "%s"
	CurrentAction uint32     // ?? CURRENTACTION %d
	Location      [6]float32 // LOCATION %0.7f %0.7f %0.7f %d %d %d
	Unk1          uint32     // ?? UNK1 %d
	// NUMACTIONS %d
	Actions []ActorAction // ACTION
	// NUMFRAGMENTS %d
	FragmentRefs []uint32 // FRAGMENTREF %d
}

// Ascii returns the ascii representation of an ActorDef
func (a *ActorDef) Ascii() string {
	out := "ACTORDEF\n"
	out += fmt.Sprintf("\tACTORTAG \"%s\"\n", a.Tag)
	out += fmt.Sprintf("\tCALLBACK \"%s\"\n", a.Callback)
	//out += fmt.Sprintf("\tBOUNDSREF %d\n", a.BoundsRef)
	//if a.CurrentAction != 0 {
	//	out += fmt.Sprintf("\tCURRENTACTION %d\n", a.CurrentAction)
	//}
	//out += fmt.Sprintf("\tOFFSET %0.7f %0.7f %0.7f\n", a.Offset[0], a.Offset[1], a.Offset[2])
	//out += fmt.Sprintf("\tROTATION %0.7f %0.7f %0.7f\n", a.Rotation[0], a.Rotation[1], a.Rotation[2])
	//if a.Unk1 != 0 {
	//	out += fmt.Sprintf("\tUNK1 %d\n", a.Unk1)
	//}
	out += fmt.Sprintf("\tNUMACTIONS %d\n", len(a.Actions))
	for _, action := range a.Actions {
		out += "\tACTION\n"
		for _, lod := range action.LevelOfDetails {
			//out += fmt.Sprintf("\t\tHIERARCHIALSPRITE \"%s\"\n", lod.HierarchialSpriteDefTag)
			out += fmt.Sprintf("\t\tMINDISTANCE %0.7f\n", lod.MinDistance)
		}
		out += "\tENDACTION\n"
	}
	//out += fmt.Sprintf("\tNUMFRAGMENTS %d\n", len(a.FragmentRefs))
	//for _, frag := range a.FragmentRefs {
	//	out += fmt.Sprintf("\tFRAGMENTREF %d\n", frag)
	//}
	out += "ENDACTORDEF\n\n"
	return out
}

// ActorAction is a declaration of ACTION
type ActorAction struct {
	//NUMLEVELSOFDETAIL %d
	LevelOfDetails []ActorLevelOfDetail // LEVELOFDETAIL
}

// ActorLevelOfDetail is a declaration of LEVELOFDETAIL
type ActorLevelOfDetail struct {
	Unk1        uint32  // ?? HIERARCHIALSPRITE "%s"
	MinDistance float32 // MINDISTANCE %0.7f
}

// ActorInst is a declaration of ACTORINST
type ActorInst struct {
	Tag            string     // ?? ACTORTAG "%s"
	Flags          uint32     // ?? FLAGS %d
	SphereTag      string     // ?? SPHERETAG "%s"
	CurrentAction  uint32     // ?? CURRENTACTION %d
	DefinitionTag  string     // DEFINITION "%s"
	Location       [6]float32 // LOCATION %0.7f %0.7f %0.7f %d %d %d
	Unk1           uint32     // ?? UNK1 %d
	BoundingRadius float32    // BOUNDINGRADIUS %0.7f
	Scale          float32    // SCALEFACTOR %0.7f
	Unk2           int32      // ?? UNK2 %d
}

// Ascii returns the ascii representation of an ActorInst
func (a *ActorInst) Ascii() string {
	out := "ACTORINST\n"
	out += fmt.Sprintf("\tACTORTAG \"%s\"\n", a.Tag)

	if a.Flags&0x20 != 0 {
		out += "\tACTIVE\n"
	}
	out += fmt.Sprintf("\tSPHERETAG \"%s\"\n", a.SphereTag)
	if a.CurrentAction != 0 {
		out += fmt.Sprintf("\tCURRENTACTION %d\n", a.CurrentAction)
	}
	out += fmt.Sprintf("\tDEFINITION \"%s\"\n", a.DefinitionTag)
	out += fmt.Sprintf("\tLOCATION %0.7f %0.7f %0.7f\n", a.Location[0], a.Location[1], a.Location[2])
	if a.Unk1 != 0 {
		out += fmt.Sprintf("\tUNK1 %d\n", a.Unk1)
	}
	out += fmt.Sprintf("\tBOUNDINGRADIUS %0.7f\n", a.BoundingRadius)
	out += fmt.Sprintf("\tSCALEFACTOR %0.7f\n", a.Scale)
	if a.Unk2 != 0 {
		out += fmt.Sprintf("\tUNK2 %d\n", a.Unk2)
	}
	out += "ENDACTORINST\n\n"
	return out
}
