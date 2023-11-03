package common

import (
	"bytes"
	"io"
)

// ByteSeeker is used primarily for tag systems in tests, it isn't efficient otherwise
type ByteSeeker struct {
	io.Seeker
	*bytes.Buffer
	offset int64
}

func NewByteSeeker() *ByteSeeker {
	return &ByteSeeker{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (b *ByteSeeker) Bytes() []byte {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	return b.Buffer.Bytes()
}

func (b *ByteSeeker) Len() int {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	return b.Buffer.Len()
}

func (b *ByteSeeker) Cap() int {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	return b.Buffer.Cap()
}

func (b *ByteSeeker) Reset() {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	b.Buffer.Reset()
	b.offset = 0
}

func (b *ByteSeeker) Write(p []byte) (n int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.Write(p)
	b.offset += int64(n)
	return
}

func (b *ByteSeeker) WriteByte(c byte) error {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	err := b.Buffer.WriteByte(c)
	if err == nil {
		b.offset++
	}
	return err
}

func (b *ByteSeeker) WriteString(s string) (n int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.WriteString(s)
	b.offset += int64(n)
	return
}

func (b *ByteSeeker) Seek(offset int64, whence int) (int64, error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	var err error
	switch whence {
	case io.SeekStart:
		b.offset = offset
	case io.SeekCurrent:
		b.offset += offset
	case io.SeekEnd:
		b.offset = int64(b.Buffer.Len()) + offset
	default:
		err = io.ErrUnexpectedEOF
	}
	return b.offset, err
}

func (b *ByteSeeker) Read(p []byte) (n int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.Read(p)
	b.offset += int64(n)
	return
}

func (b *ByteSeeker) ReadByte() (c byte, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	c, err = b.Buffer.ReadByte()
	if err == nil {
		b.offset++
	}
	return
}

func (b *ByteSeeker) ReadRune() (r rune, size int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	r, size, err = b.Buffer.ReadRune()
	if err == nil {
		b.offset += int64(size)
	}
	return
}

func (b *ByteSeeker) ReadFrom(r io.Reader) (n int64, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.ReadFrom(r)
	b.offset += n
	return
}

func (b *ByteSeeker) ReadString(delim byte) (line string, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	line, err = b.Buffer.ReadString(delim)
	b.offset += int64(len(line))
	return
}
