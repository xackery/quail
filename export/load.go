package export

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail/lay"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/prt"
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
)

func (e *Export) LoadArchive() error {
	var err error

	events := []struct {
		invoke func() error
		name   string
	}{
		{invoke: e.loadZon, name: "zon"},
		{invoke: e.loadMds, name: "mds"},
		{invoke: e.loadMod, name: "mod"},
		{invoke: e.loadTer, name: "ter"},
	}

	for _, evt := range events {
		err = evt.invoke()
		if err != nil {
			return fmt.Errorf("load %s: %w", evt.name, err)
		}
		if e.model != nil {
			return nil
		}
	}

	return nil
}

func (e *Export) loadZon() error {
	var err error
	var data []byte
	for _, entry := range e.archive.Files() {
		if !strings.HasSuffix(entry.Name(), ".zon") {
			continue
		}

		data, err = e.archive.File(entry.Name())
		if err != nil {
			return fmt.Errorf("file %s: %w", entry.Name(), err)
		}
	}
	if len(data) == 0 {
		return nil
	}

	e.model, err = zon.New(e.name, e.archive)
	if err != nil {
		return fmt.Errorf("zon new: %w", err)
	}

	err = e.model.Load(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("zon load: %w", err)
	}

	err = e.loadParticles()
	if err != nil {
		return fmt.Errorf("zon loadParticles: %w", err)
	}
	return nil
}

func (e *Export) loadMds() error {
	var err error
	var data []byte
	for _, entry := range e.archive.Files() {
		if !strings.HasSuffix(entry.Name(), ".mds") {
			continue
		}

		data, err = e.archive.File(entry.Name())
		if err != nil {
			return fmt.Errorf("file %s: %w", entry.Name(), err)
		}
	}
	if len(data) == 0 {
		return nil
	}

	e.model, err = mds.New(e.name, e.archive)
	if err != nil {
		return fmt.Errorf("mds new: %w", err)
	}

	err = e.model.Load(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("mds load: %w", err)
	}

	err = e.loadLayer()
	if err != nil {
		return fmt.Errorf("mds loadLayer: %w", err)
	}

	err = e.loadParticles()
	if err != nil {
		return fmt.Errorf("mds loadParticles: %w", err)
	}
	return nil
}

func (e *Export) loadMod() error {
	var err error
	var data []byte
	for _, entry := range e.archive.Files() {
		if !strings.HasSuffix(entry.Name(), ".mod") {
			continue
		}

		data, err = e.archive.File(entry.Name())
		if err != nil {
			return fmt.Errorf("file %s: %w", entry.Name(), err)
		}
	}
	if len(data) == 0 {
		return nil
	}

	e.model, err = mod.New(e.name, e.archive)
	if err != nil {
		return fmt.Errorf("mod new: %w", err)
	}

	err = e.model.Load(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("mod load: %w", err)
	}

	err = e.loadLayer()
	if err != nil {
		return fmt.Errorf("mod loadLayer: %w", err)
	}

	err = e.loadParticles()
	if err != nil {
		return fmt.Errorf("mod loadParticles: %w", err)
	}
	return nil
}

func (e *Export) loadTer() error {
	var err error
	var data []byte
	for _, entry := range e.archive.Files() {
		if !strings.HasSuffix(entry.Name(), ".ter") {
			continue
		}

		data, err = e.archive.File(entry.Name())
		if err != nil {
			return fmt.Errorf("file %s: %w", entry.Name(), err)
		}
	}
	if len(data) == 0 {
		return nil
	}

	e.model, err = ter.New(e.name, e.archive)
	if err != nil {
		return fmt.Errorf("ter new: %w", err)
	}

	err = e.model.Load(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("ter load: %w", err)
	}

	err = e.loadLayer()
	if err != nil {
		return fmt.Errorf("ter loadLayer: %w", err)
	}

	err = e.loadParticles()
	if err != nil {
		return fmt.Errorf("ter loadParticles: %w", err)
	}
	return nil
}

func (e *Export) loadLayer() error {
	layName := fmt.Sprintf("%s.lay", e.name)
	layEntry, err := e.archive.File(layName)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return fmt.Errorf("file '%s': %w", layName, err)
	}

	if len(layEntry) == 0 {
		return nil
	}

	l, err := lay.New(layName, e.archive)
	if err != nil {
		return fmt.Errorf("lay new: %w", err)
	}
	err = l.Load(bytes.NewReader(layEntry))
	if err != nil {
		return fmt.Errorf("lay load: %w", err)
	}

	err = e.model.SetLayers(l.Layers())
	if err != nil {
		return fmt.Errorf("lay setlayers: %w", err)
	}
	return nil
}

func (e *Export) loadParticles() error {
	prtName := fmt.Sprintf("%s.prt", e.name)
	prtEntry, err := e.archive.File(prtName)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return fmt.Errorf("file '%s': %w", prtName, err)
	}

	if len(prtEntry) == 0 {
		return nil
	}

	p, err := prt.New(prtName, e.archive)
	if err != nil {
		return fmt.Errorf("prt new: %w", err)
	}
	err = p.Load(bytes.NewReader(prtEntry))
	if err != nil {
		return fmt.Errorf("prt load: %w", err)
	}

	err = e.model.SetParticles(p.Particles())
	if err != nil {
		return fmt.Errorf("prt setparticles: %w", err)
	}
	return nil
}
