//go:build !wasm && !tinygo.wasm

package os

import (
	"os"
)

var (
	Args        = os.Args
	Stdin       = os.Stdin
	Stdout      = os.Stdout
	ErrNotExist = os.ErrNotExist
	ModePerm    = os.ModePerm
	IsExist     = os.IsExist
)

type (
	File     = os.File
	DirEntry = os.DirEntry
)

var (
	Stat       = os.Stat
	ReadDir    = os.ReadDir
	WriteFile  = os.WriteFile
	ReadFile   = os.ReadFile
	Open       = os.Open
	MkdirAll   = os.MkdirAll
	Create     = os.Create
	IsNotExist = os.IsNotExist
	RemoveAll  = os.RemoveAll
	Getwd      = os.Getwd
	Remove     = os.Remove
	Exit       = os.Exit
	Getenv     = os.Getenv
)
