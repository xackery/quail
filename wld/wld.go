// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

import (
	"fmt"
	"io"
	"sync"

	"github.com/xackery/quail/helper"
)

var AsciiVersion = "v0.0.1"

// Wld is a struct representing a Wld file
type Wld struct {
	FileName               string
	GlobalAmbientLight     string
	Version                uint32
	SimpleSpriteDefs       []*SimpleSpriteDef
	MaterialDefs           []*MaterialDef
	MaterialPalettes       []*MaterialPalette
	DMSpriteDef2s          []*DMSpriteDef2
	ActorDefs              []*ActorDef
	ActorInsts             []*ActorInst
	LightDefs              []*LightDef
	PointLights            []*PointLight
	Sprite3DDefs           []*Sprite3DDef
	TrackInstances         []*TrackInstance
	TrackDefs              []*TrackDef
	HierarchicalSpriteDefs []*HierarchicalSpriteDef
	PolyhedronDefs         []*PolyhedronDefinition
	WorldTrees             []*WorldTree
	Regions                []*Region
	AmbientLights          []*AmbientLight
	Zones                  []*Zone

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
	Tag          string
	Flags        uint32
	DmTrackTag   string
	Fragment3Ref int32
	Fragment4Ref int32
	Params2      [3]uint32
	MaxDistance  float32
	Min          [3]float32
	Max          [3]float32

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
	Faces []*Face // DMFACE
	// NUMMESHOPS %d
	MeshOps              []*MeshOp // MESHOP
	FaceMaterialGroups   [][2]uint16
	VertexMaterialGroups [][2]int16 // VERTEXMATERIALGROUPS %d %d
	BoundingRadius       float32    // BOUNDINGRADIUS %0.7f
	FPScale              uint16     // FPScale %d
	PolyhedronTag        string     // POLYHEDRON "%s"
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

func (d *DMSpriteDef2) Definition() string {
	return "DMSPRITEDEF2"
}

func (d *DMSpriteDef2) Write(w io.Writer) error {
	fmt.Fprintf(w, "DMSPRITEDEF2\n")
	fmt.Fprintf(w, "\tTAG \"%s\"\n", d.Tag)
	if d.Flags != 0 {
		fmt.Fprintf(w, "\tFLAGS %d\n", d.Flags)
	}
	fmt.Fprintf(w, "\tCENTEROFFSET %0.7f %0.7f %0.7f\n", d.CenterOffset[0], d.CenterOffset[1], d.CenterOffset[2])
	fmt.Fprintf(w, "\n")
	if len(d.Vertices) > 0 {
		fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(d.Vertices))
		for _, vert := range d.Vertices {
			fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
		}
		fmt.Fprintf(w, "\n")
	}
	if len(d.UVs) > 0 {
		fmt.Fprintf(w, "\tNUMUVS %d\n", len(d.UVs))
		for _, uv := range d.UVs {
			fmt.Fprintf(w, "\tUV %0.7f, %0.7f\n", uv[0], uv[1])
		}
		fmt.Fprintf(w, "\n")
	}
	if len(d.VertexNormals) > 0 {
		fmt.Fprintf(w, "\tNUMVERTEXNORMALS %d\n", len(d.VertexNormals))
		for _, vn := range d.VertexNormals {
			fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vn[0], vn[1], vn[2])
		}
		fmt.Fprintf(w, "\n")
	}
	if len(d.SkinAssignmentGroups) > 0 {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\tSKINASSIGNMENTGROUPS %d", len(d.SkinAssignmentGroups))
		for i, sa := range d.SkinAssignmentGroups {
			endComma := ","
			if i == len(d.SkinAssignmentGroups)-1 {
				endComma = ""
			}
			fmt.Fprintf(w, " %d, %d%s", sa[0], sa[1], endComma)
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "\tMATERIALPALETTE \"%s\"\n", d.MaterialPaletteTag)
	fmt.Fprintf(w, "\n")
	if d.PolyhedronTag != "" {
		fmt.Fprintf(w, "\tPOLYHEDRON\n")
		fmt.Fprintf(w, "\t\tDEFINITION \"%s\"\n", d.PolyhedronTag)
		fmt.Fprintf(w, "\tENDPOLYHEDRON\n\n")
	}
	if len(d.Faces) > 0 {
		fmt.Fprintf(w, "\tNUMFACE2S %d\n", len(d.Faces))
		fmt.Fprintf(w, "\n")
		for i, face := range d.Faces {
			fmt.Fprintf(w, "\tDMFACE2 //%d\n", i+1)
			if face.Flags != 0 {
				fmt.Fprintf(w, "\t\tFLAGS %d\n", face.Flags)
			}
			fmt.Fprintf(w, "\t\tTRIANGLE   %d, %d, %d\n", face.Triangle[0], face.Triangle[1], face.Triangle[2])
			fmt.Fprintf(w, "\tENDDMFACE2 //%d\n\n", i+1)
		}
		fmt.Fprintf(w, "\n")
	}
	if len(d.MeshOps) > 0 {
		//fmt.Fprintf(w, "\tNUMMESHOPS 0\n")
		fmt.Fprintf(w, "\t//TODO: NUMMESHOPS %d\n", len(d.MeshOps))
		for _, meshOp := range d.MeshOps {
			fmt.Fprintf(w, "\t// TODO: MESHOP %d %d %0.7f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
			// MESHOP_VA %d
		}
		fmt.Fprintf(w, "\n")
	}
	if len(d.FaceMaterialGroups) > 0 {
		fmt.Fprintf(w, "\tFACEMATERIALGROUPS %d", len(d.FaceMaterialGroups))
		for _, group := range d.FaceMaterialGroups {
			endComma := ","
			if group == d.FaceMaterialGroups[len(d.FaceMaterialGroups)-1] {
				endComma = ""
			}
			fmt.Fprintf(w, " %d, %d%s", group[0], group[1], endComma)
		}
		fmt.Fprintf(w, "\n")
	}
	if len(d.VertexMaterialGroups) > 0 {
		fmt.Fprintf(w, "\tVERTEXMATERIALGROUPS %d", len(d.VertexMaterialGroups))
		for _, group := range d.VertexMaterialGroups {
			endComma := ","
			if group == d.VertexMaterialGroups[len(d.VertexMaterialGroups)-1] {
				endComma = ""
			}
			fmt.Fprintf(w, " %d, %d%s", group[0], group[1], endComma)
		}
		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7e\n", d.BoundingRadius)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFPSCALE %d\n", d.FPScale)
	fmt.Fprintf(w, "ENDDMSPRITEDEF2\n\n")
	return nil
}

func (d *DMSpriteDef2) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG")
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	d.Tag = records[1]

	records, err = r.ReadProperty("CENTEROFFSET")
	if err != nil {
		return fmt.Errorf("CENTEROFFSET: %w", err)
	}
	d.CenterOffset[0], err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("center offset x: %w", err)
	}
	d.CenterOffset[1], err = helper.ParseFloat32(records[2])
	if err != nil {
		return fmt.Errorf("center offset y: %w", err)
	}
	d.CenterOffset[2], err = helper.ParseFloat32(records[3])
	if err != nil {
		return fmt.Errorf("center offset z: %w", err)
	}

	records, err = r.ReadProperty("MATERIALPALETTE")
	if err != nil {
		return fmt.Errorf("MATERIALPALETTE: %w", err)
	}
	d.MaterialPaletteTag = records[1]

	records, err = r.ReadProperty("")
	if err != nil {
		return fmt.Errorf("POLYHEDRON: %w", err)
	}
	if records[0] == "POLYHEDRON" {
		records, err = r.ReadProperty("DEFINITION")
		if err != nil {
			return fmt.Errorf("POLYHEDRON DEFINITION: %w", err)
		}
		d.PolyhedronTag = records[1]
		records, err = r.ReadProperty("ENDPOLYHEDRON")
		if err != nil {
			return fmt.Errorf("ENDPOLYHEDRON: %w", err)
		}
	}

	records, err = r.ReadProperty("BOUNDINGRADIUS")
	if err != nil {
		return fmt.Errorf("BOUNDINGRADIUS: %w", err)
	}
	d.BoundingRadius, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("FPSCALE")
	if err != nil {
		return fmt.Errorf("FPSCALE: %w", err)
	}
	d.FPScale, err = helper.ParseUint16(records[1])
	if err != nil {
		return fmt.Errorf("fpscale: %w", err)
	}

	records, err = r.ReadProperty("ENDDMSPRITEDEF2")
	if err != nil {
		return fmt.Errorf("ENDDMSPRITEDEF2: %w", err)
	}

	return nil
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

