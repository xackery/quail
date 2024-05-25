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
	NameRef  int32    `yaml:"name_ref"`
	Flags    uint32   `yaml:"flags"`
	Regions  []uint32 `yaml:"regions"`
	UserData string   `yaml:"user_data"`
}

func (e *WldFragZone) FragCode() int {
	return FragCodeZone
}

func (e *WldFragZone) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	userData := helper.WriteStringHash(e.UserData)

	paddingSize := (4 - (len(userData) % 4)) % 4

	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Regions)))
	for _, region := range e.Regions {
		enc.Uint32(region)
	}
	if len(e.UserData) > 0 {
		enc.Uint32(uint32(len(userData)))
		enc.Bytes(userData)
		enc.Bytes(make([]byte, paddingSize))
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragZone) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
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
