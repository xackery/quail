package zon

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/g3n/engine/math32"
)

// Import takes paths for various files and attempts to import them into a zon definition
func (e *ZON) Import(lightPath string, modPath string, regionPath string) error {
	err := e.importLight(lightPath)
	if err != nil {
		return fmt.Errorf("importLight: %w", err)
	}
	err = e.importMod(modPath)
	if err != nil {
		return fmt.Errorf("importMod: %w", err)
	}
	err = e.importRegion(regionPath)
	if err != nil {
		return fmt.Errorf("importRegion: %w", err)
	}
	return nil
}

func (e *ZON) importLight(lightPath string) error {
	r, err := os.Open(lightPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil
		}
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		records := strings.Split(line, " ")
		if len(records) != 8 {
			return fmt.Errorf("expected 8 arguments, got %d", len(records))
		}
		position := math32.Vector3{}
		val, err := strconv.ParseFloat(records[1], 32)
		if err != nil {
			return fmt.Errorf("parse pos x: %w", err)
		}
		position.X = float32(val)

		val, err = strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse pos y: %w", err)
		}
		position.Y = float32(val)

		val, err = strconv.ParseFloat(records[3], 32)
		if err != nil {
			return fmt.Errorf("parse pos z: %w", err)
		}
		position.Z = float32(val)

		color := math32.Color{}
		val, err = strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse color r: %w", err)
		}
		color.R = float32(val)

		val, err = strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse color g: %w", err)
		}
		color.G = float32(val)

		val, err = strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse color b: %w", err)
		}
		color.B = float32(val)

		val, err = strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse color b: %w", err)
		}

		err = e.AddLight(records[0], position, color, float32(val))
		if err != nil {
			return fmt.Errorf("addLight line %d: %w", lineNumber, err)
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read light %s: %w", lightPath, err)
	}
	return nil
}

func (e *ZON) importMod(modPath string) error {
	r, err := os.Open(modPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil
		}
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()
		records := strings.Split(line, " ")
		if len(records) != 9 {
			return fmt.Errorf("expected 9 arguments, got %d", len(records))
		}

		err = e.AddModel(records[0])
		if err != nil {
			return fmt.Errorf("addModel: %w", err)
		}

		position := math32.Vector3{}
		val, err := strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse pos x: %w", err)
		}
		position.X = float32(val)

		val, err = strconv.ParseFloat(records[3], 32)
		if err != nil {
			return fmt.Errorf("parse pos y: %w", err)
		}
		position.Y = float32(val)

		val, err = strconv.ParseFloat(records[4], 32)
		if err != nil {
			return fmt.Errorf("parse pos z: %w", err)
		}
		position.Z = float32(val)
		rotation := math32.Vector3{}
		val, err = strconv.ParseFloat(records[5], 32)
		if err != nil {
			return fmt.Errorf("parse rotation x: %w", err)
		}
		rotation.X = float32(val)

		val, err = strconv.ParseFloat(records[6], 32)
		if err != nil {
			return fmt.Errorf("parse rotation y: %w", err)
		}
		rotation.Y = float32(val)

		val, err = strconv.ParseFloat(records[7], 32)
		if err != nil {
			return fmt.Errorf("parse rotation z: %w", err)
		}
		rotation.Z = float32(val)

		err = e.AddObject(records[0], records[1], position, rotation, float32(val))
		if err != nil {
			return fmt.Errorf("addObject line %d: %w", lineNumber, err)
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read mod %s: %w", modPath, err)
	}
	return nil
}

func (e *ZON) importRegion(regionPath string) error {
	r, err := os.Open(regionPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil
		}
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		records := strings.Split(line, " ")
		if len(records) != 10 {
			return fmt.Errorf("expected 10 arguments, got %d", len(records))
		}

		center := math32.Vector3{}
		val, err := strconv.ParseFloat(records[2], 32)
		if err != nil {
			return fmt.Errorf("parse center x: %w", err)
		}
		center.X = float32(val)

		val, err = strconv.ParseFloat(records[3], 32)
		if err != nil {
			return fmt.Errorf("parse center y: %w", err)
		}
		center.Y = float32(val)

		val, err = strconv.ParseFloat(records[4], 32)
		if err != nil {
			return fmt.Errorf("parse center z: %w", err)
		}
		center.Z = float32(val)
		extent := math32.Vector3{}
		val, err = strconv.ParseFloat(records[5], 32)
		if err != nil {
			return fmt.Errorf("parse extent x: %w", err)
		}
		extent.X = float32(val)

		val, err = strconv.ParseFloat(records[6], 32)
		if err != nil {
			return fmt.Errorf("parse extent y: %w", err)
		}
		extent.Y = float32(val)

		val, err = strconv.ParseFloat(records[7], 32)
		if err != nil {
			return fmt.Errorf("parse extent z: %w", err)
		}
		extent.Z = float32(val)

		unknown := math32.Vector3{}
		val, err = strconv.ParseFloat(records[8], 32)
		if err != nil {
			return fmt.Errorf("parse unknown x: %w", err)
		}
		unknown.X = float32(val)

		val, err = strconv.ParseFloat(records[9], 32)
		if err != nil {
			return fmt.Errorf("parse unknown y: %w", err)
		}
		unknown.Y = float32(val)

		err = e.AddRegion(records[0], center, unknown, extent)
		if err != nil {
			return fmt.Errorf("addObject line %d: %w", lineNumber, err)
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read region %s: %w", regionPath, err)
	}
	return nil
}
