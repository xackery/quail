package def

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/wce"
	"gopkg.in/yaml.v3"
)

func TestWceGenPython(t *testing.T) {

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
		&wce.EqgModDef{},
		&wce.EqgMdsDef{},
		&wce.EqgTerDef{},
		&wce.EqgAniDef{},
		&wce.EqgLayDef{},
		&wce.EqgParticlePointDef{},
		&wce.EqgParticleRenderDef{},
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

		w, err := os.Create(fmt.Sprintf("%s/%s.py", dirTest, defName))
		if err != nil {
			t.Fatalf("create failed: %s", err)
		}
		defer w.Close()

		w.WriteString(`# Generated from quail, DO NOT EDIT
import io
import wce.parse as parse

class ` + defName + `:
	@staticmethod
	def definition():
		return "` + strings.ToUpper(defName) + `"

`)
		if yamlDef.HasTag {
			w.WriteString("\ttag:str\n")
		}

		decBuf := &bytes.Buffer{}
		initBuf := &bytes.Buffer{}
		initBuf.WriteString("\tdef __init__(self, tag:str, r:io.TextIOWrapper):\n")
		if yamlDef.HasTag {
			initBuf.WriteString("\t\tself.tag = tag\n")
		}
		writeBuf := &bytes.Buffer{}
		writeBuf.WriteString("\tdef write(self, w:io.TextIOWrapper):\n")
		if yamlDef.HasTag {
			writeBuf.WriteString("\t\tw.write(f\"{self.definition()} \\\"{self.tag}\\\"\\n\")\n")
		} else {
			writeBuf.WriteString("\t\tw.write(f\"{self.definition()}\\n\")\n")
		}
		err = wcePyGen(decBuf, initBuf, writeBuf, yamlDef)
		if err != nil {
			t.Fatalf("wceGen %s: %v", defName, err)
		}
		w.Write(decBuf.Bytes())
		w.WriteString("\n")
		w.Write(initBuf.Bytes())
		w.WriteString("\n")
		w.Write(writeBuf.Bytes())
		w.WriteString("\n")
	}

	fmt.Println("Generated", len(defs), "definitions")

}

