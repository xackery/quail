package lit

import (
	"bytes"
	"fmt"
	"os"
)

// LIT is a litrain file struct
type LIT struct {
	name string
}

func New(name string) (*LIT, error) {
	t := &LIT{
		name: name,
	}
	return t, nil
}

func (e *LIT) Name() string {
	return e.name
}

func (e *LIT) Data() []byte {
	w := bytes.NewBuffer(nil)

	err := e.Save(w)
	if err != nil {
		fmt.Println("failed to save litrain data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}
