package quail

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

// PfsRead imports the quail target file
func (q *Quail) PfsRead(path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".eqg" {
		archive, err := pfs.NewFile(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer archive.Close()

		baseName := filepath.Base(path)
		baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

		q.Wld = wce.New(baseName)
		err = q.Wld.ReadEqgRaw(archive)
		if err != nil {
			return fmt.Errorf("wld read: %w", err)
		}
	}
	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("pfs load: %w", err)
	}
	defer pfs.Close()

	for _, file := range pfs.Files() {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".lit" {
			q.assetAdd(file.Name(), file.Data())
			continue
		}
		reader, err := raw.Read(ext, bytes.NewReader(file.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		reader.SetFileName(file.Name())
		err = q.RawRead(reader)
		if err != nil {
			return fmt.Errorf("rawRead %s: %w", file.Name(), err)
		}
	}

	summary := ""
	if len(q.Wld.TerDefs) > 0 {
		summary = fmt.Sprintf("%s%d terrain, ", summary, len(q.Wld.TerDefs))
	}
	if len(q.Wld.WorldTrees) > 0 {
		summary = fmt.Sprintf("%s%d trees, ", summary, len(q.Wld.WorldTrees))
	}

	if len(q.Wld.ActorDefs) > 0 {
		summary = fmt.Sprintf("%d actors, ", len(q.Wld.ActorDefs))
	}
	if len(q.Wld.ModDefs) > 0 {
		summary = fmt.Sprintf("%s%d models, ", summary, len(q.Wld.ModDefs))
	}
	if len(q.Wld.RGBTrackDefs) > 0 {
		summary = fmt.Sprintf("%s%d rgb tracks, ", summary, len(q.Wld.RGBTrackDefs))
	}

	if len(q.Wld.MdsDefs) > 0 {
		summary = fmt.Sprintf("%s%d skinned models, ", summary, len(q.Wld.MdsDefs))
	}
	if len(q.Wld.AmbientLights) > 0 {
		summary = fmt.Sprintf("%s%d lights, ", summary, len(q.Wld.AmbientLights))
	}
	if len(q.Wld.AniDefs) > 0 {
		summary = fmt.Sprintf("%s%d animations, ", summary, len(q.Wld.AniDefs))
	}
	if len(q.Wld.TrackDefs) > 0 {
		summary = fmt.Sprintf("%s%d tracks, ", summary, len(q.Wld.TrackDefs))
	}
	if len(q.Wld.DMTrackDef2s) > 0 {
		summary = fmt.Sprintf("%s%d dmtracks, ", summary, len(q.Wld.DMTrackDef2s))
	}

	if len(summary) > 0 {
		summary = summary[:len(summary)-2]
		fmt.Printf("Converted %s.\n", summary)
	}

	return nil
}

func (q *Quail) assetAdd(name string, data []byte) error {
	if q.Assets == nil {
		q.Assets = make(map[string][]byte)
	}
	q.Assets[name] = data
	return nil
}
