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
					bone.translation.X, bone.translation.Y, bone.translation.Z,
					bone.rotation.X, bone.rotation.Y, bone.rotation.Z, bone.rotation.W,
					bone.scale.X, bone.scale.Y, bone.scale.Z)
			}

		}
	}
}
