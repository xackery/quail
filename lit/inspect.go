package lit

import "fmt"

// Inspect prints out details
func (e *LIT) Inspect() {
	fmt.Println(len(e.lights), "lights:")
	for i, light := range e.lights {
		fmt.Printf("	%d %0.2f\n", i, light)
	}
}