// MaterialPalette is a declaration of MATERIALPALETTE
type MaterialPalette struct {
	Tag          string // TAG "%s"
	numMaterials int    // NUMMATERIALS %d
	flags        uint32
	Materials    []string // MATERIAL "%s"
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
}

func (m *MaterialDef) Definition() string {
	return "MATERIALDEFINITION"
}

func (m *MaterialDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", m.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", m.Tag)
	fmt.Fprintf(w, "\tFLAGS %d\n", m.Flags)
	fmt.Fprintf(w, "\tRENDERMETHOD %s\n", m.RenderMethod)
	fmt.Fprintf(w, "\tRGBPEN %d %d %d\n", m.RGBPen[0], m.RGBPen[1], m.RGBPen[2])
	fmt.Fprintf(w, "\tBRIGHTNESS %0.7f\n", m.Brightness)
	fmt.Fprintf(w, "\tSCALEDAMBIENT %0.7f\n", m.ScaledAmbient)
	if m.SimpleSpriteInstTag != "" {
		fmt.Fprintf(w, "\tSIMPLESPRITEINST\n")
		fmt.Fprintf(w, "\t\tTAG \"%s\"\n", m.SimpleSpriteInstTag)
		if m.SimpleSpriteInstFlag != 0 {
			fmt.Fprintf(w, "\t\tFLAGS %d\n", m.SimpleSpriteInstFlag)
		}
		fmt.Fprintf(w, "\tENDSIMPLESPRITEINST\n")
	}
	fmt.Fprintf(w, "ENDMATERIALDEFINITION\n\n")
	return nil
}

func (m *MaterialDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG")
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	m.Tag = records[1]

	records, err = r.ReadProperty("FLAGS")
	if err != nil {
		return fmt.Errorf("FLAGS: %w", err)
	}
	m.Flags, err = helper.ParseUint32(records[1])
	if err != nil {
		return fmt.Errorf("flags: %w", err)
	}

	records, err = r.ReadProperty("RENDERMETHOD")
	if err != nil {
		return fmt.Errorf("RENDERMETHOD: %w", err)
	}
	m.RenderMethod = records[1]

	records, err = r.ReadProperty("RGBPEN")
	if err != nil {
		return fmt.Errorf("RGBPEN: %w", err)
	}
	m.RGBPen[0], err = helper.ParseUint8(records[1])
	if err != nil {
		return fmt.Errorf("rgbpen r: %w", err)
	}

	m.RGBPen[1], err = helper.ParseUint8(records[2])
	if err != nil {
		return fmt.Errorf("rgbpen g: %w", err)
	}

	m.RGBPen[2], err = helper.ParseUint8(records[3])
	if err != nil {
		return fmt.Errorf("rgbpen b: %w", err)
	}

	records, err = r.ReadProperty("BRIGHTNESS")
	if err != nil {
		return fmt.Errorf("BRIGHTNESS: %w", err)
	}
	m.Brightness, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("brightness: %w", err)
	}

	records, err = r.ReadProperty("SCALEDAMBIENT")
	if err != nil {
		return fmt.Errorf("SCALEDAMBIENT: %w", err)
	}
	m.ScaledAmbient, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("scaled ambient: %w", err)
	}

	records, err = r.ReadProperty("")
	if err != nil {
		return fmt.Errorf("SIMPLESPRITEINST: %w", err)
	}
	if records[0] == "SIMPLESPRITEINST" {
		records, err = r.ReadProperty("TAG")
		if err != nil {
			return fmt.Errorf("SIMPLESPRITEINST TAG: %w", err)
		}
		m.SimpleSpriteInstTag = records[1]
	}

	_, err = r.ReadProperty("ENDMATERIALDEFINITION")
	if err != nil {
		return fmt.Errorf("ENDMATERIALDEFINITION: %w", err)
	}
	return nil
}

// SimpleSpriteDef is a declaration of SIMPLESPRITEDEF
type SimpleSpriteDef struct {
	Tag string // SIMPLESPRITETAG "%s"
	// NUMFRAMES %d
	BMInfos [][2]string // BMINFO "%s" "%s"
}

func (s *SimpleSpriteDef) Definition() string {
	return "SIMPLESPRITEDEF"
}

func (s *SimpleSpriteDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", s.Definition())
	fmt.Fprintf(w, "\tSIMPLESPRITETAG \"%s\"\n", s.Tag)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(s.BMInfos))
	for _, bm := range s.BMInfos {
		fmt.Fprintf(w, "\tFRAME \"%s\" \"%s\"\n", bm[0], bm[1])
	}
	if len(s.BMInfos) > 0 {
		fmt.Fprintf(w, "\tBMINFO \"%s\" \"%s\"\n", s.BMInfos[0][0], s.BMInfos[0][1])
	}
	fmt.Fprintf(w, "ENDSIMPLESPRITEDEF\n\n")
	return nil
}

