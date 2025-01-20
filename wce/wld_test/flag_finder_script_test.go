package wce_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func TestFragFlags(t *testing.T) {
	if os.Getenv("SCRIPT_TEST") != "1" {
		t.Skip("skipping test: SCRIPT_TEST not set")
	}

	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	type flagEntry struct {
		fragID      int
		baseName    string
		wldName     string
		tagName     string
		hexFlagDump string
	}
	flagSorts := []int{}

	flagsMap := map[int]*flagEntry{}

	dirs, err := os.ReadDir(eqPath)
	if err != nil {
		t.Fatalf("Failed to read eq path: %s", err.Error())
	}

	fmt.Printf("Checking %d files\n", len(dirs))
	parseCount := 0
	for _, dir := range dirs {
		parseCount++
		if parseCount%100 == 0 {
			fmt.Printf("Checking %d/%d\n", parseCount, len(dirs))
		}

		ext := filepath.Ext(dir.Name())
		if ext != ".s3d" {
			continue
		}
		if dir.IsDir() {
			continue
		}

		s3dName := dir.Name()
		s3dPath := fmt.Sprintf("%s/%s", eqPath, s3dName)

		pfs, err := pfs.NewFile(s3dPath)
		if err != nil {
			t.Fatalf("Failed to open s3d %s: %s", s3dName, err.Error())
		}

		for _, file := range pfs.Files() {
			if filepath.Ext(file.Name()) != ".wld" {
				continue
			}
			wldName := file.Name()
			if wldName == "objects.wld" {
				continue
			}
			if wldName == "lights.wld" {
				continue
			}

			rawWld := &raw.Wld{}
			rawWld.NameClear()
			err = rawWld.Read(bytes.NewReader(file.Data()))
			if err != nil {
				t.Fatalf("Failed to read wld: %s", err.Error())
			}

			for i := 0; i < len(rawWld.Fragments); i++ {
				fragRaw := rawWld.Fragments[i]

				frag, ok := fragRaw.(*rawfrag.WldFragBlitSpriteDef)
				if !ok {
					continue
				}

				flags := uint32(frag.Flags)
				tagName := rawWld.Name(int32(frag.NameRef()))
				hexFlagDump := hexFlagDump(int(flags))
				flagsMap[int(flags)] = &flagEntry{
					baseName:    s3dName,
					wldName:     wldName,
					fragID:      i,
					tagName:     tagName,
					hexFlagDump: hexFlagDump,
				}
				flagSorts = append(flagSorts, int(flags))
				fmt.Printf("Flag %d found in %s/%s/%s fragID %d %s\n", flags, s3dName, wldName, tagName, i, hexFlagDump)

				/* for _, face := range frag.Faces {
					flags := uint32(face.Flags)
					if flagsMap[int(flags)] != nil {
						continue
					}

					hexFlagDump := hexFlagDump(int(flags))
					flagsMap[int(flags)] = &flagEntry{
						baseName:    s3dName,
						wldName:     wldName,
						fragID:      i,
						tagName:     tagName,
						hexFlagDump: hexFlagDump,
					}
					flagSorts = append(flagSorts, int(flags))
					fmt.Printf("Flag %d found in %s/%s/%s fragID %d %s\n", flags, s3dName, wldName, tagName, i, hexFlagDump)
				} */
			}
		}
	}

	w, err := os.Create("flags.txt")
	if err != nil {
		t.Fatalf("Failed to create flags.txt: %s", err.Error())
	}
	defer w.Close()

	sort.Ints(flagSorts)
	maxFlagFound := 0
	for _, flag := range flagSorts {
		frag := flagsMap[flag]
		fmt.Fprintf(w, "Flag %d found in %s/%s/%s fragID %d %s\n", flag, frag.baseName, frag.wldName, frag.tagName, frag.fragID, frag.hexFlagDump)
		if maxFlag(flag) > maxFlagFound {
			maxFlagFound = maxFlag(flag)
		}
	}
	fmt.Printf("Max flag found: %d 0x%02x\n", maxFlagFound, maxFlagFound)
	fmt.Fprintf(w, "Max flag found: %d 0x%02x\n", maxFlagFound, maxFlagFound)
}

func hexFlagDump(flag int) string {
	out := fmt.Sprintf("Flag %d (0x%x) -- ", flag, flag)
	for i := 1; i < 32; i++ {
		if flag&(1<<i) > 0 {
			out += fmt.Sprintf("0x%02x ", 1<<i)
		}
	}
	return out
}

func maxFlag(flag int) int {
	for i := 1; i < 64; i++ {
		if flag&(1<<i) != 0 {
			return 1 << i
		}
	}
	return 0
}
