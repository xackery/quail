package common

type ParticleRender struct {
	Header  *Header                `yaml:"header,omitempty"`
	Entries []*ParticleRenderEntry `yaml:"entries,omitempty"`
}

func NewParticleRender(name string) *ParticleRender {
	return &ParticleRender{
		Header: &Header{
			Name: name,
		},
	}
}

// ParticleRender defines what particle to emit on a particle point
type ParticleRenderEntry struct {
	ID                  uint32 `yaml:"id"` //id is actorsemittersnew.edd
	ID2                 uint32 `yaml:"id2"`
	ParticlePoint       string `yaml:"particle_point"`
	ParticlePointSuffix []byte `yaml:"particle_point_suffix,omitempty"`
	UnknownA1           uint32 `yaml:"unknowna1"`
	UnknownA2           uint32 `yaml:"unknowna2"`
	UnknownA3           uint32 `yaml:"unknowna3"`
	UnknownA4           uint32 `yaml:"unknowna4"`
	UnknownA5           uint32 `yaml:"unknowna5"`
	Duration            uint32 `yaml:"duration"`
	UnknownB            uint32 `yaml:"unknownb"`
	UnknownFFFFFFFF     int32  `yaml:"unknownffffffff"`
	UnknownC            uint32 `yaml:"unknownc"`
}
