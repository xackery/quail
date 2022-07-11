package tog

import (
	"fmt"
	"io"
	"text/template"
)

func (e *TOG) Save(w io.Writer) error {
	tmpl, err := template.New("toggle").Parse(togTemplate)
	if err != nil {
		return fmt.Errorf("parse togTemplate: %w", err)
	}
	type foo struct {
		Objects []*Object
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
		*POSITION 	{{.Position.X}} 	{{.Position.Y}}	 	{{.Position.Z}}	
		*ROTATION 	{{.Rotation.X}}	 	{{.Rotation.Y}}	 	{{.Rotation.Z}}
		*SCALE    	{{.Scale}}	
		*FILE     	{{.FileType}}     	{{.FileName}}
	*END_OBJECT{{end}}
*END_OBJECTGROUP
`
