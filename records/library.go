package records

import (
	"fmt"
)

type Id string

type Library struct {
	Artifact string
	Version  string
	LScope   LibScope
	LType    LibType

	// the path to the library files
	SystemPath string
	// if it's a multi-module project wrapping
	IsProject bool
	// optional
	Commit string
	// mainly related to Java libraries
	GroupId string
}

/* *** private *** */

func (lib Library) text(isDeduplicated bool) (res string) {
	vORc := lib.Commit
	if len(vORc) == 0 {
		vORc = lib.Version
	}

	if len(lib.GroupId) > 0 {
		res += fmt.Sprintf("%s:", lib.GroupId)
	}

	res += fmt.Sprintf("%s:%s", lib.Artifact, vORc)
	if isDeduplicated {
		res += " (d)"
	}

	return res
}
