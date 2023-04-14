package zon

import "fmt"

func (e *ZON) Inspect() {
	e.objectManager.Inspect()

	fmt.Printf("%d models\n", len(e.Models()))
	for i, model := range e.Models() {
		fmt.Printf("	%d %+v\n", i, model)
	}

	fmt.Printf("%d lights\n", len(e.Lights()))
	for i, light := range e.Lights() {
		fmt.Printf("	%d %+v\n", i, light)
	}

}
