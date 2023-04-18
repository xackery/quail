package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Encode writes a wld file
func (e *WLD) Encode(w io.Writer) error {
	if e == nil {
		return fmt.Errorf("wld nil")
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)

	// TODO: build fragments
	fragmentCount := 0
	hashRaw := []byte{}

	enc.Uint32(0x54503D02) // wld header
	enc.Uint32(0x00015500) // old world identifier
	enc.Uint32(uint32(fragmentCount))
	enc.Uint32(0x00000000)           // unk1
	enc.Uint32(0x00000000)           // unk2
	enc.Uint32(uint32(len(hashRaw))) // hashSize

	//enc.Uint32(uint32(e.meshManager.BspRegionCount(e.name))) // TODO: Needed?
	enc.Uint32(0x680D4)    // after bsp region offset
	enc.Uint32(0x00000000) //TODO: hash size
	enc.Uint32(0x00000000) //TODO: after hash size offset

	if enc.Error() != nil {
		return fmt.Errorf("encode: %w", enc.Error())
	}
	return nil
}
