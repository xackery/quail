package def

// ParticlePoint is a particle point
type ParticlePoint struct {
	Name    string
	Entries []ParticlePointEntry
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
