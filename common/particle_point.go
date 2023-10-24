package common

// ParticlePoint is a particle point
type ParticlePoint struct {
	Header  *Header              `yaml:"header,omitempty"`
	Entries []ParticlePointEntry `yaml:"entries,omitempty"`
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
	Name        string  `yaml:"name"`
	BoneName    string  `yaml:"bone_name"`
	Translation Vector3 `yaml:"translation"`
	Rotation    Vector3 `yaml:"rotation"`
	Scale       Vector3 `yaml:"scale"`
	NameSuffix  []byte  `yaml:"name_suffix,omitempty"`
	BoneSuffix  []byte  `yaml:"bone_suffix,omitempty"`
}
