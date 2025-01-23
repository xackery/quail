package quail

import (
	"os"
	"testing"
)

func TestQuailRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		pfsName string
		wantErr bool
	}{
		{"ala.eqg", false},
	}

	for _, tt := range tests {
		t.Run(tt.pfsName, func(t *testing.T) {
			q := &Quail{}
			err := q.PfsRead(eqPath + "/" + tt.pfsName)
			if err != nil {
				t.Fatalf("failed to read pfs %s: %s", tt.pfsName, err.Error())
			}

		})
	}
}
