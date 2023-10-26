package zon

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

func TestDecode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .zon|1|anguish.zon|anguish.eqg
		//{name: "anguish.eqg"},
		// .zon|1|bazaar.zon|bazaar.eqg
		//{name: "bazaar.eqg"},
		// .zon|1|bloodfields.zon|bloodfields.eqg
		{name: "bloodfields.eqg"},
		// .zon|1|broodlands.zon|broodlands.eqg
		//{name: "broodlands.eqg"},
		// .zon|1|catacomba.zon|dranikcatacombsa.eqg
		//{name: "dranikcatacombsa.eqg"},
		// .zon|1|wallofslaughter.zon|wallofslaughter.eqg
		//{name: "wallofslaughter.eqg"},
		// .zon|2|arginhiz.zon|arginhiz.eqg
		//{name: "arginhiz.eqg"},
		// .zon|2|guardian.zon|guardian.eqg
		//{name: "guardian.eqg"},
		// .zon|4|arthicrex_te.zon|arthicrex.eqg
		{name: "arthicrex.eqg"},
		// .zon|4|ascent.zon|direwind.eqg
		{name: "direwind.eqg"},
		// .zon|4|atiiki.zon|atiiki.eqg
		{name: "atiiki.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".zon" {
					continue
				}
				zone := common.NewZone("")

				err = Decode(zone, bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestDecodeV3(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//{name: "fallen.zon"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open %s: %s", tt.name, err.Error())
			}
			zone := common.NewZone("")

			err = Decode(zone, bytes.NewReader(data))
			if err != nil {
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, tt.name), data, 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, tt.name))
				t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
			}
		})
	}
}
