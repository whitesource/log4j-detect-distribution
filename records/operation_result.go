package records

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type OperationResult struct {
	// the manifest file from which the library discovered
	// e.g. go.mod, pom.xml, package-lock.json
	ManifestFile string

	// direct libraries found
	Direct *[]Id

	// mapping from each library to its direct transitive libraries
	LibraryToChildren *map[Id][]Id

	// info about all libraries found
	Libraries *map[Id]Library

	LType LibType

	// in case an error occurred for the specific scan
	Err error

	// how long it took to produce this operation result
	Elapsed time.Duration

	// Organ represents the organ type for this operation result
	Organ Organ
}

/* *** public *** */

func (op *OperationResult) Print(w io.Writer) {
	tree := op.tree()
	_, _ = fmt.Fprintf(w, "%s (%s)\n", op.ManifestFile, op.LType)
	for i, tn := range tree {
		if i > 0 {
			_, _ = fmt.Fprintln(w, "|")
		}
		op.printOprRec(*tn, 0, w)
	}
}

// TrackTime sets the total time that elapsed to
func (op *OperationResult) TrackTime(start time.Time) {
	op.Elapsed = time.Since(start)
}

/* *** private *** */

func (op *OperationResult) printOprRec(tn treeNode, level int, w io.Writer) {
	lib := (*op.Libraries)[tn.id]

	txt := fmt.Sprintf("%s|-- %s", strings.Repeat("|  ", level), lib.text(tn.deduped))
	_, _ = fmt.Fprintln(w, txt)

	for _, child := range tn.children {
		op.printOprRec(*child, level+1, w)
	}
}

func (op *OperationResult) tree() []*treeNode {
	if op.Direct == nil || op.LibraryToChildren == nil {
		return []*treeNode{}
	}

	// short name
	l2c := *op.LibraryToChildren
	// prevent circular tree + identify deduplicated
	processed := make(map[Id]bool)
	var current, next []levelNode
	root := treeNode{
		children: []*treeNode{},
	}

	for _, id := range *op.Direct {
		current = append(current, levelNode{id: id, parent: &root})
	}

	for len(current) > 0 {
		// reset next level
		next = []levelNode{}

		for _, ln := range current {
			if _, found := (*op.Libraries)[ln.id]; !found {
				continue
			}
			_, isProcessed := processed[ln.id]
			// mark as processed regardless of previous status
			processed[ln.id] = true

			newNode := new(treeNode)
			newNode.children = []*treeNode{}
			newNode.deduped = isProcessed
			newNode.id = ln.id

			ln.parent.add(newNode)

			if !isProcessed {
				// children to be processed next iteration
				for _, v := range l2c[ln.id] {
					next = append(next, levelNode{id: v, parent: newNode})
				}
			}
		}

		current = next
	}

	return root.children
}

/* *** helper types *** */

type treeNode struct {
	id         Id
	deduped    bool
	vulnerable bool
	children   []*treeNode
}

func (tn *treeNode) add(c *treeNode) {
	tn.children = append(tn.children, c)
}

type levelNode struct {
	// id of child node to process
	id     Id
	parent *treeNode
}
