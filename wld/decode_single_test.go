package wld

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/s3d"
)

func TestDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	category := "crushbone"
	path := "test/eq/" + category + ".s3d"
	file := category + ".wld"

	archive, err := s3d.NewFile(path)
	if err != nil {
		t.Fatalf("s3d new: %s", err)
	}

	dump.New(file)
	defer dump.WriteFileClose(path + "_" + file)
	e, err := NewFile(category, archive, file)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	if len(e.materials) != 42 {
		t.Fatalf("wanted 42 materials, got %d", len(e.materials))
	}

	if len(e.meshes) != 2694 {
		t.Fatalf("wanted 2694 meshes, got %d", len(e.meshes))
	}

}

func TestDecodeSingleZonePoints(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	category := "crushbone"
	path := "test/eq/" + category + ".s3d"
	archive, err := s3d.NewFile(path)
	if err != nil {
		t.Fatalf("s3d new: %s", err)
	}

	data, err := archive.File(category + ".zon")
	if err != nil {
		t.Fatalf("s3d.file: %s", err)
	}

	//dump.New(path)
	//defer dump.WriteFileClose(path)

	e, err := New("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode: %s", err)
	}

	fmt.Println("bazaar points:")
	/*for _, region := range e.regions {
		fmt.Println(region.name, region.center)
	}*/
}

func TestDecodeITModels(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	var archive common.ArchiveReadWriter

	names := make(map[string]bool)

	filesWithIT := ""

	nameRegex := regexp.MustCompile(`IT([0-9]+).*`)
	name2Regex := regexp.MustCompile(`it([0-9]+).*.mod`)

	path := "test/eq/"
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatalf("entries: %s", err)
	}
	for _, entry := range entries {
		isFound := false
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".eqg" && ext != ".s3d" {
			continue
		}

		noExtName := entry.Name()
		noExtName = noExtName[0 : len(noExtName)-4]

		fmt.Println("parsing", entry.Name())
		if ext == ".s3d" {
			archive, err = s3d.NewFile(path + entry.Name())
			if err != nil {
				t.Fatalf("s3d newFile: %s", err)
			}
			e, err := NewFile(noExtName, archive, noExtName+".wld")
			if err != nil {
				t.Fatalf("newfile wld: %s", err)
			}

			for _, name := range e.NameCache {
				if !strings.HasPrefix(name, "IT") {
					continue
				}
				if strings.Contains(name, "_") {
					continue
				}
				//fmt.Println(name)
				matches := nameRegex.FindAllStringSubmatch(name, -1)
				if len(matches) > 0 {
					names[matches[0][1]] = true
					if !isFound {
						isFound = true
						filesWithIT += entry.Name() + ", "
					}
				}
			}
			continue
		}
		if ext == ".eqg" {
			archive, err = eqg.NewFile(path + entry.Name())
			if err != nil {
				t.Fatalf("eqg newFile: %s", err)
			}
			for _, entry := range archive.Files() {
				name := strings.ToLower(entry.Name())
				if !strings.HasPrefix(name, "it") {
					continue
				}
				//fmt.Println(name)
				matches := name2Regex.FindAllStringSubmatch(name, -1)
				if len(matches) > 0 {
					names[matches[0][1]] = true
					if !isFound {
						isFound = true
						filesWithIT += entry.Name() + ", "
					}
				}
			}
		}
	}

	w, err := os.Create("weapons.txt")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	w.WriteString(fmt.Sprintf("found at: %s\n", filesWithIT[0:len(filesWithIT)-2]))
	for name := range names {
		w.WriteString(fmt.Sprintf("IT" + name + "\n"))
	}
}
