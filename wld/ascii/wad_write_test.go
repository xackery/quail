package ascii

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func TestWadMockWrite(t *testing.T) {
	wad := &Wad{
		bitmaps: []*BMInfo{
			{
				Tag: "test",
				Textures: []string{
					"test1",
					"test2",
				},
			},
		},
		sprites: []*Sprite{
			{
				Tag:          "test",
				CurrentFrame: 1,
				Sleep:        1,
				Instances: []*SpriteInstance{
					{
						Tag:   "test",
						Flags: 1,
					},
				},
			},
		},
		dmsprites: []*DMSpriteInfo{
			{
				Tag: "test",
				Tracks: []*DMSpriteInfoTrack{
					{
						Tag: "test",
					},
				},
			},
		},
	}
	buf := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buf)
	err := wad.WADWrite(w)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	w.Flush()
	fmt.Println(buf.String())
}

func TestWadWldWrite(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		path      string
		file      string
		fragIndex int
		want      raw.FragmentReadWriter
		wantErr   bool
	}{
		//{"btp_chr.s3d", "btp_chr.wld", 0, nil, false},
		//{"bac_chr.s3d", "bac_chr.wld", 0, nil, false},
		//{"avi_chr.s3d", "avi_chr.wld", 0, nil, false},
	}
	if !common.IsTestExtensive() {
		tests = []struct {
			path      string
			file      string
			fragIndex int
			want      raw.FragmentReadWriter
			wantErr   bool
		}{
			//"globalfroglok_chr.s3d", "globalfroglok_chr.wld", 0, nil, false},
			{"gequip4.s3d", "gequip4.wld", 0, nil, false},
		}
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.path))
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.file, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(tt.file)
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.file, err.Error())
			}
			wld := &raw.Wld{}
			err = wld.Read(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read wld %s: %s", tt.file, err.Error())
			}

			wad := &Wad{}
			err = wad.WLDRead(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read wld %s: %s", tt.file, err.Error())
			}

			err = wad.WADWrite(bufio.NewWriter(os.Stdout))
			if err != nil {
				t.Fatalf("failed to write wad %s: %s", tt.file, err.Error())
			}

		})
	}
}