func (s *SimpleSpriteDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("SIMPLESPRITETAG")
	if err != nil {
		return fmt.Errorf("SIMPLESPRITETAG: %w", err)
	}
	s.Tag = records[1]

	records, err = r.ReadProperty("NUMFRAMES")
	if err != nil {
		return fmt.Errorf("NUMFRAMES: %w", err)
	}
	numFrames, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = r.ReadProperty("FRAME")
		if err != nil {
			return fmt.Errorf("FRAME: %w", err)
		}
		s.BMInfos = append(s.BMInfos, [2]string{records[1], records[2]})
	}

	_, err = r.ReadProperty("ENDSIMPLESPRITEDEF")
	if err != nil {
		return fmt.Errorf("ENDSIMPLESPRITEDEF: %w", err)
	}
	return nil
}

// ActorDef is a declaration of ACTORDEF
type ActorDef struct {
	Tag           string // ACTORTAG "%s"
	Callback      string // CALLBACK "%s"
	BoundsRef     int32  // ?? BOUNDSTAG "%s"
	CurrentAction uint32 // ?? CURRENTACTION %d
	Location      [6]float32
	Unk1          uint32 // ?? UNK1 %d
	// NUMACTIONS %d
	Actions []ActorAction // ACTION
	// NUMFRAGMENTS %d
	FragmentRefs []uint32 // FRAGMENTREF %d
}

func (a *ActorDef) Definition() string {
	return "ACTORDEF"
}

func (a *ActorDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", a.Definition())
	fmt.Fprintf(w, "\tACTORTAG \"%s\"\n", a.Tag)
	fmt.Fprintf(w, "\tCALLBACK \"%s\"\n", a.Callback)
	fmt.Fprintf(w, "\t//TODO: BOUNDSREF %d\n", a.BoundsRef)
	if a.CurrentAction != 0 {
		fmt.Fprintf(w, "\t//TODO: CURRENTACTION %d\n", a.CurrentAction)
	}
	fmt.Fprintf(w, "\tLOCATION %0.7f %0.7f %0.7f\n", a.Location[0], a.Location[1], a.Location[2])
	fmt.Fprintf(w, "\tNUMACTIONS %d\n", len(a.Actions))
	for _, action := range a.Actions {
		fmt.Fprintf(w, "\tACTION\n")
		fmt.Fprintf(w, "\t\tNUMLEVELSOFDETAIL %d\n", len(action.LevelOfDetails))
		for _, lod := range action.LevelOfDetails {
			fmt.Fprintf(w, "\t\tLEVELOFDETAIL\n")
			fmt.Fprintf(w, "\t\t\tHIERRACHICALSPRITE \"%s\"\n", lod.HierarchicalSpriteTag)
			fmt.Fprintf(w, "\t\t\tMINDISTANCE %0.7f\n", lod.MinDistance)
			fmt.Fprintf(w, "\t\tENDLEVELOFDETAIL\n")
		}
		fmt.Fprintf(w, "\tENDACTION\n")
	}
	fmt.Fprintf(w, "\t// TODO: NUMFRAGMENTS %d\n", len(a.FragmentRefs))
	for _, frag := range a.FragmentRefs {
		fmt.Fprintf(w, "\t//TODO: FRAGMENTREF %d\n", frag)
	}
	fmt.Fprintf(w, "ENDACTORDEF\n\n")
	return nil
}

func (a *ActorDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("ACTORTAG")
	if err != nil {
		return fmt.Errorf("ACTORTAG: %w", err)
	}
	a.Tag = records[1]

	records, err = r.ReadProperty("CALLBACK")
	if err != nil {
		return fmt.Errorf("CALLBACK: %w", err)
	}
	a.Callback = records[1]

	records, err = r.ReadProperty("")
	if err != nil {
		return fmt.Errorf("after CALLBACK: %w", err)
	}
	switch records[0] {
	case "LOCATION":
		a.Location[0], err = helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("location x: %w", err)
		}
		a.Location[1], err = helper.ParseFloat32(records[2])
		if err != nil {
			return fmt.Errorf("location y: %w", err)
		}
		a.Location[2], err = helper.ParseFloat32(records[3])
		if err != nil {
			return fmt.Errorf("location z: %w", err)
		}
	case "CURRENTACTION":
		a.CurrentAction, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("current action: %w", err)
		}

	case "NUMACTIONS": // leak to next section
	default:
		return fmt.Errorf("unknown property: %s", records[0])
	}

	numActions, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num actions: %w", err)
	}

	for i := 0; i < numActions; i++ {
		action := ActorAction{}
		_, err = r.ReadProperty("ACTION")
		if err != nil {
			return fmt.Errorf("ACTION: %w", err)
		}
		records, err = r.ReadProperty("NUMLEVELSOFDETAIL")
		if err != nil {
			return fmt.Errorf("NUMLEVELSOFDETAIL: %w", err)
		}
		numLod, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num lod: %w", err)
		}

		for j := 0; j < numLod; j++ {
			lod := ActorLevelOfDetail{}
			_, err = r.ReadProperty("LEVELOFDETAIL")
			if err != nil {
				return fmt.Errorf("LEVELOFDETAIL: %w", err)
			}

		lodLoop:
			for {
				records, err = r.ReadProperty("")
				if err != nil {
					return fmt.Errorf("after LEVELOFDETAIL: %w", err)
				}
				switch records[0] {
				case "MINDISTANCE":
					lod.MinDistance, err = helper.ParseFloat32(records[1])
					if err != nil {
						return fmt.Errorf("min distance: %w", err)
					}
				case "3DSPRITE":
					lod.HierarchicalSpriteTag = records[1]
				case "HIERRACHICALSPRITE":
					lod.HierarchicalSpriteTag = records[1]
				case "ENDLEVELOFDETAIL":
					break lodLoop
				default:
					return fmt.Errorf("unknown property inside LEVELOFDETAIL: %s", records[0])
				}
				action.LevelOfDetails = append(action.LevelOfDetails, lod)
			}
		}
		_, err = r.ReadProperty("ENDACTION")
		if err != nil {
			return fmt.Errorf("ENDACTION: %w", err)
		}

		a.Actions = append(a.Actions, action)
	}

	_, err = r.ReadProperty("ENDACTORDEF")
	if err != nil {
		return fmt.Errorf("ENDACTORDEF: %w", err)
	}

	return nil
}

// ActorAction is a declaration of ACTION
type ActorAction struct {
	//NUMLEVELSOFDETAIL %d
	LevelOfDetails []ActorLevelOfDetail // LEVELOFDETAIL
}

// ActorLevelOfDetail is a declaration of LEVELOFDETAIL
type ActorLevelOfDetail struct {
	HierarchicalSpriteTag string
	MinDistance           float32
}

