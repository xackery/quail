package wce_test

import (
	"fmt"
	"testing"
	"time"
)

var tests = []struct {
	baseName string
	wldName  string
}{
	//{baseName: "akanon"},
	//{baseName: "akanon", wldName: "objects.wld"},
	//{baseName: "akanon", wldName: "lights.wld"},
	//{baseName: "qeynos_chr"},
	{baseName: "freporte_chr"},
	//{baseName: "freporte_chr-pre-partial-fix", wldName: "freporte_chr.wld"},
	//{baseName: "overthere_chr"},
	//{baseName: "gfaydark", wldName: "objects.wld"},
	//{baseName: "gfaydark_obj"},
	//{baseName: "globalelf_chr"}, // BROKE
	//{baseName: "global6_chr"}, // FIX ME
	//{baseName: "global_chr"},
	//{baseName: "twilight", wldName: "objects.wld"},
	//{baseName: "twilight", wldName: "lights.wld"},
	//{baseName: "frontiermtns_chr"},
	//{baseName: "crushbone"},
	//{baseName: "gukbottom"},
	//{baseName: "qeynos"},
	//{baseName: "qeynos_chr"},
	//{baseName: "qeynos_obj"},
	//{baseName: "global_chr"},  // boat_actordef unknown sprite type
	//{baseName: "global2_chr"}, // PASS
	//{baseName: "global3_chr"}, // track O02DWF_TRACK model too short
	//{baseName: "global4_chr"}, // PASS
	//{baseName: "global5_chr"}, // PASS
	//{baseName: "global6_chr"}, // PASS
	//{baseName: "global7_chr"}, // PASS
	//{baseName: "globalelf_chr"}, // PASS
	//{baseName: "globaldaf_chr"}, // PASS
	//{baseName: "globaldam_chr"}, // PASS
	//{baseName: "arena"}, // OK
	//{baseName: "abysmal"},
	//{baseName: "global_chr"}, // original human etc
	//{baseName: "load2"}, // OK
	//{baseName: "beetle_chr"}, // fails since openzone alignemnts
	//{baseName: "load2", wldName: "objects.wld"}, // OK
	//{baseName: "crushbone", wldName: "objects.wld"}, // OK
	//{baseName: "crushbone", wldName: "lights.wld"},
	//{baseName: "qeynos_chr"}, // Needs work
	//{baseName: "crushbone_chr"},
	//{baseName: "freporte_chr"},
	//{baseName: "chequip"},
	//{baseName: "gequip"}, // TODO: fix numeric prefixed tags
	//{baseName: "gequip2"},
	//{baseName: "global2_chr"},
	//{baseName: "globaldaf_chr"},
	//{baseName: "globalhum_chr"},
	//{baseName: "freporte"},
	//{baseName: "neriakc"}, // OK
	//{baseName: "qeynos"}, // OK
	//{baseName: "load2", wldName: "lights.wld"}, // EMPTY
	//{baseName: "qeynos", wldName: "objects.wld"},
	//{baseName: "qeynos", wldName: "lights.wld"},
	//{baseName: "globalogm_chr"},

	//{baseName: "acrylia_obj"}, //acrylia_obj failed to write acrylia_obj: hierarchicalsprite ACTORCH301_HS_DEF: collision volume not found: I_L301_SPB
	// {baseName: "ael_chr"}, //failed to write ael_chr: actordef AEL_ACTORDEF: sprite AEL_HS_DEF to raw: collision volume not found: I_AELCLOUD01_SPB
	//{baseName: "gequip.takp", wldName: "gequip.wld"},
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

func TestStep5(t *testing.T) {
	TestWceDoubleReadWrite(t)
}

func TestStepAll(t *testing.T) {
	TestStep1(t)
	TestStep2(t)
	TestStep3(t)
	TestStep4(t)
	TestStep5(t)
}

func TestBit(t *testing.T) {
	fmt.Printf("0x%x\n", 1<<15)
}

func TestPastGoodTests(t *testing.T) {

	start := time.Now()

	tests = []struct {
		baseName string
		wldName  string
	}{
		{baseName: "crushbone"},
		{baseName: "crushbone", wldName: "objects.wld"},
		{baseName: "crushbone", wldName: "lights.wld"},
		{baseName: "arena"},
		{baseName: "neriakc"},
		{baseName: "twilight"},
		//{baseName: "global_chr"},
		{baseName: "global2_chr"},
		//{baseName: "global3_chr"}, // track O02DWF_TRACK model too short
		//{baseName: "global4_chr"}, // BAMHE0003_SPRITE not found
		{baseName: "global5_chr"}, //  PASS
		{baseName: "global6_chr"}, // PASS
		//{baseName: "global7_chr"}, // KEFTASK11 not found
		{baseName: "globalelf_chr"}, // PASS
		{baseName: "globaldaf_chr"}, // PASS
		{baseName: "globaldam_chr"}, // PASS
	}

	TestStep4(t)

	fmt.Printf("Took %0.2f seconds for %d total tests\n", time.Since(start).Seconds(), len(tests))
}