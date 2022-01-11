package eqg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/helper"
)

// Save will write a EQG to writer
func (e *EQG) Save(w io.WriteSeeker) error {
	var err error

	type dirEntry struct {
		crc    uint32
		offset uint32
		sz     uint32
	}
	dirEntries := []*dirEntry{}

	pos := int64(0)

	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return fmt.Errorf("write header prefix: %w", err)
	}
	pos += 4

	err = binary.Write(w, binary.LittleEndian, [4]byte{'P', 'F', 'S', ' '})
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}
	pos += 4

	err = binary.Write(w, binary.LittleEndian, uint32(131072))
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	pos += 4

	fileBuffer := bytes.NewBuffer(nil)
	err = binary.Write(fileBuffer, binary.LittleEndian, uint32(len(e.files)+1))
	if err != nil {
		return fmt.Errorf("write file count: %w", err)
	}
	pos += 4

	filePos := 0
	err = binary.Write(fileBuffer, binary.LittleEndian, uint32(len(e.files)))
	if err != nil {
		return fmt.Errorf("write name count: %w", err)
	}
	filePos += 4

	e.files = byCRC(e.files)
	for _, file := range e.files {
		pos, err = w.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("%s seek: %w", file.Name, err)
		}

		dirEntries = append(dirEntries, &dirEntry{
			crc:    helper.FilenameCRC32(file.Name),
			sz:     uint32(len(file.Data)),
			offset: uint32(pos),
		})

		//write compressed data
		cData, err := helper.Deflate(file.Data)
		if err != nil {
			return fmt.Errorf("deflate %s: %w", file.Name, err)
		}

		err = binary.Write(w, binary.LittleEndian, uint32(len(cData)))
		if err != nil {
			return fmt.Errorf("%s write compressed size: %w", file.Name, err)
		}

		// prep filebuffer
		err = binary.Write(fileBuffer, binary.LittleEndian, uint32(len(file.Name)+1))
		if err != nil {
			return fmt.Errorf("%s write name length: %w", file.Name, err)
		}
		filePos += 4

		err = helper.WriteString(fileBuffer, file.Name)
		if err != nil {
			return fmt.Errorf("%s write name: %w", file.Name, err)
		}
		filePos += len(file.Name) + 1

		err = binary.Write(fileBuffer, binary.LittleEndian, uint8(filePos))
		if err != nil {
			return fmt.Errorf("%s write file pos: %w", file.Name, err)
		}

		err = binary.Write(fileBuffer, binary.LittleEndian, uint32(len(file.Data)))
		if err != nil {
			return fmt.Errorf("%s write uncompressed size: %w", file.Name, err)
		}
		filePos += 4

		err = binary.Write(fileBuffer, binary.LittleEndian, cData)
		if err != nil {
			return fmt.Errorf("%s write compressed data: %w", file.Name, err)
		}
		filePos += len(cData)
	}

	fileOffset, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("dirOffset seek: %w", err)
	}

	cData, err := helper.Deflate(fileBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("deflate fileBuffer: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(cData)))
	if err != nil {
		return fmt.Errorf("fileBuffer write compressed size: %w", err)
	}

	pos, err = w.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("dirOffset seek: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(pos))
	if err != nil {
		return fmt.Errorf("dirOffset write: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(dirEntries)+1))
	if err != nil {
		return fmt.Errorf("dirOffset write: %w", err)
	}
	for i, dir := range dirEntries {
		err = binary.Write(w, binary.LittleEndian, dir.crc)
		if err != nil {
			return fmt.Errorf("dir %d crc write: %w", i, err)
		}
		err = binary.Write(w, binary.LittleEndian, dir.offset)
		if err != nil {
			return fmt.Errorf("dir %d offset write: %w", i, err)
		}
		err = binary.Write(w, binary.LittleEndian, dir.sz)
		if err != nil {
			return fmt.Errorf("dir %d size write: %w", i, err)
		}
	}
	err = binary.Write(w, binary.LittleEndian, uint32(0x61580AC9))
	if err != nil {
		return fmt.Errorf("magic number 0x61580AC9: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(fileOffset))
	if err != nil {
		return fmt.Errorf("fileOffset: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(fileBuffer.Bytes())))
	if err != nil {
		return fmt.Errorf("fileBuffer count: %w", err)
	}

	return nil
}
