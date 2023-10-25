package quail

import (
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

type Quail struct {
	Header                 *common.Header
	Models                 []*common.Model
	Animations             []*common.Animation
	Zone                   *common.Zone
	materialCache          map[string]*common.Material
	IsExtensionVersionDump bool
}

// New returns a new Quail instance
func New() *Quail {
	return &Quail{
		materialCache: make(map[string]*common.Material),
	}
}

// Close flushes any memory and closes any open files
func (e *Quail) Close() error {
	e.Models = nil
	e.Animations = nil
	e.Zone = nil
	e.materialCache = make(map[string]*common.Material)
	return nil
}

// SetLogLevel sets the log level
func SetLogLevel(level int) {
	log.SetLogLevel(level)
}
