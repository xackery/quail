package common

// Fragment is what every fragment object type adheres to
type WldFragmenter interface {
	// FragmentType identifies the fragment type
	FragmentType() string
	Data() []byte
}
