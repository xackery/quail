package raw

import (
	"fmt"
	"io"
	"text/template"
)

func (tog *Tog) Write(w io.Writer) error {
	tmpl, err := template.New("toggle").Parse(togTemplate)
	if err != nil {
		return fmt.Errorf("parse togTemplate: %w", err)
	}
	err = tmpl.Execute(w, tog)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	return nil
}

var togTemplate = `*BEGIN_OBJECTGROUP{{range .Entries}}
	*BEGIN_OBJECT
		*NAME     	{{.Name}}
		*POSITION 	{{index .Position 0}} 	{{index .Position 1}}	 	{{index .Position 2}}
		*ROTATION 	{{index .Rotation 0}}	 	{{index .Rotation 1}}	 	{{index .Rotation 2}}
		*SCALE    	{{.Scale}}
		*FILE     	{{.FileType}}     	{{.FileName}}
	*END_OBJECT{{end}}
*END_OBJECTGROUP
`