// ActorInst is a declaration of ACTORINST
type ActorInst struct {
	Tag            string
	Flags          uint32
	SphereTag      string
	CurrentAction  uint32
	DefinitionTag  string
	Location       [6]float32
	Unk1           uint32
	BoundingRadius float32
	Scale          float32
	Unk2           int32
}

func (a *ActorInst) Definition() string {
	return "ACTORINST"
}

func (a *ActorInst) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", a.Definition())
	fmt.Fprintf(w, "\tACTORTAG \"%s\"\n", a.Tag)
	if a.Flags&0x20 != 0 {
		fmt.Fprintf(w, "\tACTIVE\n")
	}
	fmt.Fprintf(w, "\tSPHERETAG \"%s\"\n", a.SphereTag)
	if a.CurrentAction != 0 {
		fmt.Fprintf(w, "\tCURRENTACTION %d\n", a.CurrentAction)
	}
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", a.DefinitionTag)
	fmt.Fprintf(w, "\tLOCATION %0.7f %0.7f %0.7f\n", a.Location[0], a.Location[1], a.Location[2])
	if a.Unk1 != 0 {
		fmt.Fprintf(w, "\tUNK1 %d\n", a.Unk1)
	}
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", a.BoundingRadius)
	fmt.Fprintf(w, "\tSCALEFACTOR %0.7f\n", a.Scale)
	if a.Unk2 != 0 {
		fmt.Fprintf(w, "\tUNK2 %d\n", a.Unk2)
	}
	fmt.Fprintf(w, "ENDACTORINST\n\n")
	return nil
}

func (a *ActorInst) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("DEFINITION")
	if err != nil {
		return fmt.Errorf("DEFINITION: %w", err)
	}
	a.DefinitionTag = records[1]

	records, err = r.ReadProperty("")
	if err != nil {
		return fmt.Errorf("after DEFINITION: %w", err)
	}
	if records[0] == "ACTIVE" {
		a.Flags |= 0x20
		records, err = r.ReadProperty("LOCATION")
		if err != nil {
			return fmt.Errorf("LOCATION: %w", err)
		}
	}
	if records[0] != "LOCATION" {
		return fmt.Errorf("expected LOCATION, got %s", records[0])
	}
	a.Location[0], err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("location x: %w", err)
	}
	a.Location[1], err = helper.ParseFloat32(records[2])
	if err != nil {
		return fmt.Errorf("location y: %w", err)
	}
	a.Location[2], err = helper.ParseFloat32(records[3])
	if err != nil {
		return fmt.Errorf("location z: %w", err)
	}

	records, err = r.ReadProperty("SPHERE")
	if err != nil {
		return fmt.Errorf("SPHERE: %w", err)
	}
	a.SphereTag = records[1]

	records, err = r.ReadProperty("BOUNDINGRADIUS")
	if err != nil {
		return fmt.Errorf("BOUNDINGRADIUS: %w", err)
	}
	a.BoundingRadius, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("bounding radius: %w", err)
	}

	records, err = r.ReadProperty("SCALEFACTOR")
	if err != nil {
		return fmt.Errorf("SCALEFACTOR: %w", err)
	}
	a.Scale, err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("scale factor: %w", err)
	}

	_, err = r.ReadProperty("ENDACTORINST")
	if err != nil {
		return fmt.Errorf("ENDACTORINST: %w", err)
	}
	return nil
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

func (l *LightDef) Definition() string {
	return "LIGHTDEFINITION"
}

func (l *LightDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", l.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", l.Tag)
	if l.Flags&0x01 != 0 {
		fmt.Fprintf(w, "\tCURRENT_FRAME \"%d\"\n", l.FrameCurrentRef)
	}
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", len(l.LightLevels))
	if l.Flags&0x04 != 0 {
		for _, level := range l.LightLevels {
			fmt.Fprintf(w, "\tLIGHTLEVELS %0.6f\n", level)
		}
	}
	if l.Flags&0x02 != 0 {
		fmt.Fprintf(w, "\tSLEEP %d\n", l.Sleep)
	}
	if l.Flags&0x08 != 0 {
		fmt.Fprintf(w, "\tSKIPFRAMES ON\n")
	}
	if l.Flags&0x10 != 0 {
		fmt.Fprintf(w, "\tNUMCOLORS %d\n", len(l.Colors))
		for _, color := range l.Colors {
			fmt.Fprintf(w, "\tCOLOR %0.6f %0.6f %0.6f\n", color[0], color[1], color[2])
		}
	}
	fmt.Fprintf(w, "ENDLIGHTDEFINITION\n\n")
	return nil
}

func (l *LightDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG")
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	l.Tag = records[1]

	records, err = r.ReadProperty("NUMFRAMES")
	if err != nil {
		return fmt.Errorf("NUMFRAMES: %w", err)
	}
	numFrames, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num frames: %w", err)
	}

	for i := 0; i < numFrames; i++ {
		records, err = r.ReadProperty("LIGHTLEVELS")
		if err != nil {
			return fmt.Errorf("LIGHTLEVELS: %w", err)
		}
		level, err := helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("level: %w", err)
		}
		l.LightLevels = append(l.LightLevels, level)
	}

	_, err = r.ReadProperty("ENDLIGHTDEFINITION")
	if err != nil {
		return fmt.Errorf("ENDLIGHTDEFINITION: %w", err)
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
	fmt.Fprintf(w, "\tXYZ %0.6f %0.6f %0.6f\n", p.Location[0], p.Location[1], p.Location[2])
	fmt.Fprintf(w, "\tLIGHT \"%s\"\n", p.LightDefTag)
	if p.Flags != 0 {
		fmt.Fprintf(w, "\tFLAGS %d\n", p.Flags)
	}
	fmt.Fprintf(w, "\tRADIUSOFINFLUENCE %0.7f\n", p.Radius)
	fmt.Fprintf(w, "ENDPOINTLIGHT\n\n")
	return nil
}

func (p *PointLight) Read(r *AsciiReadToken) error {
	return nil
}

// Sprite3DDef is a declaration of SPRITE3DDEF
type Sprite3DDef struct {
	Tag      string
	Vertices [][3]float32
	BSPNodes []*BSPNode
}

func (s *Sprite3DDef) Definition() string {
	return "3DSPRITEDEF"
}

