package eqg

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// Load will load an EQG
func (e *EQG) Load(r io.ReadSeeker) error {
	var err error
	type dirEntry struct {
		crc    uint32
		offset uint32
		size   uint32
	}
	//	dirEntries := []*dirEntry{}

	dirOffset := 0
	err = binary.Read(r, binary.LittleEndian, uint32(dirOffset))
	if err != nil {
		return fmt.Errorf("write header prefix: %w", err)
	}

	pfsHeader := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, pfsHeader)
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}
	if pfsHeader != [4]byte{'P', 'F', 'S', ' '} {
		return fmt.Errorf("header mismatch")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, version)
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	if uint32(0x00020000) != version {
		return fmt.Errorf("unknown version")
	}

	fileCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, fileCount)
	if err != nil {
		return fmt.Errorf("read file count: %w", err)
	}

	e.files = []common.Filer{}
	/*
		for i := uint32(0); i < fileCount; i++ {


		dirEntries = append(dirEntries, &dirEntry{
				crc:    helper.FilenameCRC32(file.Name),
				size:   uint32(len(file.Data)),
				offset: uint32(pos),
			})

			//write compressed data
			cData, err := helper.Deflate(file.Data)
			if err != nil {
				return fmt.Errorf("deflate %s: %w", file.Name, err)
			}

			err = binary.Read(r, binary.LittleEndian, cData)
			if err != nil {
				return fmt.Errorf("%s write data: %w", file.Name, err)
			}

			// prep filebuffer
			err = binary.Read(fileBuffer, binary.LittleEndian, uint32(len(file.Name)+1))
			if err != nil {
				return fmt.Errorf("%s write name length: %w", file.Name, err)
			}

			_, err = fileBuffer.Write([]byte(file.Name))
			if err != nil {
				return fmt.Errorf("write name %s: %w", file.Name, err)
			}
			//fmt.Println(len(file.Name)+1, hex.Dump([]byte(file.Name)))

			err = binary.Read(fileBuffer, binary.LittleEndian, uint8(fileBuffer.Len()))
			if err != nil {
				return fmt.Errorf("%s write file pos: %w", file.Name, err)
			}
		}

		fileOffset, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("seek fileOffset: %w", err)
		}

		fmt.Println("filebuffer deflate\n", hex.Dump(fileBuffer.Bytes()))
		cData, err := helper.Deflate(fileBuffer.Bytes())
		if err != nil {
			return fmt.Errorf("deflate fileBuffer: %w", err)
		}
		//fmt.Println("after", hex.Dump(cData))

		err = binary.Read(r, binary.LittleEndian, cData)
		if err != nil {
			return fmt.Errorf("write fileBuffer: %w", err)
		}

		dirOffset, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("seek dirOffset: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, uint32(len(dirEntries)+1))
		if err != nil {
			return fmt.Errorf("write dir count: %w", err)
		}

		for i, file := range dirEntries {
			err = binary.Read(r, binary.LittleEndian, file.crc)
			if err != nil {
				return fmt.Errorf("write direntry %d crc: %w", i, err)
			}
			err = binary.Read(r, binary.LittleEndian, file.offset)
			if err != nil {
				return fmt.Errorf("write direntry %d offset: %w", i, err)
			}
			err = binary.Read(r, binary.LittleEndian, file.size)
			if err != nil {
				return fmt.Errorf("write direntry %d size: %w", i, err)
			}
		}

		err = binary.Read(r, binary.LittleEndian, uint32(0x61580AC9))
		if err != nil {
			return fmt.Errorf("magic number 0x61580AC9: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, uint32(fileOffset))
		if err != nil {
			return fmt.Errorf("fileOffset: %w", err)
		}
		pos += 4

		err = binary.Read(r, binary.LittleEndian, uint32(len(fileBuffer.Bytes())))
		if err != nil {
			return fmt.Errorf("fileBuffer count: %w", err)
		}
		pos += int64(len(fileBuffer.Bytes()))

		err = binary.Read(r, binary.LittleEndian, [5]byte{'S', 'T', 'E', 'V', 'E'})
		if err != nil {
			return fmt.Errorf("write header magic: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, uint32(time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("write header magic: %w", err)
		}
	*/
	return fmt.Errorf("not implemented")
}
