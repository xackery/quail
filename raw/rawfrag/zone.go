package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// WldFragZone is Zone in libeq, Region Flag in openzone, ZONE in wld, BspRegionType in lantern
type WldFragZone struct {
	nameRef  int32
	Flags    uint32
	Regions  []uint32
	UserData string
}

func (e *WldFragZone) FragCode() int {
	return FragCodeZone
}

func (e *WldFragZone) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	userData := helper.WriteStringHash(e.UserData)

	paddingSize := (4 - (len(userData) % 4)) % 4

	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Regions)))
	for _, region := range e.Regions {
		enc.Uint32(region)
	}
	enc.Uint32(uint32(len(userData)))
	if len(e.UserData) > 0 {
		enc.Bytes(userData)
		enc.Bytes(make([]byte, paddingSize))
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragZone) Read(r io.ReadSeeker, isNewWorld bool) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	regionCount := dec.Uint32()
	e.Regions = make([]uint32, 0)
	for i := uint32(0); i < regionCount; i++ {
		region := dec.Uint32()
		e.Regions = append(e.Regions, region)
	}
	userDataSize := dec.Uint32()
	if userDataSize > 0 {
		e.UserData = helper.ReadStringHash([]byte(dec.StringFixed(int(userDataSize))))
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	return nil
}

func (e *WldFragZone) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragZone) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
