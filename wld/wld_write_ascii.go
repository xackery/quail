package wld

import (
	"fmt"
	"io"
)

func (wld *Wld) WriteAscii(w io.Writer) error {
	var err error
	wld.mu.Lock()

	wld.writtenMaterials = make(map[string]bool)
	wld.writtenSpriteDefs = make(map[string]bool)
	wld.writtenPalettes = make(map[string]bool)
	wld.writtenActorDefs = make(map[string]bool)
	wld.writtenActorInsts = make(map[string]bool)
	defer func() {
		wld.writtenMaterials = nil
		wld.writtenSpriteDefs = nil
		wld.writtenPalettes = nil
		wld.writtenActorDefs = nil
		wld.writtenActorInsts = nil

		wld.mu.Unlock()
	}()
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

	for i := 0; i < len(wld.ActorDefs); i++ {
		actorDef := wld.ActorDefs[i]
		err = wld.writeActorDef(w, actorDef.Tag)
		if err != nil {
			return fmt.Errorf("actor def %s: %w", actorDef.Tag, err)
		}
	}

	for i := 0; i < len(wld.ActorInsts); i++ {
		actorInst := wld.ActorInsts[i]
		err = wld.writeActorInst(w, actorInst.Tag)
		if err != nil {
			return fmt.Errorf("actor inst %s: %w", actorInst.Tag, err)
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

func (wld *Wld) writeActorDef(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenActorDefs[tag]
	if ok {
		return nil
	}

	for _, actorDef := range wld.ActorDefs {
		if actorDef.Tag != tag {
			continue
		}
		_, err = w.Write([]byte(actorDef.Ascii()))
		if err != nil {
			return err
		}

		wld.writtenActorDefs[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}

func (wld *Wld) writeActorInst(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenActorInsts[tag]
	if ok {
		return nil
	}

	for _, actorInst := range wld.ActorInsts {
		if actorInst.Tag != tag {
			continue
		}
		_, err = w.Write([]byte(actorInst.Ascii()))
		if err != nil {
			return err
		}

		wld.writtenActorInsts[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}
