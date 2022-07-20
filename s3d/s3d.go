// s3d is an EverQuest pfs archive
package s3d

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/common"
)

// S3D represents a classic everquest zone archive format
type S3D struct {
	name                     string
	ShortName                string
	files                    []common.Filer
	fileCount                int
	fileEntries              []*FileEntry
	directoryChunks          []*ChunkEntry
	directoryChunksTotalSize uint32
}

type FileEntry struct {
	Name            string
	Data            []byte
	CRC             uint32
	Offset          uint32
	chunks          []*ChunkEntry
	chunksTotalSize uint32
	filePointer     uint32
}

func (e *FileEntry) String() string {
	return fmt.Sprintf("[%s (%d bytes)]", e.Name, e.chunksTotalSize)
}

type ChunkEntry struct {
	deflatedSize int32
	inflatedSize int32
	data         []byte
}

type ByOffset []*FileEntry

func (s ByOffset) Len() int {
	return len(s)
}

func (s ByOffset) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByOffset) Less(i, j int) bool {
	return s[i].Offset < s[j].Offset
}

type ByCRC []*FileEntry

func (s ByCRC) Len() int {
	return len(s)
}

func (s ByCRC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCRC) Less(i, j int) bool {
	return s[i].CRC < s[j].CRC
}

func New(name string) (*S3D, error) {
	e := &S3D{
		name: name,
	}
	return e, nil
}

// NewFile takes path and loads it as an eqg archive
func NewFile(path string) (*S3D, error) {
	e := &S3D{
		name: filepath.Base(path),
	}
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	err = e.Load(r)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return e, nil
}
