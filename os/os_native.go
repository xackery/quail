//go:build !tinygo.wasm

package os

import (
	"io/fs"
	nativeos "os"
)

type File = nativeos.File

var (
	ErrNotExist = nativeos.ErrNotExist
	ModePerm    = nativeos.ModePerm
)

func Stat(name string) (nativeos.FileInfo, error) {
	return nativeos.Stat(name)
}

func ReadDir(name string) ([]nativeos.DirEntry, error) {
	return nativeos.ReadDir(name)
}

func WriteFile(name string, buffer []byte, perm fs.FileMode) error {
	return nativeos.WriteFile(name, buffer, perm)
}

func ReadFile(name string) ([]byte, error) {
	return nativeos.ReadFile(name)
}

func Open(name string) (*nativeos.File, error) {
	return nativeos.Open(name)
}

func MkdirAll(path string, perm fs.FileMode) error {
	return nativeos.MkdirAll(path, perm)
}

func Create(name string) (*nativeos.File, error) {
	return nativeos.Create(name)
}

func IsNotExist(err error) bool {
	return nativeos.IsNotExist(err)
}

func RemoveAll(path string) error {
	return nativeos.RemoveAll(path)
}

func Getwd() (string, error) {
	return nativeos.Getwd()
}

func Remove(name string) error {
	return nativeos.Remove(name)
}

func Exit(code int) {
	nativeos.Exit(code)
}

func Getenv(key string) string {
	return nativeos.Getenv(key)
}
