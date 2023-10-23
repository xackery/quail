package common

// ParticlePoint is a particle point
type ParticlePoint struct {
	Header  *Header
	Entries []ParticlePointEntry
}

func NewParticlePoint(name string) *ParticlePoint {
	return &ParticlePoint{
		Header: &Header{
			Name: name,
		},
	}
}

func (p *ParticlePoint) init() {
	if p.Header == nil {
		p.Header = &Header{}
	}
}

// ParticlePointEntry is a single entry in a particle point
type ParticlePointEntry struct {
	Name        string `json:"name"`
	NameSuffix  []byte
	Bone        string `json:"bone"`
	BoneSuffix  []byte
	Translation Vector3 `json:"translation"`
	Rotation    Vector3 `json:"rotation"`
	Scale       Vector3 `json:"scale"`
}
