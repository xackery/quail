package mod

import (
	"bytes"
	"fmt"
	"image/png"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/spate/glimage/dds"
	"github.com/xackery/quail/common"
)

func TestGLTFExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	err := os.Mkdir("test", fs.ModeDir)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("mkdir test: %s", err)
	}

	e, err := New("obj_gears.mod", "test/")
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open("test/obj_gears.mod")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer r.Close()
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load %s", err)
	}

	w, err := os.Create("test/obj_gears.gltf")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}

func TestCube(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	doc := gltf.NewDocument()
	doc.Scenes[0] = &gltf.Scene{Name: "cube"}

	buf := &bytes.Buffer{}
	r, err := os.Open("test/metal_rustyb.dds")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer r.Close()

	textureDiffuseName := "cube.dds"

	if filepath.Ext(textureDiffuseName) == ".dds" {
		img, err := dds.Decode(r)
		if err != nil {
			t.Fatalf("dds.Decode %s: %s", textureDiffuseName, err)
		}
		buf = bytes.NewBuffer(nil)
		err = png.Encode(buf, img)
		if err != nil {
			t.Fatalf("png.Encode %s: %s", textureDiffuseName, err)
		}
		textureDiffuseName = strings.ReplaceAll(textureDiffuseName, ".dds", ".png")
	}

	imageIdx, err := modeler.WriteImage(doc, textureDiffuseName, "image/png", buf)
	if err != nil {
		t.Fatalf("writeImage to gtlf: %s", err)
	}
	doc.Textures = append(doc.Textures, &gltf.Texture{Source: gltf.Index(imageIdx)})

	doc.Materials = append(doc.Materials, &gltf.Material{
		Name: "cube",
		PBRMetallicRoughness: &gltf.PBRMetallicRoughness{
			BaseColorTexture: &gltf.TextureInfo{
				Index: uint32(len(doc.Textures) - 1),
			},
			MetallicFactor: gltf.Float(0),
		},
	})

	mesh := &gltf.Mesh{
		Name: textureDiffuseName,
	}

	prim := &gltf.Primitive{
		Mode:     gltf.PrimitiveTriangles,
		Material: gltf.Index(uint32(len(doc.Materials) - 1)),
	}
	mesh.Primitives = append(mesh.Primitives, prim)

	positions := [][3]float32{}
	normals := [][3]float32{
		{0, 0, 1},
		{1, 0, 0},
		{0, 0, -1},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
	}
	uvs := [][2]float32{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}
	indices := []uint16{
		0, 1, 3, 3, 1, 2,
		1, 5, 2, 2, 5, 6,
		5, 4, 6, 6, 4, 7,
		4, 0, 7, 7, 0, 3,
		3, 2, 7, 7, 2, 6,
		4, 5, 0, 0, 5, 1,
	}
	vertices := [][3]float32{
		{-1, -1, -1},
		{1, -1, -1},
		{1, 1, -1},
		{-1, 1, -1},
		{-1, -1, 1},
		{1, -1, 1},
		{1, 1, 1},
		{-1, 1, 1},
	}

	for _, vert := range vertices {
		positions = append(positions, vert)
		//normals = append(normals, normals[i])
		//uvs = append(uvs, [2]float32{0, 0})
	}
	/*for _, o := range e.triangles {
		indices = append(indices, uint16(o.Index.X))
		indices = append(indices, uint16(o.Index.Y))
		indices = append(indices, uint16(o.Index.Z))
	}*/

	fmt.Println(len(normals), len(positions), len(uvs))
	prim.Attributes, err = modeler.WriteAttributesInterleaved(doc, modeler.Attributes{
		Position:       positions,
		Normal:         normals,
		TextureCoord_0: uvs,
	})
	if err != nil {
		t.Fatalf("writeAttributes: %s", err)
	}
	prim.Indices = gltf.Index(modeler.WriteIndices(doc, indices))
	doc.Meshes = append(doc.Meshes, mesh)
	doc.Nodes = append(doc.Nodes, &gltf.Node{Name: textureDiffuseName, Mesh: gltf.Index(uint32(len(doc.Meshes) - 1))})
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, uint32(len(doc.Nodes)-1))

	for _, buff := range doc.Buffers {
		buff.EmbeddedResource()
	}

	w, err := os.Create("test/cube.gltf")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	enc := gltf.NewEncoder(w)
	enc.AsBinary = false
	err = enc.Encode(doc)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
}

func TestTriangle(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/"
	inFile := "test/triangle.gltf"
	outFile := "test/triangle_out.gltf"

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	e.materials = append(e.materials, &common.Material{Name: "metal_rustyb.dds", Properties: common.Properties{{Name: "e_texturediffuse0", Value: "metal_rustyb.dds", Category: 2}}})
	data, err := ioutil.ReadFile("test/metal_rustyb.dds")
	if err != nil {
		t.Fatalf("%s", err)
	}
	fe, err := common.NewFileEntry("metal_rustyb.dds", data)
	if err != nil {
		t.Fatalf("NewFileEntry: %s", err)
	}
	e.files = append(e.files, fe)
	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.GLTFExport(w)
	//err = e.Save(w)
	if err != nil {
		t.Fatalf("gltfExport: %s", err)
	}
}
