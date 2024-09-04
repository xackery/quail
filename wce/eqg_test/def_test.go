package wce_test

import (
	"fmt"
	"testing"
	"time"
)

var tests = []struct {
	baseName string
	fileName string
}{
	//{baseName: "ork"},
	//{baseName: "arena"},
	{baseName: "commonlands"},
}

func TestStep1(t *testing.T) {
	TestWceReadWrite(t)
}

func TestStep2(t *testing.T) {
	TestWceDoubleReadWrite(t)
}

func TestStepAll(t *testing.T) {
	TestStep1(t)
	TestStep2(t)
}

func TestBit(t *testing.T) {
	fmt.Printf("0x%x\n", 1<<15)
}

func TestPastGoodTests(t *testing.T) {

	start := time.Now()

	tests = []struct {
		baseName string
		fileName string
	}{
		//{baseName: "globaldam_chr"}, // PASS
	}

	TestStep1(t)

	fmt.Printf("Took %0.2f seconds for %d total tests\n", time.Since(start).Seconds(), len(tests))
}
