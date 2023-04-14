package geo

import (
	"fmt"
	"io"
)

// ParticleRender defines what particle to emit on a particle point
type ParticleRender struct {
	ID              uint32    `json:"id"` //id is actorsemittersnew.edd
	ID2             uint32    `json:"id2"`
	ParticlePoint   string    `json:"particlePoint"`
	UnknownA        [5]uint32 `json:"unknowna"` //Pretty sure last 3 have something to do with durations
	Duration        uint32    `json:"duration"`
	UnknownB        uint32    `json:"unknownb"`
	UnknownFFFFFFFF int32     `json:"unknownffffffff"`
	UnknownC        uint32    `json:"unknownc"`
}

// NewParticleRender returns a new particle render
func NewParticleRender() *ParticleRender {
	return &ParticleRender{}
}

// WriteHeader writes the header for a ParticleRender
func (e *ParticleRender) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("id|id2|particlePoint|unknowna|duration|unknownb|unknownffffffff|unknownc\n")
	return err
}

// Write writes a ParticleRender to a string writer
func (e *ParticleRender) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%d|%d|%s|%d|%d|%d|%d|%d|%d\n",
		e.ID,
		e.ID2,
		e.ParticlePoint,
		e.UnknownA[0],
		e.UnknownA[1],
		e.UnknownA[2],
		e.UnknownA[3],
		e.UnknownA[4],
		e.Duration,
		e.UnknownB,
		e.UnknownFFFFFFFF,
		e.UnknownC))
	return err
}
