package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw/rawfrag"
	"github.com/xackery/quail/tag"
)

// Write writes wld.Fragments to a .wld writer. Use quail.WldMarshal to convert a Wld to wld.Fragments
func (wld *Wld) Write(w io.Writer) error {
	var err error
	if wld.Fragments == nil {
		wld.Fragments = []model.FragmentReadWriter{&rawfrag.WldFragDefault{}}
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	tag.NewWithCoder(enc)
	enc.Bytes([]byte{0x02, 0x3D, 0x50, 0x54}) // header
	enc.Uint32(wld.Version)
	tag.Mark("red", "header")

	enc.Uint32(uint32(len(wld.Fragments)))
	tag.Mark("blue", "fragcount")

	enc.Uint32(wld.BspRegionCount) //bspRegionCount
	tag.Mark("green", "bspRegionCount")
	enc.Uint32(wld.Unk2) //unk2
	tag.Mark("lime", "unk2")

	totalFragSize := 0
	totalFragBuf := bytes.NewBuffer(nil)
	for i := range wld.Fragments {
		frag := wld.Fragments[i]
		fragBuf := bytes.NewBuffer(nil)
		chunkBuf := bytes.NewBuffer(nil)
		chunkEnc := encdec.NewEncoder(chunkBuf, binary.LittleEndian)

		err := frag.Write(fragBuf)
		if err != nil {
			return fmt.Errorf("write fragment id %d 0x%x (%s): %w", i, frag.FragCode(), FragName(frag.FragCode()), err)
		}
		chunkEnc.Uint32(uint32(fragBuf.Len()))
		chunkEnc.Uint32(uint32(frag.FragCode()))
		chunkEnc.Bytes(fragBuf.Bytes())

		totalFragSize += fragBuf.Len()

		totalFragBuf.Write(chunkBuf.Bytes())
	}

	nameData := NameData()
	enc.Uint32(uint32(len(nameData))) //hashSize
	tag.Mark("green", "hashsize")

	enc.Uint32(wld.Unk3) //unk3
	tag.Mark("lime", "unk3")
	hashRaw := helper.WriteStringHash(string(nameData))
	enc.Bytes(hashRaw) //hashRaw
	tag.Mark("red", "namehash")
	enc.Bytes(totalFragBuf.Bytes())
	tag.Mark("blue", "fragments")
	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
