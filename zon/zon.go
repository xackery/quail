package zon

import (
	"bytes"
	"fmt"
	"os"

	"github.com/g3n/engine/math32"
)

// ZON is a zon file struct
type ZON struct {
	name    string
	models  []*model
	objects []*object
	regions []*region
	lights  []*light
}

type model struct {
	name string
}

type object struct {
	modelName string
	name      string
	position  math32.Vector3
	rotation  math32.Vector3
	scale     float32
}

type region struct {
	name    string
	center  math32.Vector3
	unknown math32.Vector3
	extent  math32.Vector3
}

type light struct {
	name     string
	position math32.Vector3
	color    math32.Color
	radius   float32
}

func New(name string) (*ZON, error) {
	z := &ZON{
		name: name,
	}
	return z, nil
}

func (e *ZON) Name() string {
	return e.name
}

func (e *ZON) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Save(w)
	if err != nil {
		fmt.Println("failed to save zon data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}

// Models returns a slice of names
func (e *ZON) ModelNames() []string {
	names := []string{}
	for _, m := range e.models {
		names = append(names, m.name)
	}
	return names
}
