package lit

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/pfs/eqg"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	eqgPath := "test/eq/steamfontmts.eqg"

	archive, err := eqg.NewFile(eqgPath)
	if err != nil {
		t.Fatalf("eqg new: %s", err)
	}

	inFile := "steamfontcavesc_obj_sc_csupporta31.lit"
	dump.New(inFile)
	defer dump.WriteFileClose("test/eq/steamfontmts_eqg_steamfontcavesc_obj_sc_csupporta31.lit.png")

	e, err := NewFile("light", archive, inFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	if e == nil {
		t.Fatalf("e is nil")
	}
}
