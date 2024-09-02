package wld

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexLine = regexp.MustCompile(`"([^"]*)"|(\S+)`)

type AsciiReadToken struct {
	basePath       string
	lineNumber     int
	buf            *bytes.Buffer
	wld            *Wld
	totalLineCount int // will be higher than lineNumber due to includes
}

// LoadAsciiFile returns a new AsciiReader that reads from r.
func LoadAsciiFile(path string, wld *Wld) (*AsciiReadToken, error) {
	buf, err := caseInsensitiveOpen(path)
	if err != nil {
		return nil, err
	}
	a := &AsciiReadToken{
		lineNumber: 0,
		buf:        buf,
		wld:        wld,
	}
	a.basePath = filepath.Dir(strings.ToLower(path))

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
func caseInsensitiveOpen(path string) (*bytes.Buffer, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	entries, err := os.ReadDir(strings.ToLower(dir))
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !strings.EqualFold(entry.Name(), base) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(strings.ToLower(dir), entry.Name()))
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(data), nil
		//			return os.Open(filepath.Join(strings.ToLower(dir), entry.Name()))
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
		if match[2] == "//" {
			break
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
		return args, fmt.Errorf("property %s needs %d arguments, got %d", name, minNumArgs, len(args)-1)
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
		&AmbientLight{},
		&DMSpriteDef{},
		&DMSpriteDef2{},
		&GlobalAmbientLightDef{},
		&HierarchicalSpriteDef{},
		&LightDef{},
		&MaterialDef{},
		&MaterialPalette{},
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
		&WorldTree{},
		&Zone{},
		&WorldDef{},
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
				return fmt.Errorf("include: %w", err)
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
				if a.wld.GlobalAmbientLightDef != nil {
					return fmt.Errorf("duplicate global ambient light definition")
				}
				a.wld.GlobalAmbientLightDef = frag
				definitions[i] = &GlobalAmbientLightDef{}
			case *DMSpriteDef2:
				a.wld.DMSpriteDef2s = append(a.wld.DMSpriteDef2s, frag)
				definitions[i] = &DMSpriteDef2{}
			case *HierarchicalSpriteDef:
				a.wld.HierarchicalSpriteDefs = append(a.wld.HierarchicalSpriteDefs, frag)
				definitions[i] = &HierarchicalSpriteDef{}
			case *MaterialDef:
				a.wld.MaterialDefs = append(a.wld.MaterialDefs, frag)
				definitions[i] = &MaterialDef{}
			case *MaterialPalette:
				a.wld.MaterialPalettes = append(a.wld.MaterialPalettes, frag)
				definitions[i] = &MaterialPalette{}
			case *PolyhedronDefinition:
				a.wld.PolyhedronDefs = append(a.wld.PolyhedronDefs, frag)
				definitions[i] = &PolyhedronDefinition{}
			case *SimpleSpriteDef:
				a.wld.SimpleSpriteDefs = append(a.wld.SimpleSpriteDefs, frag)
				definitions[i] = &SimpleSpriteDef{}
			case *TrackDef:
				a.wld.TrackDefs = append(a.wld.TrackDefs, frag)
				definitions[i] = &TrackDef{}
			case *TrackInstance:
				a.wld.TrackInstances = append(a.wld.TrackInstances, frag)
				definitions[i] = &TrackInstance{}
			case *LightDef:
				a.wld.LightDefs = append(a.wld.LightDefs, frag)
				definitions[i] = &LightDef{}
			case *Sprite3DDef:
				a.wld.Sprite3DDefs = append(a.wld.Sprite3DDefs, frag)
				definitions[i] = &Sprite3DDef{}
			case *WorldTree:
				a.wld.WorldTrees = append(a.wld.WorldTrees, frag)
				definitions[i] = &WorldTree{}
			case *Region:
				a.wld.Regions = append(a.wld.Regions, frag)
				definitions[i] = &Region{}
			case *AmbientLight:
				a.wld.AmbientLights = append(a.wld.AmbientLights, frag)
				definitions[i] = &AmbientLight{}
			case *ActorDef:
				a.wld.ActorDefs = append(a.wld.ActorDefs, frag)
				definitions[i] = &ActorDef{}
			case *ActorInst:
				a.wld.ActorInsts = append(a.wld.ActorInsts, frag)
				definitions[i] = &ActorInst{}
			case *Zone:
				a.wld.Zones = append(a.wld.Zones, frag)
				definitions[i] = &Zone{}
			case *RGBTrackDef:
				a.wld.RGBTrackDefs = append(a.wld.RGBTrackDefs, frag)
				definitions[i] = &RGBTrackDef{}
			case *ParticleCloudDef:
				a.wld.ParticleCloudDefs = append(a.wld.ParticleCloudDefs, frag)
				definitions[i] = &ParticleCloudDef{}
			case *Sprite2DDef:
				a.wld.Sprite2DDefs = append(a.wld.Sprite2DDefs, frag)
				definitions[i] = &Sprite2DDef{}
			case *PointLight:
				a.wld.PointLights = append(a.wld.PointLights, frag)
				definitions[i] = &PointLight{}
			case *DMSpriteDef:
				a.wld.DMSpriteDefs = append(a.wld.DMSpriteDefs, frag)
				definitions[i] = &DMSpriteDef{}
			case *WorldDef:
				if a.wld.WorldDef != nil {
					return fmt.Errorf("duplicate world definition")
				}
				a.wld.WorldDef = frag
				definitions[i] = &WorldDef{}
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
		a.wld.lastReadModelTag = strings.TrimSuffix(args[1], "/_ROOT.WCE")
	}
	ir, err := LoadAsciiFile(path, a.wld)
	if err != nil {
		return fmt.Errorf("new ascii reader: %w", err)
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
