package wld

import (
	"testing"
)

// TODO: no refs
func TestWLD_objectLocationRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"gequip.s3d",
		},
		21,                   //fragCode
		-1,                   //fragIndex
		e.objectLocationRead) //callback
}
