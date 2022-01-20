package eqg

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
)

// Load will load an EQG
func (e *EQG) Load(r io.ReadSeeker) error {
	var err error
	type dirEntry struct {
		crc    uint32
		offset uint32
		size   uint32
	}
	dirEntries := []*dirEntry{}

	dirOffset := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &dirOffset)
	if err != nil {
		return fmt.Errorf("read dirOffset: %w", err)
	}

	// jump to dir entries
	_, err = r.Seek(int64(dirOffset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek dirOffset: %w", err)
	}

	fileCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &fileCount)
	if err != nil {
		return fmt.Errorf("read fileCount: %w", err)
	}

	for i := 0; i < int(fileCount); i++ {
		entry := &dirEntry{}
		err = binary.Read(r, binary.LittleEndian, &entry.crc)
		if err != nil {
			return fmt.Errorf("read %d crc: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &entry.offset)
		if err != nil {
			return fmt.Errorf("read %d offset: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &entry.size)
		if err != nil {
			return fmt.Errorf("read %d size: %w", i, err)
		}
		dirEntries = append(dirEntries, entry)
	}

	// reset back to start of file
	_, err = r.Seek(4, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek start: %w", err)
	}

	dump.Hex(dirOffset, fmt.Sprintf("dirOffset=%d", dirOffset))
	pfsHeader := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &pfsHeader)
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}
	dump.Hex(pfsHeader, "header=%s", string(pfsHeader[0:]))
	if pfsHeader != [4]byte{'P', 'F', 'S', ' '} {
		return fmt.Errorf("header mismatch")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	dump.Hex(version, "version=%d", version)
	if uint32(0x00020000) != version {
		return fmt.Errorf("unknown version")
	}

	/*fileCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &fileCount)
	if err != nil {
		return fmt.Errorf("read file count: %w", err)
	}
	dump.Hex(fileCount, "fileCount=%d", fileCount)*/

	e.files = []common.Filer{}

	var deflateSize uint32
	var inflateSize uint32

	for i := 0; i < int(fileCount); i++ {
		pos, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("seek current %d: %w", i, err)
		}

		entryIndex := -1
		for j, entry := range dirEntries {
			if entry.offset != uint32(pos) {
				continue
			}
			entryIndex = j
		}
		if entryIndex == -1 {
			return fmt.Errorf("data chunk %d has malformed offset %x", i, pos)
		}

		err = binary.Read(r, binary.LittleEndian, &deflateSize)
		if err != nil {
			return fmt.Errorf("read deflate size %d: %w", i, err)
		}
		dump.Hex(deflateSize, "%ddeflateSize=%d", i, deflateSize)
		err = binary.Read(r, binary.LittleEndian, &inflateSize)
		if err != nil {
			return fmt.Errorf("read inflate size %d: %w", i, err)
		}
		dump.Hex(deflateSize, "%dinflateSize=%d", i, inflateSize)

		deflateData := make([]byte, deflateSize)
		err = binary.Read(r, binary.LittleEndian, &deflateData)
		if err != nil {
			return fmt.Errorf("read inflate size %d: %w", i, err)
		}
		dump.Hex(deflateData, "%ddata chunk", i)
	}

	dump.Hex(fileCount, "fileCount=%d", fileCount)
	for i, entry := range dirEntries {
		dump.Hex(entry.crc, "%dcrc", i)
		dump.Hex(entry.offset, "%doffset", i)
		dump.Hex(entry.size, "%dsize", i)
	}
	r.Seek(int64(4+(len(dirEntries)*12)), io.SeekCurrent)

	steveFooter := [5]byte{}
	err = binary.Read(r, binary.LittleEndian, &steveFooter)
	if err != nil {
		return fmt.Errorf("read steveFooter: %w", err)
	}
	dump.Hex(steveFooter, "steveFooter")
	if steveFooter != [5]byte{'S', 'T', 'E', 'V', 'E'} {
		return fmt.Errorf("steve footer not STEVE")
	}
	var dateFooter uint32
	err = binary.Read(r, binary.LittleEndian, &dateFooter)
	if err != nil {
		return fmt.Errorf("read dateFooter: %w", err)
	}
	dump.Hex(dateFooter, "dateFooter=%d", dateFooter)

	/*
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
	return nil
}
