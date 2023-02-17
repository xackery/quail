package zon

import "fmt"

func (e *ZON) Inspect() {
	fmt.Printf("%d objects\n", len(e.Objects()))
	for i, object := range e.Objects() {
		fmt.Printf("	%d %+v\n", i, object)
	}

	fmt.Printf("%d models\n", len(e.Models()))
	for i, model := range e.Models() {
		fmt.Printf("	%d %+v\n", i, model)
	}

	fmt.Printf("%d lights\n", len(e.Lights()))
	for i, light := range e.Lights() {
		fmt.Printf("	%d %+v\n", i, light)
	}

}
