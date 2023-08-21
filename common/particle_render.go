package common

type ParticleRender struct {
	Version int
	Name    string
	Entries []*ParticleRenderEntry
}

// ParticleRender defines what particle to emit on a particle point
type ParticleRenderEntry struct {
	ID                  uint32 `json:"id"` //id is actorsemittersnew.edd
	ID2                 uint32 `json:"id2"`
	ParticlePoint       string `json:"particlePoint"`
	ParticlePointSuffix []byte
	UnknownA1           uint32
	UnknownA2           uint32
	UnknownA3           uint32
	UnknownA4           uint32
	UnknownA5           uint32
	Duration            uint32 `json:"duration"`
	UnknownB            uint32 `json:"unknownb"`
	UnknownFFFFFFFF     int32  `json:"unknownffffffff"`
	UnknownC            uint32 `json:"unknownc"`
}
