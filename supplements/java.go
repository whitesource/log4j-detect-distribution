package supplements

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/hash"
	"github.com/whitesource/log4j-detect/records"
	"strings"
)

type Java struct {
	logger logr.Logger
}

func NewJavaEnhancer(logger logr.Logger) Java {
	return Java{
		logger: logger.WithName("enhancer.java"),
	}
}

func (j Java) LType() records.LibType {
	return records.LTJava
}

func (j Java) Enhance(opr records.OperationResult) records.EnhancedResult {
	deps := make(map[records.Id]records.DependencyInfo, 0)
	result := records.EnhancedResult{
		OperationResult: opr,
		Deps:            &deps,
	}
	for id, lib := range *opr.Libraries {
		if j.invalid(&lib) {
			continue
		}
		d := records.DependencyInfo{
			GroupID:        lib.GroupId,
			ArtifactID:     lib.Artifact,
			Version:        lib.Version,
			SystemPath:     lib.SystemPath,
			DependencyType: lib.LType.String(),
			DependencyFile: opr.ManifestFile,
			Scope:          lib.LScope.String(),
		}
		sha1 := j.calcRealSha1(lib.SystemPath)
		d.Sha1 = sha1
		d.AdditionalSha1 = j.calcAdditionalSha1(lib)
		d.Checksums = &map[string]string{
			hash.Sha1:           sha1,
			hash.AdditionalSha1: d.AdditionalSha1,
		}
		deps[id] = d
	}

	return result
}

func (j Java) invalid(l *records.Library) bool {
	if l.LType != records.LTJava || len(l.Artifact) == 0 || len(l.GroupId) == 0 {
		return true
	}
	if len(l.Version) == 0 || strings.ContainsAny(l.Version, "{}()$") {
		return true
	}
	return false
}

func (j Java) calcAdditionalSha1(l records.Library) string {
	sha1, err := sha(l, gavtlc{})
	if err != nil {
		fmt.Printf("error calculating additionalSha1 for %s:%s:%s@%s, %v",
			l.GroupId, l.Artifact, l.Version, l.LType, err)
		return ""
	}
	return sha1
}

func (j Java) calcRealSha1(path string) string {
	if len(path) == 0 {
		return ""
	}
	realSha1, err := hash.FileSha1(path)
	if err != nil {
		j.logger.Error(err, "failed to calculate sha1", "path", path)
		return ""
	}
	return realSha1
}
