package geo

import (
	"fmt"
	"io"
)

// Object is an object instance
type Object struct {
	Name      string
	ModelName string
	Position  *Vector3
	Rotation  *Vector3
	Scale     float32
	FileType  string
	FileName  string
}

// WriteHeader writes the header for a Object
func (e *Object) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("name|model_name|position|rotation|scale|file_type|file_name\n")
	return err
}

// WriteString writes a Object to a string writer
func (e *Object) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%s|%s|%s|%s|%0.2f|%s|%s\n",
		e.Name,
		e.ModelName,
		e.Position.String(),
		e.Rotation.String(),
		e.Scale,
		e.FileType,
		e.FileName))
	return err
}
