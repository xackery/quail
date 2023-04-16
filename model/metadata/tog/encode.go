package tog

import (
	"fmt"
	"io"
	"text/template"

	"github.com/xackery/quail/model/geo"
)

func (e *TOG) Encode(w io.Writer) error {
	tmpl, err := template.New("toggle").Parse(togTemplate)
	if err != nil {
		return fmt.Errorf("parse togTemplate: %w", err)
	}
	type foo struct {
		Objects []geo.Object
	}
	o := &foo{
		Objects: e.objects,
	}
	err = tmpl.Execute(w, o)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	return nil
}

var togTemplate = `*BEGIN_OBJECTGROUP{{range .Objects}}
	*BEGIN_OBJECT
		*NAME     	{{.Name}}
		*POSITION 	{{index .Position 0}} 	{{index .Position 1}}	 	{{index .Position 2}}
		*ROTATION 	{{index .Rotation 0}}	 	{{index .Rotation 1}}	 	{{index .Rotation 2}}
		*SCALE    	{{.Scale}}
		*FILE     	{{.FileType}}     	{{.FileName}}
	*END_OBJECT{{end}}
*END_OBJECTGROUP
`
