package mds

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model/geo"
)

// BlenderImport imports a blender structure to MDS
func (e *MDS) BlenderImport(dir string) error {
	e.version = 1
	e.MaterialManager = &geo.MaterialManager{}
	e.meshManager = &geo.MeshManager{}
	e.particleManager = &geo.ParticleManager{}
	path := dir

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat _%s: %w", e.Name(), err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	var r *os.File
	var scanner *bufio.Scanner
	var lineNumber int

	curPath := fmt.Sprintf("%s/info.txt", path)
	if helper.IsFile(curPath) {
		r, err = os.Open(curPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", curPath, err)
		}
		scanner = bufio.NewScanner(r)
		lineNumber = 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.Split(line, "=")
			if len(parts) < 2 {
				return fmt.Errorf("invalid version.txt (expected 2 records) line %d: %s", lineNumber, line)
			}
			switch parts[0] {
			case "version":
				e.version = helper.AtoU32(parts[1])
			default:
				return fmt.Errorf("invalid info.txt line %d: %s", lineNumber, line)
			}
		}
	}

	curPath = fmt.Sprintf("%s/material.txt", path)
	if helper.IsFile(curPath) {
		err = e.MaterialManager.BlenderImport(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", curPath, err)
		}
	}

	curPath = fmt.Sprintf("%s/particle_point.txt", path)
	if helper.IsFile(curPath) {
		err = e.MaterialManager.BlenderImport(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", curPath, err)
		}
	}

	err = e.meshManager.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", curPath, err)
	}

	return nil
}
