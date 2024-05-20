package bsp

type BspTree struct {
	nodeCount int
	nodes     []*BspNode
}

type BspNode struct {
	parent         *BspNode
	leftChild      *BspNode
	rightChild     *BspNode
	isLeftChild    bool
	isRightChild   bool
	normalX        int
	normalY        int
	normalZ        int
	splitDistance  int
	regionID       int
	left           int
	right          int
	boundingBoxMin [3]int
	boundingBoxMax [3]int
	center         [3]int
}
