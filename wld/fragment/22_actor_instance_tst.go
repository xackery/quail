package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// ActorInstanceTest information
type ActorInstanceTest struct {
	Parameter float32
}

func LoadActorInstanceTest(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &ActorInstanceTest{}
	err := parseActorInstanceTest(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse ActorInstanceTest: %w", err)
	}
	return l, nil
}

func parseActorInstanceTest(r io.ReadSeeker, l *ActorInstanceTest) error {
	if l == nil {
		return fmt.Errorf("ActorInstanceTest is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &l.Parameter)
	if err != nil {
		return fmt.Errorf("read ActorInstanceTest: %w", err)
	}
	return nil
}

func (l *ActorInstanceTest) FragmentType() string {
	return "ActorInstanceTest"
}
