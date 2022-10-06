package wld

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"

	"github.com/malashin/dds"
	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/quail/common"
)

func (e *WLD) parseMesh(frag *fragmentInfo) error {
	type mesher interface {
		Vertices() []*common.Vertex
		Triangles() []*common.Triangle
	}

	meshFragment, ok := frag.data.(mesher)
	if !ok {
		return nil
	}

	mesh := &mesh{
		triangles: meshFragment.Triangles(),
		vertices:  meshFragment.Vertices(),
	}
	e.meshes = append(e.meshes, mesh)
	return nil
}

func (e *WLD) parseMaterial(frag *fragmentInfo) error {
	type materialer interface {
		Name() string
		ShaderType() int
		MaterialType() int
	}

	material, ok := frag.data.(materialer)
	if !ok {
		return nil
	}

	baseName := strings.TrimSuffix(strings.ToLower(material.Name()), "_mdf")
	if strings.HasPrefix(baseName, "m000") {
		fmt.Println("skipping model material", baseName)
		return nil
	}
	baseName = strings.TrimSuffix(baseName, ".bmp")
	outImageName := baseName + ".png"
	inExt := ".bmp"

	err := e.MaterialAdd(material.Name(), fmt.Sprintf("%d", material.ShaderType()))
	if err != nil {
		return fmt.Errorf("materialadd: %w", err)
	}
	err = e.MaterialPropertyAdd(material.Name(), "e_texturediffuse0", 2, outImageName)
	if err != nil {
		return fmt.Errorf("materialPropertyAdd %s: %w", outImageName, err)
	}

	data, err := e.archive.File(baseName + inExt)
	if err != nil {
		inExt = ".dds"
		data, err = e.archive.File(baseName + inExt)
		if err != nil {
			//TODO: dds bmp fix
			//return fmt.Errorf("material '%s' not found in archive", inImageName)
			//fmt.Printf("material '%s' with .dds and .bmp ext not found in archive\n", baseName)
			return nil
		}
	}

	buf := bytes.NewBuffer(data)
	img, err := dds.Decode(buf)
	if err != nil {
		var err2 error
		img, err2 = bmp.Decode(buf)
		if err2 != nil {
			return fmt.Errorf("bmp (dds) decode %s: %s, (bmp) decode: %w", baseName+inExt, err.Error(), err2)
		}
	}

	err = png.Encode(buf, img)
	if err != nil {
		return fmt.Errorf("png encode %s: %w", baseName+inExt, err)
	}

	err = e.archive.WriteFile(outImageName, data)
	if err != nil {
		return fmt.Errorf("writeFile %s: %w", outImageName, err)
	}

	return nil
}
