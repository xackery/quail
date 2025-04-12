package def

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/helper"
	"gopkg.in/yaml.v3"
)

var (
	knownProps = make(map[string]bool)
)

func TestWceGenPython(t *testing.T) {
	// defs declared in def_md_test.go
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
from .parse import property

class ` + defName + `:
	@staticmethod
	def definition():
		return "` + strings.ToUpper(defName) + `"

`)
		if yamlDef.HasTag {
			w.WriteString("\ttag:str\n")
		}

		w.WriteString(initPyGen(yamlDef))

		decInitBuf := &bytes.Buffer{}
		propInitBuf := &bytes.Buffer{}
		initReaderBuf := &bytes.Buffer{}
		initWriterBuf := &bytes.Buffer{}
		initWriterBuf.WriteString("\tdef read(self, tag:str, r:io.TextIOWrapper|None) -> str:\n")
		if yamlDef.HasTag {
			initWriterBuf.WriteString("\t\tself.tag = tag\n")
		}

		writeBuf := &bytes.Buffer{}
		writeBuf.WriteString("\tdef write(self, w:io.TextIOWrapper)->str:\n")
		if yamlDef.HasTag {
			writeBuf.WriteString("\t\tw.write(f\"{self.definition()} \\\"{self.tag}\\\"\\n\")\n")
		} else {
			writeBuf.WriteString("\t\tw.write(f\"{self.definition()}\\n\")\n")
		}
		err = wcePyGen(propInitBuf, decInitBuf, initReaderBuf, initWriterBuf, writeBuf, yamlDef)
		if err != nil {
			t.Fatalf("wceGen %s: %v", defName, err)
		}
		initWriterBuf.WriteString("\t\tif r is None:\n")
		initWriterBuf.WriteString("\t\t\treturn \"no reader provided\"\n")

		w.Write(initWriterBuf.Bytes())
		w.WriteString("\n")
		w.Write(initReaderBuf.Bytes())
		w.WriteString("\t\treturn \"\"\n")
		w.WriteString("\n")
		w.Write(writeBuf.Bytes())
		w.WriteString("\t\treturn \"\"\n")
		w.WriteString("\n")
	}

	fmt.Println("Generated", len(defs), "definitions")

}

func initPyGen(yamlDef *Definition) string {
	out := ""

	out += initPyProperties(2, yamlDef.Properties, "self", yamlDef.HasTag)
	out += "\n"
	return out
}

func initPyProperties(tabIndex int, props []Property, scope string, hasTag bool) string {
	out := ""
	tabIndex--
	for _, prop := range props {
		if len(prop.Properties) > 0 {
			continue
		}
		propName := strings.TrimSuffix(prop.Name, "?")
		propName = strings.ToLower(propName)
		propType := pyPropType(prop)
		if propType == "None" {
			continue
		}
		out += fmt.Sprintf("%s%s:%s\n", strings.Repeat("\t", tabIndex), propName, propType)
	}
	out += "\n"

	tabIndex++
	out += fmt.Sprintf("%sdef __init__(self):\n", strings.Repeat("\t", tabIndex-1))
	if hasTag {
		out += fmt.Sprintf("%sself.tag = \"\"\n", strings.Repeat("\t", tabIndex))
	}

	for _, prop := range props {
		if len(prop.Properties) > 0 {
			continue
		}
		propName := strings.TrimSuffix(prop.Name, "?")
		propName = strings.ToLower(propName)
		propValue := pyPropType(prop)
		switch propValue {
		case "int":
			propValue = "0"
		case "float":
			propValue = "0.0"
		case "str":
			propValue = "\"\""
		case "None":
			continue
		}

		out += fmt.Sprintf("%sself.%s = %s #%d\n", strings.Repeat("\t", tabIndex), propName, propValue, tabIndex)
	}

	for _, prop := range props {
		if len(prop.Properties) == 0 {
			continue
		}
		propName := strings.TrimSuffix(prop.Name, "?")
		propName = strings.ToLower(propName)
		if strings.HasPrefix(propName, "num") {
			propName = strings.TrimPrefix(propName, "num")
		}

		propKey := fmt.Sprintf("self.%s", propName)
		if scope != "self" {
			//propKey = fmt.Sprintf("self.%s", propName)
			//tabIndex++
		}

		out += fmt.Sprintf("%s%s = [] #%s %d\n", strings.Repeat("\t", tabIndex), propKey, scope, tabIndex)
		if scope != "self" {
			//tabIndex--
		}
		hasTag = false

	}

	for _, prop := range props {
		if len(prop.Properties) == 0 {
			continue
		}
		propName := strings.TrimSuffix(prop.Name, "?")
		propName = strings.ToLower(propName)
		if strings.HasPrefix(propName, "num") {
			propName = strings.TrimPrefix(propName, "num")
		}

		scope = propName
		out += "\n"
		out += fmt.Sprintf("%sclass %s:\n", strings.Repeat("\t", tabIndex-1), propName[:len(propName)-1])
		out += initPyProperties(tabIndex+1, prop.Properties, scope, false)
	}

	return out
}

func wcePyGen(propInitBuf *bytes.Buffer, decInitBuf *bytes.Buffer, initReaderBuf *bytes.Buffer, initWriterBuf *bytes.Buffer, writeBuf *bytes.Buffer, yamlDef *Definition) error {
	knownProps = make(map[string]bool)
	for i, prop := range yamlDef.Properties {

		if i == 0 {
			decInitBuf.WriteString(fmt.Sprintf("\tdef __init__(self):\n"))
		}
		err := traversePyProp(propInitBuf, decInitBuf, initReaderBuf, initWriterBuf, writeBuf, prop, "self", 2, 1, "")
		if err != nil {
			return err
		}
		if i == len(yamlDef.Properties)-1 {
			propInitBuf.Write(decInitBuf.Bytes())
			decInitBuf.Reset()
		}
	}

	return nil
}

func traversePyProp(propInitBuf *bytes.Buffer, decInitBuf *bytes.Buffer, initReaderBuf *bytes.Buffer, initWriterBuf *bytes.Buffer, writeBuf *bytes.Buffer, prop Property, scope string, initTabCount int, decTabCount int, treeScope string) error {

	if knownProps[prop.Name] {
		return fmt.Errorf("duplicate property: %s", prop.Name)
	}
	knownProps[prop.Name] = true
	propBuf := ""
	initBuf := ""

	propKey := strings.TrimSuffix(prop.Name, "?")
	propKey = strings.ToLower(propKey)
	if strings.HasPrefix(propKey, "num") {
		propKey = strings.TrimPrefix(propKey, "num")
	}

	if treeScope == "" {
		treeScope = "self."
	}
	treeScope += propKey[:len(propKey)-1] + "."

	isNullable := strings.HasSuffix(prop.Name, "?")
	trimName := strings.TrimSuffix(prop.Name, "?")

	isManyArg := false
	if len(prop.Args) > 0 {
		if len(prop.Properties) == 0 {
			propBuf += strings.Repeat("\t", decTabCount) + strings.ToLower(trimName) + ":"
			initBuf += fmt.Sprintf("%s\tself.%s = ", strings.Repeat("\t", decTabCount), strings.ToLower(trimName))
		}
		if len(prop.Args) > 1 {
			propBuf += "tuple["
			initBuf += "tuple["
		}

		for _, arg := range prop.Args {

			if strings.HasSuffix(arg.Format, "...") {
				isManyArg = true
				arg.Format = strings.TrimSuffix(arg.Format, "...")
			}
		}

		argLen := len(prop.Args)

		if !isManyArg {
			initReaderBuf.WriteString(fmt.Sprintf("%srecords = property(r, \"%s\", %d)\n", strings.Repeat("\t", initTabCount), prop.Name, argLen))
			if len(prop.Properties) == 0 {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s.%s = ", strings.Repeat("\t", initTabCount), scope, strings.ToLower(trimName)))
				writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s \\\"{%s.%s}\\\"\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", initTabCount-1), prop.Name, scope, strings.ToLower(trimName)))
			} else {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s = ", strings.Repeat("\t", initTabCount), strings.ToLower(trimName)))
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
						initBuf += fmt.Sprintf("tuple[%s, None]", base)
						initReaderBuf.WriteString(fmt.Sprintf("(%s(records[%d]) if records[%d] != \"NULL\" else None)", base, i+1, i+1))
					} else {
						propBuf += base
						switch base {
						case "int":
							initBuf += "0"
						case "float":
							initBuf += "0.0"
						case "str":
							initBuf += "\"\""
						}
						initReaderBuf.WriteString(fmt.Sprintf("%s(records[%d])", base, i+1))
					}
				} else {
					initReaderBuf.WriteString(fmt.Sprintf("%s(records[%d])\n", base, i+1))
				}
				if len(prop.Args) > i+1 {
					propBuf += ", "
					initBuf += ", "
					initReaderBuf.WriteString(", ")
				}
			}
		} else { // many args
			initReaderBuf.WriteString(fmt.Sprintf("%srecords = property(r, \"%s\", -1)\n", strings.Repeat("\t", initTabCount), prop.Name))
			if len(prop.Properties) == 0 {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s.%s = ", strings.Repeat("\t", initTabCount), scope, strings.ToLower(trimName)))
				writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s \\\"{%s.%s}\\\"\\n\")\n", strings.Repeat("\t", initTabCount), prop.Name, scope, strings.ToLower(trimName)))
			} else {
				initReaderBuf.WriteString(fmt.Sprintf("%s%s = ", strings.Repeat("\t", initTabCount), strings.ToLower(trimName)))
			}
			propBuf += "list[str]"
			initBuf += "list[str]"
			initReaderBuf.WriteString("records[1:]\n")
		}

		if len(prop.Args) > 1 {
			propBuf += "]"
			initBuf += "]"
		}
		if prop.Note != "" && len(prop.Properties) == 0 {
			propBuf += " # " + prop.Note
		}
		propBuf += "\n"
		initBuf += "\n"
	} else { // no argument parse
		initReaderBuf.WriteString(fmt.Sprintf("%sproperty(r, \"%s\", 0)\n", strings.Repeat("\t", initTabCount), prop.Name))
		writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", initTabCount-1), prop.Name))
	}

	initReaderBuf.WriteString("\n")

	if len(prop.Properties) > 0 {
		if len(prop.Args) != 1 {
			return fmt.Errorf("when an array of properties, count should be declared as first arg")
		}
		if prop.Args[0].Format != "%d" {
			return fmt.Errorf("parse %s: when an array of properties, format of arg 1 should be %%d", prop.Name)
		}
	}
	decInitBuf.WriteString(initBuf)
	propInitBuf.WriteString(propBuf)

	lastScope := ""
	for i, prop2 := range prop.Properties {
		prop2Name := strings.TrimSuffix(prop2.Name, "?")
		prop2Name = strings.ToLower(prop2Name)
		if i == 0 {
			decInitBuf.WriteString(fmt.Sprintf("%sself.%ss = []\n", strings.Repeat("\t", initTabCount), prop2Name))
		}
		if i == len(prop.Properties)-1 {
			decInitBuf.Write(propInitBuf.Bytes())
			propInitBuf.Reset()
		}
		if i == 0 {
			decInitBuf.WriteString(fmt.Sprintf("%sclass %s:\n", strings.Repeat("\t", decTabCount), prop2Name))
			decInitBuf.WriteString(fmt.Sprintf("%s\tdef __init__(self):\n", strings.Repeat("\t", decTabCount)))

			lastScope = scope
			scope = prop2Name + tabCode(initTabCount)
			scopeTmp := scope
			if scope != "self" {
				scopeTmp = scopeTmp[0 : len(scopeTmp)-1]
			}

			initReaderBuf.WriteString(fmt.Sprintf("%s%s.%ss = []\n", strings.Repeat("\t", initTabCount), lastScope, prop2Name))
			initReaderBuf.WriteString(fmt.Sprintf("%sfor %s in range(%s):\n", strings.Repeat("\t", initTabCount), tabCode(initTabCount), strings.ToLower(trimName)))
			initReaderBuf.WriteString(fmt.Sprintf("%s\t%s = %s()\n", strings.Repeat("\t", initTabCount), prop2Name+tabCode(initTabCount), treeScope[:len(treeScope)-1]))

			writeBuf.WriteString(fmt.Sprintf("%sw.write(f\"%s%s \\\"{len(%s.%ss)}\\\"\\n\")\n", strings.Repeat("\t", initTabCount), strings.Repeat("\\t", initTabCount-1), prop.Name, lastScope, strings.ToLower(prop2Name)))
			writeBuf.WriteString(fmt.Sprintf("%sfor %s in %s.%ss:\n", strings.Repeat("\t", initTabCount), scope, lastScope, prop2Name))
		}
		propInitBuf.WriteString("\n")
		err := traversePyProp(propInitBuf, decInitBuf, initReaderBuf, initWriterBuf, writeBuf, prop2, scope, initTabCount+1, decTabCount+1, treeScope)
		if err != nil {
			return err
		}

		if i == len(prop.Properties)-1 {
			scopeTmp := scope
			if scope != "self" {
				scopeTmp = scopeTmp[0 : len(scopeTmp)-1]
			}
			propInitBuf.WriteString(fmt.Sprintf("\n%s%ss:list[%s]\n", strings.Repeat("\t", decTabCount), scopeTmp, scopeTmp))
			lastScopeTmp := lastScope
			if lastScope != "self" {
				lastScopeTmp = lastScopeTmp[0 : len(lastScopeTmp)-1]
			}
			initReaderBuf.WriteString(fmt.Sprintf("%s\t%s.%ss.append(%s)\n", strings.Repeat("\t", initTabCount), lastScope, scopeTmp, scope))
		}

	}
	return nil
}

func tabCode(tabCount int) string {
	return string(rune('i' + tabCount - 2))
}

func pyPropType(prop Property) string {
	if len(prop.Properties) == 0 && len(prop.Args) == 0 {
		return "None"
	}

	out := ""
	if len(prop.Args) > 1 {
		out += "tuple["
	}

	for _, arg := range prop.Args {
		if strings.HasSuffix(arg.Format, "...") {
			return "list[str]"
		}
	}

	isNullable := strings.HasSuffix(prop.Name, "?")
	for i, arg := range prop.Args {

		base := ""
		switch arg.Format {
		case `%s`:
			base = "str"
		case `%d`:
			base = "int"
		case `%0.8e`:
			base = "float"
		default:
			base = "Unknown"
		}

		if len(prop.Properties) == 0 {
			if isNullable {
				out += fmt.Sprintf("tuple[%s, None]", base)
			} else {
				out += base
			}
		}

		if len(prop.Args) > i+1 {
			out += ", "
		}
	}

	if len(prop.Args) > 1 {
		out += "]"
	}
	return out
}
