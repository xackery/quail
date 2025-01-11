package tree

import (
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/raw"
)

type Node struct {
	FragID   int32
	FragType string
	Tag      string
	Children map[int32]*Node
}

// Dump dumps a tree to a writer
func Dump(src interface{}, w io.Writer) error {
	var err error
	switch val := src.(type) {
	case *raw.Wld:
		err = wldDump(val, w)
		if err != nil {
			return fmt.Errorf("wld dump: %w", err)
		}
	case *raw.Bmp:
		return nil
	default:
		return fmt.Errorf("unknown type %T", val)
	}
	return nil
}

func BuildFragReferenceTree(wld *raw.Wld) (map[int32]*Node, error) {
	// Map to store all nodes
	nodes := make(map[int32]*Node)

	// Map to track which nodes are linkedNodes as children
	linkedNodes := make(map[int32]bool)

	// Create nodes and establish relationships
	for i := 1; i < len(wld.Fragments); i++ { // Start at index 1
		frag := wld.Fragments[i]
		if frag == nil {
			continue
		}

		fragID := int32(i)
		tag := ""
		if frag.NameRef() < 0 {
			tag = wld.Name(frag.NameRef())
		}

		// Find or create the node for this fragment
		node := upsertNode(nodes, fmt.Sprintf("%T", frag), fragID, strings.TrimSpace(tag))

		// Extract references from the fragment
		fragRefs := fragRefs(frag)
		for _, refID := range fragRefs {
			if refID <= 0 {
				continue
			}
			// Mark this refID as being referenced
			linkedNodes[refID] = true

			// Find or create the child node
			child := upsertNode(nodes, fmt.Sprintf("%T", frag), refID, tag)

			// Establish the parent-child relationship
			node.Children[refID] = child
		}
	}

	// Identify root nodes (nodes that are not referenced as children)
	roots := make(map[int32]*Node)
	for fragID, node := range nodes {
		if !linkedNodes[fragID] {
			roots[fragID] = node
		}
	}

	return roots, nil
}

// upsertNode finds or creates a node in the map
func upsertNode(nodes map[int32]*Node, fragType string, fragID int32, tag string) *Node {
	fragType = strings.TrimPrefix(fragType, "*rawfrag.WldFrag")

	node, ok := nodes[fragID]
	if ok {
		node.Tag = tag
		node.FragType = fragType
		return node
	}
	node = &Node{
		FragID:   fragID,
		Tag:      tag,
		FragType: fragType,
		Children: make(map[int32]*Node),
	}
	nodes[fragID] = node
	return node
}

func PrintNode(node *Node, level int) {
	fmt.Printf("%s%s: %d (%s)\n", strings.Repeat("  ", level), node.FragType, node.FragID, node.Tag)
	for _, child := range node.Children {
		PrintNode(child, level+1)
	}
}
