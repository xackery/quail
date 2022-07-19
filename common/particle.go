package common

type ParticleEntry struct {
	Id              uint32    `json:"id"` //id is actorsemittersnew.edd
	Id2             uint32    `json:"id2"`
	Bone            string    `json:"bone"`
	UnknownA        [5]uint32 `json:"unknowna"` //Pretty sure last 3 have something to do with durations
	Duration        uint32    `json:"duration"`
	UnknownB        uint32    `json:"unknownb"`
	UnknownFFFFFFFF int32     `json:"unknownffffffff"`
	UnknownC        uint32    `json:"unknownc"`
}
