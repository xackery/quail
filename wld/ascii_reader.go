package wld

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/xackery/quail/helper"
)

type AsciiReadToken struct {
	basePath       string
	lineNumber     int
	reader         io.Reader
	wld            *Wld
	totalLineCount int // will be higher than lineNumber due to includes
}

// LoadAsciiFile returns a new AsciiReader that reads from r.
func LoadAsciiFile(path string, wld *Wld) (*AsciiReadToken, error) {
	r, err := caseInsensitiveOpen(path)
	if err != nil {
		return nil, err
	}
	a := &AsciiReadToken{
		lineNumber: 0,
		reader:     r,
		wld:        wld,
	}
	a.basePath = filepath.Dir(path)

	err = a.readDefinitions()
	if err != nil {
		return nil, fmt.Errorf("%s:%d: %w", path, a.lineNumber, err)
	}
	return a, nil
}

func (a *AsciiReadToken) Close() error {
	if c, ok := a.reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// caseInsensitiveOpen attempts to open a file in a case-insensitive manner.
func caseInsensitiveOpen(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if strings.EqualFold(entry.Name(), base) {
			return os.Open(filepath.Join(dir, entry.Name()))
		}
	}

	return nil, fmt.Errorf("file %s not found", path)
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
func (a *AsciiReadToken) Read(p []byte) (n int, err error) {
	n, err = a.reader.Read(p)
	if n > 0 {
		for _, b := range p {
			if b == '\n' {
				a.lineNumber++
			}
		}
	}
	return
}

type PropOpt struct {
	Name string
	Min  int
}

func (a *AsciiReadToken) ReadProperty(name string, minNumArgs int) ([]string, error) {
	if name == "" {
		return nil, fmt.Errorf("property name is empty")
	}
	property := ""
	args := []string{}
	for {
		buf := make([]byte, 1)
		_, err := a.Read(buf)
		if err != nil {
			return args, err
		}
		if buf[0] == '/' {
			_, err = a.Read(buf)
			if err != nil {
				return args, fmt.Errorf("read comment: %w", err)
			}
			if buf[0] != '/' {
				property += "/"
				continue
			}
			err = a.readComment()
			if err != nil {
				return args, fmt.Errorf("read comment: %w", err)
			}
			buf[0] = '\n'
		}

		if buf[0] == '\t' {
			buf[0] = ' '
		}
		if buf[0] == '\n' {
			//fmt.Println(a.lineNumber, property)
			if len(strings.TrimSpace(property)) == 0 {
				continue
			}
			args = strings.Split(strings.TrimSpace(property), " ")
			if len(args) == 0 {
				return args, fmt.Errorf("property %s has no arguments", name)
			}
			if !strings.EqualFold(args[0], name) {
				return args, fmt.Errorf("expected property '%s' got '%s'", name, args[0])
			}
			if minNumArgs > 0 && minNumArgs != len(args)-1 {
				return args, fmt.Errorf("property %s needs at least %d arguments, got %d", name, minNumArgs, len(args)-1)
			}

			if minNumArgs == -1 && len(args) == 1 {
				return args, fmt.Errorf("property %s needs at least 1 argument, got 0", name)
			}

			for i := 1; i < len(args); i++ {
				args[i] = strings.ReplaceAll(args[i], "\"", "")
			}
			return args, nil
		}
		if len(property) > 1 && property[len(property)-1] == ' ' && buf[0] == ' ' {
			continue
		}
		property += string(buf)
	}
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
		&DMSpriteDef2{},
		&HierarchicalSpriteDef{},
		&LightDef{},
		&MaterialDef{},
		&MaterialPalette{},
		&PolyhedronDefinition{},
		&Region{},
		&SimpleSpriteDef{},
		&Sprite3DDef{},
		&TrackDef{},
		&TrackInstance{},
		&WorldTree{},
		&Zone{},
		&RGBTrackDef{},
		&RGBTrack{},
	}

	definition := ""
	for {
		buf := make([]byte, 1)
		_, err := a.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read: %w", err)
		}
		if buf[0] == '\n' {
			if definition != "" {
				break
			}
			continue
		}

		if buf[0] == '/' {
			_, err = a.Read(buf)
			if err != nil {
				return fmt.Errorf("read comment: %w", err)
			}
			if buf[0] == '/' {
				err = a.readComment()
				if err != nil {
					return fmt.Errorf("read comment: %w", err)
				}
				continue
			}
		}

		// check if buf[0] is a letter
		if !unicode.IsLetter(rune(buf[0])) && !unicode.IsNumber(rune(buf[0])) {
			continue
		}

		definition += strings.ToUpper(string(buf))
		if strings.HasPrefix(definition, "INCLUDE") {
			err = a.readInclude()
			if err != nil {
				return fmt.Errorf("include: %w", err)
			}
			definition = ""
			continue
		}
		if strings.HasPrefix(definition, "NUMREGIONS") {
			err = a.readRegions()
			if err != nil {
				return fmt.Errorf("regions: %w", err)
			}

			definition = ""
			continue
		}
		if strings.HasSuffix(definition, "CPIWORLD") {
			// read to newline
			for {
				_, err = a.Read(buf)
				if err != nil {
					return fmt.Errorf("read cpiworld: %w", err)
				}
				if buf[0] == '\n' {
					break
				}
			}
			definition = ""
			continue
		}
		if strings.HasSuffix(definition, "ENDWORLD") {
			// read to newline
			for {
				_, err = a.Read(buf)
				if err != nil {
					return fmt.Errorf("read endworld: %w", err)
				}
				if buf[0] == '\n' {
					break
				}
			}
			definition = ""
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
			switch (definitions[i]).(type) {
			case *DMSpriteDef2:
				a.wld.DMSpriteDef2s = append(a.wld.DMSpriteDef2s, definitions[i].(*DMSpriteDef2))
				definitions[i] = &DMSpriteDef2{}
			case *HierarchicalSpriteDef:
				a.wld.HierarchicalSpriteDefs = append(a.wld.HierarchicalSpriteDefs, definitions[i].(*HierarchicalSpriteDef))
				definitions[i] = &HierarchicalSpriteDef{}
			case *MaterialDef:
				a.wld.MaterialDefs = append(a.wld.MaterialDefs, definitions[i].(*MaterialDef))
				definitions[i] = &MaterialDef{}
			case *MaterialPalette:
				a.wld.MaterialPalettes = append(a.wld.MaterialPalettes, definitions[i].(*MaterialPalette))
				definitions[i] = &MaterialPalette{}
			case *PolyhedronDefinition:
				a.wld.PolyhedronDefs = append(a.wld.PolyhedronDefs, definitions[i].(*PolyhedronDefinition))
				definitions[i] = &PolyhedronDefinition{}
			case *SimpleSpriteDef:
				a.wld.SimpleSpriteDefs = append(a.wld.SimpleSpriteDefs, definitions[i].(*SimpleSpriteDef))
				definitions[i] = &SimpleSpriteDef{}
			case *TrackDef:
				a.wld.TrackDefs = append(a.wld.TrackDefs, definitions[i].(*TrackDef))
				definitions[i] = &TrackDef{}
			case *TrackInstance:
				a.wld.TrackInstances = append(a.wld.TrackInstances, definitions[i].(*TrackInstance))
				definitions[i] = &TrackInstance{}
			case *LightDef:
				a.wld.LightDefs = append(a.wld.LightDefs, definitions[i].(*LightDef))
				definitions[i] = &LightDef{}
			case *Sprite3DDef:
				a.wld.Sprite3DDefs = append(a.wld.Sprite3DDefs, definitions[i].(*Sprite3DDef))
				definitions[i] = &Sprite3DDef{}
			case *WorldTree:
				a.wld.WorldTrees = append(a.wld.WorldTrees, definitions[i].(*WorldTree))
				definitions[i] = &WorldTree{}
			case *Region:
				a.wld.Regions = append(a.wld.Regions, definitions[i].(*Region))
				definitions[i] = &Region{}
			case *AmbientLight:
				a.wld.AmbientLights = append(a.wld.AmbientLights, definitions[i].(*AmbientLight))
				definitions[i] = &AmbientLight{}
			case *ActorDef:
				a.wld.ActorDefs = append(a.wld.ActorDefs, definitions[i].(*ActorDef))
				definitions[i] = &ActorDef{}
			case *ActorInst:
				a.wld.ActorInsts = append(a.wld.ActorInsts, definitions[i].(*ActorInst))
				definitions[i] = &ActorInst{}
			case *Zone:
				a.wld.Zones = append(a.wld.Zones, definitions[i].(*Zone))
				definitions[i] = &Zone{}

			}

			break
		}
	}

	if definition != "" {
		return fmt.Errorf("unknown definition: %s", definition)
	}
	return nil
}

