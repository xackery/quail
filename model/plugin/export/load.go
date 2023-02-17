package export

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/model/mesh/wld"
	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
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
		{invoke: e.loadWld, name: "wld"},
	}

	for _, evt := range events {
		err = evt.invoke()
		if err != nil {
			return fmt.Errorf("%s load: %w", evt.name, err)
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
		return fmt.Errorf("new: %w", err)
	}

	err = e.model.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err = e.loadParticlePoints()
	if err != nil {
		return fmt.Errorf("loadParticlePoints: %w", err)
	}

	err = e.loadParticleRenders()
	if err != nil {
		return fmt.Errorf("zon loadParticleRenders: %w", err)
	}
	e.name += ".zon"
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
		return fmt.Errorf("%s new: %w", e.name, err)
	}

	err = e.model.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%s decode: %w", e.name, err)
	}

	err = e.loadLayer()
	if err != nil {
		return fmt.Errorf("loadLayer: %w", err)
	}

	err = e.loadParticleRenders()
	if err != nil {
		return fmt.Errorf("loadParticleRenders: %w", err)
	}

	err = e.loadParticlePoints()
	if err != nil {
		return fmt.Errorf("loadParticlePoints: %w", err)
	}
	e.name += ".mds"
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
		return fmt.Errorf("new: %w", err)
	}

	err = e.model.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err = e.loadLayer()
	if err != nil {
		return fmt.Errorf("loadLayer: %w", err)
	}

	err = e.loadParticleRenders()
	if err != nil {
		return fmt.Errorf("loadParticleRenders: %w", err)
	}

	err = e.loadParticlePoints()
	if err != nil {
		return fmt.Errorf("loadParticlePoints: %w", err)
	}
	e.name += ".mod"
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
		return fmt.Errorf("new: %w", err)
	}

	err = e.model.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err = e.loadLayer()
	if err != nil {
		return fmt.Errorf("loadLayer: %w", err)
	}

	err = e.loadParticleRenders()
	if err != nil {
		return fmt.Errorf("loadParticleRenders: %w", err)
	}

	err = e.loadParticlePoints()
	if err != nil {
		return fmt.Errorf("loadParticlePoints: %w", err)
	}
	e.name += ".ter"
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
		return fmt.Errorf("new: %w", err)
	}
	err = l.Decode(bytes.NewReader(layEntry))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err = e.model.SetLayers(l.Layers())
	if err != nil {
		return fmt.Errorf("setlayers: %w", err)
	}
	return nil
}

func (e *Export) loadParticleRenders() error {
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
		return fmt.Errorf("new: %w", err)
	}
	err = p.Decode(bytes.NewReader(prtEntry))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err = e.model.SetParticleRenders(p.ParticleRenders())
	if err != nil {
		return fmt.Errorf("setparticles: %w", err)
	}

	return nil
}

func (e *Export) loadParticlePoints() error {
	prtName := fmt.Sprintf("%s.pts", e.name)
	prtEntry, err := e.archive.File(prtName)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return fmt.Errorf("file '%s': %w", prtName, err)
	}

	if len(prtEntry) == 0 {
		return nil
	}

	p, err := pts.New(prtName, e.archive)
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}
	err = p.Decode(bytes.NewReader(prtEntry))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	err = e.model.SetParticlePoints(p.ParticlePoints())
	if err != nil {
		return fmt.Errorf("setparticles: %w", err)
	}

	return nil
}

func (e *Export) loadWld() error {
	var err error
	var data []byte
	for _, entry := range e.archive.Files() {
		if !strings.HasSuffix(entry.Name(), ".wld") {
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

	e.model, err = wld.New(e.name, e.archive)
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}

	err = e.model.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	e.name += ".wld"
	return nil
}
