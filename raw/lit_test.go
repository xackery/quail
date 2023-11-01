package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestLitRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name    string
		litName string
	}{
		//.lit|1|commons_inn_obj_lampc01.lit|commonlands.eqg
		//{name: "commonlands.eqg", litName: "commons_inn_obj_lampc01.lit"},
		//.lit|1|communalhut_obj_treasureb01.lit|buriedsea.eqg
		{name: "buriedsea.eqg", litName: "communalhut_obj_treasureb01.lit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".lay" {
					continue
				}
				lit := &Lit{}
				err = lit.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}
