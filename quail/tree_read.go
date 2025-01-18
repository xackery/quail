package quail

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/tree"
)

const (
	ErrorInvalidExt = "invalid extension"
)

// TreeRead imports the quail target file
func (q *Quail) TreeRead(w io.Writer, path string) error {
	var file string
	if strings.Contains(path, ":") {
		file = strings.Split(path, ":")[1]
		path = strings.Split(path, ":")[0]
	}
	isValidExt := false
	exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return q.treeReadFile(w, nil, path, file)
	}

	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("%s load: %w", ext, err)
	}

	return q.treeReadFile(w, pfs, path, file)
}

func (q *Quail) TreeCompare(path1 string, path2 string) error {
	var err error
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	err = q.TreeRead(buf1, path1)
	if err != nil {
		return fmt.Errorf("tree read: %w", err)
	}
	err = q.TreeRead(buf2, path2)
	if err != nil {
		return fmt.Errorf("tree read: %w", err)
	}

	err = treeCompare(buf1.String(), buf2.String())
	if err != nil {
		return fmt.Errorf("compare: %w", err)
	}

	return nil
}

func (q *Quail) treeReadFile(w io.Writer, pfs *pfs.Pfs, path string, file string) error {
	if pfs == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return q.treeInspectContent(w, filepath.Base(path), bytes.NewReader(data))
	}

	isFound := false
	for _, fe := range pfs.Files() {
		if len(file) > 1 && !strings.EqualFold(fe.Name(), file) {
			continue
		}

		err := q.treeInspectContent(w, fe.Name(), bytes.NewReader(fe.Data()))
		if err != nil && err.Error() != ErrorInvalidExt {
			return fmt.Errorf("inspect content: %w", err)
		}
		isFound = true
	}
	if isFound {
		return nil
	}
	if len(file) < 2 {
		return fmt.Errorf("no files found to tree")
	}

	return fmt.Errorf("%s not found in %s", file, filepath.Base(path))
}

func (q *Quail) treeInspectContent(w io.Writer, file string, r *bytes.Reader) error {
	var err error
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".wld":
		fmt.Fprintf(w, "Tree: %s\n", file)
		rawWld := &raw.Wld{}
		err = rawWld.Read(r)
		if err != nil {
			return fmt.Errorf("wld read: %w", err)
		}

		isChr := false
		if strings.Contains(file, "_chr") {
			isChr = true
		}

		err = tree.Dump(w, isChr, rawWld)
		if err != nil {
			return fmt.Errorf("tree dump: %w", err)
		}

		return nil
	}
	return fmt.Errorf("%s", ErrorInvalidExt)
}

func treeCompare(val, val2 string) error {
	if val == val2 {
		fmt.Println("Files are identical")
		return nil
	}

	type treeEntry struct {
		defType       string
		srcLineNumber int
		src           string
		srcCount      int
		dstLineNumber int
		dst           string
		tag           string
		dstCount      int
	}
	// trees have a structure like Root Foo: 123 (DEF)
	// we mainly are looking for identical definitions
	entries := make(map[string]*treeEntry)

	regex := regexp.MustCompile(`(.*): (\d+) \((.*)\)`)

	lines := strings.Split(val, "\n")
	for lineNumber, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}
		defType := strings.TrimSpace(matches[1])
		frag := matches[2]
		tag := matches[3]
		entry, ok := entries[tag]
		if !ok {
			entries[tag] = &treeEntry{
				defType:       defType,
				src:           frag,
				srcLineNumber: lineNumber,
				srcCount:      1,
				tag:           tag,
			}
			continue
		}
		entry.srcCount++
	}
	lines = strings.Split(val2, "\n")
	for lineNumber, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}
		defType := strings.TrimSpace(matches[1])
		frag := matches[2]
		tag := matches[3]
		_, ok := entries[tag]
		if !ok {
			entries[tag] = &treeEntry{
				defType:       defType,
				dst:           frag,
				dstLineNumber: lineNumber,
				dstCount:      1,
				tag:           tag,
			}
			continue
		}
		entries[tag].dst = frag
		entries[tag].dstLineNumber = lineNumber
		entries[tag].dstCount++
	}

	missingTotal := 0
	for tag, entry := range entries {
		if entry.src == "" {
			fmt.Printf("dst tree line %d missing on src: %s\n", entry.dstLineNumber, entry.dst)
			missingTotal++
			continue
		}
		if entry.dst == "" {
			fmt.Printf("src tree line %d missing on dst: %s\n", entry.srcLineNumber, entry.src)
			missingTotal++
			continue
		}
		if entry.dstCount != entry.srcCount {
			fmt.Printf("!! %s %s src count: %d, dst count: %d\n", entry.defType, tag, entry.srcCount, entry.dstCount)
			missingTotal++
			continue
		}
	}

	if missingTotal > 0 {
		return fmt.Errorf("%d mismatched entries", missingTotal)
	}

	fmt.Println("Files have same count of nodes")

	return nil
}
