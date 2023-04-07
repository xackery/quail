package zon

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

func (e *ZON) BlenderExport(dir string) error {
	mw, err := os.Create(fmt.Sprintf("%s/zon_%s_model.txt", dir, e.Name()))
	if err != nil {
		return fmt.Errorf("zon_%s_model.txt: %w", e.Name(), err)
	}
	defer mw.Close()
	for _, model := range e.models {
		mw.WriteString(model.baseName + " ")
		mw.WriteString(model.name + "\n")
	}

	ow, err := os.Create(fmt.Sprintf("%s/zon_%s_object.txt", dir, e.Name()))
	if err != nil {
		return fmt.Errorf("zon_%s_object.txt: %w", e.Name(), err)
	}
	defer ow.Close()
	for _, obj := range e.objects {
		mw.WriteString(obj.modelName + " ")
		mw.WriteString(obj.name + " ")
		mw.WriteString(dump.Str(obj.rotation) + " ")
		mw.WriteString(dump.Str(obj.translation) + " ")
		mw.WriteString(dump.Str(obj.scale) + "\n")
	}

	lw, err := os.Create(fmt.Sprintf("%s/zon_%s_light.txt", dir, e.Name()))
	if err != nil {
		return fmt.Errorf("zon_%s_light.txt: %w", e.Name(), err)
	}
	defer lw.Close()
	for _, light := range e.lights {
		mw.WriteString(light.name + " ")
		mw.WriteString(dump.Str(light.color) + " ")
		mw.WriteString(dump.Str(light.position) + " ")
		mw.WriteString(dump.Str(light.radius) + "\n")
	}
	return nil
}
