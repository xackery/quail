package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/xackery/quail/quail"
)

func TestWceFromPfs(t *testing.T) {
	path := "/src/eq/rof2/it12043.eqg"

	q := quail.New()

	err := q.PfsRead(path)
	if err != nil {
		t.Fatalf("Error pfs read: %s", err.Error())
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(q.Wld)
	if err != nil {
		t.Fatalf("Error json: %s", err.Error())
	}

	q.Close()

	fmt.Println(buf.String())

}
