package raw

import (
	"fmt"
	"io"
	"strings"
)

func (lod *Lod) Write(w io.Writer) error {
	// Write the header
	_, err := w.Write([]byte("EQLOD\r\n"))
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each entry
	for _, entry := range lod.Entries {
		distance := fmt.Sprintf("%06.3f", entry.Distance)
		if strings.HasSuffix(distance, ".000") {
			distance = strings.TrimSuffix(distance, ".000")
		}

		line := fmt.Sprintf("%s,%s", entry.Category, entry.ObjectName)
		if distance != "00" {
			line += "," + distance
		}
		line += "\r\n"
		_, err := w.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("failed to write entry: %w", err)
		}
	}
	return nil
}
