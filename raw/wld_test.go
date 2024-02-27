package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"gopkg.in/yaml.v3"
)

func TestWldRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name string
	}{
		//{"crushbone"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}

			wld := &Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

		})
	}
}

func TestWldFragOffsetDump(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		name string
	}{
		//{"gequip4"},
		//{"global_chr"}, // TODO:  anarelion asked mesh of EYE_DMSPRITEDEF check if the eye is just massive 22 units in size, where the other units in that file are just 1-2 units in size
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}

			wld := &Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			path := fmt.Sprintf("%s/%s.wld.yaml", dirTest, tt.name)
			w, err := os.Create(path)
			if err != nil {
				t.Fatalf("failed to create %s: %s", tt.name, err.Error())
			}
			enc := yaml.NewEncoder(w)
			enc.SetIndent(2)
			err = enc.Encode(wld.Fragments)
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}
			w.Close()
			fmt.Println("wrote", path)
		})
	}
}

func TestWldRewrite(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()

	tests := []struct {
		name string
	}{
		{"gequip4"},
		//{"global_chr"}, // TODO:  anarelion asked mesh of EYE_DMSPRITEDEF check if the eye is just massive 22 units in size, where the other units in that file are just 1-2 units in size
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// copy original
			copyData, err := os.ReadFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.src.s3d", dirTest, tt.name), copyData, 0644)
			if err != nil {
				t.Fatalf("failed to write s3d %s: %s", tt.name, err.Error())
			}

			archive, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer archive.Close()

			// get wld
			data, err := archive.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}

			// write wld
			err = os.WriteFile(fmt.Sprintf("%s/%s.src.wld", dirTest, tt.name), data, 0644)
			if err != nil {
				t.Fatalf("failed to write wld %s: %s", tt.name, err.Error())
			}

			wld := &Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			w, err := os.Create(fmt.Sprintf("%s/%s.dst.wld", dirTest, tt.name))
			if err != nil {
				t.Fatalf("failed to create %s: %s", tt.name, err.Error())
			}
			defer w.Close()

			err = wld.Write(w)
			if err != nil {
				t.Fatalf("failed to write %s: %s", tt.name, err.Error())
			}

			archive, err = pfs.New(tt.name)
			if err != nil {
				t.Fatalf("failed to create %s: %s", tt.name, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = wld.Write(buf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", tt.name, err.Error())
			}

			err = archive.Add(fmt.Sprintf("%s.wld", tt.name), buf.Bytes())
			if err != nil {
				t.Fatalf("failed to add %s: %s", tt.name, err.Error())
			}

			aw, err := os.Create(fmt.Sprintf("%s/%s.dst.s3d", dirTest, tt.name))
			if err != nil {
				t.Fatalf("failed to create %s: %s", tt.name, err.Error())
			}
			defer aw.Close()

			err = archive.Write(aw)
			if err != nil {
				t.Fatalf("failed to write %s: %s", tt.name, err.Error())
			}

			// write yaml
			buf = bytes.NewBuffer(nil)
			enc := yaml.NewEncoder(buf)
			enc.SetIndent(2)
			err = enc.Encode(wld)
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}
			err = os.WriteFile(fmt.Sprintf("%s/%s.src.yaml", dirTest, tt.name), buf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write yaml %s: %s", tt.name, err.Error())
			}

			// now read the dst.wld and generate a yaml
			yamlRead, err := os.Open(fmt.Sprintf("%s/%s.dst.wld", dirTest, tt.name))
			if err != nil {
				t.Fatalf("failed to open %s: %s", tt.name, err.Error())
			}
			defer yamlRead.Close()

			wld = &Wld{}
			err = wld.Read(yamlRead)
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf = bytes.NewBuffer(nil)
			enc = yaml.NewEncoder(buf)
			enc.SetIndent(2)
			err = enc.Encode(wld)
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}
			err = os.WriteFile(fmt.Sprintf("%s/%s.dst.yaml", dirTest, tt.name), buf.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write yaml %s: %s", tt.name, err.Error())
			}

		})
	}
}
