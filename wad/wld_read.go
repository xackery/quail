package wad

import (
	"fmt"
	"io"

	"github.com/xackery/quail/raw"
)

func (e *Wad) WLDRead(r io.ReadSeeker) error {
	e.bitmaps = []*BMInfo{}
	e.dmsprites = []*DMSpriteInfo{}

	wld := raw.Wld{}
	err := wld.Read(r)
	if err != nil {
		return err
	}

	fragments := wld.Fragments

	bmInfoRefs := map[int]*BMInfo{}
	spriteRefs := map[int]*Sprite{}
	spriteInstRefs := map[int]*SpriteInstance{}

	for i := 1; i < len(fragments); i++ {
		fragment := fragments[i]

		switch fragment.FragCode() {
		case 0x03:
			bmInfoDef, ok := fragment.(*raw.WldFragBMInfo)
			if !ok {
				return fmt.Errorf("decode frag %d: expected *WldFragBitmapInfo, got %T", i, fragment)
			}
			bmInfo := &BMInfo{
				Tag:      raw.Name(bmInfoDef.NameRef),
				Textures: bmInfoDef.TextureNames,
			}

			bmInfoRefs[i] = bmInfo
			e.bitmaps = append(e.bitmaps, bmInfo)
		case 0x04:
			spriteDef, ok := fragment.(*raw.WldFragSimpleSpriteDef)
			if !ok {
				return fmt.Errorf("decode frag %d: expected *WldFragBitmapInfo, got %T", i, fragment)
			}

			sprite := &Sprite{
				Tag:          raw.Name(spriteDef.NameRef),
				Flags:        spriteDef.Flags,
				CurrentFrame: spriteDef.CurrentFrame,
				Sleep:        spriteDef.Sleep,
			}
			for _, frame := range spriteDef.BitmapRefs {
				bmInfo, ok := bmInfoRefs[int(frame)]
				if !ok {
					return fmt.Errorf("decode frag %d: no bitmap ref %d found", i, int(frame))
				}

				sprite.Frames = append(sprite.Frames, bmInfo)
			}

			spriteRefs[i] = sprite
			e.sprites = append(e.sprites, sprite)
		case 0x05:
			spriteInst, ok := fragment.(*raw.WldFragSimpleSprite)
			if !ok {
				return fmt.Errorf("decode frag %d: expected *WldFragSimpleSprite, got %T", i, fragment)
			}
			ref := int(spriteInst.SpriteRef)
			sprite, ok := spriteRefs[ref]
			if !ok {
				return fmt.Errorf("decode frag %d: no sprite ref %d found", i, ref)
			}

			inst := &SpriteInstance{
				Tag:   raw.Name(spriteInst.NameRef),
				Flags: spriteInst.Flags,
			}
			spriteInstRefs[i] = inst
			fmt.Println("added sprite inst", inst.Tag, "on fragment", i)
			sprite.Instances = append(sprite.Instances, inst)
		case 0x30:
			material, ok := fragment.(*raw.WldFragMaterialDef)
			if !ok {
				return fmt.Errorf("decode frag %d: expected *WldFragMaterialDef, got %T", i, fragment)
			}

			spriteInst, ok := spriteInstRefs[int(material.SpriteInstanceRef)]
			if !ok {
				return fmt.Errorf("decode frag %d material: no sprite inst ref %d found", i, int(material.SpriteInstanceRef))
			}

			mat := &SpriteInstanceMaterial{
				Tag:           raw.Name(material.NameRef),
				Flags:         material.Flags,
				RenderMethod:  material.RenderMethod,
				RGBPen:        material.RGBPen,
				Brightness:    material.Brightness,
				ScaledAmbient: material.ScaledAmbient,
				Pairs:         material.Pairs,
			}

			spriteInst.Materials = append(spriteInst.Materials, mat)

		}

	}

	return nil
}
