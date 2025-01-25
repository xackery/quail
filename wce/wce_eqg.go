package wce

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/raw"
)

// MdsDef is an entry EQSKINNEDMODELDEF
type MdsDef struct {
	folders   []string
	Tag       string
	Version   uint32
	Materials []*EQMaterialDef
	Bones     []*MdsBone
	Models    []*MdsModel
}

type MdsBone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         [3]float32
	Quaternion    [4]float32
	Scale         [3]float32
}

type MdsModel struct {
	MainPiece       uint32 // 0: no, 1: yes, head is a mainpiece
	Name            string
	Vertices        []*MdsVertex
	Faces           []*MdsFace
	BoneAssignments [][4]*MdsBoneWeight
}

type MdsVertex struct {
	Position [3]float32
	Normal   [3]float32
	Tint     [4]uint8
	Uv       [2]float32
	Uv2      [2]float32
}

type MdsFace struct {
	Index        [3]uint32
	MaterialName string
	Flags        uint32
}

type MdsBoneWeight struct {
	BoneIndex int32
	Value     float32
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
		fmt.Fprintf(w, "\tNUMMATERIALS %d\n", len(e.Materials))
		for _, material := range e.Materials {
			fmt.Fprintf(w, "\t\tMATERIAL \"%s\"\n", material.Tag)
		}

		fmt.Fprintf(w, "\tNUMBONES %d\n", len(e.Bones))
		for _, bone := range e.Bones {
			fmt.Fprintf(w, "\t\tBONE \"%s\"\n", bone.Name)
			fmt.Fprintf(w, "\t\t\tNEXT %d\n", bone.Next)
			fmt.Fprintf(w, "\t\t\tCHILDREN %d\n", bone.ChildrenCount)
			fmt.Fprintf(w, "\t\t\tCHILDINDEX %d\n", bone.ChildIndex)
			fmt.Fprintf(w, "\t\t\tPIVOT %0.8e %0.8e %0.8e\n", bone.Pivot[0], bone.Pivot[1], bone.Pivot[2])
			fmt.Fprintf(w, "\t\t\tQUATERNION %0.8e %0.8e %0.8e %0.8e\n", bone.Quaternion[0], bone.Quaternion[1], bone.Quaternion[2], bone.Quaternion[3])
			fmt.Fprintf(w, "\t\t\tSCALE %0.8e %0.8e %0.8e\n", bone.Scale[0], bone.Scale[1], bone.Scale[2])
		}

		fmt.Fprintf(w, "\tNUMMODELS %d\n", len(e.Models))
		for _, model := range e.Models {
			fmt.Fprintf(w, "\t\tMODEL \"%s\"\n", model.Name)
			fmt.Fprintf(w, "\t\t\tMAINPIECE %d\n", model.MainPiece)
			fmt.Fprintf(w, "\t\t\tNUMVERTICES %d\n", len(model.Vertices))
			for _, vert := range model.Vertices {
				fmt.Fprintf(w, "\t\t\t\t\tXYZ %0.8e %0.8e %0.8e\n", vert.Position[0], vert.Position[1], vert.Position[2])
			}
			fmt.Fprintf(w, "\t\t\tNUMUVS %d\n", len(model.Vertices))
			for _, vert := range model.Vertices {
				fmt.Fprintf(w, "\t\t\t\t\tUV %0.8e %0.8e\n", vert.Uv[0], vert.Uv[1])
			}
			fmt.Fprintf(w, "\t\t\tNUMUV2S %d\n", len(model.Vertices))
			for _, vert := range model.Vertices {
				fmt.Fprintf(w, "\t\t\t\t\tUV2 %0.8e %0.8e\n", vert.Uv2[0], vert.Uv2[1])
			}

			fmt.Fprintf(w, "\t\t\tNUMNORMALS %d\n", len(model.Vertices))
			for _, vert := range model.Vertices {
				fmt.Fprintf(w, "\t\t\t\t\tNORMAL %0.8e %0.8e %0.8e\n", vert.Normal[0], vert.Normal[1], vert.Normal[2])
			}
			fmt.Fprintf(w, "\t\t\tNUMTINTS %d\n", len(model.Vertices))
			for _, vert := range model.Vertices {
				fmt.Fprintf(w, "\t\t\t\t\tTINT %d %d %d %d\n", vert.Tint[0], vert.Tint[1], vert.Tint[2], vert.Tint[3])
			}

			fmt.Fprintf(w, "\t\t\tNUMFACES %d\n", len(model.Faces))
			for _, face := range model.Faces {
				fmt.Fprintf(w, "\t\t\t\tFACE\n")
				fmt.Fprintf(w, "\t\t\t\t\tTRIANGLE %d %d %d\n", face.Index[0], face.Index[1], face.Index[2])
				fmt.Fprintf(w, "\t\t\t\t\tMATERIAL \"%s\"\n", face.MaterialName)
				fmt.Fprintf(w, "\t\t\t\t\tFLAGS %d\n", face.Flags)
			}
			fmt.Fprintf(w, "\t\t\tNUMWEIGHTS %d\n", len(model.BoneAssignments))
			for _, weights := range model.BoneAssignments {
				fmt.Fprintf(w, "\t\t\t\tWEIGHT %d %0.8e %d %0.8e %d %0.8e %d %0.8e \n", weights[0].BoneIndex, weights[0].Value, weights[1].BoneIndex, weights[1].Value, weights[2].BoneIndex, weights[2].Value, weights[3].BoneIndex, weights[3].Value)
			}
		}
		fmt.Fprintf(w, "\n")
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

	records, err = token.ReadProperty("NUMMATERIALS", 1)
	if err != nil {
		return err
	}

	numMaterials := 0
	err = parse(&numMaterials, records[1])
	if err != nil {
		return fmt.Errorf("num materials: %w", err)
	}

	for i := 0; i < numMaterials; i++ {
		eqMaterialDef := &EQMaterialDef{}
		err = eqMaterialDef.Read(token)
		if err != nil {
			return fmt.Errorf("material %d: %w", i, err)
		}
		e.Materials = append(e.Materials, eqMaterialDef)
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
		records, err = token.ReadProperty("BONE", 1)
		if err != nil {
			return fmt.Errorf("bone %d: %w", i, err)
		}
		bone := &MdsBone{}
		bone.Name = records[1]
		records, err = token.ReadProperty("NEXT", 1)
		if err != nil {
			return fmt.Errorf("bone %d next: %w", i, err)
		}
		err = parse(&bone.Next, records[1])
		if err != nil {
			return fmt.Errorf("bone %d next: %w", i, err)
		}

		records, err = token.ReadProperty("CHILDREN", 1)
		if err != nil {
			return fmt.Errorf("bone %d children: %w", i, err)
		}
		err = parse(&bone.ChildrenCount, records[1])
		if err != nil {
			return fmt.Errorf("bone %d children: %w", i, err)
		}

		records, err = token.ReadProperty("CHILDINDEX", 1)
		if err != nil {
			return fmt.Errorf("bone %d childindex: %w", i, err)
		}
		err = parse(&bone.ChildIndex, records[1])
		if err != nil {
			return fmt.Errorf("bone %d childindex: %w", i, err)
		}

		records, err = token.ReadProperty("PIVOT", 3)
		if err != nil {
			return fmt.Errorf("bone %d pivot: %w", i, err)
		}
		err = parse(&bone.Pivot, records[1:]...)
		if err != nil {
			return fmt.Errorf("bone %d pivot: %w", i, err)
		}

		records, err = token.ReadProperty("QUATERNION", 4)
		if err != nil {
			return fmt.Errorf("bone %d quaternion: %w", i, err)
		}
		err = parse(&bone.Quaternion, records[1:]...)
		if err != nil {
			return fmt.Errorf("bone %d quaternion: %w", i, err)
		}

		records, err = token.ReadProperty("SCALE", 3)
		if err != nil {
			return fmt.Errorf("bone %d scale: %w", i, err)
		}
		err = parse(&bone.Scale, records[1:]...)
		if err != nil {
			return fmt.Errorf("bone %d scale: %w", i, err)
		}

		e.Bones = append(e.Bones, bone)
	}

	records, err = token.ReadProperty("NUMMODELS", 1)
	if err != nil {
		return err
	}

	numModels := 0
	err = parse(&numModels, records[1])
	if err != nil {
		return fmt.Errorf("num models: %w", err)
	}

	for i := 0; i < numModels; i++ {
		records, err = token.ReadProperty("MODEL", 1)
		if err != nil {
			return fmt.Errorf("model %d: %w", i, err)
		}
		model := &MdsModel{}
		model.Name = records[1]

		records, err = token.ReadProperty("MAINPIECE", 1)
		if err != nil {
			return fmt.Errorf("model %d mainpiece: %w", i, err)
		}
		err = parse(&model.MainPiece, records[1])
		if err != nil {
			return fmt.Errorf("model %d mainpiece: %w", i, err)
		}

		records, err = token.ReadProperty("NUMVERTICES", 1)
		if err != nil {
			return fmt.Errorf("model %d numvertices: %w", i, err)
		}
		numVertices := 0

		err = parse(&numVertices, records[1])
		if err != nil {
			return fmt.Errorf("model %d numvertices: %w", i, err)
		}

		for j := 0; j < numVertices; j++ {
			records, err = token.ReadProperty("XYZ", 3)
			if err != nil {
				return fmt.Errorf("model %d vertex %d xyz: %w", i, j, err)
			}
			vertex := &MdsVertex{}
			err = parse(&vertex.Position, records[1:]...)
			if err != nil {
				return fmt.Errorf("model %d vertex %d xyz: %w", i, j, err)
			}
			model.Vertices = append(model.Vertices, vertex)
		}

		records, err = token.ReadProperty("NUMUVS", 1)
		if err != nil {
			return fmt.Errorf("model %d numuvs: %w", i, err)
		}
		numUvs := 0
		err = parse(&numUvs, records[1])
		if err != nil {
			return fmt.Errorf("model %d numuvs: %w", i, err)
		}

		if len(model.Vertices) != numUvs {
			return fmt.Errorf("model %d vertices to uv mismatch: %d != %d", i, len(model.Vertices), numUvs)
		}

		for j, vertex := range model.Vertices {

			records, err = token.ReadProperty("UV", 2)
			if err != nil {
				return fmt.Errorf("model %d vertex %d uv: %w", i, j, err)
			}
			err = parse(&vertex.Uv, records[1:]...)
			if err != nil {
				return fmt.Errorf("model %d vertex %d uv: %w", i, j, err)
			}
		}

		records, err = token.ReadProperty("NUMUV2S", 1)
		if err != nil {
			return fmt.Errorf("model %d numuv2s: %w", i, err)
		}
		numUv2s := 0
		err = parse(&numUv2s, records[1])
		if err != nil {
			return fmt.Errorf("model %d numuv2s: %w", i, err)
		}

		if len(model.Vertices) != numUv2s {
			return fmt.Errorf("model %d vertices to uv2 mismatch: %d != %d", i, len(model.Vertices), numUv2s)
		}

		for j, vertex := range model.Vertices {

			records, err = token.ReadProperty("UV2", 2)
			if err != nil {
				return fmt.Errorf("model %d vertex %d uv2: %w", i, j, err)
			}
			err = parse(&vertex.Uv2, records[1:]...)
			if err != nil {
				return fmt.Errorf("model %d vertex %d uv2: %w", i, j, err)
			}
		}

		records, err = token.ReadProperty("NUMNORMALS", 1)
		if err != nil {
			return fmt.Errorf("model %d numnormals: %w", i, err)
		}
		numNormals := 0
		err = parse(&numNormals, records[1])
		if err != nil {
			return fmt.Errorf("model %d numnormals: %w", i, err)
		}

		if len(model.Vertices) != numNormals {
			return fmt.Errorf("model %d vertices to normals mismatch: %d != %d", i, len(model.Vertices), numNormals)
		}

		for j, vertex := range model.Vertices {

			records, err = token.ReadProperty("NORMAL", 3)
			if err != nil {
				return fmt.Errorf("model %d vertex %d normal: %w", i, j, err)
			}
			err = parse(&vertex.Normal, records[1:]...)
			if err != nil {
				return fmt.Errorf("model %d vertex %d normal: %w", i, j, err)
			}
		}

		records, err = token.ReadProperty("NUMTINTS", 1)
		if err != nil {
			return fmt.Errorf("model %d numtints: %w", i, err)
		}
		numTints := 0
		err = parse(&numTints, records[1])
		if err != nil {
			return fmt.Errorf("model %d numtints: %w", i, err)
		}

		if len(model.Vertices) != numTints {
			return fmt.Errorf("model %d vertices to tints mismatch: %d != %d", i, len(model.Vertices), numTints)
		}

		for j, vertex := range model.Vertices {

			records, err = token.ReadProperty("TINT", 4)
			if err != nil {
				return fmt.Errorf("model %d vertex %d tint: %w", i, j, err)
			}
			err = parse(&vertex.Tint, records[1:]...)
			if err != nil {
				return fmt.Errorf("model %d vertex %d tint: %w", i, j, err)
			}

			model.Vertices = append(model.Vertices, vertex)
		}

		records, err = token.ReadProperty("NUMFACES", 1)
		if err != nil {
			return fmt.Errorf("model %d numfaces: %w", i, err)
		}
		numFaces := 0
		err = parse(&numFaces, records[1])
		if err != nil {
			return fmt.Errorf("model %d numfaces: %w", i, err)
		}

		for j := 0; j < numFaces; j++ {
			_, err = token.ReadProperty("FACE", 0)
			if err != nil {
				return fmt.Errorf("model %d face %d: %w", i, j, err)
			}
			face := &MdsFace{}
			records, err = token.ReadProperty("TRIANGLE", 3)
			if err != nil {
				return fmt.Errorf("model %d face %d triangle: %w", i, j, err)
			}
			err = parse(&face.Index, records[1:]...)
			if err != nil {
				return fmt.Errorf("model %d face %d triangle: %w", i, j, err)
			}

			records, err = token.ReadProperty("MATERIAL", 1)
			if err != nil {
				return fmt.Errorf("model %d face %d material: %w", i, j, err)
			}
			face.MaterialName = records[1]

			records, err = token.ReadProperty("FLAGS", 1)
			if err != nil {
				return fmt.Errorf("model %d face %d flags: %w", i, j, err)
			}
			err = parse(&face.Flags, records[1])
			if err != nil {
				return fmt.Errorf("model %d face %d flags: %w", i, j, err)
			}

			model.Faces = append(model.Faces, face)
		}

		records, err = token.ReadProperty("NUMWEIGHTS", 1)
		if err != nil {
			return fmt.Errorf("model %d numbonesassignments: %w", i, err)
		}
		numBoneAssignments := 0
		err = parse(&numBoneAssignments, records[1])
		if err != nil {
			return fmt.Errorf("model %d numbonesassignments: %w", i, err)
		}

		for j := 0; j < numBoneAssignments; j++ {

			weights := [4]*MdsBoneWeight{}
			records, err = token.ReadProperty("WEIGHT", 8)
			if err != nil {
				return fmt.Errorf("model %d boneassignment %d: %w", i, j, err)
			}

			for k := 0; k < 4; k++ {
				var val1 int32
				var val2 float32
				err = parse(&val1, records[1+k*2])
				if err != nil {
					return fmt.Errorf("model %d boneassignment %d: %w", i, j, err)
				}
				err = parse(&val2, records[2+k*2])
				if err != nil {
					return fmt.Errorf("model %d boneassignment %d: %w", i, j, err)
				}
				weights[k] = &MdsBoneWeight{
					BoneIndex: val1,
					Value:     val2,
				}
			}

			model.BoneAssignments = append(model.BoneAssignments, weights)
		}

		e.Models = append(e.Models, model)
	}

	return nil
}

