package ter

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/pfs/archive"
)

func TestObjImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &TER{}
	objPath := "test/soldungb.obj"
	mtlPath := "test/soldungb.mtl"
	matTxtPath := "test/soldungb_material.txt"

	err := e.ImportObj(objPath, mtlPath, matTxtPath)
	if err != nil {
		t.Fatalf("import: %s", err)
	}

	w, err := os.Create("test/out.ter")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}

func TestObjImportWrite(t *testing.T) {
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

	w, err := os.Create("test/out.ter")
	if err != nil {
		t.Fatalf("create: %s", err)
	}

	err = e.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}

func TestObjImportExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	path, err := archive.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	objPath := "test/soldungb.obj"
	mtlPath := "test/soldungb.mtl"
	matTxtPath := "test/soldungb_material.txt"

	err = e.ImportObj(objPath, mtlPath, matTxtPath)
	if err != nil {
		t.Fatalf("import: %s", err)
	}

	err = e.ObjExport("test/out.obj", "test/out.mtl", "test/out_material.txt")
	if err != nil {
		t.Fatalf("export: %s", err)
	}
	err = compare(objPath, "test/out.obj")
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
