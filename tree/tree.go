package tree

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/xackery/quail/raw"
)

type Node struct {
	FragID   int32
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

	// Map to track which nodes are referenced as children
	referenced := make(map[int32]bool)

	// Helper function to find or create a node
	findOrCreateNode := func(fragID int32) *Node {
		if node, exists := nodes[fragID]; exists {
			return node
		}
		node := &Node{
			FragID:   fragID,
			Tag:      "", // Tag will be assigned when creating the node
			Children: make(map[int32]*Node),
		}
		nodes[fragID] = node
		return node
	}

	// Helper function to extract the NameRef
	getNameRef := func(frag interface{}) int32 {
		v := reflect.ValueOf(frag)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// Check if the fragment has a "NameRef" field
		field := v.FieldByName("NameRef")
		if field.IsValid() && field.Kind() == reflect.Int32 {
			return int32(field.Int())
		}
		// Return -1 if NameRef does not exist
		return 0
	}

	// Create nodes and establish relationships
	for i := 1; i < len(wld.Fragments); i++ { // Start at index 1
		frag := wld.Fragments[i]
		if frag == nil {
			continue
		}

		fragID := int32(i)
		nameRef := getNameRef(frag)
		tag := ""
		if nameRef < 0 {
			tag = wld.Name(nameRef)
		}

		// Find or create the node for this fragment
		node := findOrCreateNode(fragID)
		node.Tag = strings.TrimSpace(tag)

		// Extract references from the fragment
		fragRefs := getFragRefs(frag)
		for _, refID := range fragRefs {
			if refID > 0 {
				// Mark this refID as being referenced
				referenced[refID] = true

				// Find or create the child node
				child := findOrCreateNode(refID)

				// Establish the parent-child relationship
				node.Children[refID] = child
			}
		}
	}

	// Identify root nodes (nodes that are not referenced as children)
	roots := make(map[int32]*Node)
	for fragID, node := range nodes {
		if !referenced[fragID] {
			roots[fragID] = node
		}
	}

	return roots, nil
}

func PrintNode(node *Node, level int) {
	fmt.Printf("%sNode: %d (%s)\n", strings.Repeat("  ", level), node.FragID, node.Tag)
	for _, child := range node.Children {
		PrintNode(child, level+1)
	}
}
