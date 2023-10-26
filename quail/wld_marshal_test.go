package quail

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/xackery/encdec"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
	"github.com/xackery/quail/pfs"
)

func TestWldMarshal(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		path      string
		file      string
		fragIndex int
	}{
		{"btp_chr.s3d", "btp_chr.wld", 0},
		{"bac_chr.s3d", "bac_chr.wld", 0},
		{"avi_chr.s3d", "avi_chr.wld", 0},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.path))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.file, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(tt.file)
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.file, err.Error())
			}
			world := common.NewWld("test")
			err = wld.Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode wld %s: %s", tt.file, err.Error())
			}

			srcFragments, err := tmpFragments(t, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read src fragments: %s", err.Error())
			}
			if len(srcFragments) == 0 {
				t.Fatalf("failed to read src fragments: no fragments")
			}

			q := New()
			err = q.WldUnmarshal(world)
			if err != nil {
				t.Fatalf("failed to convert wld %s: %s", tt.file, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = q.WldMarshal(buf)
			if err != nil {
				t.Fatalf("failed to marshal wld %s: %s", tt.file, err.Error())
			}

			/*
				dstFragments, err := tmpFragments(t, bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to read dst fragments: %s", err.Error())
				}

				if len(srcFragments) != len(dstFragments) {
					t.Fatalf("fragment count mismatch: %d != %d", len(srcFragments), len(dstFragments))
				}

				for i := 0; i < len(srcFragments); i++ {
					if !bytes.Equal(srcFragments[i], dstFragments[i]) {
						t.Fatalf("fragment %d mismatch", i)
					}
				}
			*/
		})
	}
}

func tmpFragments(t *testing.T, r io.ReadSeeker) (fragments [][]byte, err error) {

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
		return nil, fmt.Errorf("decode: %w", dec.Error())
	}
	return fragments, nil
}
