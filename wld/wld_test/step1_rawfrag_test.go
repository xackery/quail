package wld_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw/rawfrag"
)

func TestRawFragReadWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	for _, tt := range tests {
		t.Run(tt.baseName, func(t *testing.T) {
			s3dName := fmt.Sprintf("%s.s3d", tt.baseName)
			s3dPath := fmt.Sprintf("%s/%s", eqPath, s3dName)
			wldName := tt.wldName
			if wldName == "" {
				wldName = fmt.Sprintf("%s.wld", tt.baseName)
			}

			pfs, err := pfs.NewFile(s3dPath)
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", s3dName, err.Error())
			}
			defer pfs.Close()

			data, err := pfs.File(wldName)
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", wldName, err.Error())
			}

			fragments, err := tmpFragments(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read fragments: %s", err.Error())
			}

			for i := 0; i < len(fragments); i++ {
				srcData := fragments[i]
				fragBuf := bytes.NewReader(srcData)
				srcFragRW := rawfrag.NewFrag(fragBuf)
				if srcFragRW == nil {
					t.Fatalf("frag %d read: unsupported fragment", i)
				}

				err = srcFragRW.Read(fragBuf)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) read: %s", i+1, srcFragRW.FragCode(), rawfrag.FragName(int(srcFragRW.FragCode())), err.Error())
				}

				buf := &bytes.Buffer{}

				err = srcFragRW.Write(buf)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) write: %s", i+1, srcFragRW.FragCode(), rawfrag.FragName(int(srcFragRW.FragCode())), err.Error())
				}

				fragBuf.Seek(0, io.SeekStart)
				dstFragRW := rawfrag.NewFrag(fragBuf)
				err = dstFragRW.Read(bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) read: %s", i+1, dstFragRW.FragCode(), rawfrag.FragName(int(dstFragRW.FragCode())), err.Error())
				}

				diff := deep.Equal(srcFragRW, dstFragRW)
				if diff != nil {
					t.Fatalf("wld diff %s frag %d 0x%x (%s): %s", tt.baseName, i+1, srcFragRW.FragCode(), rawfrag.FragName(int(srcFragRW.FragCode())), diff)
				}

			}
			log.Debugf("Processed %d fragments for %s", len(fragments), tt.baseName)
		})
	}
}

func tmpFragments(r io.ReadSeeker) (fragments [][]byte, err error) {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	_ = dec.Bytes(4)
	_ = int(dec.Uint32())

	fragmentCount := dec.Uint32()
	_ = dec.Uint32() //unk1
	_ = dec.Uint32() //unk2
	hashSize := dec.Uint32()
	_ = dec.Uint32() //unk3
	_ = dec.Bytes(int(hashSize))

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {

		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, fragmentCount, err)
		}
		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data = append(fragCode, data...)

		fragments = append(fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, fragmentCount, err)
		}
	}

	if dec.Error() != nil {
		return nil, fmt.Errorf("read: %w", dec.Error())
	}
	return fragments, nil
}
