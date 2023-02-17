// lit is an EverQuest file for light data on a zone
package lit

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/common"
)

// LIT are light sources
type LIT struct {
	name    string
	archive common.ArchiveReader
	lights  []float32
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, archive common.ArchiveReader) (*LIT, error) {
	t := &LIT{
		name: name,
	}
	return t, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReader, file string) (*LIT, error) {
	e := &LIT{
		name:    name,
		archive: archive,
	}
	data, err := archive.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}

func (e *LIT) Name() string {
	return e.name
}

func (e *LIT) Data() []byte {
	w := bytes.NewBuffer(nil)

	err := e.Encode(w)
	if err != nil {
		fmt.Println("failed to encode litrain data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}
