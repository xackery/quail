package quail

import (
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/wld"
	"github.com/xackery/quail/pfs"
)

// Decode decodes a WLD file
func WLDDecode(r io.ReadSeeker, pfs *pfs.PFS) ([]*common.Model, error) {
	models := make([]*common.Model, 0)

	e, err := wld.New("test", pfs)
	if err != nil {
		return nil, err
	}

	err = e.Decode(r)
	if err != nil {
		return nil, err
	}

	materials := make([]*common.Material, 0)

	curModel := &common.Model{}
	//names := e.Names()
	for _, f := range e.Fragments {
		switch d := f.(type) {
		case *wld.Mesh:
			curModel.Name = d.Name

			curModel.Vertices = d.Vertices
			curModel.Triangles = d.Triangles

			materialGroup := 0
			materialCounter := 0
			for i := 0; i < len(curModel.Triangles); i++ {
				curTriangleMat := d.TriangleMaterials[materialGroup]
				if materialCounter < 1 {
					materialCounter = int(curTriangleMat.Count)
				}
				matID := curTriangleMat.MaterialID
				if matID != 0 {
					matID--
					log.Debugf("model %s materialID %d triangles %d len materials %d", curModel.Name, curTriangleMat.MaterialID, curTriangleMat.Count, len(materials))

					if curTriangleMat.MaterialID >= uint16(len(materials)) {
						log.Debugf("materialID %d out of bounds", curTriangleMat.MaterialID)
						continue
					}

					curModel.Triangles[i].MaterialName = materials[curTriangleMat.MaterialID].Name
					hasMaterial := false
					for _, mat := range curModel.Materials {
						if mat.Name == curModel.Triangles[i].MaterialName {
							hasMaterial = true
							break
						}
					}
					if !hasMaterial {
						curModel.Materials = append(curModel.Materials, materials[curTriangleMat.MaterialID])
					}
				}
				materialCounter--
				if materialCounter < 1 {
					materialGroup++
				}
			}

			models = append(models, curModel)
			//ref := e.Fragments[d.MaterialListRef].(*wld.MaterialList)
			curModel = &common.Model{}
		case *wld.TextureList:
			log.Debugf("texture list: %+v", d)
			for _, texture := range d.TextureNames {
				material := &common.Material{
					Name: helper.Clean(texture),
				}
				material.Properties = append(material.Properties, &common.MaterialProperty{
					Name:     "e_TextureDiffuse0",
					Value:    helper.Clean(texture),
					Category: 2,
				})
				curModel.Materials = append(curModel.Materials, material)
				materials = append(materials, material)
			}

		case *wld.MaterialList:
			for _, ref := range d.MaterialRefs {
				switch t := e.Fragments[int(ref)].(type) {
				case *wld.TextureList:
					log.Debugf("texture list: %+v", t)
					for _, texture := range t.TextureNames {
						material := &common.Material{
							Name: helper.Clean(texture),
						}
						material.Properties = append(material.Properties, &common.MaterialProperty{
							Name:     "e_TextureDiffuse0",
							Value:    helper.Clean(texture),
							Category: 2,
						})
						curModel.Materials = append(curModel.Materials, material)
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
								material := &common.Material{
									Name: helper.Clean(texture),
								}
								material.Properties = append(material.Properties, &common.MaterialProperty{
									Name:     "e_TextureDiffuse0",
									Value:    helper.Clean(texture),
									Category: 2,
								})
								curModel.Materials = append(curModel.Materials, material)
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

	return models, nil
}
