package quail

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/raw"
)

// Read takes a raw type and converts it to a quail type
func (q *Quail) RawRead(in raw.ReadWriter) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	switch val := in.(type) {
	case *raw.Lay:
		return q.layRead(val)
	case *raw.Ani:
		return q.aniRead(val)
	case *raw.Wld:
		return q.wldRead(val, in.FileName())
	case *raw.Dds:
		return q.ddsRead(val)
	case *raw.Bmp:
		return q.bmpRead(val)
	case *raw.Png:
		return q.pngRead(val)
	case *raw.Mod:
		return q.modRead(val)
	case *raw.Pts:
		return nil
	case *raw.Prt:
		return nil
	case *raw.Mds:
		return q.mdsRead(val)
	case *raw.Unk:
		return nil
	case *raw.Txt:
		return nil
	default:
		return fmt.Errorf("unknown type %T", val)
	}
}

func RawRead(in raw.ReadWriter, q *Quail) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	return q.RawRead(in)
}

func (q *Quail) aniRead(in *raw.Ani) error {
	if q.Header == nil {
		q.Header = &common.Header{}
	}
	q.Header.Version = int(in.Version)
	q.Header.Name = "animation"

	return nil
}

func (q *Quail) layRead(in *raw.Lay) error {
	if q.Header == nil {
		q.Header = &common.Header{}
	}
	q.Header.Version = int(in.Version)
	q.Header.Name = "layer"
	for _, model := range q.Models {
		for _, entry := range in.Entries {
			lay := &common.Layer{
				Material: entry.Material,
				Diffuse:  entry.Diffuse,
				Normal:   entry.Normal,
			}

			model.Layers = append(model.Layers, lay)
		}
	}
	return nil
}

func (q *Quail) ddsRead(in *raw.Dds) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write dds: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) bmpRead(in *raw.Bmp) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write bmp: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) pngRead(in *raw.Png) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write png: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) modRead(in *raw.Mod) error {
	model, err := q.modConvertMesh(in)
	if err != nil {
		return fmt.Errorf("modConvertMesh: %w", err)
	}
	if model != nil {
		q.Models = append(q.Models, model)
	}
	return nil
}

func (q *Quail) mdsRead(in *raw.Mds) error {
	model, err := q.mdsConvertMesh(in)
	if err != nil {
		return fmt.Errorf("mdsConvertMesh: %w", err)
	}
	if model != nil {
		q.Models = append(q.Models, model)
	}

	return nil
}

func (q *Quail) modConvertMesh(in *raw.Mod) (*common.Model, error) {
	if in == nil {
		return nil, fmt.Errorf("mod is nil")
	}
	model := common.NewModel(in.FileName())
	model.FileType = "mod"

	for _, triangle := range in.Triangles {
		model.Triangles = append(model.Triangles, common.Triangle{
			Index: common.UIndex3{
				X: uint32(triangle.Index.X),
				Y: uint32(triangle.Index.Y),
				Z: uint32(triangle.Index.Z),
			},
			MaterialName: triangle.MaterialName,
			Flag:         uint32(triangle.Flag),
		})
	}
	for _, vertex := range in.Vertices {
		model.Vertices = append(model.Vertices, common.Vertex{
			Position: common.Vector3{
				X: vertex.Position.X,
				Y: vertex.Position.Y,
				Z: vertex.Position.Z,
			},
			Uv: common.Vector2{
				X: vertex.Uv.X,
				Y: vertex.Uv.Y,
			},
			Tint: common.RGBA{
				R: vertex.Tint.R,
				G: vertex.Tint.G,
				B: vertex.Tint.B,
				A: vertex.Tint.A,
			},
		})
	}
	for _, material := range in.Materials {
		dstMaterial := &common.Material{}
		dstMaterial.Name = material.Name
		dstMaterial.Flag = material.Flag
		dstMaterial.ShaderName = material.ShaderName
		for _, property := range material.Properties {
			dstProperty := &common.MaterialProperty{}
			dstProperty.Category = property.Category
			dstProperty.Name = property.Name
			dstProperty.Value = property.Value
			dstMaterial.Properties = append(dstMaterial.Properties, dstProperty)
		}
		model.Materials = append(model.Materials, dstMaterial)
	}

	return model, nil
}

func (q *Quail) mdsConvertMesh(in *raw.Mds) (*common.Model, error) {
	if in == nil {
		return nil, fmt.Errorf("mod is nil")
	}
	model := common.NewModel(in.FileName())
	model.FileType = "mds"
	for _, triangle := range in.Triangles {
		model.Triangles = append(model.Triangles, common.Triangle{
			Index: common.UIndex3{
				X: uint32(triangle.Index.X),
				Y: uint32(triangle.Index.Y),
				Z: uint32(triangle.Index.Z),
			},
			MaterialName: triangle.MaterialName,
			Flag:         uint32(triangle.Flag),
		})
	}
	for _, vertex := range in.Vertices {
		model.Vertices = append(model.Vertices, common.Vertex{
			Position: common.Vector3{
				X: vertex.Position.X,
				Y: vertex.Position.Y,
				Z: vertex.Position.Z,
			},
			Uv: common.Vector2{
				X: vertex.Uv.X,
				Y: vertex.Uv.Y,
			},
			Tint: common.RGBA{
				R: vertex.Tint.R,
				G: vertex.Tint.G,
				B: vertex.Tint.B,
				A: vertex.Tint.A,
			},
		})
	}
	for _, material := range in.Materials {
		dstMaterial := &common.Material{}
		dstMaterial.Name = material.Name
		dstMaterial.Flag = material.Flag
		dstMaterial.ShaderName = material.ShaderName
		for _, property := range material.Properties {
			dstProperty := &common.MaterialProperty{}
			dstProperty.Category = property.Category
			dstProperty.Name = property.Name
			dstProperty.Value = property.Value
			dstMaterial.Properties = append(dstMaterial.Properties, dstProperty)
		}
		model.Materials = append(model.Materials, dstMaterial)
	}
	return model, nil
}
