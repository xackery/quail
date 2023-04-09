package fragment

import (
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// UserData is used by WLDCOM but is never known to be used by EQ
type UserData struct {
}

// LoadUserData loads a UserData
func LoadUserData(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &UserData{}
	err := parseUserData(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse UserData: %w", err)
	}
	return e, nil
}

func parseUserData(r io.ReadSeeker, e *UserData) error {
	if e == nil {
		return fmt.Errorf("UserData is nil")
	}

	return fmt.Errorf("UserData is not implemented")
}

func (e *UserData) FragmentType() string {
	return "UserData"
}

// Data returns the raw data of the fragment
func (e *UserData) Data() []byte {
	return nil
}
