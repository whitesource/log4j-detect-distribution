package records

import (
	"fmt"
	"github.com/whitesource/log4j-detect/utils"
	"io"
	"strings"
)

type VulnerableLib struct {
	GroupId        string      `json:"groupId"`
	Artifact       string      `json:"artifact"`
	PackageVersion interface{} `json:"packageVersion"`
	FileName       string      `json:"fileName"`
	Sha1           string      `json:"sha1"`
	CVE            string
}

type EnhancedResult struct {
	OperationResult
	Deps                *map[Id]DependencyInfo
	DepId2VulnerableLib map[Id]VulnerableLib
}

/* *** public *** */

func (er EnhancedResult) IsVulnerable() bool {
	return len(er.DepId2VulnerableLib) > 0
}

func (er EnhancedResult) Print(w io.Writer) {
	tree := er.tree()
	er.identifyVulnerablePaths(tree)
	firstLinePrinted := false

	for _, tn := range tree {
		if tn.vulnerable {
			if firstLinePrinted {
				_, _ = fmt.Fprintln(w, "|")
			}
			er.printErRec(w, *tn, 0)
			firstLinePrinted = true
		}
	}
}

func (er EnhancedResult) DependencyTree() []DependencyInfo {
	var result []DependencyInfo
	nodes := er.tree()

	for _, n := range nodes {
		di := er.buildDependencyTreeRec(*n)
		result = append(result, di)
	}

	return result
}

/* *** private *** */

func (er EnhancedResult) printErRec(w io.Writer, tn treeNode, level int) {
	if !tn.vulnerable {
		return
	}

	lib := (*er.Libraries)[tn.id]

	txt := fmt.Sprintf("%s|-- %s",
		strings.Repeat("|  ", level),
		lib.text(tn.deduped),
	)

	hasVulKids := hasVulChildren(tn)
	if !hasVulKids {
		txt += " " + utils.MakeRedText("Vulnerable")
	}

	_, _ = fmt.Fprintln(w, txt)

	if hasVulKids {
		for _, child := range tn.children {
			er.printErRec(w, *child, level+1)
		}
	}
}

func hasVulChildren(tn treeNode) bool {
	for _, child := range tn.children {
		if (*child).vulnerable {
			return true
		}
	}
	return false
}

func (er EnhancedResult) identifyVulnerablePaths(tree []*treeNode) {
	if !er.IsVulnerable() {
		return
	}

	for _, tn := range tree {
		er.identifyVulnerableNodes(tn)
	}
}

func (er EnhancedResult) identifyVulnerableNodes(n *treeNode) bool {
	vulPath := false
	for _, child := range n.children {
		vulSubPath := er.identifyVulnerableNodes(child)
		vulPath = vulPath || vulSubPath
	}

	_, isVulNode := er.DepId2VulnerableLib[n.id]
	n.vulnerable = isVulNode || vulPath

	return n.vulnerable
}

func (er EnhancedResult) vulnerabilityText() (res string) {
	if !er.IsVulnerable() {
		return
	}

	return utils.MakeRedText("HIGH")
}

func (er EnhancedResult) buildDependencyTreeRec(tn treeNode) DependencyInfo {
	di := (*er.Deps)[tn.id]
	di.Deduped = tn.deduped

	if len(tn.children) == 0 {
		return di
	}

	var children []DependencyInfo
	for _, cn := range tn.children {
		child := er.buildDependencyTreeRec(*cn)
		children = append(children, child)
	}

	di.Children = &children
	return di
}
