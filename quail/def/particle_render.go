package def

type ParticleRender struct {
	Name    string
	Entries []*ParticleRenderEntry
}

// ParticleRender defines what particle to emit on a particle point
type ParticleRenderEntry struct {
	ID                  uint32 `json:"id"` //id is actorsemittersnew.edd
	ID2                 uint32 `json:"id2"`
	ParticlePoint       string `json:"particlePoint"`
	ParticlePointSuffix []byte
	UnknownA            [5]uint32 `json:"unknowna"` //Pretty sure last 3 have something to do with durations
	Duration            uint32    `json:"duration"`
	UnknownB            uint32    `json:"unknownb"`
	UnknownFFFFFFFF     int32     `json:"unknownffffffff"`
	UnknownC            uint32    `json:"unknownc"`
}
