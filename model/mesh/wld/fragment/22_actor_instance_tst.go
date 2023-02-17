package fragment

import (
	"bytes"
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
	e := &ActorInstanceTest{}
	err := parseActorInstanceTest(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse ActorInstanceTest: %w", err)
	}
	return e, nil
}

func parseActorInstanceTest(r io.ReadSeeker, e *ActorInstanceTest) error {
	if e == nil {
		return fmt.Errorf("ActorInstanceTest is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &e.Parameter)
	if err != nil {
		return fmt.Errorf("read ActorInstanceTest: %w", err)
	}
	return nil
}

func (e *ActorInstanceTest) FragmentType() string {
	return "ActorInstanceTest"
}

func (e *ActorInstanceTest) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
