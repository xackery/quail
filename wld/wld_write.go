package wld

import (
	"fmt"
	"io"

	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func (wld *Wld) Write(w io.Writer) error {
	var err error

	raw.NameClear()

	out := &raw.Wld{
		MetaFileName: wld.FileName,
		Version:      wld.Version,
		Fragments:    []model.FragmentReadWriter{},
	}
	for i := 0; i < len(wld.SpriteInstances); i++ {
		spriteInstance := wld.SpriteInstances[i]
		sprite := wld.spriteByTag(spriteInstance.Sprite)
		if sprite == nil {
			return fmt.Errorf("spriteInstance %s refers to sprite %s which does not exist", spriteInstance.Tag, spriteInstance.Sprite)
		}
		if sprite.fragID == 0 { // if sprite hasn't been inesrted yet
			bitmapRefs := []uint32{}

			for j := 0; j < len(sprite.Bitmaps); j++ {
				bitmap := sprite.Bitmaps[j]
				bmInfo := wld.bitmapByTag(bitmap)
				if bmInfo == nil {
					return fmt.Errorf("spriteInstance %s refers sprite %s which refers to bitmap %s which does not exist", spriteInstance.Tag, sprite.Tag, bitmap)
				}
				if bmInfo.fragID > 0 {
					bitmapRefs = append(bitmapRefs, bmInfo.fragID)
					continue
				}

				nameRef := raw.NameAdd(bmInfo.Tag)
				out.Fragments = append(out.Fragments, &rawfrag.WldFragBMInfo{
					NameRef:      nameRef,
					TextureNames: bmInfo.Textures,
				})
				bitmapRefs = append(bitmapRefs, uint32(len(out.Fragments)))
			}

			nameRef := raw.NameAdd(sprite.Tag)
			out.Fragments = append(out.Fragments, &rawfrag.WldFragSimpleSpriteDef{
				NameRef:      nameRef,
				Flags:        sprite.Flags,
				CurrentFrame: sprite.CurrentFrame,
				Sleep:        sprite.Sleep,
				BitmapRefs:   bitmapRefs,
			})
		}

		nameRef := raw.NameAdd(spriteInstance.Tag)
		out.Fragments = append(out.Fragments, &rawfrag.WldFragSimpleSprite{
			NameRef:   nameRef,
			SpriteRef: int16(sprite.fragID),
			Flags:     spriteInstance.Flags,
		})
	}

	err = out.Write(w)
	if err != nil {
		return err
	}

	return nil
}
