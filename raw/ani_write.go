package raw

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
)

func (ani *Ani) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.StringFixed("EQGA", 4)
	enc.Uint32(ani.Version)
	enc.Uint32(uint32(len(ani.Bones)))

	return nil
}
