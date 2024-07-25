package wld

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
)

func (wld *Wld) WriteAscii(path string, isDir bool) error {
	var err error
	wld.mu.Lock()

	wld.reset()
	defer func() {
		wld.reset()
		wld.mu.Unlock()
	}()
	//var err error

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	baseName := wld.FileName
	ext := filepath.Ext(wld.FileName)
	if ext != "" {
		baseName = baseName[:len(baseName)-len(ext)]
	}

	spkPath := path + "/" + baseName + ".spk"
	if !isDir {
		spkPath = path + "/" + wld.FileName
	}

	var w *os.File
	spkBuf, err := os.Create(spkPath)
	if err != nil {
		return err
	}
	defer spkBuf.Close()
	writeHeader(spkBuf)

	w = spkBuf
	for i := 0; i < len(wld.SimpleSpriteDefs); i++ {
		spriteDef := wld.SimpleSpriteDefs[i]
		if isDir {
			if len(spriteDef.Tag) < 1 {
				return fmt.Errorf("sprite def %d has no tag", i)
			}

			fileName := strings.ToLower(spriteDef.Tag)
			index := strings.LastIndex(fileName, "_")
			if index > 0 {
				fileName = fileName[:index]
			}

			spsBuf, err := os.Create(path + "/" + fileName + ".sps")
			if err != nil {
				return err
			}
			defer spsBuf.Close()
			writeHeader(spsBuf)
			w = spsBuf
		}
		err = spriteDef.Write(w)
		if err != nil {
			return fmt.Errorf("sprite def %s: %w", spriteDef.Tag, err)
		}
	}

	isIncluded := false
	for i := 0; i < len(wld.MaterialDefs); i++ {
		material := wld.MaterialDefs[i]
		if isDir {
			if len(material.Tag) < 1 {
				return fmt.Errorf("material %d has no tag", i)
			}

			fileName := strings.ToLower(material.Tag)
			index := strings.LastIndex(fileName, "_")
			if index > 0 {
				fileName = fileName[:index]
			}

			mdfBuf, err := os.Create(path + "/" + fileName + ".mdf")
			if err != nil {
				return err
			}
			defer mdfBuf.Close()
			writeHeader(mdfBuf)

			w = mdfBuf
			fmt.Fprintf(w, "INCLUDE \"%s\"\n\n", strings.ToUpper(fileName+".SPS"))

			fmt.Fprintf(spkBuf, "INCLUDE \"%s\"\n", strings.ToUpper(fileName+".MDF"))
			isIncluded = true
		}
		err = material.Write(w)
		if err != nil {
			return fmt.Errorf("material %s: %w", material.Tag, err)
		}
	}
	if isDir && isIncluded {
		fmt.Fprintf(spkBuf, "\n")
	}

	w = spkBuf

	for i := 0; i < len(wld.MaterialPalettes); i++ {
		palette := wld.MaterialPalettes[i]
		err = palette.Write(w)
		if err != nil {
			return fmt.Errorf("palette %s: %w", palette.Tag, err)
		}
	}

	for i := 0; i < len(wld.PolyhedronDefs); i++ {
		polyhedron := wld.PolyhedronDefs[i]
		err = polyhedron.Write(w)
		if err != nil {
			return fmt.Errorf("polyhedron %s: %w", polyhedron.Tag, err)
		}
	}

	for i := 0; i < len(wld.DMSpriteDef2s); i++ {
		dmSprite := wld.DMSpriteDef2s[i]
		err = dmSprite.Write(w)
		if err != nil {
			return fmt.Errorf("dm sprite def %s: %w", dmSprite.Tag, err)
		}
	}

	for i := 0; i < len(wld.TrackDefs); i++ {
		trackDef := wld.TrackDefs[i]
		err = trackDef.Write(w)
		if err != nil {
			return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
		}
		if len(wld.TrackInstances) > i {
			trackInst := wld.TrackInstances[i]
			err = trackInst.Write(w)
			if err != nil {
				return fmt.Errorf("track inst %s: %w", trackInst.Tag, err)
			}
		}
	}

	for i := 0; i < len(wld.HierarchicalSpriteDefs); i++ {
		hierarchicalSprite := wld.HierarchicalSpriteDefs[i]
		err = hierarchicalSprite.Write(w)
		if err != nil {
			return fmt.Errorf("hierarchical sprite def %s: %w", hierarchicalSprite.Tag, err)
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
		/* for _, materialTag := range palette.Materials {
			err = wld.writeMaterial(w, materialTag)
			if err != nil {
				return fmt.Errorf("material %s: %w", materialTag, err)
			}
		} */
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
		err = material.Write(w)
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
		err = spriteDef.Write(w)
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
		err = actorDef.Write(w)
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
		err = actorInst.Write(w)
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
		err = lightDef.Write(w)
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

		err = pointLight.Write(w)
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
		err = sprite3DDef.Write(w)
		if err != nil {
			return err
		}

		wld.writtenSprite3DDefs[tag] = true
		return nil
	}

	return fmt.Errorf("not found")
}

func writeHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail v%s\n\n", common.Version)
}
