package wld

import (
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

func (wld *Wld) WriteAscii(path string, isDir bool) error {
	var err error
	//var err error

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	rootPath := path + "/_root.wce"
	if wld.FileName == "objects.wld" {
		rootPath = path + "/_objects.wce"
	}
	if wld.FileName == "lights.wld" {
		rootPath = path + "/_lights.wce"
	}
	if !isDir {
		rootPath = path + "/" + wld.FileName
	}

	var w *os.File
	rootBuf, err := os.Create(rootPath)
	if err != nil {
		return err
	}
	defer rootBuf.Close()
	writeAsciiHeader(rootBuf)

	// now we can write

	lightsMap := map[string]bool{}
	zoneMaterials := map[string]bool{}

	if wld.GlobalAmbientLight != nil {
		err = wld.GlobalAmbientLight.Write(rootBuf)
		if err != nil {
			return fmt.Errorf("global ambient light: %w", err)
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.DMSpriteDef2s); i++ {
		dmSprite := wld.DMSpriteDef2s[i]

		baseTag := strings.ToLower(strings.TrimSuffix(strings.ToUpper(dmSprite.Tag), "_DMSPRITEDEF"))

		isZoneChunk := false
		// if baseTag is r### then foo
		if strings.HasPrefix(baseTag, "r") {
			regionChunk := 0
			chunkCount, err := fmt.Sscanf(baseTag, "r%d", &regionChunk)
			isZoneChunk = err == nil && chunkCount == 1
		}

		var dmBuf *os.File
		if !isZoneChunk {
			dmBuf, err = os.Create(path + "/" + baseTag + ".mod")
			if err != nil {
				return err
			}
			defer dmBuf.Close()
			writeAsciiHeader(dmBuf)

			w = dmBuf
		} else {
			w = rootBuf
			dmBuf = rootBuf
		}

		if dmSprite.MaterialPaletteTag != "" && !zoneMaterials[dmSprite.MaterialPaletteTag] {

			isMaterialPaletteFound := false
			for _, materialPal := range wld.MaterialPalettes {
				if materialPal.Tag != dmSprite.MaterialPaletteTag {
					continue
				}

				for _, materialTag := range materialPal.Materials {
					isMaterialDefFound := false
					for _, materialDef := range wld.MaterialDefs {
						if materialDef.Tag != materialTag {
							continue
						}

						if materialDef.SimpleSpriteTag != "" {
							isSimpleSpriteFound := false
							for _, simpleSprite := range wld.SimpleSpriteDefs {
								if simpleSprite.Tag != materialDef.SimpleSpriteTag {
									continue
								}
								isSimpleSpriteFound = true
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

		if dmSprite.PolyhedronTag != "" {
			isPolyFound := false
			for _, polyhedron := range wld.PolyhedronDefs {
				if polyhedron.Tag != dmSprite.PolyhedronTag {
					continue
				}
				err = polyhedron.Write(w)
				if err != nil {
					return fmt.Errorf("polyhedron %s: %w", polyhedron.Tag, err)
				}
				isPolyFound = true
				break
			}
			if !isPolyFound {
				fmt.Printf("polyhedron %s not found\n", dmSprite.PolyhedronTag)
				//	return fmt.Errorf("polyhedron %s not found", dmSprite.PolyhedronTag)
			}
		}

		err = dmSprite.Write(w)
		if err != nil {
			return fmt.Errorf("dm sprite def %s: %w", dmSprite.Tag, err)
		}
		if !isZoneChunk {
			fmt.Fprintf(rootBuf, "INCLUDE \"%s.MOD\"\n", strings.ToUpper(baseTag))
		}
		tracksWritten := map[string]bool{}
		var aniBuf *os.File
		for _, hierarchySprite := range wld.HierarchicalSpriteDefs {
			isFound := false

			if hierarchySprite.DMSpriteTag == dmSprite.Tag {
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

			err = hierarchySprite.Write(w)
			if err != nil {
				return fmt.Errorf("hierarchical sprite %s: %w", hierarchySprite.Tag, err)
			}

			for _, dag := range hierarchySprite.Dags {
				if dag.Track == "" {
					continue
				}
				if tracksWritten[dag.Track] {
					continue
				}

				isTrackFound := false
				for _, track := range wld.TrackInstances {
					if track.Tag != dag.Track {
						continue
					}

					isTrackDefFound := false

					var trackBuf *os.File
					for _, trackDef := range wld.TrackDefs {
						if trackDef.Tag != track.DefinitionTag {
							continue
						}
						isTrackDefFound = true

						if tracksWritten[trackDef.Tag] {
							break
						}

						trackBuf = dmBuf
						if isAnimationPrefix(trackDef.Tag) {
							if aniBuf == nil {
								aniBuf, err = os.Create(path + "/" + baseTag + ".ani")
								if err != nil {
									return err
								}
								defer aniBuf.Close()
								writeAsciiHeader(aniBuf)
							}

							trackBuf = aniBuf
						}
						err = trackDef.Write(trackBuf)
						if err != nil {
							return fmt.Errorf("track def %s: %w", trackDef.Tag, err)
						}

						tracksWritten[trackDef.Tag] = true
						break
					}
					if !isTrackDefFound {
						return fmt.Errorf("hierarchy %s track %s definition not found", hierarchySprite.Tag, track.DefinitionTag)
					}

					isTrackFound = true

					tracksWritten[dag.Track] = true

					trackBuf = dmBuf
					if isAnimationPrefix(dag.Track) {
						if aniBuf == nil {
							aniBuf, err = os.Create(path + "/" + baseTag + ".ani")
							if err != nil {
								return err
							}
							defer aniBuf.Close()
							writeAsciiHeader(aniBuf)
						}

						trackBuf = aniBuf
					}
					err = track.Write(trackBuf)
					if err != nil {
						return fmt.Errorf("track %s: %w", track.Tag, err)
					}
				}
				if !isTrackFound {
					return fmt.Errorf("hierarchy %s track %s not found", hierarchySprite.Tag, dag.Track)
				}
			}

			break
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.PolyhedronDefs); i++ {
		polyhedron := wld.PolyhedronDefs[i]
		err = polyhedron.Write(w)
		if err != nil {
			return fmt.Errorf("polyhedron %s: %w", polyhedron.Tag, err)
		}
	}

	w = rootBuf
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

	w = rootBuf
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

	w = rootBuf
	for i := 0; i < len(wld.WorldTrees); i++ {
		worldTree := wld.WorldTrees[i]

		err = worldTree.Write(w)
		if err != nil {
			return fmt.Errorf("world tree %s: %w", worldTree.Tag, err)
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.Regions); i++ {
		region := wld.Regions[i]

		err = region.Write(w)
		if err != nil {
			return fmt.Errorf("region %s: %w", region.Tag, err)
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.AmbientLights); i++ {
		ambientLight := wld.AmbientLights[i]

		err = ambientLight.Write(w)
		if err != nil {
			return fmt.Errorf("ambient light %s: %w", ambientLight.Tag, err)
		}
	}

	w = rootBuf
	for i := 0; i < len(wld.ActorInsts); i++ {
		actor := wld.ActorInsts[i]

		if actor.DMRGBTrackTag.Valid {
			isTrackFound := false
			isTrackDefFound := false
			for _, track := range wld.RGBTrackInsts {
				for _, trackDef := range wld.RGBTrackDefs {
					if trackDef.Tag != track.DefinitionTag {
						continue
					}
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
					continue
				}

				err = track.Write(w)
				if err != nil {
					return fmt.Errorf("track %s: %w", track.Tag, err)
				}

				isTrackFound = true
				break
			}

			if !isTrackFound {
				return fmt.Errorf("actor '%s' track %s not found", actor.Tag, actor.DMRGBTrackTag.String)
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

				for _, action := range actorDef.Actions {
					for _, lod := range action.LevelOfDetails {
						if lod.SpriteTag == "" {
							continue
						}
						isActorSpriteFound := false
						for _, sprite := range wld.Sprite3DDefs {
							if sprite.Tag != lod.SpriteTag {
								continue
							}

							err = sprite.Write(w)
							if err != nil {
								return fmt.Errorf("sprite %s: %w", sprite.Tag, err)
							}
							isActorSpriteFound = true
						}
						if !isActorSpriteFound {
							return fmt.Errorf("actor %s sprite %s not found", actorDef.Tag, lod.SpriteTag)
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

	w = rootBuf
	for i := 0; i < len(wld.Zones); i++ {
		zone := wld.Zones[i]

		err = zone.Write(w)
		if err != nil {
			return fmt.Errorf("zone %s: %w", zone.Tag, err)
		}
	}

	return nil
}

func writeAsciiHeader(w io.Writer) {
	fmt.Fprintf(w, "// wcemu %s\n", AsciiVersion)
	fmt.Fprintf(w, "// This file was created by quail v%s\n\n", common.Version)
}
