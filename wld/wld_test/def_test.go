package wld_test

import (
	"fmt"
	"testing"
)

var tests = []struct {
	baseName string
	wldName  string
}{
	//{baseName: "global_chr"}, // original human etc
	{baseName: "abysmal"},

	//{baseName: "load2"}, // OK
	//{baseName: "beetle_chr"}, // fails since openzone alignemnts
	//{baseName: "load2", wldName: "objects.wld"}, // OK
	//{baseName: "crushbone"}, // OK
	//{baseName: "crushbone", wldName: "objects.wld"}, // OK
	//{baseName: "crushbone", wldName: "lights.wld"},
	//{baseName: "qeynos_chr"}, // Needs work
	//{baseName: "crushbone_chr"}, // OK
	//{baseName: "freporte_chr"},
	//{baseName: "chequip"},
	//{baseName: "gequip"}, // TODO: fix numeric prefixed tags
	//{baseName: "gequip2"},
	//{baseName: "global2_chr"},
	//{baseName: "globaldaf_chr"},
	//{baseName: "globalhum_chr"},
	//{baseName: "freporte"},
	//{baseName: "qeynos"}, // OK
	//{baseName: "load2", wldName: "lights.wld"}, // EMPTY
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

func TestBit(t *testing.T) {
	fmt.Printf("0x%x\n", 1<<15)
}
