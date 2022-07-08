package ter

import (
	"fmt"
	"io"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

// ExportGLTF exports a provided ter file to gltf format
func (e *TER) ExportGLTF(w io.Writer) error {
	var err error
	doc := gltf.NewDocument()

	for _, mat := range e.materials {
		doc.Materials = append(doc.Materials, &gltf.Material{
			Name: mat.Name,
		})
	}

	mesh := &gltf.Mesh{
		Name: e.name,
	}

	prim := &gltf.Primitive{
		Mode: gltf.PrimitiveTriangles,
	}
	mesh.Primitives = append(mesh.Primitives, prim)

	positions := [][3]float32{}
	normals := [][3]float32{}
	uvs := [][2]float32{}
	indices := []uint16{}

	for _, vert := range e.vertices {
		positions = append(positions, [3]float32{vert.Position.X, vert.Position.Y, vert.Position.Z})
		normals = append(normals, [3]float32{vert.Normal.X, vert.Normal.Y, vert.Normal.Z})
		uvs = append(uvs, [2]float32{vert.Uv.X, vert.Uv.Y})
	}
	for _, o := range e.triangles {
		indices = append(indices, uint16(o.Index.X))
		indices = append(indices, uint16(o.Index.Y))
		indices = append(indices, uint16(o.Index.Z))
	}

	prim.Attributes, err = modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position:       positions,
		Normal:         normals,
		TextureCoord_0: uvs,
	})
	if err != nil {
		return fmt.Errorf("writeAttributes: %w", err)
	}
	prim.Indices = gltf.Index(modeler.WriteIndices(doc, indices))
	doc.Meshes = append(doc.Meshes, mesh)

	for _, buff := range doc.Buffers {
		buff.EmbeddedResource()
	}
	enc := gltf.NewEncoder(w)
	enc.AsBinary = false
	err = enc.Encode(doc)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
