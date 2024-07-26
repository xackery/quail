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

	for i := 0; i < len(wld.AmbientLights); i++ {
		ambientLight := wld.AmbientLights[i]
		err = ambientLight.Write(w)
		if err != nil {
			return fmt.Errorf("ambient light %s: %w", ambientLight.Tag, err)
		}
	}

	for i := 0; i < len(wld.ActorDefs); i++ {
		actorDef := wld.ActorDefs[i]
		err = actorDef.Write(w)
		if err != nil {
			return fmt.Errorf("actor def %s: %w", actorDef.Tag, err)
		}
	}

	for i := 0; i < len(wld.ActorInsts); i++ {
		actorInst := wld.ActorInsts[i]
		err = actorInst.Write(w)
		if err != nil {
			return fmt.Errorf("actor inst %s: %w", actorInst.Tag, err)
		}

	}

	for i := 0; i < len(wld.Zones); i++ {
		zone := wld.Zones[i]
		err = zone.Write(w)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	for i := 0; i < len(wld.LightDefs); i++ {
		lightDef := wld.LightDefs[i]
		err = lightDef.Write(w)
		if err != nil {
			return fmt.Errorf("light def %s: %w", lightDef.Tag, err)
		}
	}

	for i := 0; i < len(wld.PointLights); i++ {
		pointLight := wld.PointLights[i]
		err = pointLight.Write(w)
		if err != nil {
			return fmt.Errorf("point light %s: %w", pointLight.Tag, err)
		}

	}

	for i := 0; i < len(wld.Sprite3DDefs); i++ {
		sprite3DDef := wld.Sprite3DDefs[i]
		err = sprite3DDef.Write(w)
		if err != nil {
			return fmt.Errorf("sprite 3d def %s: %w", sprite3DDef.Tag, err)
		}
	}

	for i := 0; i < len(wld.WorldTrees); i++ {
		worldTree := wld.WorldTrees[i]
		err = worldTree.Write(w)
		if err != nil {
			return fmt.Errorf("world tree %d: %w", i, err)
		}
	}

	for i := 0; i < len(wld.Regions); i++ {
		region := wld.Regions[i]
		err = region.Write(w)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.RegionTag, err)
		}
	}

	return nil
}

func writeHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail v%s\n\n", common.Version)
}
