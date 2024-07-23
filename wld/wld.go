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
	LightDefs          []*LightDef
	PointLights        []*PointLight
	Sprite3DDefs       []*Sprite3DDef

	//writing temporary files
	mu                  sync.RWMutex
	writtenPalettes     map[string]bool
	writtenMaterials    map[string]bool
	writtenSpriteDefs   map[string]bool
	writtenActorDefs    map[string]bool
	writtenActorInsts   map[string]bool
	writtenLightDefs    map[string]bool
	writtenPointLights  map[string]bool
	writtenSprite3DDefs map[string]bool
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
	// NUMCOLORS %d
	Colors [][4]uint8 // RGBA %d %d %d %d
	// NUMFACE2S %d
	Faces []*Face // DMFACE2S
	// NUMMESHOPS %d
	MeshOps              []*MeshOp   // MESHOP
	FaceMaterialGroups   [][2]uint16 // FACEMATERIALGROUPS %d %d
	VertexMaterialGroups [][2]int16  // VERTEXMATERIALGROUPS %d %d
	BoundingRadius       float32     // BOUNDINGRADIUS %0.7f
	FPScale              uint16      // FPScale %d
}

func (wld *Wld) reset() {
	wld.writtenMaterials = make(map[string]bool)
	wld.writtenSpriteDefs = make(map[string]bool)
	wld.writtenPalettes = make(map[string]bool)
	wld.writtenActorDefs = make(map[string]bool)
	wld.writtenActorInsts = make(map[string]bool)
	wld.writtenLightDefs = make(map[string]bool)
	wld.writtenPointLights = make(map[string]bool)
	wld.writtenSprite3DDefs = make(map[string]bool)
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
	out += fmt.Sprintf("\t//TODO: NUMMESHOPS %d\n", len(d.MeshOps))
	for _, meshOp := range d.MeshOps {
		out += fmt.Sprintf("\t// TODO: MESHOP %d %d %0.7f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
		// MESHOP_VA %d
	}
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
	Index1    uint16
	Index2    uint16
	Offset    float32
	Param1    uint8
	TypeField uint8
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
		out += fmt.Sprintf("\tFRAME \"%s\" \"%s\"\n", bm[0], bm[1])
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
	out += fmt.Sprintf("\t//TODO: BOUNDSREF %d\n", a.BoundsRef)
	if a.CurrentAction != 0 {
		out += fmt.Sprintf("\t//TODO: CURRENTACTION %d\n", a.CurrentAction)
	}
	out += fmt.Sprintf("\tLOCATION %0.7f %0.7f %0.7f\n", a.Location[0], a.Location[1], a.Location[2])
	out += fmt.Sprintf("\tNUMACTIONS %d\n", len(a.Actions))
	for _, action := range a.Actions {
		out += "\tACTION\n"
		for _, lod := range action.LevelOfDetails {
			out += fmt.Sprintf("\t\t// TODO: UNK1 \"%d\"\n", lod.Unk1)
			out += fmt.Sprintf("\t\tMINDISTANCE %0.7f\n", lod.MinDistance)
		}
		out += "\tENDACTION\n"
	}
	out += fmt.Sprintf("\t// TODO: NUMFRAGMENTS %d\n", len(a.FragmentRefs))
	for _, frag := range a.FragmentRefs {
		out += fmt.Sprintf("\t//TODO: FRAGMENTREF %d\n", frag)
	}
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

// LightDef is a declaration of LIGHTDEF
type LightDef struct {
	Tag             string // TAG "%s"
	Flags           uint32 // ?? FLAGS %d
	FrameCurrentRef uint32 // ?? FRAMECURRENT "%d"
	Sleep           uint32 // SLEEP %d
	// NUMFRAMES %d
	LightLevels []float32 // LIGHTLEVELS %0.7f
	// NUMCOLORS %d
	Colors [][3]float32 // COLORS %0.7f %0.7f %0.7f
}

// Ascii returns the ascii representation of a LightDef
func (l *LightDef) Ascii() string {
	out := "LIGHTDEFINITION\n"
	out += fmt.Sprintf("\tTAG \"%s\"\n", l.Tag)

	if l.Flags&0x01 != 0 {
		out += fmt.Sprintf("\tCURRENT_FRAME \"%d\"\n", l.FrameCurrentRef)
	}
	out += fmt.Sprintf("\tNUMFRAMES %d\n", len(l.LightLevels))
	if l.Flags&0x04 != 0 {
		for _, level := range l.LightLevels {
			out += fmt.Sprintf("\tLIGHTLEVELS %0.6f\n", level)
		}
	}
	if l.Flags&0x02 != 0 {
		out += fmt.Sprintf("\tSLEEP %d\n", l.Sleep)
	}
	if l.Flags&0x08 != 0 {
		out += "\tSKIPFRAMES ON\n"
	}
	if l.Flags&0x10 != 0 {
		for _, color := range l.Colors {
			out += fmt.Sprintf("\tCOLOR %0.6f %0.6f %0.6f\n", color[0], color[1], color[2])
		}
	}
	out += "ENDLIGHTDEFINITION\n\n"
	return out
}

// PointLight is a declaration of POINTLIGHT
type PointLight struct {
	Tag         string     // TAG "%s"
	LightDefTag string     // LIGHT "%s"
	Flags       uint32     // FLAGS %d
	Location    [3]float32 // XYZ %0.7f %0.7f %0.7f
	Radius      float32    // RADIUSOFINFLUENCE %0.7f
}

// Ascii returns the ascii representation of a PointLight
func (p *PointLight) Ascii() string {
	out := "POINTLIGHT\n"
	//out += fmt.Sprintf("\tTAG \"%s\"\n", p.Tag)
	out += fmt.Sprintf("\tXYZ %0.6f %0.6f %0.6f\n", p.Location[0], p.Location[1], p.Location[2])
	out += fmt.Sprintf("\tLIGHT \"%s\"\n", p.LightDefTag)
	if p.Flags != 0 {
		out += fmt.Sprintf("\t//TODO: FLAGS %d\n", p.Flags)
	}
	out += fmt.Sprintf("\tRADIUSOFINFLUENCE %0.7f\n", p.Radius)
	out += "ENDPOINTLIGHT\n\n"
	return out
}

// Sprite3DDef is a declaration of SPRITE3DDEF
type Sprite3DDef struct {
	Tag string // 3DSPRITETAG "%s"
	// NUMVERTICES %d
	Vertices [][3]float32 // XYZ %0.7f %0.7f %0.7f
	// NUMBSPNODES %d
	BSPNodes []*BSPNode // BSPNODE
}

// Ascii returns the ascii representation of a Sprite3DDef
func (s *Sprite3DDef) Ascii() string {
	out := "3DSPRITEDEF\n"
	out += fmt.Sprintf("\t3DSPRITETAG \"%s\"\n", s.Tag)
	out += fmt.Sprintf("\tNUMVERTICES %d\n", len(s.Vertices))
	for _, vert := range s.Vertices {
		out += fmt.Sprintf("\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	out += fmt.Sprintf("\tNUMBSPNODES %d\n", len(s.BSPNodes))
	for i, node := range s.BSPNodes {
		out += fmt.Sprintf("\tBSPNODE //%d\n", i+1)
		out += fmt.Sprintf("\tNUMVERTICES %d\n", len(node.Vertices))
		vertStr := ""
		for _, vert := range node.Vertices {
			vertStr += fmt.Sprintf("%d ", vert)
		}
		if len(vertStr) > 0 {
			vertStr = vertStr[:len(vertStr)-1]
		}
		out += fmt.Sprintf("\tVERTEXLIST %s\n", vertStr)
		out += fmt.Sprintf("\tRENDERMETHOD %s\n", node.RenderMethod)
		out += "\tRENDERINFO\n"
		out += fmt.Sprintf("\t\tPEN %d\n", node.RenderPen)
		out += "\tENDRENDERINFO\n"
		if node.FrontTree != 0 {
			out += fmt.Sprintf("\tFRONTTREE %d\n", node.FrontTree)
		}
		if node.BackTree != 0 {
			out += fmt.Sprintf("\tBACKTREE %d\n", node.BackTree)
		}
		out += "ENDBSPNODE\n"
	}
	out += "END3DSPRITEDEF\n\n"
	return out
}

// BSPNode is a declaration of BSPNODE
type BSPNode struct {
	// NUMVERTICES %d
	Vertices                    []uint32   // VERTEXLIST %d %d %d %d
	RenderMethod                string     // RENDERMETHOD %s
	RenderFlags                 uint8      // FLAGS %d
	RenderPen                   uint32     // PEN %d
	RenderBrightness            float32    // BRIGHTNESS %0.7f
	RenderScaledAmbient         float32    // SCALEDAMBIENT %0.7f
	RenderSimpleSpriteReference uint32     // SIMPLESPRITEINSTREF %d
	RenderUVInfoOrigin          [3]float32 // ORIGIN %0.7f %0.7f %0.7f
	RenderUVInfoUAxis           [3]float32 // UAXIS %0.7f %0.7f %0.7f
	RenderUVInfoVAxis           [3]float32 // VAXIS %0.7f %0.7f %0.7f
	RenderUVMapEntries          []BspNodeUVInfo
	FrontTree                   uint32 // FRONTTREE %d
	BackTree                    uint32 // BACKTREE %d
}

// BspNodeUVInfo is a declaration of UV
type BspNodeUVInfo struct {
	UvOrigin [3]float32 // UV %0.7f %0.7f %0.7f
	UAxis    [3]float32 // UAXIS %0.7f %0.7f %0.7f
	VAxis    [3]float32 // VAXIS %0.7f %0.7f %0.7f
}

// Ascii returns the ascii representation of a BSPNode
func (b *BSPNode) Ascii() string {
	out := "BSPNODE\n"
	out += fmt.Sprintf("\tNUMVERTICES %d\n", len(b.Vertices))
	for _, vert := range b.Vertices {
		out += fmt.Sprintf("\tVERTEXLIST %d\n", vert)
	}
	out += fmt.Sprintf("\tRENDERMETHOD %s\n", b.RenderMethod)
	if b.RenderFlags != 0 {
		out += fmt.Sprintf("\tFLAGS %d\n", b.RenderFlags)
	}
	if b.RenderPen != 0 {
		out += fmt.Sprintf("\tPEN %d\n", b.RenderPen)
	}
	if b.RenderBrightness != 0 {
		out += fmt.Sprintf("\tBRIGHTNESS %0.7f\n", b.RenderBrightness)
	}
	if b.RenderScaledAmbient != 0 {
		out += fmt.Sprintf("\tSCALEDAMBIENT %0.7f\n", b.RenderScaledAmbient)
	}
	if b.RenderSimpleSpriteReference != 0 {
		out += fmt.Sprintf("\tSIMPLESPRITEINSTREF %d\n", b.RenderSimpleSpriteReference)
	}
	if b.RenderUVInfoOrigin != [3]float32{} {
		out += fmt.Sprintf("\tORIGIN %0.7f %0.7f %0.7f\n", b.RenderUVInfoOrigin[0], b.RenderUVInfoOrigin[1], b.RenderUVInfoOrigin[2])
		out += fmt.Sprintf("\tUAXIS %0.7f %0.7f %0.7f\n", b.RenderUVInfoUAxis[0], b.RenderUVInfoUAxis[1], b.RenderUVInfoUAxis[2])
		out += fmt.Sprintf("\tVAXIS %0.7f %0.7f %0.7f\n", b.RenderUVInfoVAxis[0], b.RenderUVInfoVAxis[1], b.RenderUVInfoVAxis[2])
	}
	return out
}

/*
	3DSPRITEDEF
		3DSPRITETAG "merPolyset#22"
		NUMVERTICES 13
		XYZ  0.853552 -0.387270 0.132829
		XYZ  0.853552 -0.387270 -0.0764756
		XYZ  0.669943 -0.387270 -0.0764756
		XYZ  0.669943 -0.387270 0.132829
		XYZ  0.520097 0.625145 0.263453
		XYZ  0.741984 0.789095 -0.104524
		XYZ  0.520097 0.625145 -0.213762
		XYZ  0.956755 0.625145 0.263453
		XYZ  0.956755 0.625145 -0.213762
		XYZ  1.03165 -0.391053 0.272626
		XYZ  1.03165 -0.391053 -0.251556
		XYZ  0.507472 -0.391053 0.272626
		XYZ  0.507472 -0.391053 -0.251556
		NUMBSPNODES 10
		BSPNODE  1
			NUMVERTICES 4
			VERTEXLIST  1, 2, 3, 4
			RENDERMETHOD  SOLIDFILLAMBIENTGOURAUD1
			RENDERINFO
				PEN  11  // RGBPEN  0, 0, 0
				BRIGHTNESS  0.750000
				SCALEDAMBIENT  1.00000
			ENDRENDERINFO
			FRONTTREE  2
			BACKTREE  3
		ENDBSPNODE  1
		BSPNODE  2
			NUMVERTICES 4
			VERTEXLIST  10, 11, 13, 12
			RENDERMETHOD  TEXTURE4AMBIENT
			RENDERINFO
				PEN  156  // RGBPEN  0, 150, 255
				BRIGHTNESS  0.750000
				SCALEDAMBIENT  1.00000
				SIMPLESPRITEINST
					TAG "msleve_SPRITE"
				ENDSIMPLESPRITEINST
				UVORIGIN  0.219274 -0.391053 -0.0851645
				UAXIS  1.0       0.986040        0.00000      -0.569290
				VAXIS  1.0        0.00000        0.00000        0.00000
				UV  0.460789 0.00000
				UV  0.690980 0.00000
				UV  0.292278 0.00000
				UV  0.0620874 0.00000
			ENDRENDERINFO
			FRONTTREE  0
			BACKTREE  0
		ENDBSPNODE  2
		SPHERELIST
		DEFINITION  "merPolyset#22_BOUNDS"
	ENDSPHERELIST
	BOUNDINGRADIUS 1.17286
END3DSPRITEDEF */
