package common

import (
	"bytes"
	"io"
)

// ByteSeekerTest is used primarily for tag systems in tests, it isn't efficient otherwise
type ByteSeekerTest struct {
	io.Seeker
	*bytes.Buffer
	offset int64
}

func NewByteSeekerTest() *ByteSeekerTest {
	return &ByteSeekerTest{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (b *ByteSeekerTest) Bytes() []byte {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	return b.Buffer.Bytes()
}

func (b *ByteSeekerTest) Len() int {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	return b.Buffer.Len()
}

func (b *ByteSeekerTest) Cap() int {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	return b.Buffer.Cap()
}

func (b *ByteSeekerTest) Reset() {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	b.Buffer.Reset()
	b.offset = 0
}

func (b *ByteSeekerTest) Write(p []byte) (n int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.Write(p)
	b.offset += int64(n)
	return
}

func (b *ByteSeekerTest) WriteByte(c byte) error {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	err := b.Buffer.WriteByte(c)
	if err == nil {
		b.offset++
	}
	return err
}

func (b *ByteSeekerTest) WriteString(s string) (n int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.WriteString(s)
	b.offset += int64(n)
	return
}

func (b *ByteSeekerTest) Seek(offset int64, whence int) (int64, error) {
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

func (b *ByteSeekerTest) Read(p []byte) (n int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.Read(p)
	b.offset += int64(n)
	return
}

func (b *ByteSeekerTest) ReadByte() (c byte, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	c, err = b.Buffer.ReadByte()
	if err == nil {
		b.offset++
	}
	return
}

func (b *ByteSeekerTest) ReadRune() (r rune, size int, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	r, size, err = b.Buffer.ReadRune()
	if err == nil {
		b.offset += int64(size)
	}
	return
}

func (b *ByteSeekerTest) ReadFrom(r io.Reader) (n int64, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	n, err = b.Buffer.ReadFrom(r)
	b.offset += n
	return
}

func (b *ByteSeekerTest) ReadString(delim byte) (line string, err error) {
	if b.Buffer == nil {
		b.Buffer = bytes.NewBuffer(nil)
	}
	line, err = b.Buffer.ReadString(delim)
	b.offset += int64(len(line))
	return
}
