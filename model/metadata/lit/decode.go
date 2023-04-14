package lit

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/model/geo"
)

func (e *LIT) Decode(r io.ReadSeeker) error {
	var err error

	fileCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &fileCount)
	if err != nil {
		return fmt.Errorf("read fileCount: %w", err)
	}

	lightCount := fileCount
	fmt.Println("lightCount", lightCount)

	for i := 0; i < int(lightCount); i++ {
		color := &geo.RGBA{}
		err = binary.Read(r, binary.LittleEndian, color)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read light %d: %w", i, err)
		}
		e.lights = append(e.lights, color)
	}

	return nil
}
