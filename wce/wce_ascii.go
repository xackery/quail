package wce

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	/*
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
		} */

	err = wce.writeAsciiData(path)
	if err != nil {
		return err
	}

	return nil
}

func (wce *Wce) writeAsciiData(path string) error {

	token := NewAsciiWriteToken(path, wce)
	defer token.Close()

	if wce.WorldDef == nil {
		return fmt.Errorf("worlddef not found")
	}
	var err error
	if wce.WorldDef != nil {
		err = wce.WorldDef.Write(token)
		if err != nil {
			return fmt.Errorf("world def: %w", err)
		}
	}

	if wce.GlobalAmbientLightDef != nil {
		err = wce.GlobalAmbientLightDef.Write(token)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	for _, matDef := range wce.MaterialDefs {
		err = matDef.Write(token)
		if err != nil {
			return fmt.Errorf("materialdef %s: %w", matDef.Tag, err)
		}
	}

	for _, lightDef := range wce.LightDefs {
		err = lightDef.Write(token)
		if err != nil {
			return fmt.Errorf("lightdef %s: %w", lightDef.Tag, err)
		}
	}

	for _, polyDef := range wce.PolyhedronDefs {
		err = polyDef.Write(token)
		if err != nil {
			return fmt.Errorf("polyhedron %s: %w", polyDef.Tag, err)
		}
	}

	for _, dSprite := range wce.DMSpriteDefs {
		err = dSprite.Write(token)
		if err != nil {
			return fmt.Errorf("dmspritedef %s: %w", dSprite.Tag, err)
		}
	}

	for _, dSprite2 := range wce.DMSpriteDef2s {
		err = dSprite2.Write(token)
		if err != nil {
			return fmt.Errorf("dmspritedef2 %s: %w", dSprite2.Tag, err)
		}
	}

	for _, tree := range wce.WorldTrees {
		err = tree.Write(token)
		if err != nil {
			return fmt.Errorf("world tree %s: %w", tree.Tag, err)
		}
	}

	for _, region := range wce.Regions {
		err = region.Write(token)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	for _, pLight := range wce.PointLights {
		err = pLight.Write(token)
		if err != nil {
			return fmt.Errorf("point light (%s): %w", pLight.Tag, err)
		}
	}

	for _, blitSpriteDef := range wce.BlitSpriteDefs {
		err = blitSpriteDef.Write(token)
		if err != nil {
			return fmt.Errorf("blitspritedef %s: %w", blitSpriteDef.Tag, err)
		}
	}

	for _, particleCloudDef := range wce.ParticleCloudDefs {
		err = particleCloudDef.Write(token)
		if err != nil {
			return fmt.Errorf("particleclouddef %s: %w", particleCloudDef.Tag, err)
		}
	}

	for _, track := range wce.TrackInstances {
		err = track.Write(token)
		if err != nil {
			return fmt.Errorf("track %s_%d: %w", track.Tag, track.TagIndex, err)
		}
	}

	for _, hierarchicalSpriteDef := range wce.HierarchicalSpriteDefs {
		err = hierarchicalSpriteDef.Write(token)
		if err != nil {
			return fmt.Errorf("hierarchicalspritedef %s: %w", hierarchicalSpriteDef.Tag, err)
		}
	}

	for _, actorDef := range wce.ActorDefs {
		err = actorDef.Write(token)
		if err != nil {
			return fmt.Errorf("actordef %s: %w", actorDef.Tag, err)
		}
	}

	for _, actor := range wce.ActorInsts {
		err = actor.Write(token)
		if err != nil {
			return fmt.Errorf("actor %s: %w", actor.Tag, err)
		}
	}

	for _, zone := range wce.Zones {
		err = zone.Write(token)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	// EQG

	for _, matDef := range wce.EQMaterialDefs {
		err = matDef.Write(token)
		if err != nil {
			return fmt.Errorf("eqmaterialdef %s: %w", matDef.Tag, err)
		}
	}

	for _, mdsDef := range wce.MdsDefs {
		err = mdsDef.Write(token)
		if err != nil {
			return fmt.Errorf("mdsdef %s: %w", mdsDef.Tag, err)
		}
	}

	for _, modDef := range wce.ModDefs {
		err = modDef.Write(token)
		if err != nil {
			return fmt.Errorf("moddef %s: %w", modDef.Tag, err)
		}
	}

	for _, terDef := range wce.TerDefs {
		err = terDef.Write(token)
		if err != nil {
			return fmt.Errorf("terdef %s: %w", terDef.Tag, err)
		}
	}

	for _, aniDef := range wce.AniDefs {
		err = aniDef.Write(token)
		if err != nil {
			return fmt.Errorf("anidef %s: %w", aniDef.Tag, err)
		}
	}

	token.Close()

	type folderType struct {
		hasBase bool
		hasAni  bool
	}
	folders := make(map[string]*folderType)
	for key, w := range token.writers {
		if key == "world" || key == "region" {
			w.Close()
			continue
		}
		isAni := false
		if strings.Contains(key, "_ani") {
			key = strings.Replace(key, "_ani", "", 1)
			isAni = true
		}
		_, ok := folders[key]
		if !ok {
			folders[key] = &folderType{}
		}
		if isAni {
			folders[key].hasAni = true
			continue
		}
		folders[key].hasBase = true
	}

	rootW, err := os.Create(fmt.Sprintf("%s/_root.wce", path))
	if err != nil {
		return err
	}
	wce.writeAsciiHeader(rootW)

	if token.IsWriterUsed("region") {
		rootW.WriteString("INCLUDE \"REGION.WCE\"\n")
	}

	if token.IsWriterUsed("world") {
		rootW.WriteString("INCLUDE \"WORLD.WCE\"\n")
	}

	includes := make(map[string]string)

	sortedFolders := make([]string, 0)
	for folder := range folders {
		sortedFolders = append(sortedFolders, folder)
	}
	sort.Strings(sortedFolders)

	writtenSubfolders := make(map[string]bool)
	for _, folder := range sortedFolders {
		folderInfo, ok := folders[folder]
		if !ok {
			return fmt.Errorf("folder %s not found", folder)
		}

		if strings.Contains(folder, "/") {
			rootFolder := strings.Split(folder, "/")[0]
			tag := strings.Split(folder, "/")[1]
			if _, ok := writtenSubfolders[rootFolder]; !ok {
				rootW.WriteString(fmt.Sprintf("INCLUDE \"%s/_ROOT.WCE\"\n", strings.ToUpper(rootFolder)))
			}

			if folderInfo.hasBase {
				includes[rootFolder] += fmt.Sprintf("INCLUDE \"%s.WCE\"\n", strings.ToUpper(tag))
			}

			if folderInfo.hasAni {
				includes[rootFolder] += fmt.Sprintf("INCLUDE \"%s_ANI.WCE\"\n", strings.ToUpper(tag))
			}

			writtenSubfolders[rootFolder] = true

			continue
		}

		rootW.WriteString(fmt.Sprintf("INCLUDE \"%s/_ROOT.WCE\"\n", strings.ToUpper(folder)))
		if folderInfo.hasBase {
			includes[folder] += fmt.Sprintf("INCLUDE \"%s.WCE\"\n", strings.ToUpper(folder))
		}
		if folderInfo.hasAni {
			includes[folder] += fmt.Sprintf("INCLUDE \"%s_ANI.WCE\"\n", strings.ToUpper(folder))
		}
	}
	rootW.Close()

	for folder, out := range includes {
		w, err := os.Create(fmt.Sprintf("%s/%s/_root.wce", path, strings.ToLower(folder)))
		if err != nil {
			return err
		}
		wce.writeAsciiHeader(w)
		w.WriteString(out)
		w.Close()
	}

	return nil
}

func (wce *Wce) writeAsciiHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail %s\n", common.Version)
	fmt.Fprintf(w, "// Original file: %s\n\n", wce.FileName)
}
