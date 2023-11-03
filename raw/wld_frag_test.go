package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestFragment(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	tests := []struct {
		path      string
		file      string
		fragIndex int
		isDump    bool
	}{
		//{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 557}, // tex coord count misaligned
		//{path: "gequip.s3d", file: "gequip.wld", fragIndex: 0}, // Mesh
		//{path: "gfaydark.s3d", file: "gfaydark.wld", fragIndex: 0}, // Mesh
		//{path: "frozenshadow.s3d", file: "frozenshadow.wld", fragIndex: 0}, // Mesh
		//{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 2916, isDump: true}, // PASS
		//{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 0}, // PASS
		//{path: "poknowledge.s3d", file: "poknowledge.wld", fragIndex: 0}, // PASS
		// {path: "gequip4.s3d", file: "gequip4.wld", fragIndex: 0}, // PASS
		// {path: "gequip3.s3d", file: "gequip3.wld", fragIndex: 0}, // PASS
		//{path: "gfaydark_obj.s3d", file: "gfaydark_obj.wld", fragIndex: 0}, // PASS
		//{path: "gequip2.s3d", file: "gequip2.wld", fragIndex: 22280}, // PASS
		//{path: "zel_v2_chr.s3d", file: "zel_v2_chr.wld", fragIndex: 0}, // PASS
		//{path: "wol_v3_chr.s3d", file: "wol_v3_chr.wld", fragIndex: 0}, // PASS
		//{path: "globalhuf_chr.s3d", file: "globalhuf_chr.wld", fragIndex: 0}, // PASS
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

			fragments, err := tmpFragments(t, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read fragments: %s", err.Error())
			}

			for i := 1; i <= len(fragments); i++ {
				if tt.fragIndex != 0 && i != tt.fragIndex {
					continue
				}
				srcData := fragments[i-1]
				r := bytes.NewReader(srcData)
				reader := NewFrag(r)
				if reader == nil {
					t.Fatalf("frag %d 0x%x (%s) read: unsupported fragment", i, reader.FragCode(), FragName(int(reader.FragCode())))
				}

				err = reader.Read(r)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) read: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.src.hex", dirTest, tt.file), srcData, 0644)
					tag.Write(fmt.Sprintf("%s/%s.src.hex.tags", dirTest, tt.file))
				}

				buf := common.NewByteSeeker()
				buf.Write(srcData[:4])

				err = reader.Write(buf)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) write: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				dstData := buf.Bytes()

				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.dst.hex", dirTest, tt.file), dstData, 0644)
					tag.Write(fmt.Sprintf("%s/%s.dst.hex.tags", dirTest, tt.file))
				}

				//if !reflect.DeepEqual(data, dstData) {
				for j := 0; j < len(dstData); j++ {
					if tt.fragIndex != 0 && i != tt.fragIndex {
						break
					}
					if len(srcData) <= j {
						min := 0
						max := len(srcData)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
						max = len(dstData)
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))

						t.Fatalf("frag %d 0x%x (%s) src eof at offset %d (dst is too large by %d bytes)", i, reader.FragCode(), FragName(int(reader.FragCode())), i, len(dstData)-len(srcData))
					}
					if len(dstData) <= j {
						t.Fatalf("frag %d 0x%x (%s) dst eof at offset %d (dst is too small by %d bytes)", i, reader.FragCode(), FragName(int(reader.FragCode())), i, len(srcData)-len(dstData))
					}
					if dstData[j] == srcData[j] {
						continue
					}
					fmt.Printf("frag %d 0x%x (%s) mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", i, reader.FragCode(), FragName(int(reader.FragCode())), i, srcData[j], dstData[j], dstData[j])
					max := j + 16
					if max > len(srcData) {
						max = len(srcData)
					}

					min := j - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
					if max > len(dstData) {
						max = len(dstData)
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))
					t.Fatalf("frag %d 0x%x (%s) write: data mismatch", i, reader.FragCode(), FragName(int(reader.FragCode())))
				}
			}
			if tt.fragIndex != 0 {
				log.Debugf("Processed 1 fragment @ %d", tt.fragIndex)
				return
			}
			log.Debugf("Processed %d fragments", len(fragments))
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
		return nil, fmt.Errorf("read: %w", dec.Error())
	}
	return fragments, nil
}
