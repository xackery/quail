package tog

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/xackery/quail/eqg"
)

func TestSave(t *testing.T) {
	archive, err := eqg.New("test")
	if err != nil {
		t.Fatalf("eqg new: %s", err)
	}
	e, err := New("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	e.objects = []*Object{
		{
			Name: "test",
		},
		{
			Name: "test2",
		},
	}

	buf := bytes.NewBuffer(nil)
	err = e.Save(buf)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Println(buf.String())
}
