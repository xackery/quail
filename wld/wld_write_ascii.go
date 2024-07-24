package wld

import (
	"fmt"
	"io"
)

func (wld *Wld) WriteAscii(w io.Writer) error {
	var err error
	wld.mu.Lock()

	wld.reset()
	defer func() {
		wld.reset()
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

	for i := 0; i < len(wld.PointLights); i++ {
		pointLight := wld.PointLights[i]
		err = wld.writePointLight(w, pointLight.Tag)
		if err != nil {
			return fmt.Errorf("point light %s: %w", pointLight.Tag, err)
		}
	}

	for i := 0; i < len(wld.Sprite3DDefs); i++ {
		sprite3DDef := wld.Sprite3DDefs[i]

		err = wld.writeSprite3DDef(w, sprite3DDef.Tag)
		if err != nil {
			return fmt.Errorf("sprite 3d def %s: %w", sprite3DDef.Tag, err)
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
		err = palette.Write(w)
		if err != nil {
			return fmt.Errorf("palette %s: %w", palette.Tag, err)
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

func (wld *Wld) writeLightDef(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenLightDefs[tag]
	if ok {
		return nil
	}

	for _, lightDef := range wld.LightDefs {
		if lightDef.Tag != tag {
			continue
		}
		_, err = w.Write([]byte(lightDef.Ascii()))
		if err != nil {
			return err
		}

		wld.writtenLightDefs[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}

func (wld *Wld) writePointLight(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenPointLights[tag]
	if ok {
		return nil
	}

	for _, pointLight := range wld.PointLights {
		if pointLight.Tag != tag {
			continue
		}
		if pointLight.LightDefTag != "" {
			err = wld.writeLightDef(w, pointLight.LightDefTag)
			if err != nil {
				return fmt.Errorf("light def %s: %w", pointLight.LightDefTag, err)
			}
		}

		_, err = w.Write([]byte(pointLight.Ascii()))
		if err != nil {
			return err
		}

		wld.writtenPointLights[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}

func (wld *Wld) writeSprite3DDef(w io.Writer, tag string) error {
	var err error
	if tag == "" {
		return nil
	}
	_, ok := wld.writtenSprite3DDefs[tag]
	if ok {
		return nil
	}

	for _, sprite3DDef := range wld.Sprite3DDefs {
		if sprite3DDef.Tag != tag {
			continue
		}
		_, err = w.Write([]byte(sprite3DDef.Ascii()))
		if err != nil {
			return err
		}

		wld.writtenSprite3DDefs[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}
