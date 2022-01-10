package zon

import "github.com/g3n/engine/math32"

// ZON is a zon file struct
type ZON struct {
	models  []*model
	objects []*object
	regions []*region
	lights  []*light
}

type model struct {
	name string
}

type object struct {
	modelName string
	name      string
	position  math32.Vector3
	rotation  math32.Vector3
	scale     float32
}

type region struct {
	name    string
	center  math32.Vector3
	unknown math32.Vector3
	extent  math32.Vector3
}

type light struct {
	name     string
	position math32.Vector3
	color    math32.Color
	radius   float32
}
