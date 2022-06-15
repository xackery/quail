package ter

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestObjImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &TER{}
	objPath := "../eq/soldungb/cache/soldungb.obj"
	mtlPath := "../eq/soldungb/cache/soldungb.mtl"
	matTxtPath := "../eq/soldungb/cache/soldungb_material.txt"

	err := e.ImportObj(objPath, mtlPath, matTxtPath)
	if err != nil {
		t.Fatalf("import: %s", err)
	}

	w, err := os.Create("../eq/tmp/out.ter")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}

func TestObjImportSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &TER{}
	objPath := "../eq/soldungb/soldungb.obj"
	mtlPath := "../eq/soldungb/soldungb.mtl"
	matTxtPath := "../eq/soldungb/soldungb_material.txt"

	err := e.ImportObj(objPath, mtlPath, matTxtPath)
	if err != nil {
		t.Fatalf("import: %s", err)
	}

	w, err := os.Create("../eq/tmp/out.ter")
	if err != nil {
		t.Fatalf("create: %s", err)
	}

	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}

func TestObjImportExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	objPath := "../eq/soldungb/cache/soldungb.obj"
	mtlPath := "../eq/soldungb/cache/soldungb.mtl"
	matTxtPath := "../eq/soldungb/cache/soldungb_material.txt"

	err = e.ImportObj(objPath, mtlPath, matTxtPath)
	if err != nil {
		t.Fatalf("import: %s", err)
	}

	err = e.ExportObj("../eq/tmp/out.obj", "../eq/tmp/out.mtl", "../eq/tmp/out_material.txt")
	if err != nil {
		t.Fatalf("export: %s", err)
	}
	err = compare(objPath, "../eq/tmp/out.obj")
	if err != nil {
		t.Fatalf("compare: %s", err)
	}
	t.Fatalf("comp done")
}

func compare(src string, dst string) error {
	lineNumber := 0
	rs, err := os.Open(src)
	if err != nil {
		return err
	}
	defer rs.Close()

	rd, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer rd.Close()

	rsr := bufio.NewScanner(rs)
	rdr := bufio.NewScanner(rd)
	for rsr.Scan() {
		if !rdr.Scan() {
			return fmt.Errorf("rdr ended")
		}
		lineNumber++
		if lineNumber < 5 {
			continue
		}
		lineS := rsr.Text()
		lineD := rdr.Text()
		if lineS == lineD {
			continue
		}
		fmt.Printf("line %d mismatch: %s vs %s\n", lineNumber, lineS, lineD)
		return fmt.Errorf("ended early")
	}
	err = rsr.Err()
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	return nil
}
