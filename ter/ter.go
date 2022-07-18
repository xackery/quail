// Ter package
package ter

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/common"
)

// TER is a terrain file struct
type TER struct {
	name      string
	materials []*common.Material
	vertices  []*common.Vertex
	faces     []*common.Face
	files     []common.Filer
	eqg       common.Archiver
}

func New(name string, archive common.Archiver) (*TER, error) {
	t := &TER{
		name: name,
		eqg:  archive,
	}
	return t, nil
}

func (e *TER) Name() string {
	return e.name
}

func (e *TER) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Save(w)
	if err != nil {
		fmt.Println("failed to save terrain data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}

func (e *TER) SetName(value string) {
	e.name = value
}
