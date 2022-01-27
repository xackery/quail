package tog

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	e.objects = []*Object{
		{
			name: "test",
		},
		{
			name: "test2",
		},
	}

	buf := bytes.NewBuffer(nil)
	err = e.Save(buf)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Println(buf.String())
}
