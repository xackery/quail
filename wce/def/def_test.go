package def

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/wce"
	"gopkg.in/yaml.v3"
)

type defReadWriter interface {
	Definition() string
	Write(token *wce.AsciiWriteToken) error
	Read(token *wce.AsciiReadToken) error
}

func TestWceGenTypescript(t *testing.T) {

	defs := []defReadWriter{
		&wce.ActorDef{},
		&wce.ActorInst{},
		&wce.BlitSpriteDef{},
		&wce.DMSpriteDef{},
		&wce.DMSpriteDef2{},
		&wce.GlobalAmbientLightDef{},
		&wce.MaterialDef{},
		&wce.MaterialPalette{},
		&wce.SimpleSpriteDef{},
		&wce.WorldDef{},
		&wce.LightDef{},
		&wce.PointLight{},
		&wce.Sprite3DDef{},
		&wce.PolyhedronDefinition{},
		&wce.TrackInstance{},
		&wce.TrackDef{},
		&wce.HierarchicalSpriteDef{},
		&wce.WorldTree{},
		&wce.Region{},
		&wce.AmbientLight{},
		&wce.Zone{},
		&wce.RGBTrackDef{},
		&wce.ParticleCloudDef{},
		&wce.Sprite2DDef{},
		&wce.DMTrackDef2{},
	}

	dirTest := helper.DirTest()

	for _, def := range defs {
		defName := strings.ToLower(def.Definition())

		r, err := os.Open(fmt.Sprintf("%s/../wce/def/%s.yaml", dirTest, strings.ToLower(defName)))
		if err != nil {
			t.Fatalf("open %s: %v", defName, err)
		}
		defer r.Close()

		yamlDef := &Definition{}
		err = yaml.NewDecoder(r).Decode(yamlDef)
		if err != nil {
			t.Fatalf("decode %s: %v", defName, err)
		}

		w, err := os.Create(fmt.Sprintf("%s/%s.ts", dirTest, defName))
		if err != nil {
			t.Fatalf("create failed: %s", err)
		}
		defer w.Close()

		w.WriteString(`// Generated
import * as data from "./data";

export const ` + defName + `: data.DefinitionInfo = {
`)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "    ")
		err = enc.Encode(yamlDef)
		if err != nil {
			t.Fatalf("json encode %s: %v", defName, err)
		}

		lineNumber := 0
		scanner := bufio.NewScanner(buf)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			lineNumber++
			if lineNumber == 1 {
				continue
			}
			line := scanner.Text()
			colonIndex := strings.Index(line, ":")
			if colonIndex == -1 {
				w.WriteString(line + "\n")
				continue
			}
			line = strings.ReplaceAll(line[0:colonIndex], `"`, "") + line[colonIndex:]
			w.WriteString(line + "\n")

		}

		w.WriteString(fmt.Sprintf(""))

		wceBuf := &bytes.Buffer{}
		err = wceGen(wceBuf, yamlDef)
		if err != nil {
			t.Fatalf("wceGen %s: %v", defName, err)
		}

		wceInst := wce.New("test")

		wceOut := ""
		if yamlDef.Comment != "" {
			wceOut += "// " + yamlDef.Comment + "\n"
		}
		wceOut += yamlDef.Name + "\n"
		wceOut += wceBuf.String()

		err = os.WriteFile(fmt.Sprintf("%s/%s.wce", dirTest, strings.ToLower(defName)), []byte(wceOut), os.ModePerm)
		if err != nil {
			t.Fatalf("write wce %s: %v", defName, err)
		}

		a := wce.AsciiReadTokenNew(wceBuf, wceInst)
		err = def.Read(a)
		if err != nil {
			t.Fatalf("parse wce %s: %v", defName, err)
		}
	}

	fmt.Println("Tested", len(defs), "definitions")

}

func wceGen(buf *bytes.Buffer, yamlDef *Definition) error {
	for _, prop := range yamlDef.Properties {
		err := traverseProp(buf, prop, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

func traverseProp(buf *bytes.Buffer, prop Property, tabCount int) error {

	commentBuf := ""
	propBuf := ""

	if prop.Comment != "" {
		commentBuf += strings.Repeat("\t", tabCount) + "// " + prop.Comment
	}

	propBuf += strings.Repeat("\t", tabCount) + prop.Name
	for i, arg := range prop.Args {
		argIndex := i + 1
		if arg.Comment != "" {
			commentBuf += fmt.Sprintf("\n %s // Argument %d (%s): %s", strings.Repeat("\t", tabCount), argIndex, arg.Format, arg.Comment)
		}
		propBuf += " "
		if arg.Example != "" {
			propBuf += arg.Example
			continue
		}
		switch arg.Format {
		case `%s`:
			propBuf += `"` + fmt.Sprintf("%d", argIndex) + `"`
		case `%d`:
			propBuf += fmt.Sprintf("%d", argIndex)
		case `%0.8e`:
			propBuf += fmt.Sprintf("%0.8f", float64(argIndex)+0.12345678)
		default:
			return fmt.Errorf("unhandled type: %s", arg.Format)
		}
	}

	if len(commentBuf) > 0 {
		buf.WriteString(commentBuf + "\n")
	}
	buf.WriteString(propBuf + "\n")
	for _, prop2 := range prop.Properties {
		err := traverseProp(buf, prop2, tabCount+1)
		if err != nil {
			return err
		}
	}

	return nil
}
