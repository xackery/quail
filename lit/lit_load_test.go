package lit

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/_steamfontmts.eqg/steamfontcavesc_obj_sc_csupporta31.lit"
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump: %s", err)
	}
	defer d.Save("../eq/tmp/out.png")

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}

	e, err := New("light")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	defer f.Close()

}
