package mod

// Inspect prints out details
func (e *MOD) Inspect() {
	e.MaterialManager.Inspect()
	e.meshManager.Inspect()
	e.particleManager.Inspect()
}
