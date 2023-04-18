package geo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
)

// ObjectManager is an object manager
type ObjectManager struct {
	objects []Object
}

// Object returns an object by name
func (e ObjectManager) Object(name string) (Object, bool) {
	for _, object := range e.objects {
		if object.Name == name {
			return object, true
		}
	}
	return Object{}, false
}

// BlenderExport writes all materials to a file
func (e ObjectManager) BlenderExport(objectPath string) error {

	if len(e.objects) > 0 {
		ow, err := os.Create(objectPath)
		if err != nil {
			return fmt.Errorf("create file %s: %w", objectPath, err)
		}
		defer ow.Close()
		object := &Object{}
		err = object.WriteHeader(ow)
		if err != nil {
			return fmt.Errorf("write header: %w", err)
		}

		for _, o := range e.objects {
			err = o.Write(ow)
			if err != nil {
				return fmt.Errorf("write object: %w", err)
			}
		}
	}

	return nil
}

// ReadFile reads all objects from a file
func (e ObjectManager) ReadFile(objectPath string) error {
	r, err := os.Open(objectPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", objectPath, err)
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
		if len(parts) < 3 {
			return fmt.Errorf("invalid triangle.txt (expected 4 records) line %d: %s", lineNumber, line)
		}

		e.objects = append(e.objects, Object{
			Name:      parts[0],
			ModelName: parts[1],
			Position:  AtoVector3(parts[2]),
			Rotation:  AtoVector3(parts[3]),
			Scale:     helper.AtoF32(parts[4]),
			FileType:  parts[5],
			FileName:  parts[6],
		})
	}
	r.Close()

	return nil
}

// Add adds an object to the manager
func (e ObjectManager) Add(object Object) {
	e.objects = append(e.objects, object)
}

// Inspect prints out all objects
func (e ObjectManager) Inspect() {
	fmt.Printf("%d objects\n", len(e.objects))
	for i, object := range e.objects {
		fmt.Printf("	%d %+v\n", i, object)
	}
}

// Objects returns all objects
func (e ObjectManager) Objects() []Object {
	return e.objects
}

// Count returns the number of objects
func (e ObjectManager) Count() int {
	return len(e.objects)
}
