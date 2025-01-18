package pfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	//"strings"

	"github.com/xackery/quail/helper"
)

// Read will read a Pfs archive
func (e *Pfs) Read(r io.ReadSeeker) error {
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

	pfsHeader := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &pfsHeader)
	if err != nil {
		return fmt.Errorf("write header magic: %w", err)
	}
	if pfsHeader != [4]byte{'P', 'F', 'S', ' '} {
		return fmt.Errorf("header mismatch")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	if uint32(0x00020000) != version {
		return fmt.Errorf("unknown version")
	}

	e.files = []*FileEntry{}
	fileByCRCs := make(map[uint32][]byte)
	dirNameByCRCs := make(map[uint32]string)

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
		entry := dirEntries[entryIndex]

		data := []byte{}
		currentSize := 0
		for entry.size > uint32(currentSize) {
			err = binary.Read(r, binary.LittleEndian, &deflateSize)
			if err != nil {
				return fmt.Errorf("read deflate size %d: %w", i, err)
			}

			err = binary.Read(r, binary.LittleEndian, &inflateSize)
			if err != nil {
				return fmt.Errorf("read inflate size %d: %w", i, err)
			}

			deflateData := make([]byte, deflateSize)
			err = binary.Read(r, binary.LittleEndian, &deflateData)
			if err != nil {
				return fmt.Errorf("read inflate size %d: %w", i, err)
			}

			chunkData, err := helper.Inflate(deflateData, int(inflateSize))
			if err != nil {
				return fmt.Errorf("inflate %d: %w", i, err)
			}
			currentSize += int(inflateSize)
			data = append(data, chunkData...)
		}

		nameBuf := bytes.NewBuffer(data)
		var fileNameCount uint32
		err = binary.Read(nameBuf, binary.LittleEndian, &fileNameCount)
		if err != nil {
			return fmt.Errorf("read fileNameCount %w", err)
		}

		if fileNameCount != fileCount-1 {
			fileByCRCs[entry.crc] = data
			continue
		}

		for j := 0; j < int(fileNameCount); j++ {
			var fileNameLength uint32
			err = binary.Read(nameBuf, binary.LittleEndian, &fileNameLength)
			if err != nil {
				return fmt.Errorf("read fileNameCount %w", err)
			}

			nameData := make([]byte, fileNameLength)
			err = binary.Read(nameBuf, binary.LittleEndian, &nameData)
			if err != nil {
				return fmt.Errorf("read nameData %w", err)
			}
			name := string(nameData[0 : len(nameData)-1])
			//name = strings.ToLower(name)
			dirNameByCRCs[helper.FilenameCRC32(name)] = name
		}
	}

	for crc, data := range fileByCRCs {
		dirName, ok := dirNameByCRCs[crc]
		if !ok {
			//log.Warnf("dirName for crc %d not found", crc)
			continue
			//return fmt.Errorf("dirName for crc %d not found", crc)
		}

		//force spaces in pfs archives to _
		//dirName = strings.ReplaceAll(dirName, " ", "_")
		e.files = append(e.files, NewFileEntry(dirName, data))
	}

	_, err = r.Seek(int64(4+(len(dirEntries)*12)), io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("seek to footer: %w", err)
	}

	steveFooter := [5]byte{}
	err = binary.Read(r, binary.LittleEndian, &steveFooter)
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("read steveFooter: %w", err)
		}
		return nil
	}
	if steveFooter != [5]byte{'S', 'T', 'E', 'V', 'E'} {
		return fmt.Errorf("steve footer not STEVE")
	}
	var dateFooter uint32
	err = binary.Read(r, binary.LittleEndian, &dateFooter)
	if err != nil {
		return fmt.Errorf("read dateFooter: %w", err)
	}

	e.fileCount = len(e.files)
	return nil
}
