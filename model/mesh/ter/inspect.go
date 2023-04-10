package ter

import "fmt"

// Inspect prints out details
func (e *TER) Inspect() {
	fmt.Println(len(e.materials), "materials:")
	for i, material := range e.materials {
		fmt.Printf("  %d %s\n", i, material.Name)
	}

	fmt.Println(len(e.triangles), "triangles")
	fmt.Println(len(e.vertices), "vertices")
	fmt.Println(len(e.particlePoints), "particle points")
	for i, particle := range e.particlePoints {
		fmt.Printf("  %d %s %s\n", i, particle.Bone, particle.Name)
	}
	fmt.Println(len(e.particleRenders), "particle renders")
	for i, particle := range e.particleRenders {
		fmt.Printf("  %d %s %d\n", i, particle.ParticlePoint, particle.ID)
	}
}
