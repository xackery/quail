package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Write writes wld.Fragments to a .wld writer. Use quail.WldMarshal to convert a Wld to wld.Fragments
func (wld *Wld) Write(w io.Writer) error {
	var err error
	if wld.Fragments == nil {
		wld.Fragments = make(map[int]FragmentReadWriter)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Bytes([]byte{0x02, 0x3D, 0x50, 0x54}) // header
	enc.Uint32(wld.Version)
	enc.Uint32(uint32(len(wld.Fragments)))
	enc.Uint32(0) //unk1
	enc.Uint32(0) //unk2

	fragBuf := bytes.NewBuffer(nil)
	for i, frag := range wld.Fragments {
		err := frag.Write(w)
		if err != nil {
			return fmt.Errorf("write fragment id %d 0x%x (%s): %w", i, frag.FragCode(), FragName(frag.FragCode()), err)
		}
	}

	nameData := NameData()
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(0) //unk3
	enc.Bytes(nameData)
	enc.Bytes(fragBuf.Bytes())

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