func (e *MdsDef) ToRaw(wce *Wce, dst *raw.Mds) error {
	dst.Version = e.Version

	for _, material := range e.Materials {
		rawMaterial := &raw.Material{}
		err := material.ToRaw(wce, rawMaterial)
		if err != nil {
			return fmt.Errorf("material %s: %w", material.Tag, err)
		}
		dst.Materials = append(dst.Materials, rawMaterial)
	}

	for _, bone := range e.Bones {
		rawBone := &raw.Bone{
			Name:          bone.Name,
			Next:          bone.Next,
			ChildrenCount: bone.ChildrenCount,
			ChildIndex:    bone.ChildIndex,
			Pivot:         bone.Pivot,
			Quaternion:    bone.Quaternion,
			Scale:         bone.Scale,
		}
		dst.Bones = append(dst.Bones, rawBone)
	}

	for _, model := range e.Models {
		rawModel := &raw.MdsModel{
			MainPiece: model.MainPiece,
			Name:      model.Name,
		}
		for _, vert := range model.Vertices {
			rawVertex := &raw.Vertex{
				Position: vert.Position,
				Normal:   vert.Normal,
				Tint:     vert.Tint,
				Uv:       vert.Uv,
				Uv2:      vert.Uv2,
			}
			rawModel.Vertices = append(rawModel.Vertices, rawVertex)
		}
		for _, face := range model.Faces {
			rawFace := &raw.Face{
				Index:        face.Index,
				MaterialName: face.MaterialName,
				Flags:        face.Flags,
			}
			rawModel.Faces = append(rawModel.Faces, rawFace)
		}
		for _, weights := range model.BoneAssignments {
			rawWeights := [4]*raw.MdsBoneWeight{}
			for i := 0; i < 4; i++ {
				rawWeights[i] = &raw.MdsBoneWeight{}
				if i >= len(weights) {
					continue
				}
				weight := weights[i]
				rawWeights[i].BoneIndex = weight.BoneIndex
				rawWeights[i].Value = weight.Value
			}

			rawModel.BoneAssignments = append(rawModel.BoneAssignments, rawWeights)
		}
		dst.Models = append(dst.Models, rawModel)
	}

	return nil
}

