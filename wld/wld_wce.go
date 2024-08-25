package wld

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xackery/quail/common"
)

// ReadAscii reads the ascii file at path
func (wld *Wld) ReadAscii(path string) error {

	wld.reset()
	wld.maxMaterialHeads = make(map[string]int)
	wld.maxMaterialTextures = make(map[string]int)

	asciiReader, err := LoadAsciiFile(path, wld)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = asciiReader.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, asciiReader.lineNumber, err)
	}
	fmt.Println(asciiReader.TotalLineCountRead(), "total lines parsed for", filepath.Base(path))
	return nil
}

func (wld *Wld) WriteAscii(path string) error {
	var err error

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	baseTags := []string{}
	for _, actorDef := range wld.ActorDefs {
		if len(actorDef.Tag) < 3 {
			return fmt.Errorf("actorDef %s tag too short", actorDef.Tag)
		}
		baseTag := baseTagTrim(actorDef.Tag)
		isFound := false
		for _, tag := range baseTags {
			if tag == baseTag {
				isFound = true
				break
			}
		}
		if !isFound {
			baseTags = append(baseTags, baseTag)
		}
	}
	if wld.WorldDef != nil && wld.WorldDef.Zone == 1 {
		baseTags = append(baseTags, "R")
	}

	err = wld.writeAsciiData(path, baseTags)
	if err != nil {
		return err
	}

	return nil
}

