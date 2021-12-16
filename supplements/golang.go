package supplements

import (
	"fmt"
	"github.com/whitesource/log4j-detect/hash"
	"github.com/whitesource/log4j-detect/records"
)

type Golang struct{}

func NewGolangEnhancer() Golang {
	return Golang{}
}

func (g Golang) LType() records.LibType {
	return records.LTGolang
}

func (g Golang) Enhance(opResult records.OperationResult) records.EnhancedResult {
	deps := make(map[records.Id]records.DependencyInfo, 0)
	result := records.EnhancedResult{
		OperationResult: opResult,
		Deps:            &deps,
	}
	for id, lib := range *opResult.Libraries {
		if g.invalid(&lib) {
			continue
		}
		d := records.DependencyInfo{
			GroupID:        lib.Artifact,
			ArtifactID:     lib.Artifact,
			Version:        lib.Version,
			SystemPath:     lib.SystemPath,
			DependencyType: lib.LType.String(),
			Commit:         lib.Commit,
			DependencyFile: opResult.ManifestFile,
		}
		d.Sha1 = g.sha1(lib)
		d.Checksums = &map[string]string{
			hash.Sha1: d.Sha1,
		}
		deps[id] = d
	}

	return result
}

func (g Golang) invalid(l *records.Library) bool {
	if l.LType != records.LTGolang {
		return true
	}
	if len(l.Artifact) == 0 {
		return true
	}
	if len(l.Commit) == 0 && (len(l.Version) == 0 || l.Version == "v0.0.0") {
		return true
	}
	return false
}

func (g Golang) sha1(l records.Library) string {
	if len(l.Commit) > 0 {
		sha1, err := sha(l, act{})
		if err == nil {
			return sha1
		} else {
			fmt.Printf("error calculating sha1 for %s@%s:%s, %v", l.Artifact, l.Commit, l.LType, err)
		}
	}

	if len(l.Version) > 0 && l.Version != "v0.0.0" {
		sha1, err := sha(l, avt{})
		if err == nil {
			return sha1
		} else {
			fmt.Printf("error calculating sha1 for %s@%s:%s, %v", l.Artifact, l.Version, l.LType, err)
		}
	}
	return ""
}
