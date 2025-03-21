package raw

import (
	"fmt"
	"io"
)

func (lod *Lod) Write(w io.Writer) error {
	// Write the header
	_, err := w.Write([]byte("EQLOD\n"))
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each entry
	for _, entry := range lod.Entries {
		line := fmt.Sprintf("%s,%s,%0.3f\n", entry.Category, entry.ObjectName, entry.Distance)
		_, err := w.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("failed to write entry: %w", err)
		}
	}
	return nil
}
