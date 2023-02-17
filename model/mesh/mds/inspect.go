package mds

import "fmt"

// Inspect prints out details
func (e *MDS) Inspect() {
	fmt.Println(len(e.materials), "materials:")
	for i, material := range e.materials {
		fmt.Printf("	%d %s\n", i, material.Name)
	}
}
