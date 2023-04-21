package wld

import (
	"testing"
)

func TestWLD_worldTreeRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"acrylia.s3d",
		},
		33,              //fragCode
		-1,              //fragIndex
		e.worldTreeRead) //callback
}
