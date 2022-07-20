package eqg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// Save will write a EQG to writer
func (e *EQG) Save(w io.WriteSeeker) error {
	var err error

	type dirEntry struct {
		crc    uint32
		offset uint32
		size   uint32
	}
	dirEntries := []*dirEntry{}

	pos := int64(0)

	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return fmt.Errorf("write header prefix primer: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, [4]byte{'P', 'F', 'S', ' '})
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(0x00020000))
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}

	fileBuffer := bytes.NewBuffer(nil)
	err = binary.Write(fileBuffer, binary.LittleEndian, uint32(len(e.files)))
	if err != nil {
		return fmt.Errorf("write file count: %w", err)
	}

	e.files = common.FilerByCRC(e.files)
	for _, file := range e.files {
		pos, err = w.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("%s seek: %w", file.Name(), err)
		}

		dirEntries = append(dirEntries, &dirEntry{
			crc:    helper.FilenameCRC32(file.Name()),
			size:   uint32(len(file.Data())),
			offset: uint32(pos),
		})

		//write compressed data
		cData, err := helper.Deflate(file.Data())
		if err != nil {
			return fmt.Errorf("deflate %s: %w", file.Name(), err)
		}

		err = binary.Write(w, binary.LittleEndian, cData)
		if err != nil {
			return fmt.Errorf("%s write data: %w", file.Name(), err)
		}

		// prep filebuffer
		err = binary.Write(fileBuffer, binary.LittleEndian, uint32(len(file.Name())+1))
		if err != nil {
			return fmt.Errorf("%s write name length: %w", file.Name(), err)
		}

		err = helper.WriteString(fileBuffer, file.Name())
		//_, err = fileBuffer.Write([]byte(file.Name()))
		if err != nil {
			return fmt.Errorf("write name %s: %w", file.Name(), err)
		}
		//fmt.Println(len(file.Name)+1, hex.Dump([]byte(file.Name)))

		/*err = binary.Write(fileBuffer, binary.LittleEndian, uint8(fileBuffer.Len()))
		if err != nil {
			return fmt.Errorf("%s write file pos: %w", file.Name(), err)
		}*/
	}

	fileOffset, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("seek fileOffset: %w", err)
	}

	//fmt.Println("filebuffer deflate\n", hex.Dump(fileBuffer.Bytes()))
	cData, err := helper.Deflate(fileBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("deflate fileBuffer: %w", err)
	}
	//fmt.Println("after", hex.Dump(cData))

	err = binary.Write(w, binary.LittleEndian, cData)
	if err != nil {
		return fmt.Errorf("write fileBuffer: %w", err)
	}

	dirOffset, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("seek dirOffset: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(dirEntries)+1))
	if err != nil {
		return fmt.Errorf("write dir count: %w", err)
	}

	for i, file := range dirEntries {

		err = binary.Write(w, binary.LittleEndian, file.crc)
		if err != nil {
			return fmt.Errorf("write direntry %d crc: %w", i, err)
		}
		err = binary.Write(w, binary.LittleEndian, file.offset)
		if err != nil {
			return fmt.Errorf("write direntry %d offset: %w", i, err)
		}
		err = binary.Write(w, binary.LittleEndian, file.size)
		if err != nil {
			return fmt.Errorf("write direntry %d size: %w", i, err)
		}
	}

	err = binary.Write(w, binary.LittleEndian, uint32(0x61580AC9))
	if err != nil {
		return fmt.Errorf("crc direntry: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(fileOffset))
	if err != nil {
		return fmt.Errorf("fileOffset: %w", err)
	}
	pos += 4

	err = binary.Write(w, binary.LittleEndian, uint32(len(fileBuffer.Bytes())))
	if err != nil {
		return fmt.Errorf("fileBuffer count: %w", err)
	}
	pos += int64(len(fileBuffer.Bytes()))

	err = binary.Write(w, binary.LittleEndian, [5]byte{'S', 'T', 'E', 'V', 'E'})
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}

	_, err = w.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek start: %w", err)
	}

	// err = binary.Write(w, binary.LittleEndian, uint32(pos))
	err = binary.Write(w, binary.LittleEndian, uint32(dirOffset))
	if err != nil {
		return fmt.Errorf("write header prefix proper: %w", err)
	}

	return nil
}
