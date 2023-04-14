package ter

// Inspect prints out details
func (e *TER) Inspect() {
	e.MaterialManager.Inspect()
	e.meshManager.Inspect()
	e.particleManager.Inspect()
}
