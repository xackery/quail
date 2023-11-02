package raw

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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
		{name: "guardian.eqg"}, // FIXME: nameIndex out of range
		// TODO: zone4 support
		// .zon|4|arthicrex_te.zon|arthicrex.eqg
		//{name: "arthicrex.eqg"},
		// .zon|4|ascent.zon|direwind.eqg
		//{name: "direwind.eqg"},
		// .zon|4|atiiki.zon|atiiki.eqg
		//{name: "atiiki.eqg"},
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

func TestZonReadV3(t *testing.T) {
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
		//{name: "fallen.zon"}, // FIXME: nameIndex out of range
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open %s: %s", tt.name, err.Error())
			}
			zon := &Zon{}
			err = zon.Read(bytes.NewReader(data))
			if err != nil {
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, tt.name), data, 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, tt.name))
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
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
		//{name: "anguish.eqg"},
		// .zon|1|bazaar.zon|bazaar.eqg
		//{name: "bazaar.eqg"}, // FIXME: modelNameOffset not found
		// .zon|1|bloodfields.zon|bloodfields.eqg
		//{name: "bloodfields.eqg"},
		// .zon|1|broodlands.zon|broodlands.eqg
		//{name: "broodlands.eqg"},
		// .zon|1|catacomba.zon|dranikcatacombsa.eqg
		//{name: "dranikcatacombsa.eqg"},
		// .zon|1|wallofslaughter.zon|wallofslaughter.eqg
		//{name: "wallofslaughter.eqg"},
		// .zon|2|arginhiz.zon|arginhiz.eqg
		{name: "arginhiz.eqg"}, // FIXME: nameIndex out of range
		// .zon|2|guardian.zon|guardian.eqg
		//{name: "guardian.eqg"},
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

				//srcData := file.Data()
				//dstData := buf.Bytes()
				/*for i := 0; i < len(srcData); i++ {
					if len(dstData) <= i {
						min := 0
						max := len(srcData)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
						max = len(dstData)
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))

						t.Fatalf("%s src eof at offset %d (dst is too large by %d bytes)", tt.name, i, len(dstData)-len(srcData))
					}
					if len(dstData) <= i {
						t.Fatalf("%s dst eof at offset %d (dst is too small by %d bytes)", tt.name, i, len(srcData)-len(dstData))
					}
					if srcData[i] == dstData[i] {
						continue
					}

					fmt.Printf("%s mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", tt.name, i, srcData[i], dstData[i], dstData[i])
					max := i + 16
					if max > len(srcData) {
						max = len(srcData)
					}

					min := i - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
					if max > len(dstData) {
						max = len(dstData)
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))
					//os.WriteFile(fmt.Sprintf("%s/_src_%s", dirTest, file.Name()), file.Data(), 0644)
					//os.WriteFile(fmt.Sprintf("%s/_dst_%s", dirTest, file.Name()), buf.Bytes(), 0644)
					t.Fatalf("%s encode: data mismatch", tt.name)
				}*/
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
		args    args
		wantErr bool
	}{
		// .zon|4|arthicrex_te.zon|arthicrex.eqg
		//{name: "arthicrex.eqg"}, // FIXME: v4 encode is broken
		// .zon|4|ascent.zon|direwind.eqg
		//{name: "direwind.eqg"},
		// .zon|4|atiiki.zon|atiiki.eqg
		//{name: "atiiki.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".dat" {
					continue
				}
				zon := &Zon{}
				err = zon.ReadV4(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = zon.WriteV4(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				srcData := file.Data()
				dstData := buf.Bytes()
				for i := 0; i < len(srcData); i++ {
					if len(dstData) <= i {
						min := 0
						max := len(srcData)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
						max = len(dstData)
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))

						t.Fatalf("%s src eof at offset %d (dst is too large by %d bytes)", tt.name, i, len(dstData)-len(srcData))
					}
					if len(dstData) <= i {
						t.Fatalf("%s dst eof at offset %d (dst is too small by %d bytes)", tt.name, i, len(srcData)-len(dstData))
					}
					if srcData[i] == dstData[i] {
						continue
					}

					fmt.Printf("%s mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", tt.name, i, srcData[i], dstData[i], dstData[i])
					max := i + 16
					if max > len(srcData) {
						max = len(srcData)
					}

					min := i - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
					if max > len(dstData) {
						max = len(dstData)
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))
					t.Fatalf("%s encode: data mismatch", tt.name)
				}
			}
		})
	}
}
