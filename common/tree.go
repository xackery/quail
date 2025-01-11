package common

// TreeLinker is an interface for linking tree nodes
type TreeLinker interface {
	Tag() string
	FragID() int
	Children() []TreeLinker
	FragType() string
	AddParent(parent TreeLinker)
	//RemoveParent(parent TreeLinker) error
	//ListParents() []TreeLinker
}
