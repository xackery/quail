package lod

import "fmt"

// Inspect prints out details
func (e *LOD) Inspect() {
	fmt.Println(len(e.lods), "lods:")
	for i, le := range e.lods {
		fmt.Printf("	%d %s %s %d\n", i, le.Category, le.ObjectName, le.Distance)
	}
}
