package wld_test

import "testing"

var tests = []struct {
	baseName string
	wldName  string
}{
	//{baseName: "load2"}, // OK
	//{baseName: "load2", wldName: "objects.wld"}, // OK
	//{baseName: "load2", wldName: "lights.wld"}, // OK
	//{baseName: "qeynos"},  // OK
	//{baseName: "qeynos", wldName: "objects.wld"}, // OK
	//{baseName: "qeynos", wldName: "lights.wld"},  // OK
	//{baseName: "globalogm_chr"}, // OK
	//{baseName: "qeynos_chr"}, // Needs work
	{baseName: "beetle_chr"},
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