func (s *Sprite3DDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", s.Definition())
	fmt.Fprintf(w, "\t3DSPRITETAG \"%s\"\n", s.Tag)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(s.Vertices))
	for _, vert := range s.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMBSPNODES %d\n", len(s.BSPNodes))
	for i, node := range s.BSPNodes {
		fmt.Fprintf(w, "\tBSPNODE //%d\n", i+1)
		fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(node.Vertices))
		vertStr := ""
		for _, vert := range node.Vertices {
			vertStr += fmt.Sprintf("%d ", vert)
		}
		if len(vertStr) > 0 {
			vertStr = vertStr[:len(vertStr)-1]
		}
		fmt.Fprintf(w, "\tVERTEXLIST %s\n", vertStr)
		fmt.Fprintf(w, "\tRENDERMETHOD %s\n", node.RenderMethod)
		fmt.Fprintf(w, "\tRENDERINFO\n")
		fmt.Fprintf(w, "\t\tPEN %d\n", node.RenderPen)
		fmt.Fprintf(w, "\tENDRENDERINFO\n")
		if node.FrontTree != 0 {
			fmt.Fprintf(w, "\tFRONTTREE %d\n", node.FrontTree)
		}
		if node.BackTree != 0 {
			fmt.Fprintf(w, "\tBACKTREE %d\n", node.BackTree)
		}
		fmt.Fprintf(w, "ENDBSPNODE\n")
	}
	fmt.Fprintf(w, "END3DSPRITEDEF\n\n")
	return nil
}

func (s *Sprite3DDef) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("3DSPRITETAG")
	if err != nil {
		return fmt.Errorf("3DSPRITETAG: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("3DSPRITETAG: missing tag name")
	}

	s.Tag = records[1]

	records, err = r.ReadProperty("NUMVERTICES")
	if err != nil {
		return fmt.Errorf("NUMVERTICES: %w", err)
	}
	numVertices, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = r.ReadProperty("XYZ")
		if err != nil {
			return fmt.Errorf("XYZ: %w", err)
		}
		if len(records) < 4 {
			return fmt.Errorf("XYZ: missing coordinates")
		}
		var vert [3]float32
		vert[0], err = helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("vertex x: %w", err)
		}
		vert[1], err = helper.ParseFloat32(records[2])
		if err != nil {
			return fmt.Errorf("vertex y: %w", err)
		}
		vert[2], err = helper.ParseFloat32(records[3])
		if err != nil {
			return fmt.Errorf("vertex z: %w", err)
		}
		s.Vertices = append(s.Vertices, vert)
	}

	records, err = r.ReadProperty("NUMBSPNODES")
	if err != nil {
		return fmt.Errorf("NUMBSPNODES: %w", err)
	}
	numNodes, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num nodes: %w", err)
	}

	for i := 0; i < numNodes; i++ {
		node := &BSPNode{}
		_, err = r.ReadProperty("BSPNODE")
		if err != nil {
			return fmt.Errorf("BSPNODE: %w", err)
		}

		records, err = r.ReadProperty("NUMVERTICES")
		if err != nil {
			return fmt.Errorf("NUMVERTICES: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("NUMVERTICES: missing number of vertices")
		}
		numVertices, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("num vertices: %w", err)
		}

		records, err = r.ReadProperty("VERTEXLIST")
		if err != nil {
			return fmt.Errorf("VERTEXLIST: %w", err)
		}
		if len(records) < numVertices+1 {
			return fmt.Errorf("VERTEXLIST: expected %d vertices, got %d", numVertices, len(records)-1)
		}

		node.Vertices = make([]uint32, numVertices)
		for j := 0; j < numVertices; j++ {
			node.Vertices[j], err = helper.ParseUint32(records[j+1])
			if err != nil {
				return fmt.Errorf("vertex %d: %w", j, err)
			}
		}

		records, err = r.ReadProperty("RENDERMETHOD")
		if err != nil {
			return fmt.Errorf("RENDERMETHOD: %w", err)
		}
		node.RenderMethod = records[1]

		_, err = r.ReadProperty("RENDERINFO")
		if err != nil {
			return fmt.Errorf("RENDERINFO: %w", err)
		}

		records, err = r.ReadProperty("PEN")
		if err != nil {
			return fmt.Errorf("PEN: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("PEN: missing pen value")
		}
		node.RenderPen, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("pen: %w", err)
		}

		_, err = r.ReadProperty("ENDRENDERINFO")
		if err != nil {
			return fmt.Errorf("ENDRENDERINFO: %w", err)
		}

		records, err = r.ReadProperty("FRONTTREE")
		if err != nil {
			return fmt.Errorf("FRONTTREE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("FRONTTREE: missing tree value")
		}
		node.FrontTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = r.ReadProperty("BACKTREE")
		if err != nil {
			return fmt.Errorf("BACKTREE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("BACKTREE: missing tree value")
		}
		node.BackTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = r.ReadProperty("ENDBSPNODE")
		if err != nil {
			return fmt.Errorf("ENDBSPNODE: %w", err)
		}

		s.BSPNodes = append(s.BSPNodes, node)
	}

	_, err = r.ReadProperty("END3DSPRITEDEF")
	if err != nil {
		return fmt.Errorf("END3DSPRITEDEF: %w", err)
	}

	return nil
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

type PolyhedronDefinition struct {
	Tag            string
	BoundingRadius float32
	ScaleFactor    float32
	numVertices    int // NUMVERTICES %d
	Vertices       [][3]float32
	numFaces       int // NUMFACES %d
	Faces          []*PolyhedronDefinitionFace
}

type PolyhedronDefinitionFace struct {
	numVertices int // NUMVERTICES %d
	Vertices    []uint32
}

func (p *PolyhedronDefinition) Definition() string {
	return "POLYHEDRONDEFINITION"
}

func (p *PolyhedronDefinition) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", p.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", p.Tag)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", p.BoundingRadius)
	fmt.Fprintf(w, "\tSCALEFACTOR %0.7f\n", p.ScaleFactor)
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(p.Vertices))
	for _, vert := range p.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7e %0.7e %0.7e\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\tNUMFACES %d\n", len(p.Faces))
	for i, face := range p.Faces {
		fmt.Fprintf(w, "\tFACE %d\n", i+1)
		fmt.Fprintf(w, "\t\tNUMVERTICES %d\n", len(face.Vertices))
		vertStr := ""
		for _, vert := range face.Vertices {
			vertStr += fmt.Sprintf("%d, ", vert)
		}
		if len(vertStr) > 0 {
			vertStr = vertStr[:len(vertStr)-2]
		}
		fmt.Fprintf(w, "\t\tVERTEXLIST %s\n", vertStr)
		fmt.Fprintf(w, "\tENDFACE %d\n", i+1)
	}
	fmt.Fprintf(w, "ENDPOLYHEDRONDEFINITION\n\n")
	return nil
}

func (p *PolyhedronDefinition) Read(r *AsciiReadToken) error {
	return nil
}

type TrackInstance struct {
	Tag            string // TAG "%s"
	DefiniationTag string // DEFINITION "%s"
	isInterpolated bool   // INTERPOLATE
	Sleep          uint32 // SLEEP %d
}

func (t *TrackInstance) Definition() string {
	return "TRACKINSTANCE"
}

