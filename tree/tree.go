package tree

import (
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

type Node struct {
	FragID   int32
	FragType string
	Tag      string
	Children map[int32]*Node
}

// Dump dumps a tree to a writer
func Dump(w io.Writer, isChr bool, src interface{}) error {
	var err error
	var nodes map[int32]*Node
	switch val := src.(type) {
	case *raw.Wld:
		nodes, _, err = BuildFragReferenceTree(isChr, val)
		if err != nil {
			return fmt.Errorf("build frag reference tree: %w", err)
		}
		for _, root := range nodes {
			fmt.Fprintf(w, "Root ")
			PrintNode(w, root, 0)
		}
	case *raw.Bmp:
		return nil
	default:
		return fmt.Errorf("unknown type %T", val)
	}
	return nil
}

func BuildFragReferenceTree(isChr bool, wld *raw.Wld) (map[int32]*Node, map[int32]*Node, error) {
	// Map to store all nodes
	nodes := make(map[int32]*Node)

	// Map to track which nodes are linkedNodes as children
	linkedNodes := make(map[int32]bool)

	//var zoneNode *Node

	actorNodes := make(map[string]*Node)
	// find actornodes and build them first
	for i := 1; i < len(wld.Fragments); i++ { // Start at index 1
		frag := wld.Fragments[i]
		if frag == nil {
			continue
		}
		switch frag.(type) {
		// case *rawfrag.WldFragTrackDef:
		// 	// there's a special case where orphaned tracks happen
		// 	tag := wld.Name(frag.NameRef())
		// 	_, model := helper.TrackAnimationParse(isChr, tag)
		// 	switch model {
		// 	case "POINT":
		// 		// this is an dummy actordef
		// 		tag = "POINT_HS_DEF"
		// 		actorNodes[tag] = upsertNode(nodes, "", fmt.Sprintf("%T", frag), int32(i), strings.TrimSpace(tag))
		// 		continue
		// 	}

		case *rawfrag.WldFragHierarchicalSpriteDef:
		default:
			continue
		}

		tag := ""
		if frag.NameRef() < 0 {
			tag = wld.Name(frag.NameRef())
		}
		actorNodes[tag] = upsertNode(nodes, fmt.Sprintf("%T", frag), int32(i), strings.TrimSpace(tag))
	}

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

		switch frag.(type) {
		// case *rawfrag.WldFragGlobalAmbientLightDef:
		// 	zoneNode = node
		// 	if zoneNode.Tag == "" {
		// 		zoneNode.Tag = "zone"
		// 	}
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

			child := upsertNode(nodes, fmt.Sprintf("%T", childFrag), refID, strings.TrimSpace(childTag))

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

	return roots, nodes, nil
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

func PrintNode(w io.Writer, node *Node, level int) {
	fmt.Fprintf(w, "%s%s: %d (%s)\n", strings.Repeat("  ", level), node.FragType, node.FragID, node.Tag)
	for _, child := range node.Children {
		PrintNode(w, child, level+1)
	}
}
