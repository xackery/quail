package raw

import "io"

// Edd contations particle definitions used by prt
// examples are in eq root, actoremittersnew.edd, environmentemittersnew.edd, spellsnew.edd
type Edd struct {
}

func (edd *Edd) Read(r io.ReadSeeker) error {
	return nil
}
