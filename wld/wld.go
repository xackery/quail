// virtual is Virtual World file format, it is used to make binary world more human readable and editable
package wld

import (
	"fmt"
	"io"
	"strconv"
	"strings"
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
	TrackInstances     []*TrackInstance
	TrackDefs          []*TrackDef

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
	fmt.Fprintf(w, "\tNUMVERTICES %d\n", len(d.Vertices))
	for _, vert := range d.Vertices {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMUVS %d\n", len(d.UVs))
	for _, uv := range d.UVs {
		fmt.Fprintf(w, "\tUV %0.7f %0.7f\n", uv[0], uv[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMVERTEXNORMALS %d\n", len(d.VertexNormals))
	for _, vn := range d.VertexNormals {
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vn[0], vn[1], vn[2])
	}
	fmt.Fprintf(w, "\n")
	assigments := ""
	for _, sa := range d.SkinAssignmentGroups {
		assigments += fmt.Sprintf("%d %d, ", sa[0], sa[1])
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tSKINASSIGNMENTGROUPS %s\n", assigments)

	fmt.Fprintf(w, "\tMATERIALPALETTE \"%s\"\n", d.MaterialPaletteTag)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMFACE2S %d\n", len(d.Faces))
	fmt.Fprintf(w, "\n")
	for i, face := range d.Faces {
		fmt.Fprintf(w, "\tDMFACE2S //%d\n", i+1)
		if face.Flags != 0 {
			fmt.Fprintf(w, "\t\tFLAGS %d\n", face.Flags)
		}
		fmt.Fprintf(w, "\t\tTRIANGLE   %d, %d, %d\n", face.Triangle[0], face.Triangle[1], face.Triangle[2])
		fmt.Fprintf(w, "\tENDFACE //%d\n\n", i+1)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tNUMMESHOPS 0\n")
	fmt.Fprintf(w, "\t//TODO: NUMMESHOPS %d\n", len(d.MeshOps))
	for _, meshOp := range d.MeshOps {
		fmt.Fprintf(w, "\t// TODO: MESHOP %d %d %0.7f %d %d\n", meshOp.Index1, meshOp.Index2, meshOp.Offset, meshOp.Param1, meshOp.TypeField)
		// MESHOP_VA %d
	}
	fmt.Fprintf(w, "\n")
	groups := ""
	for _, group := range d.FaceMaterialGroups {
		groups += fmt.Sprintf("%d %d, ", group[0], group[1])
	}
	if len(groups) > 0 {
		groups = groups[:len(groups)-2]
	}
	fmt.Fprintf(w, "\tFACEMATERIALGROUPS %s\n", groups)
	groups = ""
	for _, group := range d.VertexMaterialGroups {
		groups += fmt.Sprintf("%d %d, ", group[0], group[1])
	}
	if len(groups) > 0 {
		groups = groups[:len(groups)-2]
	}
	fmt.Fprintf(w, "\tVERTEXMATERIALGROUPS %s\n", groups)
	fmt.Fprintf(w, "\tBOUNDINGRADIUS %0.7f\n", d.BoundingRadius)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tFPSCALE %d\n", d.FPScale)
	fmt.Fprintf(w, "ENDDMSPRITEDEF2\n\n")
	return nil
}

func (d *DMSpriteDef2) Read(r *AsciiReadToken) error {
	definition := "DMSPRITEDEF2"
	for {
		line, err := r.ReadProperty(definition)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if line == "ENDDMSPRITEDEF2" {
			return nil
		}
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "TAG"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "TAG %s", &d.Tag)
			if err != nil {
				return fmt.Errorf("tag: %w", err)
			}
		case strings.HasPrefix(line, "FLAGS"):
			_, err = fmt.Sscanf(line, "FLAGS %d", &d.Flags)
			if err != nil {
				return fmt.Errorf("flags: %w", err)
			}
		case strings.HasPrefix(line, "CENTEROFFSET"):
			_, err = fmt.Sscanf(line, "CENTEROFFSET %f %f %f", &d.CenterOffset[0], &d.CenterOffset[1], &d.CenterOffset[2])
			if err != nil {
				return fmt.Errorf("center offset: %w", err)
			}
		case strings.HasPrefix(line, "NUMVERTICES"):
			var numVertices int
			_, err = fmt.Sscanf(line, "NUMVERTICES %d", &numVertices)
			if err != nil {
				return fmt.Errorf("num vertices: %w", err)
			}
			d.Vertices = make([][3]float32, numVertices)
			for i := 0; i < numVertices; i++ {
				line, err = r.ReadProperty(definition)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(line, "XYZ") {
					return fmt.Errorf("expected XYZ, got %s", line)
				}
				_, err = fmt.Sscanf(line, "XYZ %f %f %f", &d.Vertices[i][0], &d.Vertices[i][1], &d.Vertices[i][2])
				if err != nil {
					return fmt.Errorf("vertex %d: %w", i, err)
				}
			}
		case strings.HasPrefix(line, "NUMUVS"):
			var numUVs int
			_, err = fmt.Sscanf(line, "NUMUVS %d", &numUVs)
			if err != nil {
				return fmt.Errorf("num uvs: %w", err)
			}
			d.UVs = make([][2]float32, numUVs)

			for i := 0; i < numUVs; i++ {
				line, err = r.ReadProperty(definition)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(line, "UV") {
					return fmt.Errorf("expected UV, got %s", line)
				}
				line = strings.TrimPrefix(line, "UV")
				line = strings.TrimSpace(line)
				uvs := strings.Split(line, ", ")
				if len(uvs) != 2 {
					return fmt.Errorf("expected 2 uvs, got %d", len(uvs))
				}
				u, err := strconv.ParseFloat(uvs[0], 32)
				if err != nil {
					return fmt.Errorf("uv %d u: %w", i, err)
				}
				v, err := strconv.ParseFloat(uvs[1], 32)
				if err != nil {
					return fmt.Errorf("uv %d v: %w", i, err)
				}
				d.UVs[i] = [2]float32{float32(u), float32(v)}
			}
		case strings.HasPrefix(line, "NUMVERTEXNORMALS"):
			var numNormals int
			_, err = fmt.Sscanf(line, "NUMVERTEXNORMALS %d", &numNormals)
			if err != nil {
				return fmt.Errorf("num normals: %w", err)
			}
			d.VertexNormals = make([][3]float32, numNormals)
			for i := 0; i < numNormals; i++ {
				line, err = r.ReadProperty(definition)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(line, "XYZ") {
					return fmt.Errorf("expected XYZ, got %s", line)
				}

				_, err = fmt.Sscanf(line, "XYZ %f %f %f", &d.VertexNormals[i][0], &d.VertexNormals[i][1], &d.VertexNormals[i][2])
				if err != nil {
					return fmt.Errorf("normal %d: %w", i, err)
				}
			}
		case strings.HasPrefix(line, "SKINASSIGNMENTGROUPS"):
			line = strings.TrimPrefix(line, "SKINASSIGNMENTGROUPS")
			line = strings.TrimSpace(line)
			index := strings.Index(line, " ")
			if index == -1 {
				return fmt.Errorf("expected space in skin assignment groups")
			}
			numGroups, err := strconv.Atoi(line[:index])
			if err != nil {
				return fmt.Errorf("num groups: %w", err)
			}
			d.SkinAssignmentGroups = make([][2]uint16, numGroups)
			line = line[index+1:]
			line = strings.ReplaceAll(line, ",", "")
			for i := 0; i < numGroups; i++ {
				index = strings.Index(line, " ")
				if index == -1 {
					return fmt.Errorf("expected space for val0 in skin assignment group %d", i)
				}
				val0, err := strconv.ParseUint(line[:index], 10, 16)
				if err != nil {
					return fmt.Errorf("group %d val0: %w", i, err)
				}
				line = line[index+1:]
				index = strings.Index(line, " ")
				if i == numGroups-1 {
					index = len(line)
				}
				if index == -1 {
					return fmt.Errorf("expected space for val1 in skin assignment group %d", i)
				}
				val1, err := strconv.ParseUint(line[:index], 10, 16)
				if err != nil {
					return fmt.Errorf("group %d val1: %w", i, err)
				}
				if i < numGroups-1 {
					line = line[index+1:]
				}
				d.SkinAssignmentGroups[i] = [2]uint16{uint16(val0), uint16(val1)}
			}

		case strings.HasPrefix(line, "MATERIALPALETTE"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "MATERIALPALETTE %s", &d.MaterialPaletteTag)
			if err != nil {
				return fmt.Errorf("material palette tag: %w", err)
			}
		case strings.HasPrefix(line, "NUMCOLORS"):
			var numColors int
			_, err = fmt.Sscanf(line, "NUMCOLORS %d", &numColors)
			if err != nil {
				return fmt.Errorf("num colors: %w", err)
			}
			d.Colors = make([][4]uint8, numColors)
			for i := 0; i < numColors; i++ {
				line, err = r.ReadProperty(definition)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(line, "RGBA") {
					return fmt.Errorf("expected RGBA, got %s", line)
				}
				_, err = fmt.Sscanf(line, "RGBA %d %d %d %d", &d.Colors[i][0], &d.Colors[i][1], &d.Colors[i][2], &d.Colors[i][3])
				if err != nil {
					return fmt.Errorf("color %d: %w", i, err)
				}
			}
		case strings.HasPrefix(line, "NUMFACE2S"):
			var numFaces int

			_, err = fmt.Sscanf(line, "NUMFACE2S %d", &numFaces)
			if err != nil {
				return fmt.Errorf("num faces: %w", err)
			}
			d.Faces = make([]*Face, numFaces)
			for i := 0; i < numFaces; i++ {
				face := &Face{}
				line, err = r.ReadProperty(definition)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(line, "DMFACE2") {
					return fmt.Errorf("expected DMFACE2, got %s", line)
				}
				for {
					line, err = r.ReadProperty(definition)
					if err != nil {
						return err
					}
					if strings.HasPrefix(line, "ENDDMFACE2") {
						break
					}
					if strings.HasPrefix(line, "FLAGS") {
						_, err = fmt.Sscanf(line, "FLAGS %d", &face.Flags)
						if err != nil {
							return fmt.Errorf("face %d flags: %w", i, err)
						}
					} else if strings.HasPrefix(line, "TRIANGLE") {
						_, err = fmt.Sscanf(line, "TRIANGLE   %d, %d, %d", &face.Triangle[0], &face.Triangle[1], &face.Triangle[2])
						if err != nil {
							return fmt.Errorf("face %d triangle: %w", i, err)
						}
					}
				}
				d.Faces[i] = face
			}
		case strings.HasPrefix(line, "NUMMESHOPS"):
			var numMeshOps int
			_, err = fmt.Sscanf(line, "NUMMESHOPS %d", &numMeshOps)
			if err != nil {
				return fmt.Errorf("num mesh ops: %w", err)
			}
			d.MeshOps = make([]*MeshOp, numMeshOps)
			for i := 0; i < numMeshOps; i++ {
				meshOp := &MeshOp{}
				line, err = r.ReadProperty(definition)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(line, "MESHOP") {
					return fmt.Errorf("expected MESHOP, got %s", line)
				}

				line = strings.TrimPrefix(line, "MESHOP")
				line = strings.TrimSpace(line)
				index := strings.Index(line, " ")
				if index == -1 {
					return fmt.Errorf("expected space in mesh op %d", i)
				}
				//context := line[1:index]
				//fmt.Println("context:", context)

				d.MeshOps[i] = meshOp
			}
		case strings.HasPrefix(line, "FACEMATERIALGROUPS"):
			line = strings.TrimPrefix(line, "FACEMATERIALGROUPS")
			line = strings.TrimSpace(line)
			index := strings.Index(line, " ")
			if index == -1 {
				return fmt.Errorf("expected space in face material groups")
			}
			numGroups, err := strconv.Atoi(line[:index])
			if err != nil {
				return fmt.Errorf("num groups: %w", err)
			}
			d.FaceMaterialGroups = make([][2]uint16, numGroups)
			line = line[index+1:]
			line = strings.ReplaceAll(line, ",", "")
			for i := 0; i < numGroups; i++ {
				index = strings.Index(line, " ")
				if index == -1 {
					return fmt.Errorf("expected space for val0 in face material group %d", i)
				}
				val0, err := strconv.ParseUint(line[:index], 10, 16)
				if err != nil {
					return fmt.Errorf("group %d val0: %w", i, err)
				}
				line = line[index+1:]
				index = strings.Index(line, " ")
				if i == numGroups-1 {
					index = len(line)
				}
				if index == -1 {
					return fmt.Errorf("expected space for val1 in face material group %d", i)
				}
				val1, err := strconv.ParseUint(line[:index], 10, 16)
				if err != nil {
					return fmt.Errorf("group %d val1: %w", i, err)
				}
				if i < numGroups-1 {
					line = line[index+1:]
				}
				d.FaceMaterialGroups[i] = [2]uint16{uint16(val0), uint16(val1)}
			}

		case strings.HasPrefix(line, "VERTEXMATERIALGROUPS"):
			line = strings.TrimPrefix(line, "VERTEXMATERIALGROUPS")
			line = strings.TrimSpace(line)
			index := strings.Index(line, " ")
			if index == -1 {
				return fmt.Errorf("expected space in vertex material groups")
			}
			numGroups, err := strconv.Atoi(line[:index])
			if err != nil {
				return fmt.Errorf("num groups: %w", err)
			}
			d.VertexMaterialGroups = make([][2]int16, numGroups)
			line = line[index+1:]
			line = strings.ReplaceAll(line, ",", "")
			for i := 0; i < numGroups; i++ {
				index = strings.Index(line, " ")
				if index == -1 {
					return fmt.Errorf("expected space for val0 in vertex material group %d", i)
				}
				val0, err := strconv.ParseInt(line[:index], 10, 16)
				if err != nil {
					return fmt.Errorf("group %d val0: %w", i, err)
				}
				line = line[index+1:]
				index = strings.Index(line, " ")
				if i == numGroups-1 {
					index = len(line)
				}
				if index == -1 {
					return fmt.Errorf("expected space for val1 in vertex material group %d", i)
				}
				val1Str := line[:index]

				val1, err := strconv.ParseInt(val1Str, 10, 16)
				if err != nil {
					return fmt.Errorf("group %d val1: %w", i, err)
				}
				if i < numGroups-1 {
					line = line[index+1:]
				}
				d.VertexMaterialGroups[i] = [2]int16{int16(val0), int16(val1)}
			}

		case strings.HasPrefix(line, "BOUNDINGRADIUS"):
			_, err = fmt.Sscanf(line, "BOUNDINGRADIUS %f", &d.BoundingRadius)
			if err != nil {
				return fmt.Errorf("bounding radius: %w", err)
			}
		case strings.HasPrefix(line, "FPSCALE"):
			_, err = fmt.Sscanf(line, "FPSCALE %d", &d.FPScale)
			if err != nil {
				return fmt.Errorf("fpscale: %w", err)
			}
		default:
			return fmt.Errorf("unknown property: %s", line)
		}
	}
	return nil
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
	for {
		line, err := r.ReadProperty(m.Definition())
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fmt.Println("line", line)
		if line == "ENDMATERIALPALETTE" {
			break
		}
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "TAG"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "TAG %s", &m.Tag)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "NUMMATERIALS"):
			_, err = fmt.Sscanf(line, "NUMMATERIALS %d", &m.numMaterials)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "MATERIAL"):
			line = strings.ReplaceAll(line, "\"", "")
			var mat string
			_, err = fmt.Sscanf(line, "MATERIAL %s", &mat)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			m.Materials = append(m.Materials, mat)
		}
	}

	if m.Tag == "" {
		return fmt.Errorf("missing tag")
	}

	if m.numMaterials != len(m.Materials) {
		return fmt.Errorf("expected %d materials, got %d", m.numMaterials, len(m.Materials))
	}
	return io.EOF
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

