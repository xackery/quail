package geo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
)

// ParticleManager is a particle manager
type ParticleManager struct {
	renders []ParticleRender
	points  []ParticlePoint
}

// NewParticleManager creates a new particle manager
func NewParticleManager() *ParticleManager {
	return &ParticleManager{}
}

// WriteFile writes all materials to a file
func (e *ParticleManager) BlenderExport(path string) error {
	if len(e.renders) > 0 {
		renderPath := fmt.Sprintf("%s/particle_render.txt", path)
		pw, err := os.Create(renderPath)
		if err != nil {
			return fmt.Errorf("create file %s: %w", renderPath, err)
		}
		defer pw.Close()

		pr := &ParticleRender{}
		err = pr.WriteHeader(pw)
		if err != nil {
			return fmt.Errorf("write particle render: %w", err)
		}

		for _, pr := range e.renders {
			err = pr.Write(pw)
			if err != nil {
				return fmt.Errorf("write particle render: %w", err)
			}
		}
	}

	if len(e.points) > 0 {
		pointPath := fmt.Sprintf("%s/particle_point.txt", path)
		mw, err := os.Create(pointPath)
		if err != nil {
			return fmt.Errorf("create particle point %s: %w", pointPath, err)
		}
		defer mw.Close()
		pp := &ParticlePoint{}
		err = pp.WriteHeader(mw)
		if err != nil {
			return fmt.Errorf("write header: %w", err)
		}

		for _, pp := range e.points {
			err = pp.Write(mw)
			if err != nil {
				return fmt.Errorf("write particle point: %w", err)
			}

		}
	}
	return nil
}

// ReadFile reads a material file
func (e *ParticleManager) ReadFile(pointPath string, renderPath string) error {
	r, err := os.Open(pointPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", pointPath, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			return fmt.Errorf("invalid particle points (expected 5 records) line %d: %s", lineNumber, line)
		}

		e.points = append(e.points, ParticlePoint{
			Name:        parts[0],
			Bone:        parts[1],
			Translation: AtoVector3(parts[2]),
			Rotation:    AtoVector3(parts[3]),
			Scale:       AtoVector3(parts[4]),
		})
	}
	r.Close()

	r, err = os.Open(renderPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", renderPath, err)
	}
	scanner = bufio.NewScanner(r)
	lineNumber = 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 7 {
			return fmt.Errorf("invalid particle_render.txt (expected 7 records) line %d: %s", lineNumber, line)
		}
		e.renders = append(e.renders, ParticleRender{
			Duration:      helper.AtoU32(parts[0]),
			ID:            helper.AtoU32(parts[1]),
			ID2:           helper.AtoU32(parts[2]),
			ParticlePoint: parts[3],
			//UnknownA:        helper.AtoU32(parts[4]),
			UnknownB: helper.AtoU32(parts[5]),
			//UnknownFFFFFFFF: helper.AtoU32(parts[6]),
		})
		return fmt.Errorf("todo: blender import fix for particles")
	}
	r.Close()

	return nil
}

// Inspect prints the particles
func (e *ParticleManager) Inspect() {
	fmt.Println(len(e.points), "particle points")
	for i, particle := range e.points {
		fmt.Printf("  %d %s %s\n", i, particle.Bone, particle.Name)
	}
	fmt.Println(len(e.renders), "particle renders")
	for i, particle := range e.renders {
		fmt.Printf("  %d %s %d\n", i, particle.ParticlePoint, particle.ID)
	}
}

// PointCount returns the number of points
func (e *ParticleManager) PointCount() int {
	return len(e.points)
}

// RenderCount returns the number of renders
func (e *ParticleManager) RenderCount() int {
	return len(e.renders)
}

// PointAdd adds a point
func (e *ParticleManager) PointAdd(point ParticlePoint) {
	e.points = append(e.points, point)
}

// RenderAdd adds a render
func (e *ParticleManager) RenderAdd(render ParticleRender) {
	e.renders = append(e.renders, render)
}
