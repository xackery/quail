package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/metadata/mat"
	"github.com/xackery/quail/tag"
)

// Decode decodes a wld file that was prepped by Load
func Decode(wld *common.Wld) error {
	tag.New()
	curMaterial := &common.Material{
		ShaderName: "Opaque_MaxCB1.fx",
	}

	curModel := &common.Model{}
	materials := make(map[uint32]*common.Material)
	materialList := []uint32{}
	for i := uint32(1); i <= wld.FragmentCount; i++ {
		data, err := wld.Fragment(int(i - 1))
		if err != nil {
			return fmt.Errorf("%d decode: %w", i, err)
		}

		r := bytes.NewReader(data)
		dec := encdec.NewDecoder(r, binary.LittleEndian)

		fragCode := dec.Int32()

		name := ""
		nameRef := int32(0)
		switch fragCode {
		case 0x00: // Empty
		case 0x01: // Default Palette File
		case 0x02: // UserData
		case 0x03: // FrameAndBMInfo (TextureImages)
			// TextureImages always follows up with Texture, TextureRef
			err = mat.DecodeTextureImages(curMaterial, &nameRef, r)
			if err != nil {
				return fmt.Errorf("%d 0x03 TextureList decode %s: %w", i, name, err)
			}
			curMaterial.Name = wld.Names[-nameRef]

		case 0x04: // SimpleSpriteDef (Texture)
			textureRefs := []*int32{}
			err = mat.DecodeTexture(curMaterial, &nameRef, textureRefs, r)
			if err != nil {
				return fmt.Errorf("%d 0x04 TextureAnimation decode %s: %w", i, name, err)
			}
			curMaterial.Name = wld.Names[-nameRef]
			for _, ref := range textureRefs {
				curMaterial.Animation.Textures = append(curMaterial.Animation.Textures, wld.Names[*ref])
			}
		case 0x05: // SimpleSpriteInst (TextureReference)
			dec := encdec.NewDecoder(r, binary.LittleEndian)
			nameRef := dec.Int32() // nameRef
			curMaterial.Name = wld.Names[-nameRef]
			_ = dec.Int32() // textureRef
			if dec.Error() != nil {
				return fmt.Errorf("%d 0x05 TextureReference decode %s: %w", i, name, dec.Error())
			}

		case 0x06: // 2DSpriteDef (TwoDimensionalObject)
		case 0x07: // 2DSpriteReference (TwoDimensionalObjectReference)
		case 0x08: // 3DSpriteDef (Camera)
		case 0x09: // 3DSpriteReference (CameraReference)
		case 0x0a: // 4DSpriteDef
		case 0x0b: // 4DSpriteReference
		case 0x0c: // ParticleSpriteDef
		case 0x0d: // ParticleSpriteReference
		case 0x0e: // CompositeSpriteDef
		case 0x0f: // CompositeSpriteReference
		case 0x10: // HierarchicalSpriteDef (SkeletonTrackSet)
		case 0x11: // HierarchicalSpriteReference (SkeletonTrackSetReference)
		case 0x12: // TrackDefinition (MobSkeletonPieceTrack)
		case 0x13: // TrackInstance (MobSkeletonPieceTrackReference)
		case 0x14: // ActorDef (Model)
		case 0x15: // ActorReference (ObjectLocation)
		case 0x16: // SphereReference (ZoneUknown)
		case 0x17: // PolyhedronDef (PolygonAnimation)
		case 0x18: // PolyhedronReference (PolygonAnimationReference)
		case 0x19: // SphereListDef
		case 0x1a: // SphereListReference
		case 0x1b: // LightDef (LightSource)
		case 0x1c: // LightReference (LightSourceReference)
		case 0x1d: // PointLight
		case 0x1e: // Unknown
		case 0x1f: // SoundDef
		case 0x20: // SoundInstance
		case 0x21: // WorldTree (BspTree)
		case 0x22: // Region (BspRegion)
		case 0x23: // ActiveGeometryRegion
		case 0x24: // SkyRegion
		case 0x25: // DirectionalLight
		case 0x26: // BlitSpriteDefinition (BlitSpriteDefinition)
		case 0x27: // BlitSpriteReference (BlitSpriteReference)
		case 0x28: // PointLight (LightInfo)
		case 0x29: // Zone (RegionFlag)
		case 0x2a: // AmbientLight (AmbientLight)
		case 0x2b: // DirectionalLight
		case 0x2c: // DMSpriteDef (AlternateMesh)
		case 0x2d: // MeshReference (MeshReference)
		case 0x2e: // Unknown
		case 0x2f: // Unknown (MeshAnimatedVerticesReference)
		case 0x30: // MaterialDef
			err = mat.DecodeMaterialDef(curMaterial, &nameRef, r)
			if err != nil {
				return fmt.Errorf("%d 0x03 TextureList decode %s: %w", i, name, err)
			}

			curMaterial.Name = wld.Names[-nameRef]
			wld.Materials = append(wld.Materials, curMaterial)
			curMaterial = &common.Material{
				ShaderName: "Opaque_MaxCB1.fx",
			}
			materials[i] = curMaterial
			fmt.Printf("adding materialRef %d %s\n", i, curMaterial.Name)
		case 0x31: // MaterialPalette (MaterialList)
			dec := encdec.NewDecoder(r, binary.LittleEndian)
			dec.Int32()  // nameRef
			dec.Uint32() // flags
			materialListCount := dec.Uint32()
			for i := 0; i < int(materialListCount); i++ {
				materialList = append(materialList, dec.Uint32())
			}

			if dec.Error() != nil {
				return fmt.Errorf("%d 0x05 TextureReference decode %s: %w", i, name, dec.Error())
			}
		case 0x32: // ? (VertexColor)
		case 0x33: // ? (VertexColorReference)
		case 0x34: // Unknown
		case 0x35: // ? (FirstFragment)
		case 0x36: // DMSpriteDef2 (Mesh)
			err = mod.DecodeMesh(curModel, &nameRef, wld.IsOldWorld, r)
			if err != nil {
				return fmt.Errorf("%d 0x36 Mesh decode %s: %w", i, name, err)
			}
			curModel.Name = wld.Names[-nameRef]
			wld.Models = append(wld.Models, curModel)
			curModel = &common.Model{}
		}
	}

	newModels := []*common.Model{}
	terModel := &common.Model{
		FileType: "ter",
	}
	customModels := []*common.Model{}

	for _, model := range wld.Models {
		var curModel *common.Model

		if model.FileType != "ter" {
			curModel = model
		}
		if model.FileType == "ter" {
			curModel = terModel
			if terModel.Name == "" {
				terModel.Name = model.Name
				if terModel.Name == "" {
					terModel.Name = "ter"
				}
				newModels = append(newModels, terModel)
			}

			// sample first triangle to check for flags
			triangle := model.Triangles[0]
			material := materialLookup(triangle, materials, materialList)

			/*if material != nil && material.ShaderName == "mesh_invisible" {
				curModel = terInvisModel
				if terInvisModel.Name == "" {
					terInvisModel.Name = model.Name + "_invis"
					if terInvisModel.Name == "" {
						terInvisModel.Name = "ter_invis"
					}
					newModels = append(newModels, terInvisModel)
				}
			}*/
			if material != nil {
				isFound := false
				for _, cust := range customModels {
					if cust.Name != material.ShaderName {
						continue
					}
					isFound = true
					curModel = cust
				}
				if !isFound {
					curModel = &common.Model{
						Name:     material.ShaderName,
						FileType: "ter",
					}
					customModels = append(customModels, curModel)
					newModels = append(newModels, curModel)
				}
			}

		}

		vertOffset := len(curModel.Vertices)
		curModel.Vertices = append(curModel.Vertices, model.Vertices...)

		for _, triangle := range model.Triangles {
			newTriangle := common.Triangle{
				Flag: triangle.Flag,
			}
			material := materialLookup(triangle, materials, materialList)
			if material != nil {
				isNew := true
				for _, curMat := range curModel.Materials {
					if curMat.Name == material.Name {
						isNew = false
						break
					}
				}
				if isNew {
					curModel.Materials = append(curModel.Materials, material)
				}
				newTriangle.MaterialName = material.Name
			}
			newTriangle.Index.X = triangle.Index.X + uint32(vertOffset)
			newTriangle.Index.Y = triangle.Index.Y + uint32(vertOffset)
			newTriangle.Index.Z = triangle.Index.Z + uint32(vertOffset)

			curModel.Triangles = append(curModel.Triangles, newTriangle)
		}
	}
	wld.Models = newModels

	return nil
}

func materialLookup(triangle common.Triangle, materials map[uint32]*common.Material, materialList []uint32) *common.Material {
	if !strings.Contains(triangle.MaterialName, "material_") {
		return nil
	}

	materialNameRef, err := strconv.Atoi(strings.TrimPrefix(triangle.MaterialName, "material_"))
	if err != nil {
		return nil
	}

	material, ok := materials[materialList[materialNameRef]]
	if !ok {
		return nil
	}

	return material
}
