package wce

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

var regexLine = regexp.MustCompile(`"([^"]*)"|(\S+)`)

type AsciiReadToken struct {
	folder         string
	basePath       string
	lineNumber     int
	buf            *bytes.Buffer
	wce            *Wce
	totalLineCount int // will be higher than lineNumber due to includes
}

func AsciiReadTokenNew(buf *bytes.Buffer, wce *Wce) *AsciiReadToken {
	return &AsciiReadToken{
		lineNumber: 0,
		buf:        buf,
		wce:        wce,
	}
}

// LoadAsciiFile returns a new AsciiReader that reads from r.
func LoadAsciiFile(path string, wce *Wce) (*AsciiReadToken, error) {
	buf, err := caseInsensitiveOpen(path, wce)
	if err != nil {
		return nil, err
	}
	a := &AsciiReadToken{
		lineNumber: 0,
		buf:        buf,
		wce:        wce,
	}
	a.basePath = filepath.Dir(strings.ToLower(path))
	a.folder = filepath.Base(a.basePath)

	err = a.readDefinitions()
	if err != nil {
		return nil, fmt.Errorf("%s:%d: %w", path, a.lineNumber, err)
	}
	return a, nil
}

func (a *AsciiReadToken) Close() error {
	return nil
}

// caseInsensitiveOpen attempts to open a file in a case-insensitive manner.
func caseInsensitiveOpen(path string, wce *Wce) (*bytes.Buffer, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	entries, err := wce.FileSystem.ReadDir(strings.ToLower(dir))
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !strings.EqualFold(entry.Name(), base) {
			continue
		}
		data, err := wce.FileSystem.ReadFile(filepath.Join(strings.ToLower(dir), entry.Name()))
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(data), nil
		//			return wce.FileSystem.Open(filepath.Join(strings.ToLower(dir), entry.Name()))
	}

	return nil, fmt.Errorf("file %s not found", path)
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
func (a *AsciiReadToken) ReadLine() (string, error) {
	line := ""
	p := make([]byte, 1)
	for {
		_, err := a.buf.Read(p)
		if err != nil {
			if err == io.EOF {
				a.lineNumber++
				return line, err
			}

			return "", err
		}
		if p[0] != '\n' {
			line += string(p)
			continue
		}
		a.lineNumber++
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			line = ""
			continue
		}
		if strings.TrimSpace(line) == "" {
			line = ""
			continue
		}
		return line, nil
	}
}

func (a *AsciiReadToken) ReadSegmentedLine() ([]string, error) {
	line, err := a.ReadLine()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		if len(line) == 0 {
			return nil, err
		}
	}
	matches := regexLine.FindAllStringSubmatch(line, -1)
	args := []string{}
	for _, match := range matches {
		if strings.Contains(match[2], "//") {
			match[2] = match[2][:strings.Index(match[2], "//")]
			if match[2] == "" {
				break
			}
		}
		if match[1] != "" {
			args = append(args, match[1])
		} else {
			args = append(args, match[2])
		}
	}
	return args, nil
}

type PropOpt struct {
	Name string
	Min  int
}

func (a *AsciiReadToken) ReadProperty(name string, minNumArgs int) ([]string, error) {
	if name == "" {
		return nil, fmt.Errorf("property name is empty")
	}
	args, err := a.ReadSegmentedLine()
	if err != nil {
		return args, fmt.Errorf("read property %s: %w", name, err)
	}
	if len(args) == 0 {
		return args, fmt.Errorf("property %s has no arguments", name)
	}

	value := args[len(args)-1]
	if !strings.HasSuffix(name, "?") && value == "NULL" {
		return args, fmt.Errorf("invalid property NULL for %s", name)
	}
	if len(args) == 0 {
		return args, fmt.Errorf("property %s has no arguments", name)
	}
	if !strings.EqualFold(args[0], name) {
		return args, fmt.Errorf("expected property '%s' got '%s'", name, args[0])
	}
	if minNumArgs > 0 && minNumArgs != len(args)-1 {
		return args, fmt.Errorf("property %s needs %d arguments, got %d: %s", name, minNumArgs, len(args)-1, args)
	}

	if minNumArgs == -1 && len(args) == 1 {
		return args, fmt.Errorf("property %s needs at least 1 argument, got 0", name)
	}

	for i := 1; i < len(args); i++ {
		args[i] = strings.ReplaceAll(args[i], "\"", "")
	}
	return args, nil
}

func (a *AsciiReadToken) TotalLineCountRead() int {
	return a.totalLineCount + a.lineNumber
}

