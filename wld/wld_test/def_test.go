package wld_test

import "testing"

var tests = []struct {
	baseName string
	wldName  string
}{
	//{baseName: "load2"}, // OK
	//{baseName: "beetle_chr"}, // OK
	//{baseName: "load2", wldName: "objects.wld"}, // OK
	//{baseName: "load2", wldName: "lights.wld"}, // EMPTY
	{baseName: "qeynos_chr"}, // Needs work
	//{baseName: "qeynos"},
	//{baseName: "qeynos", wldName: "objects.wld"},
	//{baseName: "qeynos", wldName: "lights.wld"},
	//{baseName: "globalogm_chr"},
}

func TestStep1(t *testing.T) {
	TestRawFragReadWrite(t)
}

func TestStep2(t *testing.T) {
	TestRawBinReadWrite(t)
}

func TestStep3(t *testing.T) {
	TestRawWldReadWrite(t)
}

func TestStep4(t *testing.T) {
	TestWceReadWrite(t)
}

func TestStepAll(t *testing.T) {
	TestStep1(t)
	TestStep2(t)
	TestStep3(t)
	TestStep4(t)
}
