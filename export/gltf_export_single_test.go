package export

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestGLTFFlushEQPath(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	files, err := os.ReadDir("test/eq/")
	if err != nil {
		t.Fatalf("readdir: %s", err)
	}
	for _, fe := range files {
		if !strings.HasSuffix(fe.Name(), ".gltf") {
			continue
		}
		err = os.Remove(fmt.Sprintf("test/eq/%s", fe.Name()))
		if err != nil {
			t.Fatalf("remove: %s", err)
		}
	}
}

func TestGLTFEncodeSamplesSingleTest(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	tests := []struct {
		category string
	}{
		/*{category: "bxi"},
		{category: "bnf"},
		{category: "bnm"},
		{category: "bnr"},
		{category: "bnx"},*/
		//{category: "bny"},
		//{category: "dv6"},
		//{category: "lth"}, //Morell thule
		{category: "djm"}, //Djinn Male
		//{category: "zmm"}, //Zombie Male
		//{category: "wrm"}, //worm
		//{category: "dgj"}, //jade dragon
		//{category: "prt"},
	}
	for _, tt := range tests {

		fmt.Println("loading", tt.category)
		eqgFile := fmt.Sprintf("test/eq/%s.eqg", tt.category)
		//modFile := fmt.Sprintf("obj_%s.mod", tt.category)

		//outFile := fmt.Sprintf("test/eq/%s_mod.gltf", tt.category)
		//txtFile := fmt.Sprintf("test/eq/%s_mod.txt", tt.category)

		archive, err := eqg.NewFile(eqgFile)
		if err != nil {
			t.Fatalf("eqg new: %s", err)
		}

		e, err := New(tt.category, archive)
		if err != nil {
			t.Fatalf("export new: %s", err)
		}

		err = e.LoadArchive()
		if err != nil {
			t.Fatalf("loadArchive %s: %s", tt.category, err)
		}

		outFile := fmt.Sprintf("test/eq/%s.gltf", tt.category)
		w, err := os.Create(outFile)
		if err != nil {
			t.Fatalf("create %s", err)
		}
		defer w.Close()

		doc, err := gltf.New()
		if err != nil {
			t.Fatalf("gltf.New: %s", err)
		}

		err = e.GLTFEncode(doc)
		if err != nil {
			t.Fatalf("gltf: %s", err)
		}

		err = doc.Export(w)
		if err != nil {
			t.Fatalf("export: %s", err)
		}

	}
}
