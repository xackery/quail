package tree

import (
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

type Node struct {
	FragID   int32
	FragType string
	Tag      string
	Parent   string
	Children map[int32]*Node
}

// Dump dumps a tree to a writer
func Dump(isChr bool, src interface{}, w io.Writer) error {
	var err error
	var nodes map[int32]*Node
	switch val := src.(type) {
	case *raw.Wld:
		nodes, err = BuildFragReferenceTree(isChr, val)
		if err != nil {
			return fmt.Errorf("build frag reference tree: %w", err)
		}
		for _, root := range nodes {
			fmt.Printf("Root ")
			PrintNode(root, 0)
		}
	case *raw.Bmp:
		return nil
	default:
		return fmt.Errorf("unknown type %T", val)
	}
	return nil
}

func BuildFragReferenceTree(isChr bool, wld *raw.Wld) (map[int32]*Node, error) {
	// Map to store all nodes
	nodes := make(map[int32]*Node)

	// Map to track which nodes are linkedNodes as children
	linkedNodes := make(map[int32]bool)

	actorNodes := make(map[string]*Node)
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
		node := upsertNode(nodes, "", fmt.Sprintf("%T", frag), fragID, strings.TrimSpace(tag))

		_, ok := frag.(*rawfrag.WldFragHierarchicalSpriteDef)
		if ok {
			actorNodes[tag] = node
		}

		switch frag.(type) {
		case *rawfrag.WldFragTrackDef:
			_, model := helper.TrackAnimationParse(isChr, tag)
			parentNode, ok := actorNodes[model+"_HS_DEF"]
			if ok {
				parentNode.Children[fragID] = node
			}
		case *rawfrag.WldFragTrack:
			_, model := helper.TrackAnimationParse(isChr, tag)
			parentNode, ok := actorNodes[model+"_HS_DEF"]
			if ok {
				parentNode.Children[fragID] = node
			}

		}

		// Extract references from the fragment
		fragRefs := fragRefs(frag)
		for _, refID := range fragRefs {
			if refID <= 0 {
				continue
			}
			// Mark this refID as being referenced
			linkedNodes[refID] = true

			// Find or create the child node
			childFrag := wld.Fragments[refID]
			childTag := ""
			if childFrag != nil && childFrag.NameRef() < 0 {
				childTag = wld.Name(childFrag.NameRef())
			}

			child := upsertNode(nodes, node.Tag, fmt.Sprintf("%T", childFrag), refID, strings.TrimSpace(childTag))

			// Establish the parent-child relationship
			node.Children[refID] = child
		}
	}

	// Identify root nodes (nodes that are not referenced as children)
	roots := make(map[int32]*Node)
	for fragID, node := range nodes {
		if linkedNodes[fragID] {
			continue
		}
		roots[fragID] = node

	}

	return roots, nil
}

// upsertNode finds or creates a node in the map
func upsertNode(nodes map[int32]*Node, parent string, fragType string, fragID int32, tag string) *Node {
	fragType = strings.TrimPrefix(fragType, "*rawfrag.WldFrag")

	node, ok := nodes[fragID]
	if ok {
		node.Parent = parent
		node.Tag = tag
		node.FragType = fragType
		return node
	}
	node = &Node{
		Parent:   parent,
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