func (a *AsciiReadToken) readInclude() error {
	filename := ""
	for {
		buf := make([]byte, 1)
		_, err := a.Read(buf)
		if err != nil {
			return err
		}
		if buf[0] == ' ' {
			continue
		}
		if buf[0] == '\n' {
			if filename == "" {
				return fmt.Errorf("include: missing filename")
			}
			return fmt.Errorf("include: missing end quote")
		}
		if filename != "" && buf[0] == '"' {
			break
		}
		if buf[0] == '"' {
			continue
		}
		filename += string(buf)
	}
	path := a.basePath + "/" + filename
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

func (a *AsciiReadToken) readComment() error {
	for {
		buf := make([]byte, 1)
		_, err := a.Read(buf)
		if err != nil {
			return fmt.Errorf("read comment: %w", err)
		}
		if buf[0] == '\n' {
			return nil
		}
	}
}

func (a *AsciiReadToken) readRegions() error {
	var err error
	line := ""
	for {
		buf := make([]byte, 1)
		_, err = a.Read(buf)
		if err != nil {
			return fmt.Errorf("read regions: %w", err)
		}
		if buf[0] == '\n' {
			break
		}
		if buf[0] == ' ' {
			continue
		}
		line += string(buf)
	}
	numRegions, err := helper.ParseInt(line)
	if err != nil {
		return fmt.Errorf("parse numregions: %w", err)
	}

	a.wld.Regions = make([]*Region, numRegions)
	for i := 0; i < numRegions; i++ {
		r := &Region{}
		_, err = a.ReadProperty("REGION", 1)
		if err != nil {
			return fmt.Errorf("REGION: %w", err)
		}
		err = r.Read(a)
		if err != nil {
			return fmt.Errorf("read region: %w", err)
		}

		a.wld.Regions[i] = r
	}

	return nil
}
