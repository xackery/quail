package mod

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestMOD_BlenderImport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	os.MkdirAll("test", 0755)
	tests := []struct {
		name    string
		srcDir  string
		dstDir  string
		wantErr bool
	}{
		{name: "it13926.eqg", srcDir: "test/_it13926.eqg", dstDir: "test/decode.mod", wantErr: false},
		//{name: "test", eqgPath: fmt.Sprintf("%s/it13900.eqg", eqPath), dstDir: "test/", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eqPath := fmt.Sprintf("%s/%s", eqPath, tt.name)
			pfs, err := eqg.NewFile(eqPath)
			if err != nil {
				t.Fatalf("Failed to open eqg: %s", err.Error())
			}
			baseName := filepath.Base(tt.name)
			baseName = baseName[:len(baseName)-4]
			modName := baseName + ".mod"

			isFound := false
			var data []byte
			for _, file := range pfs.Files() {
				if file.Name() != modName {
					continue
				}
				isFound = true
				data = file.Data()
				break
			}
			if !isFound {
				t.Fatalf("Failed to find %s in eqg", modName)
			}

			e, err := New(modName, nil)
			if err != nil {
				t.Fatalf("Failed to create mod: %s", err.Error())
			}

			err = e.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("Failed to decode mod: %s", err.Error())
			}

			err = e.BlenderExport(tt.srcDir)
			if err != nil {
				t.Fatalf("Failed to blender export mod: %s", err.Error())
			}

			err = e.BlenderImport(tt.srcDir + "/_" + modName)
			if (err != nil) != tt.wantErr {
				t.Errorf("mod.BlenderImport() error = %v, wantErr %v", err, tt.wantErr)
			}

			e.version = 3

			w, err := os.Create(tt.dstDir)
			if err != nil {
				t.Fatalf("Failed to create file: %s", err.Error())
			}
			err = e.Encode(w)
			if err != nil {
				t.Fatalf("Failed to encode mod: %s", err.Error())
			}
			w.Close()

		})
	}
}
