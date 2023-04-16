package wld

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/xackery/quail/log"
)

func parseFragment(path string, fragCode int32, fragOffset int, parser func(r io.ReadSeeker, fragOffset int) error) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer r.Close()

	return parser(r, fragOffset)
}

func parseFragments(path string, fragCode int32, parser func(r io.ReadSeeker, fragOffset int) error) (int, error) {

	total := 0
	files, err := os.ReadDir(path)
	if err != nil {
		return total, fmt.Errorf("read dir: %w", err)
	}
	for _, fe := range files {
		if fe.IsDir() {
			continue
		}

		if !strings.HasSuffix(fe.Name(), fmt.Sprintf("0x%02x.hex", fragCode)) {
			continue
		}

		fragOffset, err := strconv.Atoi(fe.Name()[0:4])
		if err != nil {
			return total, fmt.Errorf("strconv %s: %w", fe.Name(), err)
		}

		log.Debugf(fe.Name())
		err = parseFragment(path+"/"+fe.Name(), fragCode, fragOffset, parser)
		if err != nil {
			return total, fmt.Errorf("parse fragment %s: %w", fe.Name(), err)
		}
		total++
	}

	return total, nil
}