func (a *AsciiReadToken) readDefinitions() error {
	type definitionReader interface {
		Definition() string
		Read(r *AsciiReadToken) error
	}
	definitions := []definitionReader{
		&ActorDef{},
		&ActorInst{},
		&EqgAniDef{},
		&AmbientLight{},
		&BlitSpriteDef{},
		&DMSpriteDef{},
		&DMSpriteDef2{},
		&DMTrackDef2{},
		&EQMaterialDef{},
		&GlobalAmbientLightDef{},
		&HierarchicalSpriteDef{},
		&EqgLayDef{},
		&EqgParticlePointDef{},
		&EqgParticleRenderDef{},
		&LightDef{},
		&MaterialDef{},
		&MaterialPalette{},
		&EqgMdsDef{},
		&EqgModDef{},
		&EqgTerDef{},
		&EqgZonDef{},
		&EqgLodDef{},
		&EqgLayDef{},
		&ParticleCloudDef{},
		&PointLight{},
		&PolyhedronDefinition{},
		&Region{},
		&RGBTrackDef{},
		&SimpleSpriteDef{},
		&Sprite2DDef{},
		&Sprite3DDef{},
		&TrackDef{},
		&TrackInstance{},
		&WorldDef{folders: []string{"world"}},
		&WorldTree{},
		&Zone{},
	}

	definition := ""
	for {
		args, err := a.ReadSegmentedLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read: %w", err)
		}
		if len(args) == 0 {
			continue
		}

		definition = strings.ToUpper(string(args[0]))
		if strings.HasPrefix(definition, "INCLUDE") {
			err = a.readInclude(args)
			if err != nil {
				return fmt.Errorf("-> %w", err)
			}
			definition = ""
			continue
		}
		if strings.HasPrefix(definition, "//") {
			continue
		}

		for i := 0; i < len(definitions); i++ {
			defName := definitions[i].Definition()
			defRead := definitions[i].Read
			if !strings.HasPrefix(definition, defName) {
				continue
			}
			if defName != definition {
				continue
			}

			definition = ""
			err = defRead(a)
			if err != nil {
				return fmt.Errorf("%s: %w", defName, err)
			}
			switch frag := (definitions[i]).(type) {
			case *GlobalAmbientLightDef:
				if a.wce.GlobalAmbientLightDef != nil {
					return fmt.Errorf("duplicate global ambient light definition")
				}
				a.wce.GlobalAmbientLightDef = frag
				definitions[i] = &GlobalAmbientLightDef{}
			case *BlitSpriteDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				// Check if the Tag already exists in wce.BlitSpriteDefs
				exists := false
				for _, existing := range a.wce.BlitSpriteDefs {
					if existing.Tag == frag.Tag {
						exists = true
						break
					}
				}
				// Append the current BlitSpriteDef only if no match is found
				if !exists {
					a.wce.BlitSpriteDefs = append(a.wce.BlitSpriteDefs, frag)
				}
				definitions[i] = &BlitSpriteDef{}
			case *DMSpriteDef2:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				if existingIndex, exists := a.wce.tagIndexes[frag.Tag]; exists && existingIndex == frag.TagIndex {
					definitions[i] = &DMSpriteDef2{}
					break
				}
				a.wce.DMSpriteDef2s = append(a.wce.DMSpriteDef2s, frag)
				a.wce.tagIndexes[frag.Tag] = frag.TagIndex
				definitions[i] = &DMSpriteDef2{}
			case *HierarchicalSpriteDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				exists := false
				for _, existing := range a.wce.HierarchicalSpriteDefs {
					if existing.Tag == frag.Tag {
						exists = true
						break
					}
				}
				if !exists {
					a.wce.HierarchicalSpriteDefs = append(a.wce.HierarchicalSpriteDefs, frag)
				}
				definitions[i] = &HierarchicalSpriteDef{}
			case *MaterialDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				// Check if the Tag and TagIndex combination already exists
				if existingIndex, exists := a.wce.tagIndexes[frag.Tag]; exists && existingIndex == frag.TagIndex {
					definitions[i] = &MaterialDef{}
					break
				}
				// Add to MaterialDefs and tagIndexes
				a.wce.MaterialDefs = append(a.wce.MaterialDefs, frag)
				a.wce.tagIndexes[frag.Tag] = frag.TagIndex
				definitions[i] = &MaterialDef{}
			case *MaterialPalette:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.MaterialPalettes = append(a.wce.MaterialPalettes, frag)
				definitions[i] = &MaterialPalette{}
			case *PolyhedronDefinition:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.PolyhedronDefs = append(a.wce.PolyhedronDefs, frag)
				definitions[i] = &PolyhedronDefinition{}
			case *SimpleSpriteDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.SimpleSpriteDefs = append(a.wce.SimpleSpriteDefs, frag)
				definitions[i] = &SimpleSpriteDef{}
			case *TrackDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.TrackDefs = append(a.wce.TrackDefs, frag)
				definitions[i] = &TrackDef{}
			case *TrackInstance:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				if existingIndex, exists := a.wce.tagIndexes[frag.Tag]; exists && existingIndex == frag.TagIndex {
					definitions[i] = &TrackInstance{}
					break
				}
				a.wce.TrackInstances = append(a.wce.TrackInstances, frag)
				a.wce.tagIndexes[frag.Tag] = frag.TagIndex
				definitions[i] = &TrackInstance{}
			case *DMTrackDef2:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				exists := false
				for _, existing := range a.wce.DMTrackDef2s {
					if existing.Tag == frag.Tag {
						exists = true
						break
					}
				}
				if !exists {
					a.wce.DMTrackDef2s = append(a.wce.DMTrackDef2s, frag)
				}
				definitions[i] = &DMTrackDef2{}
			case *LightDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.LightDefs = append(a.wce.LightDefs, frag)
				definitions[i] = &LightDef{}
			case *Sprite3DDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.Sprite3DDefs = append(a.wce.Sprite3DDefs, frag)
				definitions[i] = &Sprite3DDef{}
			case *WorldTree:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.WorldTrees = append(a.wce.WorldTrees, frag)
				definitions[i] = &WorldTree{}
			case *Region:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.Regions = append(a.wce.Regions, frag)
				definitions[i] = &Region{}
			case *AmbientLight:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.AmbientLights = append(a.wce.AmbientLights, frag)
				definitions[i] = &AmbientLight{}
			case *ActorDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.ActorDefs = append(a.wce.ActorDefs, frag)
				definitions[i] = &ActorDef{}
			case *ActorInst:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.ActorInsts = append(a.wce.ActorInsts, frag)
				definitions[i] = &ActorInst{}
			case *Zone:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.Zones = append(a.wce.Zones, frag)
				definitions[i] = &Zone{}
			case *RGBTrackDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.RGBTrackDefs = append(a.wce.RGBTrackDefs, frag)
				definitions[i] = &RGBTrackDef{}
			case *ParticleCloudDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.ParticleCloudDefs = append(a.wce.ParticleCloudDefs, frag)
				definitions[i] = &ParticleCloudDef{}
			case *Sprite2DDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.Sprite2DDefs = append(a.wce.Sprite2DDefs, frag)
				definitions[i] = &Sprite2DDef{}
			case *PointLight:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.PointLights = append(a.wce.PointLights, frag)
				definitions[i] = &PointLight{}
			case *DMSpriteDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				if existingIndex, exists := a.wce.tagIndexes[frag.Tag]; exists && existingIndex == frag.TagIndex {
					definitions[i] = &DMSpriteDef{}
					break
				}
				a.wce.DMSpriteDefs = append(a.wce.DMSpriteDefs, frag)
				a.wce.tagIndexes[frag.Tag] = frag.TagIndex
				definitions[i] = &DMSpriteDef{}
			case *WorldDef:

				a.wce.WorldDef = frag
				definitions[i] = &WorldDef{folders: []string{"world"}}
			case *EqgMdsDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.MdsDefs = append(a.wce.MdsDefs, frag)
				definitions[i] = &EqgMdsDef{}
			case *EqgModDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.ModDefs = append(a.wce.ModDefs, frag)
				definitions[i] = &EqgModDef{}
			case *EqgTerDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.TerDefs = append(a.wce.TerDefs, frag)
				definitions[i] = &EqgTerDef{}
			case *EqgAniDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.AniDefs = append(a.wce.AniDefs, frag)
				definitions[i] = &EqgAniDef{}
			case *EqgLodDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.LodDefs = append(a.wce.LodDefs, frag)
				definitions[i] = &EqgLodDef{}
			case *EqgLayDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.LayDefs = append(a.wce.LayDefs, frag)
				definitions[i] = &EqgLayDef{}
			case *EqgParticlePointDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.PtsDefs = append(a.wce.PtsDefs, frag)
				definitions[i] = &EqgParticlePointDef{}
			case *EqgParticleRenderDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.PrtDefs = append(a.wce.PrtDefs, frag)
				definitions[i] = &EqgParticleRenderDef{}
			case *EqgZonDef:
				if len(args) == 1 {
					return fmt.Errorf("definition %s has no arguments", defName)
				}
				frag.Tag = args[1]
				a.wce.ZonDefs = append(a.wce.ZonDefs, frag)
				definitions[i] = &EqgZonDef{}
			default:
				return fmt.Errorf("unknown definition type for rebuild: %T", definitions[i])
			}

			break
		}
	}

	if definition != "" {
		return fmt.Errorf("unknown definition: %s", definition)
	}
	return nil
}

func (a *AsciiReadToken) readInclude(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("INCLUDE needs 1 argument")
	}

	path := a.basePath + "/" + args[1]
	if strings.HasSuffix(args[1], "/_ROOT.WCE") {
		a.wce.lastReadFolder = strings.TrimSuffix(args[1], "/_ROOT.WCE")
	}
	ir, err := LoadAsciiFile(path, a.wce)
	if err != nil {
		return err
	}
	err = ir.readDefinitions()
	if err != nil {
		return fmt.Errorf("read definitions: %w", err)
	}

	a.totalLineCount += ir.TotalLineCountRead()

	err = ir.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}
