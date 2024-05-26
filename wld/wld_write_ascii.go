package wld

import (
	"fmt"
	"io"
)

func (wld *Wld) WriteAscii(w io.Writer) error {
	var err error
	wld.mu.Lock()
	defer wld.mu.Unlock()
	wld.writtenBMInfos = make(map[string]bool)
	wld.writtenMaterials = make(map[string]bool)
	wld.writtenPalettes = make(map[string]bool)
	//var err error

	for i := 0; i < len(wld.DMSpriteDef2s); i++ {
		dmSprite := wld.DMSpriteDef2s[i]
		err = wld.writePalette(w, dmSprite.MaterialPaletteTag)
		if err != nil {
			return fmt.Errorf("palette %s: %w", dmSprite.MaterialPaletteTag, err)
		}
		_, err = w.Write([]byte(dmSprite.Ascii()))
		if err != nil {
			return err
		}

	}

	return nil
}

func (wld *Wld) writePalette(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenPalettes[tag]
	if ok {
		return nil
	}

	for _, palette := range wld.MaterialPalettes {
		if palette.Tag != tag {
			continue
		}
		for _, materialTag := range palette.Materials {
			err = wld.writeMaterial(w, materialTag)
			if err != nil {
				return fmt.Errorf("material %s: %w", materialTag, err)
			}
		}
		_, err = w.Write([]byte(palette.Ascii()))
		if err != nil {
			return err
		}
		wld.writtenPalettes[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}

func (wld *Wld) writeMaterial(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenMaterials[tag]
	if ok {
		return nil
	}

	for _, material := range wld.MaterialDefs {
		if material.Tag != tag {
			continue
		}
		err = wld.writeSpriteDef(w, material.SimpleSpriteInstTag)
		if err != nil {
			return fmt.Errorf("sprite def %s: %w", material.SimpleSpriteInstTag, err)
		}
		_, err = w.Write([]byte(material.Ascii()))
		if err != nil {
			return err
		}
		wld.writtenMaterials[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}

func (wld *Wld) writeSpriteDef(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenSpriteDefs[tag]
	if ok {
		return nil
	}

	for _, spriteDef := range wld.SimpleSpriteDefs {
		if spriteDef.Tag != tag {
			continue
		}
		_, err = w.Write([]byte(spriteDef.Ascii()))
		if err != nil {
			return err
		}
		wld.writtenSpriteDefs[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}
