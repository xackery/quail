package geo

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
