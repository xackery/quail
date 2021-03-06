package ani

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestDecodeSingleTest(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	isDump := true
	tests := []struct {
		category string
	}{
		{category: "steamfontmts"},
	}
	for _, tt := range tests {

		fmt.Println("loading", tt.category)
		eqgFile := fmt.Sprintf("test/eq/%s.eqg", tt.category)

		ra, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer ra.Close()
		archive, err := eqg.NewFile(eqgFile)
		if err != nil {
			t.Fatalf("eqg new: %s", err)
		}

		files := archive.Files()
		for _, aniEntry := range files {
			if filepath.Ext(aniEntry.Name()) != ".ani" {
				continue
			}

			if isDump {
				dump.New(aniEntry.Name())
			}
			defer dump.Close()
			r := bytes.NewReader(aniEntry.Data())

			e, err := New(aniEntry.Name())
			if err != nil {
				t.Fatalf("new: %s", err)
			}

			err = e.Decode(r)
			if err != nil {
				t.Fatalf("decode %s: %s", aniEntry.Name(), err)
			}
			fmt.Println(e.name)
			for _, bone := range e.bones {
				fmt.Printf("delay %d translation %0.f %0.f %0.f rotation %0.f %0.f %0.f %0.f scale %0.f %0.f %0.f\n",
					bone.delay,
					bone.translation[0], bone.translation[1], bone.translation[2],
					bone.rotation[0], bone.rotation[1], bone.rotation[2], bone.rotation[3],
					bone.scale[0], bone.scale[1], bone.scale[2])
			}

		}
	}
}
