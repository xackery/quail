package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/tag"
)

// Write writes wld.Fragments to a .wld writer. Use quail.WldMarshal to convert a Wld to wld.Fragments
func (wld *Wld) Write(w io.Writer) error {
	var err error
	if wld.Fragments == nil {
		wld.Fragments = make(map[int]FragmentReadWriter)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	tag.NewWithCoder(enc)
	enc.Bytes([]byte{0x02, 0x3D, 0x50, 0x54}) // header
	enc.Uint32(wld.Version)
	tag.Mark("red", "header")

	enc.Uint32(uint32(len(wld.Fragments)))
	tag.Mark("blue", "fragcount")

	enc.Uint32(0) //bspRegionCount
	enc.Uint32(0) //unk2

	totalFragSize := 0
	fragBuf := bytes.NewBuffer(nil)
	for i, frag := range wld.Fragments {
		buf := bytes.NewBuffer(nil)
		err := frag.Write(buf)
		if err != nil {
			return fmt.Errorf("write fragment id %d 0x%x (%s): %w", i, frag.FragCode(), FragName(frag.FragCode()), err)
		}
		enc.Uint32(uint32(buf.Len()))       //fragSize
		enc.Uint32(uint32(frag.FragCode())) //fragCode
		totalFragSize += buf.Len()
		fragBuf.Write(buf.Bytes())
	}

	nameData := NameData()
	enc.Uint32(uint32(len(nameData))) //hashSize
	tag.Mark("green", "hashsize")

	enc.Uint32(0) //unk3
	hashRaw := helper.WriteStringHash(string(nameData))
	enc.Bytes(hashRaw) //hashRaw
	tag.Mark("green", "namedata")
	enc.Bytes(fragBuf.Bytes())
	tag.Mark("blue", "fragments")
	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