func (e *MdsDef) FromRaw(wce *Wce, src *raw.Mds) error {
	folder := strings.TrimSuffix(strings.ToLower(wce.FileName), ".eqg")
	e.folders = append(e.folders, folder)
	e.Tag = string(src.FileName())
	e.Version = src.Version
	for _, mat := range src.Materials {
		eqMaterialDef := &EQMaterialDef{}
		err := eqMaterialDef.FromRaw(wce, mat)
		if err != nil {
			return fmt.Errorf("material %s: %w", mat.Name, err)
		}
	}

	for _, bone := range src.Bones {
		mdsBone := &MdsBone{
			Name:          bone.Name,
			Next:          bone.Next,
			ChildrenCount: bone.ChildrenCount,
			ChildIndex:    bone.ChildIndex,
			Pivot:         bone.Pivot,
			Quaternion:    bone.Quaternion,
			Scale:         bone.Scale,
		}
		e.Bones = append(e.Bones, mdsBone)
	}

	for _, model := range src.Models {
		mdsModel := &MdsModel{
			MainPiece: model.MainPiece,
			Name:      model.Name,
		}
		for _, vert := range model.Vertices {
			mdsVertex := &MdsVertex{
				Position: vert.Position,
				Normal:   vert.Normal,
				Tint:     vert.Tint,
				Uv:       vert.Uv,
				Uv2:      vert.Uv2,
			}
			mdsModel.Vertices = append(mdsModel.Vertices, mdsVertex)
		}
		for _, face := range model.Faces {
			mdsFace := &MdsFace{
				Index:        face.Index,
				MaterialName: face.MaterialName,
				Flags:        face.Flags,
			}
			mdsModel.Faces = append(mdsModel.Faces, mdsFace)
		}
		for _, weights := range model.BoneAssignments {
			mdsWeights := [4]*MdsBoneWeight{}
			for i := 0; i < 4; i++ {
				mdsWeight := &MdsBoneWeight{}
				if i >= len(weights) {
					continue
				}
				weight := weights[i]
				mdsWeight.BoneIndex = weight.BoneIndex
				mdsWeight.Value = weight.Value
				mdsWeights[i] = mdsWeight
			}
			mdsModel.BoneAssignments = append(mdsModel.BoneAssignments, mdsWeights)
		}
		e.Models = append(e.Models, mdsModel)
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
	Faces    []*ModFace
}

type ModFace struct {
	Index        [3]uint32
	MaterialName string
	Passable     int
	Transparent  int
	Collision    int
	Culled       int
	Degenerate   int
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
		fmt.Fprintf(w, "\tNUMUVS %d\n", len(e.UVs))
		for _, u := range e.UVs {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}
		fmt.Fprintf(w, "\tNUMUV2S %d\n", len(e.UV2s))
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
			fmt.Fprintf(w, "\t\t\tPASSABLE %d\n", face.Passable)
			fmt.Fprintf(w, "\t\t\tTRANSPARENT %d\n", face.Transparent)
			fmt.Fprintf(w, "\t\t\tCOLLISIONREQUIRED %d\n", face.Collision)
			fmt.Fprintf(w, "\t\t\tCULLED %d\n", face.Culled)
			fmt.Fprintf(w, "\t\t\tDEGENERATE %d\n", face.Degenerate)
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
		face := &ModFace{}
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

		records, err = token.ReadProperty("PASSABLE", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Passable, records[1])
		if err != nil {
			return fmt.Errorf("passable %d: %w", i, err)
		}

		records, err = token.ReadProperty("TRANSPARENT", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Transparent, records[1])
		if err != nil {
			return fmt.Errorf("transparent %d: %w", i, err)
		}

		records, err = token.ReadProperty("COLLISIONREQUIRED", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Collision, records[1])
		if err != nil {
			return fmt.Errorf("collision %d: %w", i, err)
		}

		records, err = token.ReadProperty("CULLED", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Culled, records[1])
		if err != nil {
			return fmt.Errorf("culled %d: %w", i, err)
		}

		records, err = token.ReadProperty("DEGENERATE", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Degenerate, records[1])
		if err != nil {
			return fmt.Errorf("degenerate %d: %w", i, err)
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
		ModFace := &ModFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flags&uint32(raw.ModFaceFlagPassable) != 0 {
			ModFace.Passable = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagTransparent) != 0 {
			ModFace.Transparent = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagCollisionRequired) != 0 {
			ModFace.Collision = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagCulled) != 0 {
			ModFace.Culled = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagDegenerate) != 0 {
			ModFace.Degenerate = 1
		}

		e.Faces = append(e.Faces, ModFace)
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
	Faces    []*ModFace
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
		fmt.Fprintf(w, "\tNUMUVS %d\n", len(e.UVs))
		for _, u := range e.UVs {
			fmt.Fprintf(w, "\t\tUV %0.8e %0.8e\n", u[0], u[1])
		}
		fmt.Fprintf(w, "\tNUMUV2S %d\n", len(e.UV2s))
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
			fmt.Fprintf(w, "\t\t\tPASSABLE %d\n", face.Passable)
			fmt.Fprintf(w, "\t\t\tTRANSPARENT %d\n", face.Transparent)
			fmt.Fprintf(w, "\t\t\tCOLLISIONREQUIRED %d\n", face.Collision)
			fmt.Fprintf(w, "\t\t\tCULLED %d\n", face.Culled)
			fmt.Fprintf(w, "\t\t\tDEGENERATE %d\n", face.Degenerate)

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
		face := &ModFace{}
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

		records, err = token.ReadProperty("PASSABLE", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Passable, records[1])
		if err != nil {
			return fmt.Errorf("passable %d: %w", i, err)
		}

		records, err = token.ReadProperty("TRANSPARENT", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Transparent, records[1])
		if err != nil {
			return fmt.Errorf("transparent %d: %w", i, err)
		}

		records, err = token.ReadProperty("COLLISIONREQUIRED", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Collision, records[1])
		if err != nil {
			return fmt.Errorf("collision %d: %w", i, err)
		}

		records, err = token.ReadProperty("CULLED", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Culled, records[1])
		if err != nil {
			return fmt.Errorf("culled %d: %w", i, err)
		}

		records, err = token.ReadProperty("DEGENERATE", 1)
		if err != nil {
			return err
		}
		err = parse(&face.Degenerate, records[1])
		if err != nil {
			return fmt.Errorf("degenerate %d: %w", i, err)
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
		ModFace := &ModFace{
			MaterialName: string(face.MaterialName),
			Index:        face.Index,
		}
		if face.Flags&uint32(raw.ModFaceFlagPassable) != 0 {
			ModFace.Passable = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagTransparent) != 0 {
			ModFace.Transparent = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagCollisionRequired) != 0 {
			ModFace.Collision = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagCulled) != 0 {
			ModFace.Culled = 1
		}
		if face.Flags&uint32(raw.ModFaceFlagDegenerate) != 0 {
			ModFace.Degenerate = 1
		}
		e.Faces = append(e.Faces, ModFace)
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

			records, err = token.ReadProperty("MILLISECONDS", 1)
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