func wcePyGen(decBuf *bytes.Buffer, initBuf *bytes.Buffer, writeBuf *bytes.Buffer, yamlDef *Definition) error {
	for _, prop := range yamlDef.Properties {
		err := traversePyProp(decBuf, initBuf, writeBuf, prop, "self", 2, 1, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func traversePyProp(decBuf *bytes.Buffer, initBuf *bytes.Buffer, writeBuf *bytes.Buffer, prop Property, scope string, initTabCount int, decTabCount int, treeScope string) error {

	propBuf := ""
	scopeTmp := scope
	if scope != "self" {
		scopeTmp = scopeTmp[0 : len(scopeTmp)-1]
	}
	treeScope += strings.ToLower(scopeTmp) + "."

	isNullable := strings.HasSuffix(prop.Name, "?")
	trimName := strings.TrimSuffix(prop.Name, "?")

	isManyArg := false
	if len(prop.Args) > 0 {
		if len(prop.Properties) == 0 {
			propBuf += strings.Repeat("\t", decTabCount) + strings.ToLower(trimName) + ":"
		}
		if len(prop.Args) > 1 {
			propBuf += "tuple["
		}

		for _, arg := range prop.Args {

			if strings.HasSuffix(arg.Format, "...") {
				isManyArg = true
				arg.Format = strings.TrimSuffix(arg.Format, "...")
			}
		}

		argLen := len(prop.Args)

		if !isManyArg {
			initBuf.WriteString(fmt.Sprintf("%srecords = parse.property(r, \"%s\", %d)\n", strings.Repeat("\t", initTabCount), prop.Name, argLen))
			if len(prop.Properties) == 0 {
				initBuf.WriteString(fmt.Sprintf("%s%s.%s = ", strings.Repeat("\t", initTabCount), scope, strings.ToLower(trimName)))
				writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s \\\"{%s.%s}\\\"\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", initTabCount-1), prop.Name, scope, strings.ToLower(trimName)))
			} else {
				initBuf.WriteString(fmt.Sprintf("%s%s = ", strings.Repeat("\t", initTabCount), strings.ToLower(trimName)))
			}
			for i, arg := range prop.Args {

				if strings.HasSuffix(arg.Format, "...") {
					isManyArg = true
					arg.Format = strings.TrimSuffix(arg.Format, "...")
				}

				base := ""
				switch arg.Format {
				case `%s`:
					base = "str"
				case `%d`:
					base = "int"
				case `%0.8e`:
					base = "float"
				default:
					return fmt.Errorf("unhandled type: %s", arg.Format)
				}

				if len(prop.Properties) == 0 {
					if isNullable {
						propBuf += fmt.Sprintf("tuple[%s, None]", base)
						initBuf.WriteString(fmt.Sprintf("(%s(records[%d]) if records[%d] != \"NULL\" else None)", base, i+1, i+1))
					} else {
						propBuf += base
						initBuf.WriteString(fmt.Sprintf("%s(records[%d])", base, i+1))
					}
				} else {
					initBuf.WriteString(fmt.Sprintf("%s(records[%d])\n", base, i+1))
				}
				if len(prop.Args) > i+1 {
					propBuf += ", "
					initBuf.WriteString(", ")
				}
			}
		} else { // many args
			initBuf.WriteString(fmt.Sprintf("%srecords = parse.property(r, \"%s\", -1)\n", strings.Repeat("\t", initTabCount), prop.Name))
			if len(prop.Properties) == 0 {
				initBuf.WriteString(fmt.Sprintf("%s%s.%s = ", strings.Repeat("\t", initTabCount), scope, strings.ToLower(trimName)))
				writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s \\\"{%s.%s}\\\"\\n\")\n", strings.Repeat("\t", initTabCount), prop.Name, scope, strings.ToLower(trimName)))
			} else {
				initBuf.WriteString(fmt.Sprintf("%s%s = ", strings.Repeat("\t", initTabCount), strings.ToLower(trimName)))
			}
			propBuf += "list[str]"
			initBuf.WriteString("records[1:]\n")
		}

		if len(prop.Args) > 1 {
			propBuf += "]"
		}
		if prop.Note != "" && len(prop.Properties) == 0 {
			propBuf += " # " + prop.Note
		}
		propBuf += "\n"
	} else { // no argument parse
		initBuf.WriteString(fmt.Sprintf("%sparse.property(r, \"%s\", 0)\n", strings.Repeat("\t", initTabCount), prop.Name))
		writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", initTabCount-1), prop.Name))
	}

	initBuf.WriteString("\n")

	if len(prop.Properties) > 0 {
		if len(prop.Args) != 1 {
			return fmt.Errorf("when an array of properties, count should be declared as first arg")
		}
		if prop.Args[0].Format != "%d" {
			return fmt.Errorf("parse %s: when an array of properties, format of arg 1 should be %%d", prop.Name)
		}
	}
	decBuf.WriteString(propBuf)
	lastScope := ""
	for i, prop2 := range prop.Properties {
		prop2Name := strings.TrimSuffix(prop2.Name, "?")
		prop2Name = strings.ToLower(prop2Name)
		if i == 0 {
			decBuf.WriteString(fmt.Sprintf("%sclass %s:", strings.Repeat("\t", decTabCount), prop2Name))

			lastScope = scope
			scope = prop2Name + tabCode(initTabCount)
			scopeTmp := scope
			if scope != "self" {
				scopeTmp = scopeTmp[0 : len(scopeTmp)-1]
			}

			initBuf.WriteString(fmt.Sprintf("%s%s.%ss = []\n", strings.Repeat("\t", initTabCount), lastScope, prop2Name))
			initBuf.WriteString(fmt.Sprintf("%sfor %s in range(%s):\n", strings.Repeat("\t", initTabCount), tabCode(initTabCount), strings.ToLower(trimName)))
			initBuf.WriteString(fmt.Sprintf("%s\t%s = %s%s()\n", strings.Repeat("\t", initTabCount), prop2Name+tabCode(initTabCount), treeScope, prop2Name))

			writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s \\\"{len(%s.%ss)}\\\"\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", initTabCount-1), prop.Name, lastScope, strings.ToLower(prop2Name)))
			writeBuf.WriteString(fmt.Sprintf("%sfor %s in %s.%ss:\n", strings.Repeat("\t", initTabCount), scope, lastScope, prop2Name))
		}

		decBuf.WriteString("\n")
		err := traversePyProp(decBuf, initBuf, writeBuf, prop2, scope, initTabCount+1, decTabCount+1, treeScope)
		if err != nil {
			return err
		}

		if i == len(prop.Properties)-1 {
			scopeTmp := scope
			if scope != "self" {
				scopeTmp = scopeTmp[0 : len(scopeTmp)-1]
			}
			decBuf.WriteString(fmt.Sprintf("\n%s%ss:list[%s]\n", strings.Repeat("\t", decTabCount), scopeTmp, scopeTmp))
			lastScopeTmp := lastScope
			if lastScope != "self" {
				lastScopeTmp = lastScopeTmp[0 : len(lastScopeTmp)-1]
			}
			initBuf.WriteString(fmt.Sprintf("%s\t%s.%ss.append(%s)\n", strings.Repeat("\t", initTabCount), lastScope, scopeTmp, scope))
		}
	}

	return nil
}

func tabCode(tabCount int) string {
	return string(rune('i' + tabCount - 2))
}
