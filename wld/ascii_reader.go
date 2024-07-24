package wld

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type AsciiReader interface {
	Read(p []byte) (n int, err error)
	ReadProperty(definition string) (string, error)
	ReadString() (string, error)
	ReadInt() (int, error)
}

type AsciiReadToken struct {
	basePath     string
	lineNumber   int
	lastProperty string
	reader       io.Reader
	wld          *Wld
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

// ReadProperty reads the property of the definition.
func (a *AsciiReadToken) ReadProperty(definition string) (string, error) {
	properties := map[string][]string{
		"DMSPRITEDEF2":         {"TAG", "FLAGS", "CENTEROFFSET", "NUMVERTICES", "NUMUVS", "XYZ", "UV", "DMFACE2", "ENDDMFACE2", "TRIANGLE", "NUMVERTEXNORMALS", "SKINASSIGNMENTGROUPS", "MATERIALPALETTE", "NUMFACE2S", "NUMMESHOPS", "MESHOP_FA", "MESHOP_VA", "MESHOP_SW", "MESHOP_EL", "FACEMATERIALGROUPS", "VERTEXMATERIALGROUPS", "BOUNDINGRADIUS", "FPSCALE", "ENDDMSPRITEDEF2"},
		"INCLUDE":              {"FILENAME"},
		"MATERIALPALETTE":      {"TAG", "NUMMATERIALS", "ENDMATERIALPALETTE", "MATERIAL"},
		"POLYHEDRONDEFINITION": {"TAG", "BOUNDINGRADIUS", "SCALEFACTOR", "NUMVERTICES", "XYZ", "NUMFACES", "VERTEXLIST", "FACE", "ENDFACE", "ENDPOLYHEDRONDEFINITION"},
		"TRACKINSTANCE":        {"TAG", "DEFINITION", "INTERPOLATE", "SLEEP", "ENDTRACKINSTANCE"},
		"TRACKDEFINITION":      {"TAG", "NUMFRAMES", "FRAMETRANSFORM", "ENDTRACKDEFINITION"},
		"SIMPLESPRITEDEF":      {"SIMPLESPRITETAG", "NUMFRAMES", "BMINFO", "ENDSIMPLESPRITEDEF"},
		"MATERIALDEFINITION":   {"TAG", "RENDERMETHOD", "RGBPEN", "BRIGHTNESS", "SCALEDAMBIENT", "SIMPLESPRITEINST", "ENDSIMPLESPRITEINST", "ENDMATERIALDEFINITION"},
	}
	endMarkers := map[string]string{
		"DMSPRITEDEF2":         "ENDDMSPRITEDEF2",
		"INCLUDE":              "",
		"MATERIALPALETTE":      "ENDMATERIALPALETTE",
		"POLYHEDRONDEFINITION": "ENDPOLYHEDRONDEFINITION",
		"TRACKINSTANCE":        "ENDTRACKINSTANCE",
		"TRACKDEFINITION":      "ENDTRACKDEFINITION",
		"SIMPLESPRITEDEF":      "ENDSIMPLESPRITEDEF",
		"MATERIALDEFINITION":   "ENDMATERIALDEFINITION",
	}

	if definition == "" {
		return "", fmt.Errorf("definition: empty")
	}
	property := ""
	if a.lastProperty != "" {
		property = a.lastProperty
		a.lastProperty = ""
		endMark := endMarkers[definition]
		if strings.HasPrefix(strings.TrimSpace(property), endMark) {
			return strings.TrimSpace(property), nil
		}
	}
	for {
		buf := make([]byte, 1)
		_, err := a.Read(buf)
		if err != nil {
			return property, fmt.Errorf("read: %w", err)
		}
		if buf[0] == '\n' {
			continue
		}
		if buf[0] == '/' {
			_, err = a.Read(buf)
			if err != nil {
				return property, fmt.Errorf("read comment: %w", err)
			}
			if buf[0] != '/' {
				return property, fmt.Errorf("comment: expected second slash")
			}
			err = a.readComment()
			if err != nil {
				return property, fmt.Errorf("read comment: %w", err)
			}
			continue
		}
		property += string(buf)
		nextProperty := ""
		propertyUpper := strings.ToUpper(property)
		isComplete := false
		for _, propName := range properties[definition] {
			if strings.HasSuffix(propertyUpper, "\t"+propName) {
				isComplete = true
				nextProperty = propName
				break
			}
			if strings.HasPrefix(propName, "END") && strings.HasSuffix(propertyUpper, propName) {
				isComplete = true
				nextProperty = propName
				break
			}
		}
		if !isComplete {
			continue
		}

		property = strings.TrimSuffix(property, nextProperty)
		if property == "" {
			property = nextProperty
			continue
		}
		a.lastProperty = nextProperty + " "
		out := property
		out = strings.ReplaceAll(out, "\t", "")
		out = strings.ReplaceAll(out, "\r", "")
		out = strings.TrimSpace(out)
		fmt.Println("Property:", out)
		return out, nil
	}
}

func (a *AsciiReadToken) readDefinitions() error {
	type definitionReader interface {
		Definition() string
		Read(r *AsciiReadToken) error
	}
	definitions := []definitionReader{
		&DMSpriteDef2{},
		//&Include{},
		&MaterialPalette{},
		&PolyhedronDefinition{},
		&TrackInstance{},
		&TrackDef{},
		&SimpleSpriteDef{},
		&MaterialDef{},
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
			if buf[0] != '/' {
				return fmt.Errorf("comment: expected second slash")
			}
			err = a.readComment()
			if err != nil {
				return fmt.Errorf("read comment: %w", err)
			}
			continue
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

		for i := 0; i < len(definitions); i++ {
			defName := definitions[i].Definition()
			defRead := definitions[i].Read
			if strings.HasPrefix(definition, defName) {
				if defName != definition {
					continue
				}
				definition = ""
				err = defRead(a)
				if err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("%s: %w", defName, err)
				}
				break
			}
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
			return fmt.Errorf("read: %w", err)
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
	fmt.Println("Include:", path)
	ir, err := LoadAsciiFile(path, a.wld)
	if err != nil {
		return fmt.Errorf("new ascii reader: %w", err)
	}
	err = ir.readDefinitions()
	if err != nil {
		return fmt.Errorf("%s:%d: %w", path, a.lineNumber, err)
	}

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
			break
		}
	}
	return nil
}
