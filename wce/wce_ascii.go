package wce

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
func (wce *Wce) ReadAscii(path string) error {

	wce.reset()
	wce.maxMaterialHeads = make(map[string]int)
	wce.maxMaterialTextures = make(map[string]int)

	asciiReader, err := LoadAsciiFile(path, wce)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = asciiReader.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, asciiReader.lineNumber, err)
	}
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	relPath, err := filepath.Rel(dir, path)
	if err != nil {
		relPath = path
	}
	fmt.Println(asciiReader.TotalLineCountRead(), "total lines parsed for", relPath)
	return nil
}

func (wce *Wce) WriteAscii(path string) error {
	var err error

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	baseTags := []string{}
	for _, actorDef := range wce.ActorDefs {
		if len(actorDef.Tag) < 3 {
			return fmt.Errorf("actorDef %s tag too short", actorDef.Tag)
		}
		baseTag := baseTagTrim(wce.isObj, actorDef.Tag)
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
	if wce.WorldDef != nil && wce.WorldDef.Zone == 1 {
		baseTags = append(baseTags, "R")
	}

	err = wce.writeAsciiData(path, baseTags)
	if err != nil {
		return err
	}

	return nil
}

func (wce *Wce) writeAsciiData(path string, baseTags []string) error {

	token := NewAsciiWriteToken(path, wce)
	defer token.Close()

	if wce.WorldDef == nil {
		return fmt.Errorf("worlddef not found")
	}

	for _, track := range wce.TrackInstances {
		baseTag, _, _, _ := wce.trackTagAndSequence(track.Tag)
		if baseTag == "" {
			// return fmt.Errorf("track %s tag too short (baseTag empty)", track.Tag)
			baseTag = track.model
		}

		track.model = baseTag
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

	for _, trackDef := range wce.TrackDefs {
		baseTag, _, _, _ := wce.trackTagAndSequence(trackDef.Tag)
		if baseTag == "" {
			//return fmt.Errorf("trackDef %s tag too short (baseTag empty)", trackDef.Tag)
			baseTag = trackDef.model

		}

		trackDef.model = baseTag
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

	for _, trackDef := range wce.DMTrackDef2s {
		baseTag, _, _, _ := wce.trackTagAndSequence(trackDef.Tag)
		if baseTag == "" {
			//return fmt.Errorf("trackDef %s tag too short (baseTag empty)", trackDef.Tag)
			baseTag = trackDef.model

		}

		trackDef.model = baseTag
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

	for _, actorDef := range wce.ActorDefs {
		if len(actorDef.Tag) < 3 {
			return fmt.Errorf("actorDef %s tag too short", actorDef.Tag)
		}
		baseTag := baseTagTrim(wce.isObj, actorDef.Tag)
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

	for _, mdsDef := range wce.MdsDefs {
		if len(mdsDef.Tag) < 3 {
			return fmt.Errorf("mdsDef %s tag too short", mdsDef.Tag)
		}
		baseTag := baseTagTrim(wce.isObj, mdsDef.Tag)
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

	for _, modDef := range wce.ModDefs {
		if len(modDef.Tag) < 3 {
			return fmt.Errorf("modDef %s tag too short", modDef.Tag)
		}
		baseTag := modDef.Tag
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

	for _, terDef := range wce.TerDefs {
		if len(terDef.Tag) < 3 {
			return fmt.Errorf("terDef %s tag too short", terDef.Tag)
		}
		baseTag := terDef.Tag
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

	err := token.AddWriter("world", fmt.Sprintf("%s/world.wce", path))
	if err != nil {
		return fmt.Errorf("add world writer: %w", err)
	}

	err = token.AddWriter("region", fmt.Sprintf("%s/region.wce", path))
	if err != nil {
		return fmt.Errorf("add region writer: %w", err)
	}

	for _, baseTag := range baseTags {
		writePath := fmt.Sprintf("%s/%s/%s.wce", path, strings.ToLower(baseTag), strings.ToLower(baseTag))
		err = token.AddWriter(baseTag, writePath)
		if err != nil {
			return fmt.Errorf("add writer %s: %w", baseTag, err)
		}

		writePath = fmt.Sprintf("%s/%s/%s_ani.wce", path, strings.ToLower(baseTag), strings.ToLower(baseTag))
		err = token.AddWriter(baseTag+"_ani", writePath)
		if err != nil {
			return fmt.Errorf("add writer %s_ani: %w", baseTag, err)
		}
	}

	if wce.WorldDef != nil {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set worlddef writer zone: %w", err)
		}
		err = wce.WorldDef.Write(token)
		if err != nil {
			return fmt.Errorf("world def: %w", err)
		}
	}

	if wce.GlobalAmbientLightDef != nil {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set global ambient light writer zone: %w", err)
		}
		err = wce.GlobalAmbientLightDef.Write(token)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	clks := make(map[string]bool)
	err = token.SetWriter("world")
	if err != nil {
		return fmt.Errorf("set material palette writer zone: %w", err)
	}
	for _, matDef := range wce.MaterialDefs {
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

	for _, region := range wce.Regions {
		err = token.SetWriter("R")
		if err != nil {
			return fmt.Errorf("set R %s writer: %w", region.Tag, err)
		}

		err = region.Write(token)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	for _, actorDef := range wce.ActorDefs {
		token.TagClearIsWritten()
		baseTag := baseTagTrim(wce.isObj, actorDef.Tag)
		wce.lastReadModelTag = baseTag

		err = token.SetWriter(actorDef.Tag)
		if err != nil {
			return fmt.Errorf("set actordef %s writer: %w", actorDef.Tag, err)
		}
		err = actorDef.Write(token)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	if wce.WorldDef.Zone == 1 {

		for _, dSprite := range wce.DMSpriteDefs {
			err = token.SetWriter(dSprite.Tag)
			if err != nil {
				return fmt.Errorf("set dmspritedef %s writer: %w", dSprite.Tag, err)
			}

			err = dSprite.Write(token)
			if err != nil {
				return fmt.Errorf("dmspritedef %s: %w", dSprite.Tag, err)
			}
		}
	}

	// global tracks
	for _, track := range wce.TrackInstances {
		if len(track.Tag) < 3 {
			return fmt.Errorf("track %s tag too short", track.Tag)
		}
		if len(track.SpriteTag) < 3 {
			return fmt.Errorf("track %s model too short (%s)", track.Tag, track.SpriteTag)
		}

		tag := track.model
		if (track.Sleep.Valid && track.Sleep.Uint32 > 0) ||
			wce.isTrackAni(track.Tag) {
			tag += "_ani"
		}

		if token.TagIsWritten(fmt.Sprintf("%s_%d", track.Tag, track.TagIndex)) {
			continue
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

	if wce.WorldDef.Zone == 1 {

		for _, polyDef := range wce.PolyhedronDefs {
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

	for _, pLight := range wce.PointLights {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set point light %s writer: %w", pLight.Tag, err)
		}

		err = pLight.Write(token)
		if err != nil {
			return fmt.Errorf("point light (%s): %w", pLight.Tag, err)
		}
	}

	// for _, matDef := range wce.MaterialDefs {
	// 	tag := matDef.Tag
	// 	if wce.WorldDef.Zone == 1 {
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

	for _, lightDef := range wce.LightDefs {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set lightdef %s writer: %w", lightDef.Tag, err)
		}

		err = lightDef.Write(token)
		if err != nil {
			return fmt.Errorf("lightdef %s: %w", lightDef.Tag, err)
		}
	}

	for _, tree := range wce.WorldTrees {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set world tree %s writer: %w", tree.Tag, err)
		}

		err = tree.Write(token)
		if err != nil {
			return fmt.Errorf("world tree %s: %w", tree.Tag, err)
		}
	}

	for _, aLight := range wce.AmbientLights {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set ambient light %s writer: %w", aLight.Tag, err)
		}

		err = aLight.Write(token)
		if err != nil {
			return fmt.Errorf("ambient light %s: %w", aLight.Tag, err)
		}
	}

	for _, actor := range wce.ActorInsts {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set actor %s writer: %w", actor.Tag, err)
		}

		err = actor.Write(token)
		if err != nil {
			return fmt.Errorf("actor %s: %w", actor.Tag, err)
		}
	}

	// for _, actorDef := range wce.ActorDefs {
	// 	tag := actorDef.Tag
	// 	if wce.WorldDef.Zone == 1 {
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

	for _, zone := range wce.Zones {
		err = token.SetWriter("world")
		if err != nil {
			return fmt.Errorf("set zone %s writer: %w", zone.Tag, err)
		}

		err = zone.Write(token)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	// EQG

	for _, mdsDef := range wce.MdsDefs {
		err = token.SetWriter(mdsDef.Tag)
		if err != nil {
			return fmt.Errorf("set mdsdef %s writer: %w", mdsDef.Tag, err)
		}

		err = mdsDef.Write(token)
		if err != nil {
			return fmt.Errorf("mdsdef %s: %w", mdsDef.Tag, err)
		}
	}

	for _, modDef := range wce.ModDefs {
		err = token.SetWriter(modDef.Tag)
		if err != nil {
			return fmt.Errorf("set moddef %s writer: %w", modDef.Tag, err)
		}

		err = modDef.Write(token)
		if err != nil {
			return fmt.Errorf("moddef %s: %w", modDef.Tag, err)
		}
	}

	token.Close()

	rootW, err := os.Create(fmt.Sprintf("%s/_root.wce", path))
	if err != nil {
		return err
	}
	wce.writeAsciiHeader(rootW)

	defer rootW.Close()

	if token.IsWriterUsed("world") {
		rootW.WriteString("INCLUDE \"WORLD.WCE\"\n")
	} else {
		err = os.Remove(fmt.Sprintf("%s/world.wce", path))
		if err != nil {
			return fmt.Errorf("remove %s: %w", fmt.Sprintf("%s/world.wce", path), err)
		}
	}

	if token.IsWriterUsed("region") {
		rootW.WriteString("INCLUDE \"REGION.WCE\"\n")
	} else {
		err = os.Remove(fmt.Sprintf("%s/region.wce", path))
		if err != nil {
			return fmt.Errorf("remove %s: %w", fmt.Sprintf("%s/region.wce", path), err)
		}
	}

	for _, baseTag := range baseTags {

		if baseTag != "PLAYER" &&
			!token.IsWriterUsed(baseTag) &&
			!token.IsWriterUsed(baseTag+"_ani") &&
			!strings.Contains(path, "_obj") {
			fmt.Println("Tag", baseTag, "was never used for model or ani (can be ignored)")
			//			return fmt.Errorf("tag %s was never used for model or ani", baseTag)
		}

		modelW, err := os.Create(fmt.Sprintf("%s/%s/_root.wce", path, strings.ToLower(baseTag)))
		if err != nil {
			return err
		}
		wce.writeAsciiHeader(modelW)

		defer modelW.Close()

		rootW.WriteString(fmt.Sprintf("INCLUDE \"%s/_ROOT.WCE\"\n", strings.ToUpper(baseTag)))

		if token.IsWriterUsed(baseTag) {
			_, err = modelW.WriteString(fmt.Sprintf("INCLUDE \"%s.WCE\"\n", strings.ToUpper(baseTag)))
			if err != nil {
				return err
			}
		} else {
			removePath := fmt.Sprintf("%s/%s/%s.wce", path, strings.ToLower(baseTag), strings.ToLower(baseTag))

			err = os.Remove(removePath)
			if err != nil {
				return fmt.Errorf("remove %s: %w", removePath, err)
			}
		}

		if token.IsWriterUsed(baseTag + "_ani") {
			_, err = modelW.WriteString(fmt.Sprintf("INCLUDE \"%s_ANI.WCE\"\n", strings.ToUpper(baseTag)))
			if err != nil {
				return fmt.Errorf("write %s: %w", fmt.Sprintf("%s/%s/%s_ani.wce", path, strings.ToLower(baseTag), strings.ToLower(baseTag)), err)
			}
		} else {
			removePath := fmt.Sprintf("%s/%s/%s_ani.wce", path, strings.ToLower(baseTag), strings.ToLower(baseTag))

			err = os.Remove(removePath)
			if err != nil {
				return fmt.Errorf("remove %s: %w", removePath, err)
			}
		}
	}

	return nil
}

func (wce *Wce) writeAsciiHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail %s\n", common.Version)
	fmt.Fprintf(w, "// Original file: %s\n\n", wce.FileName)
}
