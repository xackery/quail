package lit

import (
	"fmt"
	"io"
)

// Save writes a zon file to location
func (e *LIT) Save(w io.Writer) error {
	return fmt.Errorf("lit save is not yet supported")
}
