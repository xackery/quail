package wld

import (
	"testing"
)

func TestWLD_objectLocationRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		false, //single run stop
		[]string{
			"gfaydark.s3d",
		},
		21,                   //fragCode
		-1,                   //fragIndex
		e.objectLocationRead) //callback
}
