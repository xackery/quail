package wad

import (
	"fmt"
	"io"
	"strings"
)

type wadWriter struct {
	w io.Writer
}

func (e *Wad) WADWrite(writer io.Writer) error {
	w := &wadWriter{w: writer}

	_, err := writer.Write([]byte(`CPIWORLD        "SONY3D World"` + "\n"))
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	w.writeLn(1, fmt.Sprintf("NUMBITMAPINFO %d", len(e.bitmaps)))
	for i, bitmap := range e.bitmaps {
		w.writeLn(1, fmt.Sprintf("BITMAPINFO %d", i+1))

		w.writeLn(2, fmt.Sprintf("MATERIALTAG %s", bitmap.Tag))

		w.writeLn(2, "")
		w.writeLn(2, fmt.Sprintf("NUMTEXTURES %d", len(bitmap.Textures)))
		for j, texture := range bitmap.Textures {
			w.writeLn(2, fmt.Sprintf("TEXTURE %d", j+1))
			w.writeLn(3, fmt.Sprintf("TEXTUREFILE %s", texture))
			w.writeLn(2, fmt.Sprintf("ENDTEXTURE %d", j+1))
		}

		w.writeLn(1, fmt.Sprintf("ENDBITMAPINFO %d", i+1))
	}

	w.writeLn(1, "")
	w.writeLn(1, fmt.Sprintf("NUM2DSPRITEDEF %d", len(e.sprites)))
	for j, sprite := range e.sprites {
		w.writeLn(1, fmt.Sprintf("2DSPRITEDEF %d", j+1))
		w.writeLn(2, fmt.Sprintf("TAG %s", sprite.Tag))
		w.writeLn(2, fmt.Sprintf("FLAGS %d", sprite.Flags))
		w.writeLn(2, fmt.Sprintf("CURRENTFRAME %d", sprite.CurrentFrame))
		w.writeLn(2, fmt.Sprintf("SLEEP %d", sprite.Sleep))

		w.writeLn(2, fmt.Sprintf("NUMFRAMES %d", len(sprite.Frames)))
		for k, frame := range sprite.Frames {
			w.writeLn(3, fmt.Sprintf("FRAME %d", k+1))
			w.writeLn(4, fmt.Sprintf("TAG %s", frame.Tag))
			w.writeLn(3, fmt.Sprintf("ENDFRAME %d", k+1))
		}

		w.writeLn(2, fmt.Sprintf("NUM2DSPRITEINST %d", len(sprite.Instances)))
		for k, inst := range sprite.Instances {
			w.writeLn(2, fmt.Sprintf("2DSPRITEINST %d", k+1))
			w.writeLn(3, fmt.Sprintf("2DSPRITETAG %s", inst.Tag))
			w.writeLn(3, fmt.Sprintf("FLAGS %d", inst.Flags))

			w.writeLn(3, "")
			w.writeLn(3, fmt.Sprintf("NUMMATERIAL %d", len(inst.Materials)))
			for j, mat := range inst.Materials {
				w.writeLn(3, fmt.Sprintf("MATERIAL %d", j+1))
				w.writeLn(4, fmt.Sprintf("MATERIALTAG %s", mat.Tag))
				w.writeLn(4, fmt.Sprintf("FLAGS %d", mat.Flags))
				w.writeLn(4, fmt.Sprintf("RENDERMETHOD %d", mat.RenderMethod))
				w.writeLn(4, fmt.Sprintf("RGBPEN %d", mat.RGBPen))
				w.writeLn(4, fmt.Sprintf("BRIGHTNESS %f", mat.Brightness))
				w.writeLn(4, fmt.Sprintf("SCALEDAMBIENT %f", mat.ScaledAmbient))
				w.writeLn(4, fmt.Sprintf("PAIRS %d, %d", mat.Pairs[0], mat.Pairs[1]))

				w.writeLn(3, fmt.Sprintf("ENDMATERIAL %d", j+1))
			}

			w.writeLn(2, fmt.Sprintf("END2DSPRITEINST %d", k+1))
		}

		w.writeLn(1, fmt.Sprintf("END2DSPRITEDEF %d", j+1))
	}

	w.writeLn(1, "")
	w.writeLn(1, fmt.Sprintf("NUMDMSPRITEINFO %d", len(e.dmsprites)))
	for i, dmsprite := range e.dmsprites {
		w.writeLn(1, fmt.Sprintf("DMSPRITEINFO %d", i))
		w.writeLn(2, fmt.Sprintf("DMSPRITETAG %s", dmsprite.Tag))
		w.writeLn(2, fmt.Sprintf("FLAGS %d", dmsprite.Flags))

		// track
		w.writeLn(2, "")
		w.writeLn(2, fmt.Sprintf("NUMTRACK %d", len(dmsprite.Tracks)))
		for j, track := range dmsprite.Tracks {
			w.writeLn(2, fmt.Sprintf("TRACK %d", j))
			w.writeLn(3, fmt.Sprintf("TRACKTAG %s", track.Tag))
			w.writeLn(3, fmt.Sprintf("FLAGS %d", track.Flags))
			w.writeLn(3, fmt.Sprintf("SCALE %d", track.Scale))
			w.writeLn(3, fmt.Sprintf("SIZE %d", track.Size))

			w.writeLn(3, fmt.Sprintf("NUMTRACKFRAMES %d", len(track.Frames)))
			for _, frame := range track.Frames {
				w.writeLn(4, fmt.Sprintf("XYZ %d, %d, %d", frame[0], frame[1], frame[2]))
			}

			w.writeLn(2, fmt.Sprintf("ENDTRACK %d", j))

		}

		w.writeLn(2, fmt.Sprintf("NUMUVS %d", len(dmsprite.UVs)))
		for _, uv := range dmsprite.UVs {
			w.writeLn(3, fmt.Sprintf("UV %d, %d", uv[0], uv[1]))
		}

		w.writeLn(2, fmt.Sprintf("NUMNORMALS %d", len(dmsprite.Normals)))
		for _, normal := range dmsprite.Normals {
			w.writeLn(3, fmt.Sprintf("NORMAL %d, %d, %d", normal[0], normal[1], normal[2]))
		}

		w.writeLn(1, fmt.Sprintf("ENDDMSPRITEINFO %d", i))
	}

	// 49 DmSpriteDef2 SED_DMSPRITEDEF
	// 50 TrackDef SED_TRACKDEF
	// 51 Track SED_TRACK
	// 52 TrackDef SEDPE_TRACKDEF

	return nil
}

func (e *wadWriter) writeLn(indent int, s string) error {
	var err error
	if s == "" {
		_, err = e.w.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
		return nil
	}
	_, err = e.w.Write([]byte(strings.Repeat("  ", indent) + s + "\n"))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
