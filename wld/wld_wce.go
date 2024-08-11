package wld

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
)

// ReadAscii reads the ascii file at path
func (wld *Wld) ReadAscii(path string) error {

	wld.reset()
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
	for _, dmSprite := range wld.DMSpriteDef2s {
		if len(dmSprite.Tag) < 3 {
			return fmt.Errorf("dmspritedef2 %s tag too short", dmSprite.Tag)
		}
		baseTag := baseTagTrim(dmSprite.Tag)
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
	for _, dmSprite := range wld.DMSpriteDefs {
		if len(dmSprite.Tag) < 3 {
			return fmt.Errorf("dmsprite %s tag too short", dmSprite.Tag)
		}
		baseTag := baseTagTrim(dmSprite.Tag)
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

	rootBuf, err := os.Create(fmt.Sprintf("%s/_root.wce", path))
	if err != nil {
		return err
	}
	wld.writeAsciiHeader(rootBuf)
	err = wld.writeAsciiData(path, baseTags, rootBuf)
	if err != nil {
		return err
	}
	defer rootBuf.Close()

	data, err := os.ReadFile(fmt.Sprintf("%s/zone.mod", path))
	if err != nil {
		return fmt.Errorf("read %s: %w", fmt.Sprintf("%s/zone.mod", path), err)
	}

	if len(data) < 60 {
		err = os.Remove(fmt.Sprintf("%s/zone.mod", path))
		if err != nil {
			return fmt.Errorf("remove %s: %w", fmt.Sprintf("%s/zone.mod", path), err)
		}
	} else {
		rootBuf.WriteString("INCLUDE \"ZONE.MOD\"\n")
	}
	for _, tag := range baseTags {
		// read .ani files
		aniPath := fmt.Sprintf("%s/%s/%s.ani", path, strings.ToLower(tag), strings.ToLower(tag))
		data, err := os.ReadFile(aniPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", aniPath, err)
		}
		if len(data) < 60 {
			err = os.Remove(aniPath)
			if err != nil {
				return fmt.Errorf("remove %s: %w", aniPath, err)
			}

			buf := &bytes.Buffer{}
			wld.writeAsciiHeader(buf)
			fmt.Fprintf(buf, "INCLUDE \"%s.MOD\"\n", strings.ToUpper(tag))

			err = os.WriteFile(fmt.Sprintf("%s/%s/_root.wce", path, strings.ToLower(tag)), buf.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("write %s: %w", fmt.Sprintf("%s/%s/_root.wce", path, strings.ToLower(tag)), err)
			}
		}

	}

	return nil
}

func (wld *Wld) writeAsciiData(path string, baseTags []string, rootBuf *os.File) error {
	var w io.Writer
	var ok bool

	lightsMap := map[string]bool{}
	zoneMaterials := map[string]bool{}
	modWriters := map[string]*os.File{}
	defsWritten := map[string]bool{}
	modDefsWritten := map[string]bool{}
	zoneName := filepath.Base(path)

	for _, track := range wld.TrackInstances {
		if len(track.Tag) < 3 {
			return fmt.Errorf("trackdef %s tag too short", track.Tag)
		}
		baseTag := wld.aniWriterTag(track.Tag)
		if track.modelTag != "" {
			baseTag = track.modelTag
		}
		if len(baseTag) < 1 {
			return fmt.Errorf("trackd %s tag too short (%s)", track.Tag, baseTag)
		}
		isFound := false
		for _, tag := range baseTags {
			if tag == baseTag || tag == track.modelTag {
				isFound = true
				break
			}
		}
		if !isFound {
			baseTags = append(baseTags, baseTag)
		}
	}

	zoneBuf, err := os.Create(fmt.Sprintf("%s/zone.mod", path))
	if err != nil {
		return err
	}
	wld.writeAsciiHeader(zoneBuf)

	modWriters[zoneName] = zoneBuf
	for _, baseTag := range baseTags {
		if wld.isZone {

			modWriters[baseTag+"_ani"] = zoneBuf
			modWriters[baseTag+"_root"] = rootBuf
		} else {
			rootBuf.WriteString(fmt.Sprintf("INCLUDE \"%s/_ROOT.WCE\"\n", strings.ToUpper(baseTag)))

			err = os.MkdirAll(fmt.Sprintf("%s/%s", path, strings.ToLower(baseTag)), os.ModePerm)
			if err != nil {
				return fmt.Errorf("mkdir: %w", err)
			}

			wceBuf, err := os.Create(fmt.Sprintf("%s/%s/_root.wce", path, strings.ToLower(baseTag)))
			if err != nil {
				return err
			}
			defer wceBuf.Close()
			wld.writeAsciiHeader(wceBuf)
			modWriters[baseTag+"_root"] = wceBuf

			wceBuf.WriteString(fmt.Sprintf("INCLUDE \"%s.MOD\"\n", strings.ToUpper(baseTag)))
			wceBuf.WriteString(fmt.Sprintf("INCLUDE \"%s.ANI\"\n", strings.ToUpper(baseTag)))

			modBuf, err := os.Create(fmt.Sprintf("%s/%s/%s.mod", path, strings.ToLower(baseTag), strings.ToLower(baseTag)))
			if err != nil {
				return err
			}
			defer modBuf.Close()
			wld.writeAsciiHeader(modBuf)
			modWriters[baseTag] = modBuf

			aniBuf, err := os.Create(fmt.Sprintf("%s/%s/%s.ani", path, strings.ToLower(baseTag), strings.ToLower(baseTag)))
			if err != nil {
				return err
			}
			defer aniBuf.Close()
			wld.writeAsciiHeader(aniBuf)
			modWriters[baseTag+"_ani"] = aniBuf
		}
	}

	w = modWriters[zoneName]

	if wld.WorldDef != nil {
		err = wld.WorldDef.Write(w)
		if err != nil {
			return fmt.Errorf("world def: %w", err)
		}
	}

	if wld.GlobalAmbientLightDef != nil {
		err = wld.GlobalAmbientLightDef.Write(w)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
		wld.isZone = true
	}

	for i := 0; i < len(wld.DMSpriteDef2s); i++ {
		dmSprite := wld.DMSpriteDef2s[i]
		baseTag := baseTagTrim(dmSprite.Tag)
		if !wld.isZone {
			modDefsWritten = map[string]bool{}
		}
		w, ok = modWriters[baseTag]
		if !ok {
			return fmt.Errorf("dmsprite %s writer not found (basetag %s)", dmSprite.Tag, baseTag)
		}
		defsWritten[dmSprite.Tag] = true

		isZoneChunk := false
		// if baseTag is r### then foo
		if strings.HasPrefix(baseTag, "r") {
			regionChunk := 0
			chunkCount, err := fmt.Sscanf(baseTag, "r%d", &regionChunk)
			isZoneChunk = err == nil && chunkCount == 1
		}

		if dmSprite.MaterialPaletteTag != "" && !zoneMaterials[dmSprite.MaterialPaletteTag] {

			isMaterialPaletteFound := false
			for _, materialPal := range wld.MaterialPalettes {
				if materialPal.Tag != dmSprite.MaterialPaletteTag {
					continue
				}

				if modDefsWritten[materialPal.Tag] {
					continue
				}
				modDefsWritten[materialPal.Tag] = true

				for _, materialTag := range materialPal.Materials {
					isMaterialDefFound := false
					for _, materialDef := range wld.MaterialDefs {
						if materialDef.Tag != materialTag {
							continue
						}

						isMaterialDefFound = true

						if modDefsWritten[materialDef.Tag] {
							continue
						}
						modDefsWritten[materialDef.Tag] = true
						defsWritten[materialDef.Tag] = true

						if materialDef.SimpleSpriteTag != "" {
							isSimpleSpriteFound := false
							for _, simpleSprite := range wld.SimpleSpriteDefs {
								if simpleSprite.Tag != materialDef.SimpleSpriteTag {
									continue
								}
								isSimpleSpriteFound = true
								if modDefsWritten[simpleSprite.Tag] {
									continue
								}
								modDefsWritten[simpleSprite.Tag] = true
								err = simpleSprite.Write(w)
								if err != nil {
									return fmt.Errorf("simple sprite %s: %w", simpleSprite.Tag, err)
								}
								break
							}
							if !isSimpleSpriteFound {
								return fmt.Errorf("simple sprite %s not found", materialDef.SimpleSpriteTag)
							}
						}

						err = materialDef.Write(w)
						if err != nil {
							return fmt.Errorf("material %s: %w", materialDef.Tag, err)
						}
						break
					}
					if !isMaterialDefFound {
						return fmt.Errorf("dmspritedef2 %s materialdef %s not found", dmSprite.Tag, materialTag)
					}
				}

				isMaterialPaletteFound = true
				err = materialPal.Write(w)
				if err != nil {
					return fmt.Errorf("material palette %s: %w", materialPal.Tag, err)
				}
				break
			}
			if !isMaterialPaletteFound {
				return fmt.Errorf("material palette %s not found", dmSprite.MaterialPaletteTag)
			}
			zoneMaterials[dmSprite.MaterialPaletteTag] = true
		}

		if dmSprite.PolyhedronTag != "" && dmSprite.PolyhedronTag != "NEGATIVE_TWO" {
			poly := wld.ByTag(dmSprite.PolyhedronTag)
			if poly == nil {
				return fmt.Errorf("dmsprite %s polyhedron %s not found", dmSprite.Tag, dmSprite.PolyhedronTag)
			}
			switch polyDef := poly.(type) {
			case *PolyhedronDefinition:
				err = polyDef.Write(w)
				if err != nil {
					return fmt.Errorf("polyhedron %s: %w", polyDef.Tag, err)
				}
				defsWritten[polyDef.Tag] = true
			case *Sprite3DDef:
				err = polyDef.Write(w)
				if err != nil {
					return fmt.Errorf("sprite 3d %s: %w", polyDef.Tag, err)
				}
				defsWritten[polyDef.Tag] = true
			default:
				return fmt.Errorf("dmsprite %s polyhedron %s unknown type %T", dmSprite.Tag, dmSprite.PolyhedronTag, poly)
			}

		}

		err = dmSprite.Write(w)
		if err != nil {
			return fmt.Errorf("dm sprite def %s: %w", dmSprite.Tag, err)
		}
		if !isZoneChunk {
			//		fmt.Fprintf(rootBuf, "INCLUDE \"%s.MOD\"\n", strings.ToUpper(baseTag))
		}
		tracksWritten := map[string]bool{}
		for _, hierarchySprite := range wld.HierarchicalSpriteDefs {
			isFound := false

			if hierarchySprite.PolyhedronTag == dmSprite.Tag {
				isFound = true
			}
			if !isFound {
				for _, skin := range hierarchySprite.AttachedSkins {
					if skin.DMSpriteTag != dmSprite.Tag {
						continue
					}
					isFound = true
					break
				}
			}
			if !isFound {
				for _, dag := range hierarchySprite.Dags {
					if dag.SpriteTag == dmSprite.Tag {
						isFound = true
						break
					}
				}
			}
			if !isFound {
				continue
			}

			if defsWritten[hierarchySprite.Tag] {
				break
			}
			defsWritten[hierarchySprite.Tag] = true

			err = hierarchySprite.Write(w)
			if err != nil {
				return fmt.Errorf("hierarchical sprite %s: %w", hierarchySprite.Tag, err)
			}

			for _, dag := range hierarchySprite.Dags {
				if dag.Track == "" {
					continue
				}
				trackUUID := fmt.Sprintf("%s_%d", dag.Track, dag.TrackIndex)
				if tracksWritten[trackUUID] {
					continue
				}

				isTrackFound := false
				for _, track := range wld.TrackInstances {
					if track.Tag != dag.Track && track.modelTag != dag.Track {
						continue
					}

					if track.TagIndex != dag.TrackIndex {
						continue
					}

					defsWritten[trackUUID] = true

					isTrackDefFound := false

					for _, trackDef := range wld.TrackDefs {
						trackDefUUID := fmt.Sprintf("%s_%d", trackDef.Tag, trackDef.TagIndex)

						if trackDef.Tag != track.DefinitionTag {
							continue
						}
						if trackDef.TagIndex != track.DefinitionTagIndex {
							continue
						}
						isTrackDefFound = true

						if tracksWritten[trackDefUUID] {
							break
						}

						err = trackDef.Write(w)
						if err != nil {
							return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
						}

						tracksWritten[trackDefUUID] = true
						defsWritten[trackDefUUID] = true
						break
					}
					if !isTrackDefFound {
						return fmt.Errorf("hierarchy %s track %s definition not found", hierarchySprite.Tag, track.DefinitionTag)
					}

					isTrackFound = true

					tracksWritten[trackUUID] = true
					defsWritten[trackUUID] = true

					err = track.Write(w)
					if err != nil {
						return fmt.Errorf("track %s: %w", track.Tag, err)
					}
				}
				if !isTrackFound {
					return fmt.Errorf("dnspritedef2 hierarchy %s track %s not found", hierarchySprite.Tag, dag.Track)
				}
			}

			break
		}

	}

	for i := 0; i < len(wld.DMSpriteDefs); i++ {
		dmSprite := wld.DMSpriteDefs[i]
		baseTag := baseTagTrim(dmSprite.Tag)
		if !wld.isZone {
			modDefsWritten = map[string]bool{}
		}
		w, ok = modWriters[baseTag]
		if !ok {
			return fmt.Errorf("dmsprite %s writer not found (basetag %s)", dmSprite.Tag, baseTag)
		}
		defsWritten[dmSprite.Tag] = true

		isZoneChunk := false
		// if baseTag is r### then foo
		if strings.HasPrefix(baseTag, "r") {
			regionChunk := 0
			chunkCount, err := fmt.Sscanf(baseTag, "r%d", &regionChunk)
			isZoneChunk = err == nil && chunkCount == 1
		}

		if dmSprite.MaterialPaletteTag != "" && !zoneMaterials[dmSprite.MaterialPaletteTag] {

			isMaterialPaletteFound := false
			for _, materialPal := range wld.MaterialPalettes {
				if materialPal.Tag != dmSprite.MaterialPaletteTag {
					continue
				}

				if modDefsWritten[materialPal.Tag] {
					continue
				}
				modDefsWritten[materialPal.Tag] = true

				for _, materialTag := range materialPal.Materials {
					isMaterialDefFound := false
					for _, materialDef := range wld.MaterialDefs {
						if materialDef.Tag != materialTag {
							continue
						}

						if modDefsWritten[materialDef.Tag] {
							continue
						}
						modDefsWritten[materialDef.Tag] = true
						defsWritten[materialDef.Tag] = true

						if materialDef.SimpleSpriteTag != "" {
							isSimpleSpriteFound := false
							for _, simpleSprite := range wld.SimpleSpriteDefs {
								if simpleSprite.Tag != materialDef.SimpleSpriteTag {
									continue
								}
								isSimpleSpriteFound = true
								if modDefsWritten[simpleSprite.Tag] {
									continue
								}
								modDefsWritten[simpleSprite.Tag] = true
								err = simpleSprite.Write(w)
								if err != nil {
									return fmt.Errorf("simple sprite %s: %w", simpleSprite.Tag, err)
								}
								break
							}
							if !isSimpleSpriteFound {
								return fmt.Errorf("simple sprite %s not found", materialDef.SimpleSpriteTag)
							}
						}

						isMaterialDefFound = true
						err = materialDef.Write(w)
						if err != nil {
							return fmt.Errorf("material %s: %w", materialDef.Tag, err)
						}
						break
					}
					if !isMaterialDefFound {
						return fmt.Errorf("dmsprite %s materialdef %s not found", dmSprite.Tag, materialTag)
					}
				}

				isMaterialPaletteFound = true
				err = materialPal.Write(w)
				if err != nil {
					return fmt.Errorf("material palette %s: %w", materialPal.Tag, err)
				}
				break
			}
			if !isMaterialPaletteFound {
				return fmt.Errorf("material palette %s not found", dmSprite.MaterialPaletteTag)
			}
			zoneMaterials[dmSprite.MaterialPaletteTag] = true
		}

		err = dmSprite.Write(w)
		if err != nil {
			return fmt.Errorf("dm sprite def %s: %w", dmSprite.Tag, err)
		}
		if !isZoneChunk {
			//		fmt.Fprintf(rootBuf, "INCLUDE \"%s.MOD\"\n", strings.ToUpper(baseTag))
		}
		tracksWritten := map[string]bool{}
		for _, hierarchySprite := range wld.HierarchicalSpriteDefs {
			isFound := false

			if hierarchySprite.PolyhedronTag == dmSprite.Tag {
				isFound = true
			}
			if !isFound {
				for _, skin := range hierarchySprite.AttachedSkins {
					if skin.DMSpriteTag != dmSprite.Tag {
						continue
					}
					isFound = true
					break
				}
			}
			if !isFound {
				for _, dag := range hierarchySprite.Dags {
					if dag.SpriteTag == dmSprite.Tag {
						isFound = true
						break
					}
				}
			}
			if !isFound {
				continue
			}
			if defsWritten[hierarchySprite.Tag] {
				break
			}
			defsWritten[hierarchySprite.Tag] = true

			err = hierarchySprite.Write(w)
			if err != nil {
				return fmt.Errorf("hierarchical sprite %s: %w", hierarchySprite.Tag, err)
			}

			for _, dag := range hierarchySprite.Dags {
				if dag.Track == "" {
					continue
				}

				trackUUID := fmt.Sprintf("%s_%d", dag.Track, dag.TrackIndex)
				if tracksWritten[trackUUID] {
					continue
				}

				isTrackFound := false
				for _, track := range wld.TrackInstances {
					if track.Tag != dag.Track && track.modelTag != dag.Track {
						continue
					}
					if track.TagIndex != dag.TrackIndex {
						continue
					}

					trackUUID := fmt.Sprintf("%s_%d", track.DefinitionTag, track.DefinitionTagIndex)

					if defsWritten[trackUUID] {
						isTrackFound = true
						break
					}

					defsWritten[trackUUID] = true

					isTrackDefFound := false

					for _, trackDef := range wld.TrackDefs {
						if trackDef.Tag != track.DefinitionTag {
							continue
						}
						if trackDef.TagIndex != track.DefinitionTagIndex {
							continue
						}
						isTrackDefFound = true

						trackDefUUID := fmt.Sprintf("%s_%d", trackDef.Tag, trackDef.TagIndex)

						if tracksWritten[trackDefUUID] {
							break
						}

						err = trackDef.Write(w)
						if err != nil {
							return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
						}

						tracksWritten[trackDefUUID] = true
						defsWritten[trackDefUUID] = true
						break
					}
					if !isTrackDefFound {
						return fmt.Errorf("dmsprite hierarchy %s track %s definition not found", hierarchySprite.Tag, track.DefinitionTag)
					}

					isTrackFound = true

					tracksWritten[trackUUID] = true
					defsWritten[trackUUID] = true

					err = track.Write(w)
					if err != nil {
						return fmt.Errorf("track %s: %w", track.Tag, err)
					}
				}
				if !isTrackFound {
					return fmt.Errorf("dmsprite hierarchy %s track %s not found", hierarchySprite.Tag, dag.Track)
				}
			}

			break
		}
	}

	for i := 0; i < len(wld.TrackInstances); i++ {
		track := wld.TrackInstances[i]
		if defsWritten[track.Tag] {
			continue
		}
		baseTag := wld.aniWriterTag(track.Tag)

		w, ok = modWriters[track.modelTag+"_ani"]
		if !ok {
			return fmt.Errorf("track %s writer not found (basetag %s)", track.Tag, baseTag)
		}

		isTrackFound := false
		for _, trackDef := range wld.TrackDefs {
			if trackDef.Tag != track.DefinitionTag {
				continue
			}
			if trackDef.TagIndex != track.DefinitionTagIndex {
				continue
			}

			err = trackDef.Write(w)
			if err != nil {
				return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
			}

			isTrackFound = true
			break
		}
		if !isTrackFound {
			return fmt.Errorf("track %s definition not found", track.DefinitionTag)
		}

		err = track.Write(w)
		if err != nil {
			return fmt.Errorf("track %s: %w", track.Tag, err)
		}
	}

	for i := 0; i < len(wld.PolyhedronDefs); i++ {
		polyhedron := wld.PolyhedronDefs[i]
		if defsWritten[polyhedron.Tag] {
			continue
		}
		baseTag := baseTagTrim(polyhedron.Tag)
		w, ok = modWriters[baseTag]
		if !ok {
			return fmt.Errorf("polyhedron %s writer not found (basetag %s)", polyhedron.Tag, baseTag)
		}
		err = polyhedron.Write(w)
		if err != nil {
			return fmt.Errorf("polyhedron %s: %w", polyhedron.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.PointLights); i++ {
		pointLight := wld.PointLights[i]

		if pointLight.LightTag != "" {
			isLightFound := false
			for _, lightDef := range wld.LightDefs {
				if lightDef.Tag != pointLight.LightTag {
					continue
				}

				lightsMap[lightDef.Tag] = true
				isLightFound = true
				err = lightDef.Write(w)
				if err != nil {
					return fmt.Errorf("light def %s: %w", lightDef.Tag, err)
				}
				break
			}
			if !isLightFound {
				return fmt.Errorf("point light %s light %s not found", pointLight.Tag, pointLight.LightTag)
			}
		}

		lightsMap[pointLight.Tag] = true
		err = pointLight.Write(w)
		if err != nil {
			return fmt.Errorf("point light %s: %w", pointLight.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.MaterialDefs); i++ {
		materialDef := wld.MaterialDefs[i]
		if defsWritten[materialDef.Tag] {
			continue
		}
		defsWritten[materialDef.Tag] = true
		baseTag := baseTagTrim(materialDef.Tag)
		if modWriters[baseTag] != nil {
			w, ok = modWriters[baseTag]
			if !ok {
				return fmt.Errorf("material def %s writer not found (basetag %s)", materialDef.Tag, baseTag)
			}
		}

		if materialDef.SimpleSpriteTag != "" {
			isSimpleSpriteFound := false
			for _, simpleSprite := range wld.SimpleSpriteDefs {
				if simpleSprite.Tag != materialDef.SimpleSpriteTag {
					continue
				}
				isSimpleSpriteFound = true
				if modDefsWritten[simpleSprite.Tag] {
					continue
				}
				modDefsWritten[simpleSprite.Tag] = true
				err = simpleSprite.Write(w)
				if err != nil {
					return fmt.Errorf("simple sprite %s: %w", simpleSprite.Tag, err)
				}
				break
			}
			if !isSimpleSpriteFound {
				return fmt.Errorf("simple sprite %s not found", materialDef.SimpleSpriteTag)
			}
		}
		err = materialDef.Write(w)
		if err != nil {
			return fmt.Errorf("material def %s: %w", materialDef.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.LightDefs); i++ {
		lightDef := wld.LightDefs[i]
		if lightsMap[lightDef.Tag] {
			continue
		}

		lightsMap[lightDef.Tag] = true
		err = lightDef.Write(w)
		if err != nil {
			return fmt.Errorf("light def %s: %w", lightDef.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.WorldTrees); i++ {
		worldTree := wld.WorldTrees[i]

		err = worldTree.Write(w)
		if err != nil {
			return fmt.Errorf("world tree %s: %w", worldTree.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.Regions); i++ {
		region := wld.Regions[i]

		err = region.Write(w)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.AmbientLights); i++ {
		ambientLight := wld.AmbientLights[i]

		err = ambientLight.Write(w)
		if err != nil {
			return fmt.Errorf("ambient light %s: %w", ambientLight.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.ActorInsts); i++ {
		actor := wld.ActorInsts[i]
		baseTag := baseTagTrim(actor.Tag)
		if modWriters[baseTag] != nil {
			w, ok = modWriters[baseTag]
			if !ok {
				return fmt.Errorf("actor %s writer not found (basetag %s)", actor.Tag, baseTag)
			}
		}

		if actor.DMRGBTrackTag.Valid {
			isTrackDefFound := false
			for _, trackDef := range wld.RGBTrackDefs {
				if trackDef.Tag != actor.DMRGBTrackTag.String {
					continue
				}

				isTrackDefFound = true
				err = trackDef.Write(w)
				if err != nil {
					return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
				}
				break
			}

			if !isTrackDefFound {
				return fmt.Errorf("actor '%s' track %s definition not found", actor.Tag, actor.DMRGBTrackTag.String)
			}
		}

		if actor.DefinitionTag == "!UNK" {
			return fmt.Errorf("actor %s definition not found", actor.DefinitionTag)
		}

		if actor.DefinitionTag != "" {
			isActorDefFound := false
			for j := 0; j < len(wld.ActorDefs); j++ {
				actorDef := wld.ActorDefs[j]
				if actorDef.Tag != actor.DefinitionTag {
					continue
				}

				baseTag = baseTagTrim(actorDef.Tag)
				if modWriters[baseTag] != nil {
					w, ok = modWriters[baseTag]
					if !ok {
						return fmt.Errorf("actor def %s writer not found (basetag %s)", actorDef.Tag, baseTag)
					}
				}

				defsWritten[actorDef.Tag] = true

				for _, action := range actorDef.Actions {
					for _, lod := range action.LevelOfDetails {
						if lod.SpriteTag == "" {
							continue
						}

						spriteFrag := wld.ByTag(lod.SpriteTag)
						if spriteFrag == nil {
							return fmt.Errorf("actorinst %s actor %s sprite %s not found", actor.Tag, actorDef.Tag, lod.SpriteTag)
						}

						switch sprite := spriteFrag.(type) {
						case *SimpleSpriteDef:
							err = sprite.Write(w)
							if err != nil {
								return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
							}
						case *Sprite3DDef:
							err = sprite.Write(w)
							if err != nil {
								return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
							}
						case *HierarchicalSpriteDef:
							if defsWritten[sprite.Tag] {
								continue
							}
							defsWritten[sprite.Tag] = true
							err = sprite.Write(w)
							if err != nil {
								return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
							}
						default:
							return fmt.Errorf("actorInst '%s' actorDef %s unknown sprite type %T", actor.Tag, actorDef.Tag, sprite)
						}
					}
				}

				isActorDefFound = true
				err = actorDef.Write(w)
				if err != nil {
					return fmt.Errorf("actor def %s: %w", actor.Tag, err)
				}
				break
			}
			if !isActorDefFound {
				//	fmt.Printf("actor %s definition not found\n", actor.DefinitionTag)
				// return fmt.Errorf("actor %s definition not found", actor.DefinitionTag)
			}
		}

		err = actor.Write(w)
		if err != nil {
			return fmt.Errorf("actor %s: %w", actor.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.ActorDefs); i++ {
		actorDef := wld.ActorDefs[i]
		if defsWritten[actorDef.Tag] {
			continue
		}
		baseTag := baseTagTrim(actorDef.Tag)
		if modWriters[baseTag] != nil {
			w, ok = modWriters[baseTag]
			if !ok {
				return fmt.Errorf("actor def %s writer not found (basetag %s)", actorDef.Tag, baseTag)
			}
		}
		for _, action := range actorDef.Actions {
			for _, lod := range action.LevelOfDetails {
				if lod.SpriteTag == "" {
					continue
				}

				spriteFrag := wld.ByTag(lod.SpriteTag)
				if spriteFrag == nil {
					return fmt.Errorf("actor %s sprite %s not found", actorDef.Tag, lod.SpriteTag)
				}

				switch sprite := spriteFrag.(type) {
				case *SimpleSpriteDef:
					err = sprite.Write(w)
					if err != nil {
						return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
					}
					defsWritten[sprite.Tag] = true
				case *Sprite3DDef:
					err = sprite.Write(w)
					if err != nil {
						return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
					}
				case *HierarchicalSpriteDef:
					if defsWritten[sprite.Tag] {
						continue
					}
					defsWritten[sprite.Tag] = true
					err = sprite.Write(w)
					if err != nil {
						return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
					}
				case *DMSpriteDef2:
					if !defsWritten[sprite.Tag] {
						err = sprite.Write(w)
						if err != nil {
							return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
						}
						defsWritten[sprite.Tag] = true
					}
				case *BlitSpriteDefinition:
					if !defsWritten[sprite.Tag] {
						err = sprite.Write(w)
						if err != nil {
							return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
						}
						defsWritten[sprite.Tag] = true
					}
				case *Sprite2DDef:
					if !defsWritten[sprite.Tag] {
						err = sprite.Write(w)
						if err != nil {
							return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
						}
						defsWritten[sprite.Tag] = true
					}
				default:
					return fmt.Errorf("actordef %s refs unknown sprite %s with type %T", actorDef.Tag, lod.SpriteTag, sprite)
				}

			}
		}
		err = actorDef.Write(w)
		if err != nil {
			return fmt.Errorf("actor def %s: %w", actorDef.Tag, err)
		}
	}

	w = modWriters[zoneName]
	for i := 0; i < len(wld.Zones); i++ {
		zone := wld.Zones[i]

		err = zone.Write(w)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	for _, writer := range modWriters {
		writer.Close()
	}

	return nil
}

func (wld *Wld) writeAsciiHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail v%s\n", common.Version)
	fmt.Fprintf(w, "// Original file: %s\n\n", wld.FileName)
}

func (wld *Wld) aniWriterTag(name string) string {

	base := strings.TrimSuffix(name, "_TRACKDEF")
	base = strings.TrimSuffix(base, "_TRACK")
	if isAnimationPrefix(base) {
		base = base[3:]
	}

	if len(base) < 4 {
		return base
	}

	if len(base) == 3 {
		if raceIDs[base[:3]] {
			return base[:3]
		}
	}
	for race := range raceIDs {
		if base == race {
			return base
		}
	}

	if len(base) > 3 {
		base = base[:3]
	}

	for _, sprite := range wld.DMSpriteDef2s {
		spriteName := strings.TrimSuffix(sprite.Tag, "_DMSPRITEDEF")
		if strings.HasPrefix(spriteName, base) {
			return base
		}
		if strings.HasSuffix(spriteName, base) {
			return base
		}
	}

	return base
}

var raceIDs = map[string]bool{
	"AAM":          true,
	"ABH":          true,
	"AEL":          true,
	"AHF":          true,
	"AHM":          true,
	"AIE":          true,
	"AKF":          true,
	"AKM":          true,
	"AKN":          true,
	"ALA":          true,
	"ALG":          true,
	"ALL":          true,
	"ALR":          true,
	"AMP":          true,
	"AMY":          true,
	"ANS":          true,
	"APX":          true,
	"ARM":          true,
	"ARO":          true,
	"ARROW":        true,
	"ASM":          true,
	"AVI":          true,
	"AVK":          true,
	"AXA":          true,
	"B01":          true,
	"B02":          true,
	"B03":          true,
	"B04":          true,
	"B05":          true,
	"B06":          true,
	"B07":          true,
	"B08":          true,
	"B09":          true,
	"B10":          true,
	"BAC":          true,
	"BAF":          true,
	"BAL":          true,
	"BAM":          true,
	"BAR":          true,
	"BAS":          true,
	"BAT":          true,
	"BDR":          true,
	"BEA":          true,
	"BEH":          true,
	"BEL":          true,
	"BER":          true,
	"BET":          true,
	"BFC":          true,
	"BFF":          true,
	"BFR":          true,
	"BGB":          true,
	"BGF":          true,
	"BGG":          true,
	"BGM":          true,
	"BIX":          true,
	"BKD":          true,
	"BKN":          true,
	"BLV":          true,
	"BNF":          true,
	"BNM":          true,
	"BNR":          true,
	"BNX":          true,
	"BNY":          true,
	"BOAT":         true,
	"BON":          true,
	"BOX":          true,
	"BRC":          true,
	"BRE":          true,
	"BRF":          true,
	"BRI":          true,
	"BRL":          true,
	"BRM":          true,
	"BRN":          true,
	"BRV":          true,
	"BRX":          true,
	"BSE":          true,
	"BSG":          true,
	"BTL":          true,
	"BTM":          true,
	"BTN":          true,
	"BTP":          true,
	"BTS":          true,
	"BTX":          true,
	"BUB":          true,
	"BUR":          true,
	"BUU":          true,
	"BVK":          true,
	"BWD":          true,
	"BXI":          true,
	"CAK":          true,
	"CAM":          true,
	"CAT":          true,
	"CAZ":          true,
	"CCD":          true,
	"CDF":          true,
	"CDM":          true,
	"CDR":          true,
	"CEF":          true,
	"CEM":          true,
	"CEN":          true,
	"CHM":          true,
	"CHT":          true,
	"CLA":          true,
	"CLB":          true,
	"CLF":          true,
	"CLG":          true,
	"CLM":          true,
	"CLN":          true,
	"CLQ":          true,
	"CLS":          true,
	"CLV":          true,
	"CLW":          true,
	"CNP":          true,
	"CNT":          true,
	"COC":          true,
	"COF":          true,
	"COK":          true,
	"COM":          true,
	"CPF":          true,
	"CPM":          true,
	"CPT":          true,
	"CRB":          true,
	"CRH":          true,
	"CRL":          true,
	"CRO":          true,
	"CRS":          true,
	"CRY":          true,
	"CSE":          true,
	"CST":          true,
	"CTH":          true,
	"CUB":          true,
	"CWB":          true,
	"CWG":          true,
	"CWR":          true,
	"DAF":          true,
	"DAM":          true,
	"DBP":          true,
	"DBX":          true,
	"DCF":          true,
	"DCM":          true,
	"DDM":          true,
	"DDV":          true,
	"DEN":          true,
	"DER":          true,
	"DEV":          true,
	"DIA":          true,
	"DIW":          true,
	"DJI":          true,
	"DKE":          true,
	"DKF":          true,
	"DKM":          true,
	"DLK":          true,
	"DMA":          true,
	"DML":          true,
	"DPF":          true,
	"DPM":          true,
	"DR2":          true,
	"DRA":          true,
	"DRC":          true,
	"DRE":          true,
	"DRF":          true,
	"DRG":          true,
	"DRI":          true,
	"DRK":          true,
	"DRL":          true,
	"DRM":          true,
	"DRN":          true,
	"DRP":          true,
	"DRS":          true,
	"DRU":          true,
	"DRV":          true,
	"DSB":          true,
	"DSF":          true,
	"DSG":          true,
	"DV6":          true,
	"DVF":          true,
	"DVL":          true,
	"DVM":          true,
	"DVS":          true,
	"DWF":          true,
	"DWM":          true,
	"DYN":          true,
	"EAE":          true,
	"ECS":          true,
	"EEF":          true,
	"EEL":          true,
	"EEM":          true,
	"EEY":          true,
	"EFE":          true,
	"EFR":          true,
	"EGF":          true,
	"EGM":          true,
	"ELE":          true,
	"ELF":          true,
	"ELM":          true,
	"EMP":          true,
	"ENA":          true,
	"EPF":          true,
	"EPM":          true,
	"ERF":          true,
	"ERM":          true,
	"ERO":          true,
	"ERU":          true,
	"EVE":          true,
	"EXO":          true,
	"EYE":          true,
	"FAF":          true,
	"FAM":          true,
	"FAN":          true,
	"FBF":          true,
	"FBM":          true,
	"FDR":          true,
	"FEF":          true,
	"FEL":          true,
	"FEM":          true,
	"FEN":          true,
	"FGG":          true,
	"FGH":          true,
	"FGI":          true,
	"FGP":          true,
	"FGT":          true,
	"FIE":          true,
	"FIS":          true,
	"FKN":          true,
	"FLG":          true,
	"FMO":          true,
	"FMP":          true,
	"FMT":          true,
	"FNG":          true,
	"FPF":          true,
	"FPM":          true,
	"FRA":          true,
	"FRD":          true,
	"FRF":          true,
	"FRG":          true,
	"FRM":          true,
	"FRO":          true,
	"FRT":          true,
	"FRY":          true,
	"FSG":          true,
	"FSK":          true,
	"FUD":          true,
	"FUG":          true,
	"FUN":          true,
	"G00":          true,
	"G01":          true,
	"G02":          true,
	"G03":          true,
	"G04":          true,
	"G05":          true,
	"GAL":          true,
	"GAM":          true,
	"GAR":          true,
	"GBL":          true,
	"GBM":          true,
	"GBN":          true,
	"GCB":          true,
	"GDF":          true,
	"GDM":          true,
	"GDR":          true,
	"GEF":          true,
	"GEM":          true,
	"GEN":          true,
	"GFC":          true,
	"GFF":          true,
	"GFM":          true,
	"GFO":          true,
	"GFR":          true,
	"GFS":          true,
	"GGL":          true,
	"GGY":          true,
	"GHF":          true,
	"GHM":          true,
	"GHO":          true,
	"GHU":          true,
	"GIA":          true,
	"GIG":          true,
	"GLB":          true,
	"GLC":          true,
	"GLM":          true,
	"GMF":          true,
	"GMM":          true,
	"GMN":          true,
	"GND":          true,
	"GNF":          true,
	"GNL":          true,
	"GNM":          true,
	"GNN":          true,
	"GOB":          true,
	"GOD":          true,
	"GOF":          true,
	"GOJ":          true,
	"GOL":          true,
	"GOM":          true,
	"GOO":          true,
	"GOR":          true,
	"GPF":          true,
	"GPM":          true,
	"GR1":          true,
	"GRA":          true,
	"GRD":          true,
	"GRF":          true,
	"GRG":          true,
	"GRI":          true,
	"GRL":          true,
	"GRM":          true,
	"GRN":          true,
	"GSF":          true,
	"GSM":          true,
	"GSN":          true,
	"GSP":          true,
	"GTD":          true,
	"GUA":          true,
	"GUL":          true,
	"GUM":          true,
	"GUS":          true,
	"GVK":          true,
	"GYA":          true,
	"GYO":          true,
	"HAF":          true,
	"HAG":          true,
	"HAM":          true,
	"HAR":          true,
	"HDL":          true,
	"HDV":          true,
	"HHF":          true,
	"HHM":          true,
	"HIF":          true,
	"HIM":          true,
	"HIP":          true,
	"HLF":          true,
	"HLG":          true,
	"HLM":          true,
	"HNF":          true,
	"HNM":          true,
	"HOM":          true,
	"HPF":          true,
	"HPM":          true,
	"HRP":          true,
	"HRS":          true,
	"HSF":          true,
	"HSM":          true,
	"HSN":          true,
	"HSS":          true,
	"HUF":          true,
	"HUM":          true,
	"HYC":          true,
	"HYD":          true,
	"I00":          true,
	"I01":          true,
	"I02":          true,
	"I03":          true,
	"I04":          true,
	"I05":          true,
	"I06":          true,
	"I07":          true,
	"I08":          true,
	"I09":          true,
	"I10":          true,
	"I11":          true,
	"I12":          true,
	"I13":          true,
	"I14":          true,
	"I15":          true,
	"I16":          true,
	"I17":          true,
	"I18":          true,
	"I19":          true,
	"I20":          true,
	"I21":          true,
	"I22":          true,
	"I23":          true,
	"I24":          true,
	"I25":          true,
	"I26":          true,
	"I27":          true,
	"I28":          true,
	"I29":          true,
	"I30":          true,
	"I31":          true,
	"I32":          true,
	"IBR":          true,
	"ICF":          true,
	"ICM":          true,
	"ICN":          true,
	"ICY":          true,
	"IEC":          true,
	"IFC":          true,
	"IHU":          true,
	"IKF":          true,
	"IKG":          true,
	"IKH":          true,
	"IKM":          true,
	"IKS":          true,
	"ILA":          true,
	"ILB":          true,
	"IMP":          true,
	"INN":          true,
	"INV":          true,
	"ISB":          true,
	"ISC":          true,
	"ISE":          true,
	"IVF":          true,
	"IVM":          true,
	"IWB":          true,
	"IWF":          true,
	"IWH":          true,
	"IWM":          true,
	"IZF":          true,
	"IZM":          true,
	"JKR":          true,
	"JUB":          true,
	"KAF":          true,
	"KAM":          true,
	"KAN":          true,
	"KAR":          true,
	"KBD":          true,
	"KDG":          true,
	"KED":          true,
	"KEF":          true,
	"KEM":          true,
	"KES":          true,
	"KGO":          true,
	"KHA":          true,
	"KNG":          true,
	"KOB":          true,
	"KOP":          true,
	"KOR":          true,
	"KRB":          true,
	"KRF":          true,
	"KRK":          true,
	"KRM":          true,
	"KRN":          true,
	"LAUNCH":       true,
	"LAUNCHM":      true,
	"LCR":          true,
	"LDR":          true,
	"LEE":          true,
	"LEP":          true,
	"LGA":          true,
	"LGR":          true,
	"LIF":          true,
	"LIM":          true,
	"LIZ":          true,
	"LMF":          true,
	"LMM":          true,
	"LSP":          true,
	"LSQ":          true,
	"LTH":          true,
	"LU2":          true,
	"LU3":          true,
	"LU4":          true,
	"LUC":          true,
	"LUG":          true,
	"LUJ":          true,
	"LVR":          true,
	"LYC":          true,
	"MAL":          true,
	"MAM":          true,
	"MAP":          true,
	"MAR":          true,
	"MBL":          true,
	"MBR":          true,
	"MBX":          true,
	"MCH":          true,
	"MCL":          true,
	"MCP":          true,
	"MCR":          true,
	"MCS":          true,
	"MDR":          true,
	"MEP":          true,
	"MER":          true,
	"MES":          true,
	"MFR":          true,
	"MGL":          true,
	"MHB":          true,
	"MHY":          true,
	"MIF":          true,
	"MIM":          true,
	"MIN":          true,
	"MINIPOM200":   true,
	"MKG":          true,
	"MKI":          true,
	"MMF":          true,
	"MMM":          true,
	"MMV":          true,
	"MMY":          true,
	"MNR":          true,
	"MNT":          true,
	"MOI":          true,
	"MOS":          true,
	"MPG":          true,
	"MPH":          true,
	"MPU":          true,
	"MRD":          true,
	"MRH":          true,
	"MRK":          true,
	"MRP":          true,
	"MRT":          true,
	"MSC":          true,
	"MSD":          true,
	"MSL":          true,
	"MSO":          true,
	"MTA":          true,
	"MTC":          true,
	"MTH":          true,
	"MTL":          true,
	"MTP":          true,
	"MTR":          true,
	"MUD":          true,
	"MUH":          true,
	"MUR":          true,
	"MWI":          true,
	"MWM":          true,
	"MWO":          true,
	"MWR":          true,
	"MYG":          true,
	"NBT":          true,
	"NET":          true,
	"NGF":          true,
	"NGM":          true,
	"NIN":          true,
	"NMG":          true,
	"NMH":          true,
	"NMP":          true,
	"NMW":          true,
	"NPT":          true,
	"NYD":          true,
	"NYM":          true,
	"OBJ_BLIMP":    true,
	"OBP_MELDRATH": true,
	"OGF":          true,
	"OGM":          true,
	"OKF":          true,
	"OKM":          true,
	"ONF":          true,
	"ONM":          true,
	"ONT":          true,
	"OPF":          true,
	"OPM":          true,
	"ORB":          true,
	"ORC":          true,
	"ORK":          true,
	"OTM":          true,
	"OWB":          true,
	"PAF":          true,
	"PBR":          true,
	"PEG":          true,
	"PG3":          true,
	"PGS":          true,
	"PHX":          true,
	"PIF":          true,
	"PIM":          true,
	"PIR":          true,
	"PMA":          true,
	"PPOINT":       true,
	"PRE":          true,
	"PRI":          true,
	"PRT":          true,
	"PSC":          true,
	"PUM":          true,
	"PUS":          true,
	"PYS":          true,
	"QCF":          true,
	"QCM":          true,
	"QZT":          true,
	"RAK":          true,
	"RAL":          true,
	"RAP":          true,
	"RAT":          true,
	"RAZ":          true,
	"RDG":          true,
	"REA":          true,
	"REF":          true,
	"REM":          true,
	"REN":          true,
	"RGM":          true,
	"RHI":          true,
	"RHP":          true,
	"RIF":          true,
	"RIM":          true,
	"RKP":          true,
	"RNB":          true,
	"ROB":          true,
	"ROE":          true,
	"ROM":          true,
	"RON":          true,
	"ROW":          true,
	"RPF":          true,
	"RPT":          true,
	"RTH":          true,
	"RTN":          true,
	"RZM":          true,
	"S01":          true,
	"SAR":          true,
	"SAT":          true,
	"SBU":          true,
	"SCA":          true,
	"SCC":          true,
	"SCE":          true,
	"SCH":          true,
	"SCO":          true,
	"SCR":          true,
	"SCU":          true,
	"SCW":          true,
	"SDC":          true,
	"SDE":          true,
	"SDF":          true,
	"SDM":          true,
	"SDR":          true,
	"SDV":          true,
	"SEA":          true,
	"SED":          true,
	"SEF":          true,
	"SEG":          true,
	"SEM":          true,
	"SER":          true,
	"SEY":          true,
	"SGO":          true,
	"SGR":          true,
	"SHA":          true,
	"SHD":          true,
	"SHE":          true,
	"SHF":          true,
	"SHI":          true,
	"SHIP":         true,
	"SHL":          true,
	"SHM":          true,
	"SHN":          true,
	"SHP":          true,
	"SHR":          true,
	"SHS":          true,
	"SIF":          true,
	"SIM":          true,
	"SIN":          true,
	"SIR":          true,
	"SKB":          true,
	"SKE":          true,
	"SKI":          true,
	"SKL":          true,
	"SKN":          true,
	"SKR":          true,
	"SKT":          true,
	"SKU":          true,
	"SLG":          true,
	"SMA":          true,
	"SMD":          true,
	"SNA":          true,
	"SND":          true,
	"SNE":          true,
	"SNK":          true,
	"SNN":          true,
	"SOK":          true,
	"SOL":          true,
	"SOS":          true,
	"SOW":          true,
	"SPB":          true,
	"SPC":          true,
	"SPD":          true,
	"SPE":          true,
	"SPH":          true,
	"SPI":          true,
	"SPL":          true,
	"SPQ":          true,
	"SPR":          true,
	"SPT":          true,
	"SPW":          true,
	"SPX":          true,
	"SRG":          true,
	"SRK":          true,
	"SRN":          true,
	"SRO":          true,
	"SRV":          true,
	"SRW":          true,
	"SSA":          true,
	"SSK":          true,
	"SSN":          true,
	"STA":          true,
	"STC":          true,
	"STF":          true,
	"STG":          true,
	"STM":          true,
	"STR":          true,
	"STU":          true,
	"SUC":          true,
	"SVO":          true,
	"SWC":          true,
	"SWI":          true,
	"SWO":          true,
	"SYN":          true,
	"SZK":          true,
	"T00":          true,
	"T01":          true,
	"T02":          true,
	"T03":          true,
	"T04":          true,
	"T05":          true,
	"T06":          true,
	"T07":          true,
	"T08":          true,
	"T09":          true,
	"T10":          true,
	"T11":          true,
	"T12":          true,
	"T13":          true,
	"TAC":          true,
	"TAR":          true,
	"TAZ":          true,
	"TBF":          true,
	"TBL":          true,
	"TBM":          true,
	"TBU":          true,
	"TEF":          true,
	"TEG":          true,
	"TEL":          true,
	"TEM":          true,
	"TEN":          true,
	"TGL":          true,
	"TGO":          true,
	"THO":          true,
	"TIG":          true,
	"TIN":          true,
	"TLN":          true,
	"TMB":          true,
	"TMR":          true,
	"TMT":          true,
	"TNF":          true,
	"TNM":          true,
	"TNT":          true,
	"TOT":          true,
	"TPB":          true,
	"TPF":          true,
	"TPL":          true,
	"TPM":          true,
	"TPN":          true,
	"TPO":          true,
	"TRA":          true,
	"TRE":          true,
	"TRF":          true,
	"TRG":          true,
	"TRI":          true,
	"TRK":          true,
	"TRM":          true,
	"TRN":          true,
	"TRQ":          true,
	"TRT":          true,
	"TRW":          true,
	"TSE":          true,
	"TSF":          true,
	"TSM":          true,
	"TTB":          true,
	"TUN":          true,
	"TVP":          true,
	"TWF":          true,
	"TZF":          true,
	"TZM":          true,
	"UDF":          true,
	"UDK":          true,
	"UNB":          true,
	"UNI":          true,
	"UNM":          true,
	"UVK":          true,
	"VAC":          true,
	"VAF":          true,
	"VAL":          true,
	"VAM":          true,
	"VAS":          true,
	"VAZ":          true,
	"VEG":          true,
	"VEK":          true,
	"VNM":          true,
	"VOL":          true,
	"VPF":          true,
	"VPM":          true,
	"VRM":          true,
	"VSF":          true,
	"VSG":          true,
	"VSK":          true,
	"VSM":          true,
	"VST":          true,
	"WAE":          true,
	"WAL":          true,
	"WAS":          true,
	"WBU":          true,
	"WEL":          true,
	"WER":          true,
	"WET":          true,
	"WIL":          true,
	"WLF":          true,
	"WLM":          true,
	"WMP":          true,
	"WOE":          true,
	"WOF":          true,
	"WOK":          true,
	"WOL":          true,
	"WOM":          true,
	"WOR":          true,
	"WRB":          true,
	"WRF":          true,
	"WRM":          true,
	"WRU":          true,
	"WRW":          true,
	"WUF":          true,
	"WUR":          true,
	"WWF":          true,
	"WYR":          true,
	"WYV":          true,
	"XAL":          true,
	"XEF":          true,
	"XEG":          true,
	"XEM":          true,
	"XHF":          true,
	"XHM":          true,
	"XIM":          true,
	"YAK":          true,
	"YET":          true,
	"ZBC":          true,
	"ZEB":          true,
	"ZEL":          true,
	"ZMF":          true,
	"ZMM":          true,
	"ZOF":          true,
	"ZOM":          true,
}
