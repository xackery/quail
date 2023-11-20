package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestZonRead(t *testing.T) {
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
		// .zon|2|guardian.zon|guardian.eqg
		//{name: "guardian.eqg"}, // PASS
		// .zon|1|anguish.zon|anguish.eqg
		//{name: "anguish.eqg"}, // PASS
		// .zon|1|bazaar.zon|bazaar.eqg
		//{name: "bazaar.eqg"}, // PASS
		// .zon|1|bloodfields.zon|bloodfields.eqg
		//{name: "bloodfields.eqg"}, // PASS
		// .zon|1|broodlands.zon|broodlands.eqg
		//{name: "broodlands.eqg"}, // PASS
		// .zon|1|catacomba.zon|dranikcatacombsa.eqg
		//{name: "dranikcatacombsa.eqg"}, // PASS
		// .zon|1|wallofslaughter.zon|wallofslaughter.eqg
		//{name: "wallofslaughter.eqg"}, // PASS
		// .zon|2|arginhiz.zon|arginhiz.eqg
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
				zon := &Zon{}

				err = zon.Read(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestZonWrite(t *testing.T) {
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
		// {name: "anguish.eqg"}, // TODO: mismatch
		// .zon|1|bazaar.zon|bazaar.eqg
		//{name: "bazaar.eqg"}, // TODO: mismatch
		// .zon|1|bloodfields.zon|bloodfields.eqg
		//{name: "bloodfields.eqg"}, // TODO: mismatch
		// .zon|1|broodlands.zon|broodlands.eqg
		//{name: "broodlands.eqg"}, // TODO: mismatch
		// .zon|1|catacomba.zon|dranikcatacombsa.eqg
		//{name: "dranikcatacombsa.eqg"}, // TODO: mismatch
		// .zon|1|wallofslaughter.zon|wallofslaughter.eqg
		//{name: "wallofslaughter.eqg"}, // TODO: mismatch
		// .zon|2|arginhiz.zon|arginhiz.eqg
		//{name: "arginhiz.eqg"}, // TODO: mismatch
		// .zon|2|guardian.zon|guardian.eqg
		//{name: "guardian.eqg"}, // TODO: v2 zone mismatch
		// .zon|2|guardian.zon|guardian.eqg
		//{name: "fallen.zon"}, // TODO: v2 zone mismatch
	}

	var err error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fileName := tt.name
			var data []byte
			ext := filepath.Ext(tt.name)
			if ext == ".zon" {
				data, err = os.ReadFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
				if err != nil {
					t.Fatalf("failed to open %s: %s", tt.name, err.Error())
				}
			} else {
				pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
				if err != nil {
					t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
				}
				fileName := strings.ReplaceAll(tt.name, ".eqg", "") + ".zon"
				data, err = pfs.File(fileName)
				if err != nil {
					t.Fatalf("failed to open %s.zon: %s", tt.name, err.Error())
				}
			}

			zon := &Zon{}
			err = zon.Read(bytes.NewReader(data))
			os.WriteFile(fmt.Sprintf("%s/%s", dirTest, fileName), data, 0644)
			tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, fileName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = zon.Write(buf)
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}

			zon2 := &Zon{}
			err = zon2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			if len(zon.Lights) != len(zon2.Lights) {
				t.Fatalf("lights mismatch: %d != %d", len(zon.Lights), len(zon2.Lights))
			}

			if len(zon.Models) != len(zon2.Models) {
				t.Fatalf("models mismatch: %d != %d", len(zon.Models), len(zon2.Models))
			}

			if len(zon.Objects) != len(zon2.Objects) {
				t.Fatalf("objects mismatch: %d != %d", len(zon.Objects), len(zon2.Objects))
			}

			if len(zon.Regions) != len(zon2.Regions) {
				t.Fatalf("regions mismatch: %d != %d", len(zon.Regions), len(zon2.Regions))
			}

			srcData := data
			dstData := buf.Bytes()
			err = common.ByteCompareTest(srcData, dstData)
			if err != nil {
				t.Fatalf("%s byteCompare: %s", tt.name, err)
			}
		})
	}
}

func TestZonWriteV4(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	type args struct {
	}

	tests := []struct {
		name    string
		eqgPath string
		args    args
		wantErr bool
	}{
		// .zon|4|arthicrex_te.zon|arthicrex.eqg
		//{name: "arthicrex_te.zon", eqgPath: "arthicrex.eqg"}, // FIXME: v4 write is broken
		// .zon|4|ascent.zon|direwind.eqg
		//{name: "direwind.eqg"}, // PASS
		// .zon|4|atiiki.zon|atiiki.eqg
		//{name: "atiiki.eqg"}, // PASS
		// .zon|4|arthicrex_te.zon|arthicrex.eqg
		//{name: "arthicrex.eqg"}, // PASS
		// .zon|4|ascent.zon|direwind.eqg
		//{name: "direwind.eqg"}, // PASS
		// .zon|4|atiiki.zon|atiiki.eqg
		{name: "atiiki.eqg"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error
			data, err = os.ReadFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.eqgPath))
				if err != nil {
					t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
				}
				data, err = pfs.File(tt.name)
				if err != nil {
					t.Fatalf("failed to open %s.zon: %s", tt.name, err.Error())
				}
			}

			zon := &Zon{}
			err = zon.ReadV4(bytes.NewReader(data))
			os.WriteFile(fmt.Sprintf("%s/%s", dirTest, tt.name), data, 0644)
			tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, tt.name))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = zon.WriteV4(buf)
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}

			zon2 := &Zon{}
			err = zon2.ReadV4(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			if zon.V4Info.Name != zon2.V4Info.Name {
				t.Fatalf("name mismatch: %s != %s", zon.V4Info.Name, zon2.V4Info.Name)
			}

			if zon.V4Info.MinLng != zon2.V4Info.MinLng {
				t.Fatalf("minLng mismatch: %d != %d", zon.V4Info.MinLng, zon2.V4Info.MinLng)
			}

			if zon.V4Info.MaxLng != zon2.V4Info.MaxLng {
				t.Fatalf("maxLng mismatch: %d != %d", zon.V4Info.MaxLng, zon2.V4Info.MaxLng)
			}

			if zon.V4Info.MinLat != zon2.V4Info.MinLat {
				t.Fatalf("minLat mismatch: %d != %d", zon.V4Info.MinLat, zon2.V4Info.MinLat)
			}

			if zon.V4Info.MaxLat != zon2.V4Info.MaxLat {
				t.Fatalf("maxLat mismatch: %d != %d", zon.V4Info.MaxLat, zon2.V4Info.MaxLat)
			}

			if zon.V4Info.MinExtents.X != zon2.V4Info.MinExtents.X {
				t.Fatalf("minExtents.X mismatch: %f != %f", zon.V4Info.MinExtents.X, zon2.V4Info.MinExtents.X)
			}

			if zon.V4Info.MinExtents.Y != zon2.V4Info.MinExtents.Y {
				t.Fatalf("minExtents.Y mismatch: %f != %f", zon.V4Info.MinExtents.Y, zon2.V4Info.MinExtents.Y)
			}

		})
	}
}