func (wld *Wld) writeAsciiData(path string, baseTags []string) error {

	token := NewAsciiWriteToken(path, wld)
	defer token.Close()

	if wld.WorldDef == nil {
		return fmt.Errorf("worlddef not found")
	}
	/*
		for _, track := range wld.TrackInstances {
			if len(track.Tag) < 3 {
				return fmt.Errorf("trackdef %s tag too short", track.Tag)
			}
			baseTag := track.SpriteTag
			if len(baseTag) < 1 {
				return fmt.Errorf("track sprite tag %s too short (%s)", track.Tag, baseTag)
			}
			isFound := false
			for _, tag := range baseTags {
				if tag == baseTag || tag == track.SpriteTag {
					isFound = true
					break
				}
			}
			if !isFound {
				baseTags = append(baseTags, baseTag)
			}
		} */

	for _, actorDef := range wld.ActorDefs {
		if len(actorDef.Tag) < 3 {
			return fmt.Errorf("actorDef %s tag too short", actorDef.Tag)
		}
		baseTag := baseTagTrim(actorDef.Tag)
		isFound := false
		for _, tag := range baseTags {
			if tag == baseTag {
				isFound = true
				break
			}
		}
		if !isFound {
			baseTags = append(baseTags, baseTag)
		}
	}

	err := token.AddWriter("world", fmt.Sprintf("%s/world.mod", path))
	if err != nil {
		return fmt.Errorf("add writer: %w", err)
	}

	for _, baseTag := range baseTags {
		writePath := fmt.Sprintf("%s/%s/%s.mod", path, strings.ToLower(baseTag), strings.ToLower(baseTag))
		err = token.AddWriter(baseTag, writePath)
		if err != nil {
			return fmt.Errorf("add writer %s: %w", baseTag, err)
		}

		writePath = fmt.Sprintf("%s/%s/%s.ani", path, strings.ToLower(baseTag), strings.ToLower(baseTag))
		err = token.AddWriter(baseTag+"_ani", writePath)
		if err != nil {
			return fmt.Errorf("add writer %s_ani: %w", baseTag, err)
		}
	}

	if wld.WorldDef != nil {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set worlddef writer zone: %w", err)
		}
		err = wld.WorldDef.Write(token)
		if err != nil {
			return fmt.Errorf("world def: %w", err)
		}
	}

	if wld.GlobalAmbientLightDef != nil {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set global ambient light writer zone: %w", err)
		}
		err = wld.GlobalAmbientLightDef.Write(token)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	clks := make(map[string]bool)
	err = token.SetWriter("world")
	if err != nil {
		return fmt.Errorf("set material palette writer zone: %w", err)
	}
	for _, matDef := range wld.MaterialDefs {
		if !strings.HasPrefix(matDef.Tag, "CLK") {
			continue
		}

		_, err := strconv.Atoi(matDef.Tag[3:6])
		if err != nil {
			continue
		}
		if clks[matDef.Tag] {
			continue
		}
		clks[matDef.Tag] = true

		err = matDef.Write(token)
		if err != nil {
			return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
		}
	}

	for _, region := range wld.Regions {
		err = token.SetWriter("R")
		if err != nil {
			return fmt.Errorf("set region %s writer: %w", region.Tag, err)
		}

		err = region.Write(token)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	for _, actorDef := range wld.ActorDefs {
		token.TagClearIsWritten()
		err = token.SetWriter(actorDef.Tag)
		if err != nil {
			return fmt.Errorf("set actordef %s writer: %w", actorDef.Tag, err)
		}
		err = actorDef.Write(token)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	if wld.WorldDef.Zone == 1 {
		for _, dmSpriteDef := range wld.DMSpriteDef2s {
			err = token.SetWriter(dmSpriteDef.Tag)
			if err != nil {
				return fmt.Errorf("set dmspritedef2 %s writer: %w", dmSpriteDef.Tag, err)
			}
			err = dmSpriteDef.Write(token)
			if err != nil {
				return fmt.Errorf("dmspritedef2 %s: %w", dmSpriteDef.Tag, err)
			}
		}

		for _, hSprite := range wld.HierarchicalSpriteDefs {
			err = token.SetWriter(hSprite.Tag)
			if err != nil {
				return fmt.Errorf("set hspritedef %s writer: %w", hSprite.Tag, err)
			}

			err = hSprite.Write(token)
			if err != nil {
				return fmt.Errorf("hspritedef %s: %w", hSprite.Tag, err)
			}
		}

		for _, dSprite := range wld.DMSpriteDefs {
			err = token.SetWriter(dSprite.Tag)
			if err != nil {
				return fmt.Errorf("set dmspritedef %s writer: %w", dSprite.Tag, err)
			}

			err = dSprite.Write(token)
			if err != nil {
				return fmt.Errorf("dmspritedef %s: %w", dSprite.Tag, err)
			}
		}

		// global tracks
		for _, track := range wld.TrackInstances {
			if len(track.Tag) < 3 {
				return fmt.Errorf("track %s tag too short", track.Tag)
			}
			if len(track.SpriteTag) < 3 {
				return fmt.Errorf("track %s model too short", track.Tag)
			}
			tag := track.SpriteTag
			if (track.Sleep.Valid && track.Sleep.Uint32 > 0) ||
				isAnimationPrefix(track.Tag) {
				tag += "_ani"
			}

			err = token.SetWriter(tag)
			if err != nil {
				return fmt.Errorf("set track baseTag (%s) %s writer: %w", tag, track.Tag, err)
			}

			err = track.Write(token)
			if err != nil {
				return fmt.Errorf("track %s_%d: %w", track.Tag, track.TagIndex, err)
			}
		}

		for _, polyDef := range wld.PolyhedronDefs {
			err = token.SetWriter(polyDef.Tag)
			if err != nil {
				return fmt.Errorf("set polyhedron %s writer: %w", polyDef.Tag, err)
			}

			err = polyDef.Write(token)
			if err != nil {
				return fmt.Errorf("polyhedron %s: %w", polyDef.Tag, err)
			}
		}
	}

	for _, pLight := range wld.PointLights {
		err = token.SetWriter(pLight.Tag)
		if err != nil {
			return fmt.Errorf("set point light %s writer: %w", pLight.Tag, err)
		}

		err = pLight.Write(token)
		if err != nil {
			return fmt.Errorf("point light %s: %w", pLight.Tag, err)
		}
	}

	// for _, matDef := range wld.MaterialDefs {
	// 	tag := matDef.Tag
	// 	if wld.WorldDef.Zone == 1 {
	// 		tag = "R"
	// 	}

	// 	if strings.Count(tag, "_") > 1 {
	// 		tag = strings.Split(tag, "_")[0]
	// 	}

	// 	err = token.SetWriter(tag)
	// 	if err != nil {
	// 		return fmt.Errorf("set materialdef %s writer: %w", matDef.Tag, err)
	// 	}

	// 	err = matDef.Write(token)
	// 	if err != nil {
	// 		return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
	// 	}
	// }

	for _, lightDef := range wld.LightDefs {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set lightdef %s writer: %w", lightDef.Tag, err)
		}

		err = lightDef.Write(token)
		if err != nil {
			return fmt.Errorf("lightdef %s: %w", lightDef.Tag, err)
		}
	}

	for _, tree := range wld.WorldTrees {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set world tree %s writer: %w", tree.Tag, err)
		}

		err = tree.Write(token)
		if err != nil {
			return fmt.Errorf("world tree %s: %w", tree.Tag, err)
		}
	}

	for _, aLight := range wld.AmbientLights {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set ambient light %s writer: %w", aLight.Tag, err)
		}

		err = aLight.Write(token)
		if err != nil {
			return fmt.Errorf("ambient light %s: %w", aLight.Tag, err)
		}
	}

	for _, actor := range wld.ActorInsts {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set actor %s writer: %w", actor.Tag, err)
		}

		err = actor.Write(token)
		if err != nil {
			return fmt.Errorf("actor %s: %w", actor.Tag, err)
		}
	}

	// for _, actorDef := range wld.ActorDefs {
	// 	tag := actorDef.Tag
	// 	if wld.WorldDef.Zone == 1 {
	// 		tag = "R"
	// 	}
	// 	err = token.SetWriter(tag)
	// 	if err != nil {
	// 		return fmt.Errorf("set actor def %s writer: %w", actorDef.Tag, err)
	// 	}

	// 	err = actorDef.Write(token)
	// 	if err != nil {
	// 		return fmt.Errorf("actor def %s: %w", actorDef.Tag, err)
	// 	}
	// }

	for _, zone := range wld.Zones {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set zone %s writer: %w", zone.Tag, err)
		}

		err = zone.Write(token)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	token.Close()

	rootW, err := os.Create(fmt.Sprintf("%s/_root.wce", path))
	if err != nil {
		return err
	}
	wld.writeAsciiHeader(rootW)

	defer rootW.Close()

	if token.IsWriterUsed("world") {
		rootW.WriteString("INCLUDE \"WORLD.MOD\"\n")
	} else {
		err = os.Remove(fmt.Sprintf("%s/world.mod", path))
		if err != nil {
			return fmt.Errorf("remove %s: %w", fmt.Sprintf("%s/world.mod", path), err)
		}
	}

	for _, baseTag := range baseTags {

		if baseTag != "PLAYER" &&
			!token.IsWriterUsed(baseTag) &&
			!token.IsWriterUsed(baseTag+"_ani") {
			return fmt.Errorf("tag %s was never used for model or ani", baseTag)
		}

		modelW, err := os.Create(fmt.Sprintf("%s/%s/_root.wce", path, strings.ToLower(baseTag)))
		if err != nil {
			return err
		}
		wld.writeAsciiHeader(modelW)

		defer modelW.Close()

		rootW.WriteString(fmt.Sprintf("INCLUDE \"%s/_ROOT.WCE\"\n", strings.ToUpper(baseTag)))

		if token.IsWriterUsed(baseTag) {
			_, err = modelW.WriteString(fmt.Sprintf("INCLUDE \"%s.MOD\"\n", strings.ToUpper(baseTag)))
			if err != nil {
				return err
			}
		} else {
			removePath := fmt.Sprintf("%s/%s/%s.mod", path, strings.ToLower(baseTag), strings.ToLower(baseTag))

			err = os.Remove(removePath)
			if err != nil {
				return fmt.Errorf("remove %s: %w", removePath, err)
			}
		}

		if token.IsWriterUsed(baseTag + "_ani") {
			_, err = modelW.WriteString(fmt.Sprintf("INCLUDE \"%s.ANI\"\n", strings.ToUpper(baseTag)))
			if err != nil {
				return fmt.Errorf("write %s: %w", fmt.Sprintf("%s/%s/%s.ani", path, strings.ToLower(baseTag), strings.ToLower(baseTag)), err)
			}
		} else {
			removePath := fmt.Sprintf("%s/%s/%s.ani", path, strings.ToLower(baseTag), strings.ToLower(baseTag))

			err = os.Remove(removePath)
			if err != nil {
				return fmt.Errorf("remove %s: %w", removePath, err)
			}
		}
	}

	return nil
}

func (wld *Wld) writeAsciiHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail v%s\n", common.Version)
	fmt.Fprintf(w, "// Original file: %s\n\n", wld.FileName)
}
