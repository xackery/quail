package ter

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestGLTFExportBroodlands(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	zone := "broodlands"
	path := fmt.Sprintf("test/eq/_%s.eqg", zone)
	inFile := fmt.Sprintf("test/eq/_%s.eqg/ter_%s.ter", zone, zone)
	outFile := fmt.Sprintf("test/eq/%s.gltf", zone)

	e, err := New("arena", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", inFile, err)
	}

	fw, err := os.Create(fmt.Sprintf("test/%s.txt", zone))
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fw.Close()
	fmt.Fprintf(fw, "faces:\n")
	for i, o := range e.faces {
		fmt.Fprintf(fw, "%d %+v\n", i, o)
	}

	fmt.Fprintf(fw, "vertices:\n")
	for i, o := range e.vertices {
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFExportCityOfBronze(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	zone := "cityofbronze"
	path := fmt.Sprintf("test/eq/_%s.eqg", zone)
	inFile := fmt.Sprintf("test/eq/_%s.eqg/ter_%s.ter", zone, zone)
	outFile := fmt.Sprintf("test/eq/%s.gltf", zone)
	isDumpEnabed := false

	if isDumpEnabed {
		d, err := dump.New(path)
		if err != nil {
			t.Fatalf("dump.new: %s", err)
		}
		defer d.Save(fmt.Sprintf("%s.png", inFile))
	}
	e, err := New("cityofbronze", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	fw, err := os.Create(fmt.Sprintf("test/%s.txt", zone))
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fw.Close()
	fmt.Fprintf(fw, "faces:\n")
	for i, o := range e.faces {
		fmt.Fprintf(fw, "%d %+v\n", i, o)
	}

	fmt.Fprintf(fw, "vertices:\n")
	for i, o := range e.vertices {
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
