package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
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
	if !wld.IsNewWorld {
		enc.Uint32(0x00015500)
	} else {
		enc.Uint32(0x1000C800)
	}
	tag.Mark("red", "header")

	if len(wld.Fragments) == 0 {
		return fmt.Errorf("no fragments found")
	}

	_, hasDefaultFrag := wld.Fragments[0].(*rawfrag.WldFragDefault)
	if !hasDefaultFrag {
		return fmt.Errorf("first fragment must be WldFragDefault")
	}

	enc.Uint32(uint32(len(wld.Fragments) - 1))
	tag.Mark("blue", "fragcount")

	maxFragSize := 0
	totalRegionCount := 0
	totalFragSize := 0
	totalFragBuf := bytes.NewBuffer(nil)
	for i := range wld.Fragments {
		frag := wld.Fragments[i]
		if frag.FragCode() == rawfrag.FragCodeDefault {
			if i != 0 {
				return fmt.Errorf("default fragment must be first fragment")
			}
			continue
		}
		fragBuf := bytes.NewBuffer(nil)
		chunkBuf := bytes.NewBuffer(nil)
		chunkEnc := encdec.NewEncoder(chunkBuf, binary.LittleEndian)

		err := frag.Write(fragBuf, wld.IsNewWorld)
		if err != nil {
			return fmt.Errorf("write fragment id %d 0x%x (%s): %w", i, frag.FragCode(), FragName(frag.FragCode()), err)
		}
		chunkEnc.Uint32(uint32(fragBuf.Len()))
		chunkEnc.Uint32(uint32(frag.FragCode()))
		chunkEnc.Bytes(fragBuf.Bytes())

		totalFragSize += fragBuf.Len()
		if fragBuf.Len() > maxFragSize {
			maxFragSize = fragBuf.Len() + 8
		}

		_, ok := frag.(*rawfrag.WldFragRegion)
		if ok {
			totalRegionCount++
		}

		totalFragBuf.Write(chunkBuf.Bytes())
	}

	enc.Uint32(uint32(totalRegionCount)) //aka bspRegionCount
	tag.Mark("green", "totalRegionCount")

	enc.Uint32(uint32(maxFragSize))
	tag.Mark("lime", "max_object_bytes")

	nameData := NameData()

	// pad namedata with 0's so it's divisible by 4
	for len(nameData)%4 != 0 {
		nameData = append(nameData, 0)
	}

	enc.Uint32(uint32(len(nameData))) //hashSize
	tag.Mark("green", "hashsize")

	if len(names) < 1 {
		return fmt.Errorf("no names found")
	}
	enc.Uint32(uint32(len(names) - 1)) // there's a 0x00 string at start but it's not counted
	tag.Mark("green", "string count")

	enc.Bytes(nameData)
	tag.Mark("red", "namehash")
	enc.Bytes(totalFragBuf.Bytes())
	tag.Mark("blue", "fragments")
	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