/*
POLYHEDRONDEFINITION

	TAG	"prepe_POLYHDEF"
	BOUNDINGRADIUS	1.2431762e+002
	SCALEFACTOR	1.0
	NUMVERTICES	287
	XYZ	-5.9604645e-008 1.9073486e-005 -3.8146973e-006
	NUMFACES	280
	FACE 1
		NUMVERTICES	3
		VERTEXLIST	3, 1, 2
	ENDFACE 1
	ENDPOLYHEDRONDEFINITION
*/
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
		fmt.Fprintf(w, "\tXYZ %0.7f %0.7f %0.7f\n", vert[0], vert[1], vert[2])
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
	for {
		line, err := r.ReadProperty(p.Definition())
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if line == "ENDPOLYHEDRONDEFINITION" {
			return nil
		}
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "TAG"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "TAG %s", &p.Tag)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "BOUNDINGRADIUS"):
			valStr := ""
			_, err = fmt.Sscanf(line, "BOUNDINGRADIUS %s", &valStr)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			val, err := strconv.ParseFloat(valStr, 32)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			p.BoundingRadius = float32(val)
		case strings.HasPrefix(line, "SCALEFACTOR"):
			valStr := ""
			_, err = fmt.Sscanf(line, "SCALEFACTOR %s", &valStr)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			val, err := strconv.ParseFloat(valStr, 32)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			p.ScaleFactor = float32(val)
		case strings.HasPrefix(line, "NUMVERTICES"):
			_, err = fmt.Sscanf(line, "NUMVERTICES %d", &p.numVertices)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			p.Vertices = make([][3]float32, p.numVertices)
			for i := 0; i < p.numVertices; i++ {
				line, err = r.ReadProperty(p.Definition())
				if err != nil {
					return err
				}
				valStr1, valStr2, valStr3 := "", "", ""
				_, err = fmt.Sscanf(line, "XYZ %s %s %s", &valStr1, &valStr2, &valStr3)
				if err != nil {
					return fmt.Errorf("vertex %d: %w", i, err)
				}
				val1, err := strconv.ParseFloat(valStr1, 32)
				if err != nil {
					return fmt.Errorf("vertex %d: %w", i, err)
				}
				val2, err := strconv.ParseFloat(valStr2, 32)
				if err != nil {
					return fmt.Errorf("vertex %d: %w", i, err)
				}
				val3, err := strconv.ParseFloat(valStr3, 32)
				if err != nil {
					return fmt.Errorf("vertex %d: %w", i, err)
				}
				p.Vertices[i] = [3]float32{float32(val1), float32(val2), float32(val3)}
			}
		case strings.HasPrefix(line, "NUMFACES"):
			_, err = fmt.Sscanf(line, "NUMFACES %d", &p.numFaces)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			p.Faces = make([]*PolyhedronDefinitionFace, p.numFaces)
			for i := 0; i < p.numFaces; i++ {
				line, err = r.ReadProperty(p.Definition())
				if err != nil {
					return err
				}
				if line == "" {
					continue
				}
				if !strings.HasPrefix(line, "FACE") {
					return fmt.Errorf("expected FACE %d, got %s", i+1, line)
				}
				face := &PolyhedronDefinitionFace{}
				_, err = fmt.Sscanf(line, "FACE %d", &face.numVertices)
				if err != nil {
					return fmt.Errorf("face %d: %w", i+1, err)
				}
				face.Vertices = make([]uint32, face.numVertices)
				line, err = r.ReadProperty(p.Definition())
				if err != nil {
					return err
				}
				if line == "" {
					continue
				}
				if !strings.HasPrefix(line, "NUMVERTICES") {
					return fmt.Errorf("expected FACE %d NUMVERTICES, got %s", i+1, line)
				}
				numVertices := 0
				_, err = fmt.Sscanf(line, "NUMVERTICES %d", &numVertices)
				if err != nil {
					return fmt.Errorf("face %d numvertices: %w", i+1, err)
				}
				face.Vertices = make([]uint32, numVertices)
				line, err = r.ReadProperty(p.Definition())
				if err != nil {
					return err
				}
				if line == "" {
					continue
				}
				if !strings.HasPrefix(line, "VERTEXLIST") {
					return fmt.Errorf("expected VERTEXLIST, got %s", line)
				}

				vertStr := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "VERTEXLIST")), ",")
				if len(vertStr) != numVertices {
					return fmt.Errorf("face %d: expected %d vertices, got %d", i+1, numVertices, len(vertStr))
				}
				for k, v := range vertStr {
					v = strings.TrimSpace(v)
					val, err := strconv.ParseUint(v, 10, 32)
					if err != nil {
						return fmt.Errorf("face %d element %d: %w", i+1, k, err)
					}
					face.Vertices[k] = uint32(val)
				}
				line, err = r.ReadProperty(p.Definition())
				if err != nil {
					return err
				}
				if line == "" {
					continue
				}
				if !strings.HasPrefix(line, "ENDFACE") {
					return fmt.Errorf("expected ENDFACE %d, got %s", i, line)
				}
				p.Faces[i] = face
			}
		}
	}
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
	if t.Sleep != 0 {
		fmt.Fprintf(w, "\tSLEEP %d\n", t.Sleep)
	}
	fmt.Fprintf(w, "ENDTRACKDEFINITION\n\n")
	return nil
}

