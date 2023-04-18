package pts

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/geo"
)

func (e *PTS) Decode(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	header := dec.StringFixed(4)
	if header != "EQPT" {
		return fmt.Errorf("header does not match EQPT, got %s", header)
	}

	particleCount := dec.Uint32()
	dump.Hex(particleCount, "particleCount=%d", particleCount)
	version := dec.Uint32()
	dump.Hex(version, "version=%d", version)
	if version != 1 {
		return fmt.Errorf("version is %d, wanted 1", version)
	}
	for i := 0; i < int(particleCount); i++ {
		pp := geo.ParticlePoint{}
		pp.Name = dec.StringFixed(64)
		dump.Hex(pp.Name, "%dname=%s", i, pp.Name)
		pp.Bone = dec.StringFixed(64)
		dump.Hex(pp.Bone, "%dbone=%s", i, pp.Bone)
		pp.Translation.X = dec.Float32()
		pp.Translation.Y = dec.Float32()
		pp.Translation.Z = dec.Float32()
		dump.Hex(pp.Translation, "%dtranslation=%0.0f", i, pp.Translation)
		pp.Rotation.X = dec.Float32()
		pp.Rotation.Y = dec.Float32()
		pp.Rotation.Z = dec.Float32()
		dump.Hex(pp.Rotation, "%drotation=%0.0f", i, pp.Rotation)
		pp.Scale.X = dec.Float32()
		pp.Scale.Y = dec.Float32()
		pp.Scale.Z = dec.Float32()
		dump.Hex(pp.Scale, "%dscale=%0.0f", i, pp.Scale)
		e.particleManager.PointAdd(pp)
	}
	if dec.Error() != nil {
		return fmt.Errorf("pts decode: %w", dec.Error())
	}

	return nil
}
