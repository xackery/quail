package ani

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	//path := "test/bl2h_ba_1_ala.ani"
	path := "test/obj_gears_default.ani"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.New: %s", err)
	}
	defer d.Save("test/out.png")
	e := &ANI{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

}
