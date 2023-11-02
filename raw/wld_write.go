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
	if wld.Fragments == nil {
		wld.Fragments = make(map[int]FragmentReadWriter)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Bytes([]byte{0x02, 0x3D, 0x50, 0x54})
	enc.Uint32(wld.Version)
	enc.Uint32(uint32(len(wld.Fragments)))
	enc.Uint32(0) //unk1
	enc.Uint32(0) //unk2

	fragBuf := bytes.NewBuffer(nil)
	nameData, err := writeFragments(wld, fragBuf)
	if err != nil {
		return fmt.Errorf("write fragments: %w", err)
	}
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(0) //unk3
	enc.Bytes(nameData)
	enc.Bytes(fragBuf.Bytes())

	if enc.Error() != nil {
		return fmt.Errorf("encode: %w", enc.Error())
	}
	return nil
}

// writeFragments converts fragment structs to bytes
func writeFragments(wld *Wld, w io.Writer) ([]byte, error) {
	nameBuf := bytes.NewBuffer(nil)
	for i, frag := range wld.Fragments {
		err := frag.Write(w)
		if err != nil {
			return nil, fmt.Errorf("fragment id %d 0x%x (%s): %w", i, frag.FragCode(), FragName(frag.FragCode()), err)
		}
		// Name builder?
		/*
			pos, ok := names[name]
			if !ok {
				tmpNames = append(tmpNames, model.Header.Name)
				names[name] = int32(nameBuf.Len())
				nameBuf.Write([]byte(name))
				nameBuf.Write([]byte{0})
				pos = names[name]
			}
		*/
	}
	return nameBuf.Bytes(), nil
}
