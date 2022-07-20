package s3d

import (
	"fmt"
	"io"
)

// Encode writes a s3d file to provided writer
func (e *S3D) Encode(w io.Writer) error {
	/*
		var directoryIndex uint32
		var magicNumber uint32 = 0x20534650
		var versionNumber uint32 = 2 // TODO: verify

		var value uint32

		//this is just placeholder we need to write it after we know position
		err := binary.Write(w, binary.LittleEndian, &directoryIndex)
		if err != nil {
			return fmt.Errorf("write directory index: %w", err)
		}

		err = binary.Write(w, binary.LittleEndian, &magicNumber)
		if err != nil {
			return fmt.Errorf("write magic number: %w", err)
		}

		err = binary.Write(w, binary.LittleEndian, &versionNumber)
		if err != nil {
			return fmt.Errorf("write version number: %w", err)
		}

		var fileCount uint32 = uint32(len(e.files))
		err = binary.Write(w, binary.LittleEndian, &fileCount)
		if err != nil {
			return fmt.Errorf("write file count: %w", err)
		}

		filenames := []string{}

		for i, fe := range e.fileEntries {

			crc, err := helper.GenerateCRC32(fe.Data)
			if err != nil {
				return fmt.Errorf("generate crc for %s: %w", fe.Name, err)
			}
			err = binary.Write(w, binary.LittleEndian, &crc)
			if err != nil {
				return fmt.Errorf("write crc %d/%d: %w", i, fileCount, err)
			}

			err = binary.Write(w, binary.LittleEndian, &entry.Offset)
			if err != nil {
				return fmt.Errorf("write offset %d/%d: %w", i, fileCount, err)
			}
			debugInfo := fmt.Sprintf("%d/%d 0x%x", i, fileCount, entry.Offset)
			// size is the uncompressed size of the file
			var size uint32
			err = binary.Write(w, binary.LittleEndian, &size)
			if err != nil {
				return fmt.Errorf("write size %s: %w", debugInfo, err)
			}

			buf := bytes.NewBuffer(nil)

			cachedOffset, err := r.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("seek cached offset %s: %w", debugInfo, err)
			}
			_, err = r.Seek(int64(entry.Offset), io.SeekStart)
			if err != nil {
				return fmt.Errorf("seek offset %s: %w", debugInfo, err)
			}

			for uint32(buf.Len()) != size {
				var deflatedLength uint32
				var inflatedLength uint32
				err = binary.Write(w, binary.LittleEndian, &deflatedLength)
				if err != nil {
					return fmt.Errorf("write deflated length %s: %w", debugInfo, err)
				}

				err = binary.Write(w, binary.LittleEndian, &inflatedLength)
				if err != nil {
					return fmt.Errorf("write inflated length %s: %w", debugInfo, err)
				}

				//if inflatedLength < deflatedLength {
				//	return fmt.Errorf("inflated < deflated, offset misaligned? %d/%d", i, fileCount)
				//}

				compressedData := make([]byte, deflatedLength)
				err = binary.Write(w, binary.LittleEndian, compressedData)
				if err != nil {
					return fmt.Errorf("write compressed data: %s: %w", debugInfo, err)
				}

				fr, err := zlib.NewReaderDict(bytes.NewReader(compressedData), nil)
				if err != nil {
					return fmt.Errorf("zlib new reader dict %s: %w", debugInfo, err)
				}

				inflatedWritten, err := io.Copy(buf, fr)
				if err != nil {
					return fmt.Errorf("copy: %s: %w", debugInfo, err)
				}

				if inflatedWritten != int64(inflatedLength) {
					return fmt.Errorf("inflate mismatch %s: %w", debugInfo, err)
				}
			}
			if buf.Len() < int(size) {
				return fmt.Errorf("size mismatch %s: %w", debugInfo, err)
			}
			entry.Data = buf.Bytes()

			if entry.CRC == 0x61580AC9 {
				fr := bytes.NewReader(buf.Bytes())
				var filenameCount uint32
				err = binary.Write(wr, binary.LittleEndian, &filenameCount)
				if err != nil {
					return fmt.Errorf("filename count %s: %w", debugInfo, err)
				}

				for j := uint32(0); j < filenameCount; j++ {
					err = binary.Write(wr, binary.LittleEndian, &value)
					if err != nil {
						return fmt.Errorf("filename length %s: %w", debugInfo, err)
					}
					filename, err := helper.ParseFixedString(fr, value)
					if err != nil {
						return fmt.Errorf("filename %s: %w", debugInfo, err)
					}
					filenames = append(filenames, filename)
				}

				_, err = r.Seek(cachedOffset, io.SeekStart)
				if err != nil {
					return fmt.Errorf("seek cached offset %s: %w", debugInfo, err)
				}
				continue
			}
			e.fileEntries = append(e.fileEntries, entry)
			_, err = r.Seek(cachedOffset, io.SeekStart)
			if err != nil {
				return fmt.Errorf("seek cached offset %s: %w", debugInfo, err)
			}
		}

		sort.Sort(ByOffset(e.fileEntries))
		for i, entry := range e.fileEntries {
			if len(filenames) < i {
				return fmt.Errorf("entry %d has no name", i)
			}
			entry.Name = filenames[i]
			fe, err := common.NewFileEntry(entry.Name[0:len(entry.Name)-1], entry.Data)
			if err != nil {
				return fmt.Errorf("entry %d newFileEntry: %w", i, err)
			}
			e.files = append(e.files, fe)
		}*/
	return fmt.Errorf("encode is not yet supported")
}
