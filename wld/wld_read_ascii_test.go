package wld

import "testing"

func TestReadDMSpriteDef2(t *testing.T) {
	wld := &Wld{}
	err := wld.ReadAscii("testdata/all.spk")
	if err != nil {
		t.Fatalf("failed to read: %s", err.Error())
	}

}
