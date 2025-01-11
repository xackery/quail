package tree

import (
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/common"
)

// PrintNode prints a tree node to a writer
func PrintNode(w io.Writer, node common.TreeLinker, level int) error {
	if node == nil {
		return nil
	}
	fmt.Fprintf(w, "%s%s: %d (%s)\n", strings.Repeat("  ", level), node.FragType(), node.FragID(), node.Tag())
	for _, child := range node.Children() {
		err := PrintNode(w, child, level+1)
		if err != nil {
			return fmt.Errorf("node %s: %w", child.Tag(), err)
		}
	}
	return nil
}
