package geo

import "io"

// ParticlePoint defines a location relative to a bone as where a particle should emit
type ParticlePoint struct {
	Name        string   `json:"name"`
	Bone        string   `json:"bone"`
	Translation *Vector3 `json:"translation"`
	Rotation    *Vector3 `json:"rotation"`
	Scale       *Vector3 `json:"scale"`
}

// NewParticlePoint returns a new particle point
func NewParticlePoint() *ParticlePoint {
	return &ParticlePoint{
		Translation: &Vector3{},
		Rotation:    &Vector3{},
		Scale:       &Vector3{},
	}
}

// WriteHeader writes the header for a ParticlePoint
func (e *ParticlePoint) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("name|bone|translation|rotation|scale\n")
	return err
}

// Write writes a ParticlePoint to a string writer
func (e *ParticlePoint) Write(w io.StringWriter) error {
	_, err := w.WriteString(e.Name + "|" + e.Bone + "|" + e.Translation.String() + "|" + e.Rotation.String() + "|" + e.Scale.String() + "\n")
	return err
}
