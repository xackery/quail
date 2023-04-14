package mds

// Inspect prints out details
func (e *MDS) Inspect() {
	e.MaterialManager.Inspect()
	e.meshManager.Inspect()
	e.particleManager.Inspect()
}
