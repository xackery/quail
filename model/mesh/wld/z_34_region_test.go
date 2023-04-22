package wld

import (
	"testing"
)

func TestWLD_regionRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"abysmal.s3d",
		},
		34,           //fragCode
		-1,           //fragIndex
		e.regionRead) //callback
}
