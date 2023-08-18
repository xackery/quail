package quail

import (
	"io"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/wld"
	"github.com/xackery/quail/pfs/archive"
	"github.com/xackery/quail/quail/def"
)

// Decode decodes a WLD file
func WLDDecode(r io.ReadSeeker, pfs archive.ReadWriter) ([]*def.Mesh, error) {
	meshes := make([]*def.Mesh, 0)

	e, err := wld.New("test", pfs)
	if err != nil {
		return nil, err
	}

	err = e.Decode(r)
	if err != nil {
		return nil, err
	}

	materials := make([]*def.Material, 0)

	curMesh := &def.Mesh{}
	//names := e.Names()
	for _, f := range e.Fragments {
		switch d := f.(type) {
		case *wld.Mesh:
			curMesh.Name = d.Name

			curMesh.Vertices = d.Vertices
			curMesh.Triangles = d.Triangles

			materialGroup := 0
			materialCounter := 0
			for i := 0; i < len(curMesh.Triangles); i++ {
				curTriangleMat := d.TriangleMaterials[materialGroup]
				if materialCounter < 1 {
					materialCounter = int(curTriangleMat.Count)
				}
				matID := curTriangleMat.MaterialID
				if matID != 0 {
					matID--
					log.Debugf("mesh %s materialID %d triangles %d len materials %d", curMesh.Name, curTriangleMat.MaterialID, curTriangleMat.Count, len(materials))

					if curTriangleMat.MaterialID >= uint16(len(materials)) {
						log.Debugf("materialID %d out of bounds", curTriangleMat.MaterialID)
						continue
					}

					curMesh.Triangles[i].MaterialName = materials[curTriangleMat.MaterialID].Name
					hasMaterial := false
					for _, mat := range curMesh.Materials {
						if mat.Name == curMesh.Triangles[i].MaterialName {
							hasMaterial = true
							break
						}
					}
					if !hasMaterial {
						curMesh.Materials = append(curMesh.Materials, materials[curTriangleMat.MaterialID])
					}
				}
				materialCounter--
				if materialCounter < 1 {
					materialGroup++
				}
			}

			meshes = append(meshes, curMesh)
			//ref := e.Fragments[d.MaterialListRef].(*wld.MaterialList)
			curMesh = &def.Mesh{}
		case *wld.TextureList:
			log.Debugf("texture list: %+v", d)
			for _, texture := range d.TextureNames {
				material := &def.Material{
					Name: helper.Clean(texture),
				}
				material.Properties = append(material.Properties, &def.MaterialProperty{
					Name:     "e_TextureDiffuse0",
					Value:    helper.Clean(texture),
					Category: 2,
				})
				curMesh.Materials = append(curMesh.Materials, material)
				materials = append(materials, material)
			}

		case *wld.MaterialList:
			for _, ref := range d.MaterialRefs {
				switch t := e.Fragments[int(ref)].(type) {
				case *wld.TextureList:
					log.Debugf("texture list: %+v", t)
					for _, texture := range t.TextureNames {
						material := &def.Material{
							Name: helper.Clean(texture),
						}
						material.Properties = append(material.Properties, &def.MaterialProperty{
							Name:     "e_TextureDiffuse0",
							Value:    helper.Clean(texture),
							Category: 2,
						})
						curMesh.Materials = append(curMesh.Materials, material)
						materials = append(materials, material)
					}
				case *wld.Material:
					log.Debugf("material: %+v", t)
				case *wld.MaterialList:
					for _, ref := range t.MaterialRefs {
						switch t := e.Fragments[int(ref)].(type) {
						case *wld.TextureList:
							log.Debugf("texture list: %+v", t)
							for _, texture := range t.TextureNames {
								material := &def.Material{
									Name: helper.Clean(texture),
								}
								material.Properties = append(material.Properties, &def.MaterialProperty{
									Name:     "e_TextureDiffuse0",
									Value:    helper.Clean(texture),
									Category: 2,
								})
								curMesh.Materials = append(curMesh.Materials, material)
								materials = append(materials, material)
							}
						case *wld.Material:
							log.Debugf("material: %+v", t)
						case *wld.MaterialList:
							log.Debugf("material list: %+v", t)
						default:
							log.Debugf("unknown sub material list ref: %T", t)
						}
					}
				default:
					log.Debugf("unknown material list ref: %T", t)
				}
			}
		}
	}

	return meshes, nil
}
