package main

// #include <stdlib.h>
import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/xackery/quail/quail"
)

var (
	loadedQuail *quail.Quail
)

//export QuailLoadPfs
func QuailLoadPfs(cpath *C.char) *C.char {
	path := C.GoString(cpath)

	loadedQuail = quail.New()
	err := loadedQuail.PfsRead(path)
	if err != nil {
		return C.CString(fmt.Sprintf("Error: %s", err.Error()))
	}

	return C.CString("Success")
}

//export QuailWceFromPfs
func QuailWceFromPfs(cpath *C.char) *C.char {
	path := C.GoString(cpath)

	q := quail.New()
	err := q.PfsRead(path)
	if err != nil {
		return C.CString(fmt.Sprintf("Error pfs read: %s", err.Error()))
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(q.Wld)
	if err != nil {
		return C.CString(fmt.Sprintf("Error json: %s", err.Error()))
	}

	q.Close()

	return C.CString(buf.String())
}

//export QuailGetWldSummary
func QuailGetWldSummary() *C.char {
	if loadedQuail == nil || loadedQuail.Wld == nil {
		return C.CString("No WLD loaded")
	}

	summary := fmt.Sprintf("Models: %d, Terrains: %d, Animations: %d",
		len(loadedQuail.Wld.ModDefs),
		len(loadedQuail.Wld.TerDefs),
		len(loadedQuail.Wld.AniDefs))

	return C.CString(summary)
}

//export QuailFreeString
func QuailFreeString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func main() {
	// Required for Go to build a proper DLL
}
