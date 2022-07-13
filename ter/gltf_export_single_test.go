package ter

import (
	"fmt"
	"os"
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/xackery/quail/dump"
)

func TestGLTFExportBroodlands(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	zone := "broodlands"
	path := fmt.Sprintf("test/eq/_%s.eqg", zone)
	inFile := fmt.Sprintf("test/eq/_%s.eqg/ter_%s.ter", zone, zone)
	outFile := fmt.Sprintf("test/eq/%s.gltf", zone)

	e, err := New("arena", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	fw, err := os.Create(fmt.Sprintf("test/%s.txt", zone))
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fw.Close()
	fmt.Fprintf(fw, "faces:\n")
	for i, o := range e.faces {
		fmt.Fprintf(fw, "%d %+v\n", i, o)
	}

	fmt.Fprintf(fw, "vertices:\n")
	for i, o := range e.vertices {
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFExportCityOfBronze(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	zone := "cityofbronze"
	path := fmt.Sprintf("test/eq/_%s.eqg", zone)
	inFile := fmt.Sprintf("test/eq/_%s.eqg/ter_%s.ter", zone, zone)
	outFile := fmt.Sprintf("test/eq/%s.gltf", zone)
	isDumpEnabed := false

	if isDumpEnabed {
		d, err := dump.New(path)
		if err != nil {
			t.Fatalf("dump.new: %s", err)
		}
		defer d.Save(fmt.Sprintf("%s.png", inFile))
	}
	e, err := New("cityofbronze", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	fw, err := os.Create(fmt.Sprintf("test/%s.txt", zone))
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fw.Close()
	fmt.Fprintf(fw, "faces:\n")
	for i, o := range e.faces {
		fmt.Fprintf(fw, "%d %+v\n", i, o)
	}

	fmt.Fprintf(fw, "vertices:\n")
	for i, o := range e.vertices {
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTF(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error

	outFile := "test/cube2.gltf"

	doc := &gltf.Document{}
	attrs, _ := modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position: [][3]float32{{0, 0, 0}, {0, 10, 0}, {0, 0, 10}},
		Color:    [][3]uint8{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}},
	})
	doc.Meshes = []*gltf.Mesh{
		{Extras: 8.0, Name: "mesh_1", Weights: []float32{1.2, 2}},
		{Extras: 8.0, Name: "mesh_2", Primitives: []*gltf.Primitive{
			{Extras: 8.0, Attributes: gltf.Attribute{gltf.POSITION: 1}, Indices: gltf.Index(2), Material: gltf.Index(1), Mode: gltf.PrimitiveLines},
			{Extras: 8.0, Targets: []gltf.Attribute{{gltf.POSITION: 1, "THEN": 4}, {"OTHER": 2}}, Indices: gltf.Index(2), Material: gltf.Index(1), Mode: gltf.PrimitiveLines, Attributes: attrs},
		}}}

	err = gltf.Save(doc, outFile)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFBin(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	outFile := "test/example.gtlf"
	doc := &gltf.Document{
		Accessors: []*gltf.Accessor{
			{BufferView: gltf.Index(0), ComponentType: gltf.ComponentUshort, Count: 36, Type: gltf.AccessorScalar},
			{BufferView: gltf.Index(1), ComponentType: gltf.ComponentFloat, Count: 24, Max: []float32{0.5, 0.5, 0.5}, Min: []float32{-0.5, -0.5, -0.5}, Type: gltf.AccessorVec3},
			{BufferView: gltf.Index(2), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec3},
			{BufferView: gltf.Index(3), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec4},
			{BufferView: gltf.Index(4), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec2},
		},
		Asset: gltf.Asset{Version: "2.0", Generator: "FBX2glTF"},
		BufferViews: []*gltf.BufferView{
			{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: gltf.TargetElementArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 72, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 360, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 384, ByteOffset: 648, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 192, ByteOffset: 1032, Target: gltf.TargetArrayBuffer},
		},
		Buffers: []*gltf.Buffer{{ByteLength: 1224, Data: []byte{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}}},
		Materials: []*gltf.Material{{
			Name: "Default", AlphaMode: gltf.AlphaOpaque, AlphaCutoff: gltf.Float(0.5),
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{BaseColorFactor: &[4]float32{0.8, 0.8, 0.8, 1}, MetallicFactor: gltf.Float(0.1), RoughnessFactor: gltf.Float(0.99)},
		}},
		Meshes: []*gltf.Mesh{{Name: "Cube", Primitives: []*gltf.Primitive{{Indices: gltf.Index(0), Material: gltf.Index(0), Mode: gltf.PrimitiveTriangles, Attributes: map[string]uint32{gltf.POSITION: 1, gltf.COLOR_0: 3, gltf.NORMAL: 2, gltf.TEXCOORD_0: 4}}}}},
		Nodes: []*gltf.Node{
			{Name: "RootNode", Children: []uint32{1, 2, 3}},
			{Name: "Mesh"},
			{Name: "Cube", Mesh: gltf.Index(0)},
			{Name: "Texture Group"},
		},
		Samplers: []*gltf.Sampler{{WrapS: gltf.WrapRepeat, WrapT: gltf.WrapRepeat}},
		Scene:    gltf.Index(0),
		Scenes:   []*gltf.Scene{{Name: "Root Scene", Nodes: []uint32{0}}},
	}
	if err := gltf.SaveBinary(doc, outFile); err != nil {
		panic(err)
	}
}

func TestFoo(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	outFile := "test/foo.glb"

	doc := gltf.NewDocument()
	attrs, _ := modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position:       [][3]float32{{1, 2, 3}, {0, 0, -1}},
		Normal:         [][3]float32{{1, 2, 3}, {0, 0, -1}},
		Tangent:        [][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}},
		TextureCoord_0: [][2]uint8{{0, 255}, {255, 0}},
		TextureCoord_1: [][2]float32{{1, 2}, {1, 2}},
		Joints:         [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}},
		Weights:        [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}},
		Color:          [][3]uint8{{255, 255, 255}, {0, 255, 0}},
		CustomAttributes: []modeler.CustomAttribute{
			{Name: "COLOR_1", Data: [][3]uint8{{0, 0, 255}, {100, 200, 0}}},
			{Name: "COLOR_2", Data: [][4]uint8{{23, 58, 188, 1}, {0, 155, 0, 0}}},
		},
	})
	indicesAccessor := modeler.WriteIndices(doc, []uint16{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	doc.Meshes = []*gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []*gltf.Primitive{
			{
				Indices:    gltf.Index(indicesAccessor),
				Attributes: attrs,
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	if err := gltf.SaveBinary(doc, outFile); err != nil {
		panic(err)
	}
}

func TestGLTFSaveExample(t *testing.T) {
	outFile := "test/tmp.gltf"

	doc := &gltf.Document{
		Accessors: []*gltf.Accessor{
			{BufferView: gltf.Index(0), ComponentType: gltf.ComponentUshort, Count: 36, Type: gltf.AccessorScalar},
			{BufferView: gltf.Index(1), ComponentType: gltf.ComponentFloat, Count: 24, Max: []float32{0.5, 0.5, 0.5}, Min: []float32{-0.5, -0.5, -0.5}, Type: gltf.AccessorVec3},
			{BufferView: gltf.Index(2), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec3},
			{BufferView: gltf.Index(3), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec4},
			{BufferView: gltf.Index(4), ComponentType: gltf.ComponentFloat, Count: 24, Type: gltf.AccessorVec2},
		},
		Asset: gltf.Asset{Version: "2.0", Generator: "FBX2glTF"},
		BufferViews: []*gltf.BufferView{
			{Buffer: 0, ByteLength: 72, ByteOffset: 0, Target: gltf.TargetElementArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 72, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 288, ByteOffset: 360, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 384, ByteOffset: 648, Target: gltf.TargetArrayBuffer},
			{Buffer: 0, ByteLength: 192, ByteOffset: 1032, Target: gltf.TargetArrayBuffer},
		},
		Buffers: []*gltf.Buffer{{ByteLength: 1224, Data: []byte{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}}},
		Materials: []*gltf.Material{{
			Name: "Default", AlphaMode: gltf.AlphaOpaque, AlphaCutoff: gltf.Float(0.5),
			PBRMetallicRoughness: &gltf.PBRMetallicRoughness{BaseColorFactor: &[4]float32{0.8, 0.8, 0.8, 1}, MetallicFactor: gltf.Float(0.1), RoughnessFactor: gltf.Float(0.99)},
		}},
		Meshes: []*gltf.Mesh{{Name: "Cube", Primitives: []*gltf.Primitive{{Indices: gltf.Index(0), Material: gltf.Index(0), Mode: gltf.PrimitiveTriangles, Attributes: map[string]uint32{gltf.POSITION: 1, gltf.COLOR_0: 3, gltf.NORMAL: 2, gltf.TEXCOORD_0: 4}}}}},
		Nodes: []*gltf.Node{
			{Name: "RootNode", Children: []uint32{1, 2, 3}},
			{Name: "Mesh"},
			{Name: "Cube", Mesh: gltf.Index(0)},
			{Name: "Texture Group"},
		},
		Samplers: []*gltf.Sampler{{WrapS: gltf.WrapRepeat, WrapT: gltf.WrapRepeat}},
		Scene:    gltf.Index(0),
		Scenes:   []*gltf.Scene{{Name: "Root Scene", Nodes: []uint32{0}}},
	}
	for _, buff := range doc.Buffers {
		buff.EmbeddedResource()
	}
	err := gltf.Save(doc, outFile)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFInterleaved(t *testing.T) {
	outFile := "test/tmp.gltf"

	doc := gltf.NewDocument()
	attrs, _ := modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position:       [][3]float32{{1, 2, 3}, {0, 0, -1}},
		Normal:         [][3]float32{{1, 2, 3}, {0, 0, -1}},
		Tangent:        [][4]float32{{1, 2, 3, 4}, {1, 2, 3, 4}},
		TextureCoord_0: [][2]uint8{{0, 255}, {255, 0}},
		TextureCoord_1: [][2]float32{{1, 2}, {1, 2}},
		Joints:         [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}},
		Weights:        [][4]uint8{{1, 2, 3, 4}, {1, 2, 3, 4}},
		Color:          [][3]uint8{{255, 255, 255}, {0, 255, 0}},
		CustomAttributes: []modeler.CustomAttribute{
			{Name: "COLOR_1", Data: [][3]uint8{{0, 0, 255}, {100, 200, 0}}},
			{Name: "COLOR_2", Data: [][4]uint8{{23, 58, 188, 1}, {0, 155, 0, 0}}},
		},
	})
	indicesAccessor := modeler.WriteIndices(doc, []uint16{0, 1, 2, 3, 1, 0, 0, 2, 3, 1, 4, 2, 4, 3, 2, 4, 1, 3})
	doc.Meshes = []*gltf.Mesh{{
		Name: "Pyramid",
		Primitives: []*gltf.Primitive{
			{
				Indices:    gltf.Index(indicesAccessor),
				Attributes: attrs,
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: "Root", Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)
	for _, buff := range doc.Buffers {
		buff.EmbeddedResource()
	}
	err := gltf.Save(doc, outFile)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