func (t *TrackInstance) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", t.Tag)
	fmt.Fprintf(w, "\tDEFINITION \"%s\"\n", t.DefiniationTag)
	if t.isInterpolated {
		fmt.Fprintf(w, "\tINTERPOLATE\n")
	}
	fmt.Fprintf(w, "\tSLEEP %d\n", t.Sleep)
	fmt.Fprintf(w, "ENDTRACKINSTANCE\n\n")
	return nil
}

func (t *TrackInstance) Read(r *AsciiReadToken) error {
	return nil
}

type TrackDef struct {
	Tag            string                // TAG "%s"
	numFrames      int                   // NUMFRAMES %d
	FrameTransform []TrackFrameTransform // FRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f
}

type TrackFrameTransform struct {
	LocDenom float32
	Rotation [3]int32
	Position [3]float32
}

func (t *TrackDef) Definition() string {
	return "TRACKDEFINITION"
}

func (t *TrackDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", t.Tag)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", t.numFrames)
	for _, frame := range t.FrameTransform {
		fmt.Fprintf(w, "\tFRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f\n", frame.LocDenom, frame.Rotation[0], frame.Rotation[1], frame.Rotation[2], frame.Position[0], frame.Position[1], frame.Position[2])
	}
	fmt.Fprintf(w, "ENDTRACKDEFINITION\n\n")
	return nil
}

func (t *TrackDef) Read(r *AsciiReadToken) error {

	return nil
}

type HierarchicalSpriteDef struct {
	Tag            string         // TAG "%s"
	Dags           []Dag          // DAG
	AttachedSkins  []AttachedSkin // ATTACHEDSKIN
	CenterOffset   [3]float32     // CENTEROFFSET %0.7f %0.7f %0.7f
	BoundingRadius float32        // BOUNDINGRADIUS %0.7f
	HasCollisions  bool           // DAGCOLLISIONS
}

type Dag struct {
	Tag     string   // TAG "%s"
	Flags   uint32   // NULLSPRITE, etc
	Track   string   // TRACK "%s"
	SubDags []uint32 // SUBDAGLIST %d %d
}

type AttachedSkin struct {
	Tag                       string // TAG "%s"
	LinkSkinUpdatesToDagIndex uint32 // LINKSKINUPDATES %d
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
		if dag.Flags != 0 {
			fmt.Fprintf(w, "\t\tFLAGS %d\n", dag.Flags)
		}
		if dag.Track != "" {
			fmt.Fprintf(w, "\t\tTRACK \"%s\"\n", dag.Track)
		}
		if len(dag.SubDags) > 0 {
			fmt.Fprintf(w, "\t\tNUMSUBDAGS %d\n", len(dag.SubDags))
			fmt.Fprintf(w, "\t\tSUBDAGLIST")
			for _, subDag := range dag.SubDags {
				fmt.Fprintf(w, " %d", subDag)
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "\tENDDAG // %d\n", i+1)
	}
	if len(h.Dags) > 0 {
		fmt.Fprintf(w, "\n")
	}
	if len(h.AttachedSkins) > 0 {
		fmt.Fprintf(w, "\tNUMATTACHEDSKINS %d\n", len(h.AttachedSkins))
		for _, skin := range h.AttachedSkins {
			fmt.Fprintf(w, "\tDMSPRITE \"%s\"\n", skin.Tag)
			fmt.Fprintf(w, "\tLINKSKINUPDATESTODAGINDEX %d\n", skin.LinkSkinUpdatesToDagIndex)
		}
		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintf(w, "\tCENTEROFFSET %0.1f %0.1f %0.1f\n", h.CenterOffset[0], h.CenterOffset[1], h.CenterOffset[2])
	if h.HasCollisions {
		fmt.Fprintf(w, "\tDAGCOLLISIONS\n")
	}
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7e\n", h.BoundingRadius)

	fmt.Fprintf(w, "ENDHIERARCHICALSPRITEDEF\n\n")
	return nil
}

func (h *HierarchicalSpriteDef) Read(r *AsciiReadToken) error {

	return nil
}

type WorldTree struct {
	WorldNodes []*WorldNode
}

type WorldNode struct {
	Normals        [4]float32
	WorldRegionTag string
	FrontTree      uint32
	BackTree       uint32
}

func (t *WorldTree) Definition() string {
	return "WORLDTREE"
}

func (t *WorldTree) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	for i, node := range t.WorldNodes {
		fmt.Fprintf(w, "\tWORLDNODE %d\n", i+1)
		fmt.Fprintf(w, "\t\tNORMALS %0.7f %0.7f %0.7f %0.7f\n", node.Normals[0], node.Normals[1], node.Normals[2], node.Normals[3])
		fmt.Fprintf(w, "\t\tWORLDREGIONTAG \"%s\"\n", node.WorldRegionTag)
		fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
		fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
		fmt.Fprintf(w, "\tENDWORLDNODE %d\n", i+1)
	}
	fmt.Fprintf(w, "ENDWORLDTREE\n\n")
	return nil
}

