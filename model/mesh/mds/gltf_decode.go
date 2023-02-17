package mds

import (
	"fmt"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// GLTFDecode imports a GLTF document
func (e *MDS) GLTFDecode(doc *gltf.Document) error {
	var err error
	e.skin = nil

	for _, m := range doc.Materials {
		err = e.parseMaterial(doc, m)
		if err != nil {
			return fmt.Errorf("parseMaterial %s: %w", m.Name, err)
		}
	}

	for _, n := range doc.Nodes {
		if n.Mesh != nil {
			m := doc.Meshes[*n.Mesh]
			if m == nil {
				return fmt.Errorf("accesing node '%s' mesh '%d' failed", n.Name, *n.Mesh)
			}
			err = e.parseMesh(doc, m, n.Scale)
			if err != nil {
				return fmt.Errorf("parseMesh %s %s: %w", n.Name, m.Name, err)
			}
		}

		if n.Skin != nil {
			s := doc.Skins[*n.Skin]
			err = e.parseSkin(doc, s)
			if err != nil {
				return fmt.Errorf("parseSkin %s: %w", s.Name, err)
			}
		}
	}

	//https://github.com/KhronosGroup/glTF-Tutorials/blob/master/gltfTutorial/gltfTutorial_007_Animations.md
	e.animations = append(e.animations, doc.Animations...)
	return nil
}

func (e *MDS) parseMaterial(doc *gltf.Document, m *gltf.Material) error {
	err := e.MaterialAdd(m.Name, "Opaque_MaxCB1.fx")
	if err != nil {
		return fmt.Errorf("materialAdd: %w", err)
	}

	if m.PBRMetallicRoughness.BaseColorTexture == nil {
		return nil
	}

	image := doc.Images[int(m.PBRMetallicRoughness.BaseColorTexture.Index)]
	if image == nil {
		return fmt.Errorf("expected image for '%s', but not found", m.Name)
	}

	bv := doc.BufferViews[int(*image.BufferView)]
	if bv == nil {
		return fmt.Errorf("texture '%s' expected buffer view %d, but not found", image.Name, *image.BufferView)
	}

	ext := ""
	if strings.ToLower(image.MimeType) == "image/png" {
		ext = ".png"
	}

	if ext == "" {
		return fmt.Errorf("unsupported mimetype in gltf image '%s'", image.Name)
	}

	imageName := strings.ToLower(image.Name)
	if !strings.HasSuffix(imageName, ext) {
		imageName += ext
	}

	data, err := modeler.ReadBufferView(doc, bv)
	if err != nil {
		return fmt.Errorf("readBufferView %d: %w", *image.BufferView, err)
	}

	err = e.archive.WriteFile(imageName, data)
	if err != nil {
		return fmt.Errorf("writeFile: %w", err)
	}

	err = e.MaterialPropertyAdd(m.Name, "e_TextureDiffuse0", 2, imageName)
	if err != nil {
		return fmt.Errorf("materialPropertyAdd %s: %w", imageName, err)
	}
	return nil
}

// parseMesh will parse a gltf mesh and place it in mds style
func (e *MDS) parseMesh(doc *gltf.Document, m *gltf.Mesh, scale [3]float32) error {
	for _, p := range m.Primitives {
		if p.Mode != gltf.PrimitiveTriangles {
			return fmt.Errorf("primitive in mesh '%s' is mode %d, unsupported", m.Name, p.Mode)
		}

		materialName := ""
		if p.Material != nil {
			materialName = doc.Materials[*p.Material].Name
		}

		indices, err := modeler.ReadIndices(doc, doc.Accessors[*p.Indices], []uint32{})
		if err != nil {
			return fmt.Errorf("readIndices: %w", err)
		}

		for i := 0; i < len(indices); i += 3 {
			err = e.FaceAdd([3]uint32{uint32(indices[i]), uint32(indices[i+1]), uint32(indices[i+2])}, materialName, 0)
			if err != nil {
				return fmt.Errorf("faceAdd: %w", err)
			}
		}

		posIndex, ok := p.Attributes[gltf.POSITION]
		if !ok {
			return fmt.Errorf("primitive in mesh '%s' has no position", m.Name)
		}

		positions, err := modeler.ReadPosition(doc, doc.Accessors[posIndex], [][3]float32{})
		if err != nil {
			return fmt.Errorf("readPosition: %w", err)
		}

		joints := [][4]uint16{}
		jointIndex, ok := p.Attributes[gltf.JOINTS_0]
		if ok {
			joints, err = modeler.ReadJoints(doc, doc.Accessors[jointIndex], [][4]uint16{})
			if err != nil {
				return fmt.Errorf("readJoints: %w", err)
			}
		}

		weights := [][4]float32{}
		weightIndex, ok := p.Attributes[gltf.WEIGHTS_0]
		if ok {
			weights, err = modeler.ReadWeights(doc, doc.Accessors[weightIndex], [][4]float32{})
			if err != nil {
				return fmt.Errorf("readJoints: %w", err)
			}
		}

		//fmt.Printf("pos: %+v\n", pos)
		normals := [][3]float32{}
		normalIndex, ok := p.Attributes[gltf.NORMAL]
		if ok {
			normals, err = modeler.ReadNormal(doc, doc.Accessors[normalIndex], [][3]float32{})
			if err != nil {
				return fmt.Errorf("readNormal: %w", err)
			}
		} //return fmt.Errorf("primitive in mesh '%s' has no normal", m.Name)

		tints, err := modeler.ReadColor(doc, doc.Accessors[p.Attributes[gltf.COLOR_0]], [][4]uint8{})
		if err != nil {
			return fmt.Errorf("readTint: %w", err)
		}

		//fmt.Printf("normal: %+v\n", normal)

		uvIndex, ok := p.Attributes[gltf.TEXCOORD_0]
		uvs := [][2]float32{}
		if ok {
			uvs, err = modeler.ReadTextureCoord(doc, doc.Accessors[uvIndex], [][2]float32{})
			if err != nil {
				return fmt.Errorf("readTextureCoord: %w", err)
			}
		}
		//return fmt.Errorf("primitive in mesh '%s' has no texcoord", m.Name)
		//fmt.Printf("uv: %+v\n", uv)

		for i := 0; i < len(positions); i++ {
			vertex := &common.Vertex{}
			positions[i] = helper.ApplyQuaternion(positions[i], [4]float32{-0.5, 0.5, 0.5, -0.5})
			positions[i][0] *= scale[0]
			positions[i][1] *= scale[1]
			positions[i][2] *= scale[2]
			vertex.Position = positions[i]

			vertex.Normal = normals[i]

			uvs[i][0] = uvs[i][0] * scale[0]
			uvs[i][1] = uvs[i][1] * scale[1]
			vertex.Uv = uvs[i]

			vertex.Tint = tints[i]

			if len(joints) > i {
				vertex.Joint = joints[i]
			} else {
				vertex.Joint = [4]uint16{}
			}

			if len(weights) > i {
				vertex.Weight = weights[i]
			} else {
				vertex.Weight = [4]float32{}
			}

			e.vertices = append(e.vertices, vertex)
		}
	}
	return nil
}

func (e *MDS) parseSkin(doc *gltf.Document, s *gltf.Skin) error {
	if e.skin != nil {
		return fmt.Errorf("multiple skins found, only one is supported for mds gltf decode")
	}

	e.skin = &common.Skin{
		Name:   s.Name,
		Joints: make(map[int]*common.Joint),
	}
	matrices, err := modeler.ReadAccessor(doc, doc.Accessors[*s.InverseBindMatrices], nil)
	if err != nil {
		return fmt.Errorf("read inversebindmatrices: %w", err)
	}
	e.skin.InverseBindMatrices = append(e.skin.InverseBindMatrices, matrices.([][4][4]float32)...)

	// maps gltf node lookup to jointCache array
	jointCache := make(map[uint32]*common.Joint)

	if s.Skeleton == nil {
		return fmt.Errorf("skin skeleton is nil")
	}

	jointNode := doc.Nodes[*s.Skeleton]

	joint := &common.Joint{
		Name:        jointNode.Name,
		Translation: jointNode.Translation,
		Children:    parseJointChildren(doc, jointCache, jointNode),
	}

	jointCache[*s.Skeleton] = joint
	e.skin.Joints[0] = joint

	for _, jointIndex := range s.Joints {
		jointNode := doc.Nodes[jointIndex]
		if jointNode == nil {
			return fmt.Errorf("joint node %d not found", jointIndex)
		}

		joint, ok := jointCache[jointIndex]
		if !ok {
			joint = &common.Joint{
				Name:        jointNode.Name,
				Translation: jointNode.Translation,
			}
			jointCache[jointIndex] = joint
		}
		if len(jointNode.Children) != len(joint.Children) {
			joint.Children = parseJointChildren(doc, jointCache, jointNode)
		}

		e.skin.Joints[len(e.skin.Joints)] = joint
	}

	return nil
}

func parseJointChildren(doc *gltf.Document, jointCache map[uint32]*common.Joint, joint *gltf.Node) []*common.Joint {
	joints := []*common.Joint{}
	for _, jointIndex := range joint.Children {
		jointNode := doc.Nodes[jointIndex]
		if jointNode == nil {
			fmt.Printf("warning: joint node %d not found\n", jointIndex)
			continue
		}
		joint, ok := jointCache[jointIndex]
		if !ok {
			joint = &common.Joint{
				Name:        jointNode.Name,
				Translation: jointNode.Translation,
			}
			jointCache[jointIndex] = joint
		}
		joints = append(joints, joint)
		if len(jointNode.Children) != len(joint.Children) {
			joint.Children = parseJointChildren(doc, jointCache, jointNode)
		}
	}
	return joints
}
