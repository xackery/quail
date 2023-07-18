package def

type ParticlePoint struct {
	Name    string
	Entries []*ParticlePointEntry
}

type ParticlePointEntry struct {
	Name        string `json:"name"`
	NameSuffix  []byte
	Bone        string `json:"bone"`
	BoneSuffix  []byte
	Translation Vector3 `json:"translation"`
	Rotation    Vector3 `json:"rotation"`
	Scale       Vector3 `json:"scale"`
}
