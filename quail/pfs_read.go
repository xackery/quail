package quail

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/helper"
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
			return fmt.Errorf("raw.Read %s: %w", file.Name(), err)
		}
		reader.SetFileName(file.Name())
		err = q.RawRead(reader)
		if err != nil {
			return fmt.Errorf("q rawRead %s: %w", file.Name(), err)
		}
	}

	if q.Wld == nil {
		// edge cases like grass.s3d has no wld data
		return nil
	}

	summary := ""
	if len(q.Wld.TerDefs) > 0 {
		summary = fmt.Sprintf("%s%d terrain, ", summary, len(q.Wld.TerDefs))
	}
	if len(q.Wld.WorldTrees) > 0 {
		summary = fmt.Sprintf("%s%d tree%s, ", summary, len(q.Wld.WorldTrees), helper.Pluralize(len(q.Wld.WorldTrees)))
	}

	if len(q.Wld.ActorDefs) > 0 {
		summary = fmt.Sprintf("%d actor%s, ", len(q.Wld.ActorDefs), helper.Pluralize(len(q.Wld.ActorDefs)))
	}
	if len(q.Wld.ModDefs) > 0 {
		summary = fmt.Sprintf("%s%d model%s, ", summary, len(q.Wld.ModDefs), helper.Pluralize(len(q.Wld.ModDefs)))
	}
	if len(q.Wld.RGBTrackDefs) > 0 {
		summary = fmt.Sprintf("%s%d rgb track%s, ", summary, len(q.Wld.RGBTrackDefs), helper.Pluralize(len(q.Wld.RGBTrackDefs)))
	}

	if len(q.Wld.MdsDefs) > 0 {
		summary = fmt.Sprintf("%s%d skinned model%s, ", summary, len(q.Wld.MdsDefs), helper.Pluralize(len(q.Wld.MdsDefs)))
	}
	if len(q.Wld.AmbientLights) > 0 {
		summary = fmt.Sprintf("%s%d light%s, ", summary, len(q.Wld.AmbientLights), helper.Pluralize(len(q.Wld.AmbientLights)))
	}
	if len(q.Wld.AniDefs) > 0 {
		summary = fmt.Sprintf("%s%d animation%s, ", summary, len(q.Wld.AniDefs), helper.Pluralize(len(q.Wld.AniDefs)))
	}
	if len(q.Wld.TrackDefs) > 0 {
		summary = fmt.Sprintf("%s%d track%s, ", summary, len(q.Wld.TrackDefs), helper.Pluralize(len(q.Wld.TrackDefs)))
	}
	if len(q.Wld.DMTrackDef2s) > 0 {
		summary = fmt.Sprintf("%s%d dmtrack%s, ", summary, len(q.Wld.DMTrackDef2s), helper.Pluralize(len(q.Wld.DMTrackDef2s)))
	}
	if len(q.Wld.PrtDefs) > 0 {
		summary = fmt.Sprintf("%s%d particle render%s, ", summary, len(q.Wld.PrtDefs), helper.Pluralize(len(q.Wld.PrtDefs)))
	}
	if len(q.Wld.PtsDefs) > 0 {
		summary = fmt.Sprintf("%s%d particle point%s, ", summary, len(q.Wld.PtsDefs), helper.Pluralize(len(q.Wld.PtsDefs)))
	}
	if len(q.Wld.ParticleCloudDefs) > 0 {
		summary = fmt.Sprintf("%s%d particle cloud%s, ", summary, len(q.Wld.ParticleCloudDefs), helper.Pluralize(len(q.Wld.ParticleCloudDefs)))
	}
	if len(q.Assets) > 0 {
		summary = fmt.Sprintf("%s%d asset%s, ", summary, len(q.Assets), helper.Pluralize(len(q.Assets)))
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