func (t *WorldTree) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("NUMWORLDNODES")
	if err != nil {
		return fmt.Errorf("NUMWORLDNODES: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMWORLDNODES: missing node count")
	}
	numNodes, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num nodes: %w", err)
	}

	for i := 0; i < numNodes; i++ {
		node := WorldNode{}
		_, err = r.ReadProperty("WORLDNODE")
		if err != nil {
			return fmt.Errorf("WORLDNODE: %w", err)
		}

		records, err = r.ReadProperty("NORMALABCD")
		if err != nil {
			return fmt.Errorf("NORMALABCD: %w", err)
		}
		if len(records) < 5 {
			return fmt.Errorf("NORMALABCD: missing values")
		}
		node.Normals[0], err = helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("normal a: %w", err)
		}
		node.Normals[1], err = helper.ParseFloat32(records[2])
		if err != nil {
			return fmt.Errorf("normal b: %w", err)
		}
		node.Normals[2], err = helper.ParseFloat32(records[3])
		if err != nil {
			return fmt.Errorf("normal c: %w", err)
		}
		node.Normals[3], err = helper.ParseFloat32(records[4])
		if err != nil {
			return fmt.Errorf("normal d: %w", err)
		}

		records, err = r.ReadProperty("WORLDREGIONTAG")
		if err != nil {
			return fmt.Errorf("WORLDREGIONTAG: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("WORLDREGIONTAG: missing tag")
		}
		node.WorldRegionTag = records[1]

		records, err = r.ReadProperty("FRONTTREE")
		if err != nil {
			return fmt.Errorf("FRONTTREE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("FRONTTREE: missing value")
		}

		node.FrontTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = r.ReadProperty("BACKTREE")
		if err != nil {
			return fmt.Errorf("BACKTREE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("BACKTREE: missing value")
		}

		node.BackTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = r.ReadProperty("ENDWORLDNODE")
		if err != nil {
			return fmt.Errorf("ENDWORLDNODE: %w", err)
		}

		t.WorldNodes = append(t.WorldNodes, &node)
	}

	_, err = r.ReadProperty("ENDWORLDTREE")
	if err != nil {
		return fmt.Errorf("ENDWORLDTREE: %w", err)
	}

	return nil
}

type Region struct {
	RegionTag        string
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
	if r.VisTree != nil {
		fmt.Fprintf(w, "\tNUMVISNODE %d\n", len(r.VisTree.VisNodes))
		for i, node := range r.VisTree.VisNodes {
			fmt.Fprintf(w, "\tVISNODE // %d\n", i+1)
			fmt.Fprintf(w, "\t\tNORMALABCD %0.7f %0.7f %0.7f %0.7f\n", node.Normal[0], node.Normal[1], node.Normal[2], node.Normal[3])
			fmt.Fprintf(w, "\t\tVISLISTINDEX %d\n", node.VisListIndex)
			fmt.Fprintf(w, "\t\tFRONTTREE %d\n", node.FrontTree)
			fmt.Fprintf(w, "\t\tBACKTREE %d\n", node.BackTree)
			fmt.Fprintf(w, "\tENDVISNODE // %d\n", i+1)
		}
	}
	fmt.Fprintf(w, "\tSPHERE %0.7f %0.7f %0.7f %0.7f\n", r.Sphere[0], r.Sphere[1], r.Sphere[2], r.Sphere[3])
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", r.UserData)
	fmt.Fprintf(w, "ENDREGION\n\n")
	return nil
}

func (r *Region) Read(token *AsciiReadToken) error {
	r.VisTree = &VisTree{}
	records, err := token.ReadProperty("REGIONTAG")
	if err != nil {
		return fmt.Errorf("REGIONTAG: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("REGIONTAG: missing tag")
	}
	r.RegionTag = records[1]

	records, err = token.ReadProperty("NUMREGIONVERTEX")
	if err != nil {
		return fmt.Errorf("NUMREGIONVERTEX: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMREGIONVERTEX: missing vertex count")
	}
	numVertices, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ")
		if err != nil {
			return fmt.Errorf("XYZ: %w", err)
		}
		if len(records) < 4 {
			return fmt.Errorf("XYZ: missing coordinates")
		}
		var vert [3]float32
		vert[0], err = helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("vertex x: %w", err)
		}
		vert[1], err = helper.ParseFloat32(records[2])
		if err != nil {
			return fmt.Errorf("vertex y: %w", err)
		}
		vert[2], err = helper.ParseFloat32(records[3])
		if err != nil {
			return fmt.Errorf("vertex z: %w", err)
		}
		r.RegionVertices = append(r.RegionVertices, vert)
	}

	records, err = token.ReadProperty("NUMRENDERVERTICES")
	if err != nil {
		return fmt.Errorf("NUMRENDERVERTICES: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMRENDERVERTICES: missing vertex count")
	}
	numVertices, err = helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num render vertices: %w", err)
	}

	for i := 0; i < numVertices; i++ {
		records, err = token.ReadProperty("XYZ")
		if err != nil {
			return fmt.Errorf("XYZ: %w", err)
		}
		if len(records) < 4 {
			return fmt.Errorf("XYZ: missing coordinates")
		}
		var vert [3]float32
		vert[0], err = helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("vertex x: %w", err)
		}
		vert[1], err = helper.ParseFloat32(records[2])
		if err != nil {
			return fmt.Errorf("vertex y: %w", err)
		}
		vert[2], err = helper.ParseFloat32(records[3])
		if err != nil {
			return fmt.Errorf("vertex z: %w", err)
		}
		r.RenderVertices = append(r.RenderVertices, vert)
	}

	records, err = token.ReadProperty("NUMWALLS")
	if err != nil {
		return fmt.Errorf("NUMWALLS: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMWALLS: missing wall count")
	}
	numWalls, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num walls: %w", err)
	}

	if numWalls > 0 {
		return fmt.Errorf("Walls are not supported")
	}

	records, err = token.ReadProperty("NUMOBSTACLES")
	if err != nil {
		return fmt.Errorf("NUMOBSTACLES: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMOBSTACLES: missing obstacle count")
	}
	numObstacles, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num obstacles: %w", err)
	}

	if numObstacles > 0 {
		return fmt.Errorf("Obstacles are not supported")
	}

	records, err = token.ReadProperty("NUMCUTTINGOBSTACLES")
	if err != nil {
		return fmt.Errorf("NUMCUTTINGOBSTACLES: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMCUTTINGOBSTACLES: missing cutting obstacle count")
	}
	numCuttingObstacles, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num cutting obstacles: %w", err)
	}

	if numCuttingObstacles > 0 {
		return fmt.Errorf("cutting obstacles are not supported")
	}

	_, err = token.ReadProperty("VISTREE")
	if err != nil {
		return fmt.Errorf("VISTREE: %w", err)
	}

	records, err = token.ReadProperty("NUMVISNODE")
	if err != nil {
		return fmt.Errorf("NUMVISNODE: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMVISNODE: missing vis node count")
	}
	numVisNodes, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vis nodes: %w", err)
	}

	records, err = token.ReadProperty("NUMVISLIST")
	if err != nil {
		return fmt.Errorf("NUMVISLIST: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMVISLIST: missing vis list count")
	}
	numVisLists, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num vis lists: %w", err)
	}

	for i := 0; i < numVisNodes; i++ {
		node := VisNode{}
		_, err = token.ReadProperty("VISNODE")
		if err != nil {
			return fmt.Errorf("VISNODE: %w", err)
		}

		records, err = token.ReadProperty("NORMALABCD")
		if err != nil {
			return fmt.Errorf("NORMALABCD: %w", err)
		}
		if len(records) < 5 {
			return fmt.Errorf("NORMALABCD: missing values")
		}
		node.Normal[0], err = helper.ParseFloat32(records[1])
		if err != nil {
			return fmt.Errorf("normal a: %w", err)
		}
		node.Normal[1], err = helper.ParseFloat32(records[2])
		if err != nil {
			return fmt.Errorf("normal b: %w", err)
		}
		node.Normal[2], err = helper.ParseFloat32(records[3])
		if err != nil {
			return fmt.Errorf("normal c: %w", err)
		}
		node.Normal[3], err = helper.ParseFloat32(records[4])
		if err != nil {
			return fmt.Errorf("normal d: %w", err)
		}

		records, err = token.ReadProperty("VISLISTINDEX")
		if err != nil {
			return fmt.Errorf("VISLISTINDEX: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("VISLISTINDEX: missing value")
		}
		node.VisListIndex, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("vis list index: %w", err)
		}

		records, err = token.ReadProperty("FRONTTREE")
		if err != nil {
			return fmt.Errorf("FRONTTREE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("FRONTTREE: missing value")
		}
		node.FrontTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("front tree: %w", err)
		}

		records, err = token.ReadProperty("BACKTREE")
		if err != nil {
			return fmt.Errorf("BACKTREE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("BACKTREE: missing value")
		}
		node.BackTree, err = helper.ParseUint32(records[1])
		if err != nil {
			return fmt.Errorf("back tree: %w", err)
		}

		_, err = token.ReadProperty("ENDVISNODE")
		if err != nil {
			return fmt.Errorf("ENDVISNODE: %w", err)
		}

		r.VisTree.VisNodes = append(r.VisTree.VisNodes, &node)
	}

	for i := 0; i < numVisLists; i++ {
		list := VisList{}
		_, err = token.ReadProperty("VISIBLELIST")
		if err != nil {
			return fmt.Errorf("VISIBLELIST: %w", err)
		}
		records, err = token.ReadProperty("RANGE")
		if err != nil {
			return fmt.Errorf("RANGE: %w", err)
		}
		if len(records) < 2 {
			return fmt.Errorf("RANGE: missing range")
		}
		numVisRange, err := helper.ParseInt(records[1])
		if err != nil {
			return fmt.Errorf("vis range: %w", err)
		}

		if len(records) < numVisRange+2 {
			return fmt.Errorf("RANGE: expected %d values, got %d", numVisRange+2, len(records))
		}

		for j := 0; j < numVisRange; j++ {
			val, err := helper.ParseUint32(records[j+2])
			if err != nil {
				return fmt.Errorf("vis range %d: %w", j, err)
			}
			list.Range = append(list.Range, val)
		}

		_, err = token.ReadProperty("ENDVISIBLELIST")
		if err != nil {
			return fmt.Errorf("ENDVISIBLELIST: %w", err)
		}

		r.VisTree.VisLists = append(r.VisTree.VisLists, &list)
	}

	_, err = token.ReadProperty("ENDVISTREE")
	if err != nil {
		return fmt.Errorf("ENDVISTREE: %w", err)
	}

	records, err = token.ReadProperty("SPHERE")
	if err != nil {
		return fmt.Errorf("SPHERE: %w", err)
	}
	if len(records) < 5 {
		return fmt.Errorf("SPHERE: missing values")
	}
	r.Sphere[0], err = helper.ParseFloat32(records[1])
	if err != nil {
		return fmt.Errorf("sphere x: %w", err)
	}
	r.Sphere[1], err = helper.ParseFloat32(records[2])
	if err != nil {
		return fmt.Errorf("sphere y: %w", err)
	}
	r.Sphere[2], err = helper.ParseFloat32(records[3])
	if err != nil {
		return fmt.Errorf("sphere z: %w", err)
	}
	r.Sphere[3], err = helper.ParseFloat32(records[4])
	if err != nil {
		return fmt.Errorf("sphere radius: %w", err)
	}

	records, err = token.ReadProperty("USERDATA")
	if err != nil {
		return fmt.Errorf("USERDATA: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("USERDATA: missing data")
	}
	r.UserData = records[1]

	_, err = token.ReadProperty("ENDREGION")
	if err != nil {
		return fmt.Errorf("ENDREGION: %w", err)
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
	fmt.Fprintf(w, "\tNUMREGIONS %d\n", len(a.Regions))
	regions := ""
	for _, region := range a.Regions {
		regions += fmt.Sprintf("%d ", region)
	}
	fmt.Fprintf(w, "\tREGIONS %d %s\n", len(a.Regions), regions)
	fmt.Fprintf(w, "ENDAMBIENTLIGHT\n\n")
	return nil
}

func (a *AmbientLight) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG")
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("TAG: missing tag")
	}
	a.Tag = records[1]

	records, err = r.ReadProperty("LIGHT")
	if err != nil {
		return fmt.Errorf("LIGHT: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("LIGHT: missing tag")
	}
	a.LightTag = records[1]

	records, err = r.ReadProperty("NUMREGIONS")
	if err != nil {
		return fmt.Errorf("NUMREGIONS: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMREGIONS: missing region count")
	}
	numRegions, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num regions: %w", err)
	}

	records, err = r.ReadProperty("REGIONS")
	if err != nil {
		return fmt.Errorf("REGIONS: %w", err)
	}

	if len(records) < numRegions+1 {
		return fmt.Errorf("NUMREGIONS: expected %d values, got %d", numRegions+2, len(records))
	}

	for i := 0; i < numRegions; i++ {
		region, err := helper.ParseUint32(records[i+1])
		if err != nil {
			return fmt.Errorf("region %d: %w", i, err)
		}
		a.Regions = append(a.Regions, region)
	}

	_, err = r.ReadProperty("ENDAMBIENTLIGHT")
	if err != nil {
		return fmt.Errorf("ENDAMBIENTLIGHT: %w", err)
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
	fmt.Fprintf(w, "\tNUMREGIONS %d\n", len(z.Regions))
	regions := ""
	for _, region := range z.Regions {
		regions += fmt.Sprintf("%d ", region)
	}
	fmt.Fprintf(w, "\tREGIONS %d %s\n", len(z.Regions), regions)
	fmt.Fprintf(w, "\tUSERDATA \"%s\"\n", z.UserData)
	fmt.Fprintf(w, "ENDZONE\n\n")
	return nil
}

func (z *Zone) Read(r *AsciiReadToken) error {
	records, err := r.ReadProperty("TAG")
	if err != nil {
		return fmt.Errorf("TAG: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("TAG: missing tag")
	}
	z.Tag = records[1]

	records, err = r.ReadProperty("NUMREGIONS")
	if err != nil {
		return fmt.Errorf("NUMREGIONS: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("NUMREGIONS: missing region count")
	}
	numRegions, err := helper.ParseInt(records[1])
	if err != nil {
		return fmt.Errorf("num regions: %w", err)
	}

	records, err = r.ReadProperty("REGIONS")
	if err != nil {
		return fmt.Errorf("REGIONS: %w", err)
	}

	if len(records) < numRegions+1 {
		return fmt.Errorf("NUMREGIONS: expected %d values, got %d", numRegions+2, len(records))
	}

	for i := 0; i < numRegions; i++ {
		region, err := helper.ParseUint32(records[i+1])
		if err != nil {
			return fmt.Errorf("region %d: %w", i, err)
		}
		z.Regions = append(z.Regions, region)
	}

	records, err = r.ReadProperty("USERDATA")
	if err != nil {
		return fmt.Errorf("USERDATA: %w", err)
	}
	if len(records) < 2 {
		return fmt.Errorf("USERDATA: missing data")
	}
	z.UserData = records[1]

	_, err = r.ReadProperty("ENDZONE")
	if err != nil {
		return fmt.Errorf("ENDZONE: %w", err)
	}

	return nil
}
