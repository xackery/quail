package def

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/wce"
	"gopkg.in/yaml.v3"
)

func TestWceGenMarkdown(t *testing.T) {
	defs := []defReadWriter{
		&wce.ActorDef{},
		&wce.ActorInst{},
		&wce.AmbientLight{},
		&wce.BlitSpriteDef{},
		&wce.DMSpriteDef{},
		&wce.DMSpriteDef2{},
		&wce.DMTrackDef2{},
		&wce.EqgAniDef{},
		&wce.EqgLayDef{},
		&wce.EqgMdsDef{},
		&wce.EqgModDef{},
		&wce.EqgTerDef{},
		&wce.GlobalAmbientLightDef{},
		&wce.HierarchicalSpriteDef{},
		&wce.LightDef{},
		&wce.MaterialDef{},
		&wce.MaterialPalette{},
		&wce.ParticleCloudDef{},
		&wce.PointLight{},
		&wce.PolyhedronDefinition{},
		&wce.Region{},
		&wce.RGBTrackDef{},
		&wce.SimpleSpriteDef{},
		&wce.Sprite2DDef{},
		&wce.Sprite3DDef{},
		&wce.TrackDef{},
		&wce.TrackInstance{},
		&wce.WorldDef{},
		&wce.WorldTree{},
		&wce.Zone{},
	}

	dirTest := helper.DirTest()

	w, err := os.Create(fmt.Sprintf("%s/latest.md", dirTest))
	if err != nil {
		t.Fatalf("create: %v", err)

	}
	defer w.Close()

	w.WriteString(fmt.Sprintf("# WCEmu Latest\n\n"))

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

		w.WriteString(fmt.Sprintf(""))

		err = wceMarkdownGen(w, yamlDef)
		if err != nil {
			t.Fatalf("wceGen %s: %v", defName, err)
		}

		w.WriteString(fmt.Sprintf("\n\n\n"))

	}

	fmt.Println("Latest markdown ", len(defs), "definitions")

}

func wceMarkdownGen(buf *os.File, yamlDef *Definition) error {

	buf.WriteString(fmt.Sprintf("## %s\n\n", yamlDef.Name))
	if len(yamlDef.Note) > 0 {
		buf.WriteString(fmt.Sprintf("%s\n", yamlDef.Note))
	}
	if len(yamlDef.Description) > 0 {
		buf.WriteString(fmt.Sprintf("%s\n", yamlDef.Description))
	}

	buf.WriteString("\n```c\n")

	for _, prop := range yamlDef.Properties {
		err := traverseMarkdownProp(buf, prop, 0)
		if err != nil {
			return err
		}
	}

	buf.WriteString("```\n")

	return nil
}

func traverseMarkdownProp(buf *os.File, prop Property, tabCount int) error {

	commentBuf := ""
	propBuf := ""

	if prop.Note != "" {
		commentBuf += strings.Repeat("\t", tabCount) + "// " + prop.Note
	}

	propBuf += strings.Repeat("\t", tabCount) + prop.Name
	for i, arg := range prop.Args {
		argIndex := i + 1
		if arg.Note != "" {
			commentBuf += fmt.Sprintf("\n %s // Argument %d (%s): %s", strings.Repeat("\t", tabCount), argIndex, arg.Format, arg.Note)
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
	if len(prop.Properties) > 0 {
		if len(prop.Args) != 1 {
			return fmt.Errorf("when an array of properties, count should be declared as first arg")
		}
		if prop.Args[0].Format != "%d" {
			return fmt.Errorf("parse %s: when an array of properties, format of arg 1 should be %%d", prop.Name)
		}
	}
	for _, prop2 := range prop.Properties {
		err := traverseMarkdownProp(buf, prop2, tabCount+1)
		if err != nil {
			return err
		}
	}

	return nil
}