func (t *TrackInstance) Read(r *AsciiReadToken) error {
	for {
		line, err := r.ReadProperty(t.Definition())
		if err != nil {
			if err == io.EOF {
				break
			}
			return err

		}
		if line == "ENDTRACKINSTANCE" {
			break
		}
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "TAG"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "TAG %s", &t.Tag)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "DEFINITION"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "DEFINITION %s", &t.DefiniationTag)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "INTERPOLATE"):
			t.isInterpolated = true
		case strings.HasPrefix(line, "SLEEP"):
			_, err = fmt.Sscanf(line, "SLEEP %d", &t.Sleep)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		}
	}
	return nil
}

type TrackDef struct {
	Tag            string                // TAG "%s"
	numFrames      int                   // NUMFRAMES %d
	FrameTransform []TrackFrameTransform // FRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f
}

type TrackFrameTransform struct {
}

func (t *TrackDef) Definition() string {
	return "TRACKDEFINITION"
}

func (t *TrackDef) Write(w io.Writer) error {
	fmt.Fprintf(w, "%s\n", t.Definition())
	fmt.Fprintf(w, "\tTAG \"%s\"\n", t.Tag)
	fmt.Fprintf(w, "\tNUMFRAMES %d\n", t.numFrames)
	//for _, frame := range t.FrameTransform {
	//	fmt.Fprintf(w, "\tFRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f\n", frame)
	//}
	fmt.Fprintf(w, "ENDTRACKDEFINITION\n\n")
	return nil
}

func (t *TrackDef) Read(r *AsciiReadToken) error {
	for {
		line, err := r.ReadProperty(t.Definition())
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if line == "ENDTRACKDEFINITION" {
			break
		}
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "TAG"):
			line = strings.ReplaceAll(line, "\"", "")
			_, err = fmt.Sscanf(line, "TAG %s", &t.Tag)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "NUMFRAMES"):
			_, err = fmt.Sscanf(line, "NUMFRAMES %d", &t.numFrames)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
		case strings.HasPrefix(line, "FRAMETRANSFORM"):
			frame := TrackFrameTransform{}
			_, err = fmt.Sscanf(line, "FRAMETRANSFORM %0.7f %d %d %d %0.7f %0.7f %0.7f", &frame)
			if err != nil {
				return fmt.Errorf("%s: %w", line, err)
			}
			t.FrameTransform = append(t.FrameTransform, frame)
		}
	}
	return nil
}
