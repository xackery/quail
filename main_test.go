package main

import "testing"

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
}
