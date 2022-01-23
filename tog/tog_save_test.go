package tog

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	e := &TOG{
		objects: []*Object{
			{
				Name: "test",
			},
			{
				Name: "test2",
			},
		},
	}
	buf := bytes.NewBuffer(nil)
	err := e.Save(buf)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Println(buf.String())
}
