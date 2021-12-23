package supplements

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/hash"
	"github.com/whitesource/log4j-detect/records"
)

type Ruby struct {
	logger logr.Logger
}

func NewRubyEnhancer(logger logr.Logger) Ruby {
	return Ruby{
		logger: logger.WithName("enhancer.ruby"),
	}
}

func (r Ruby) LType() records.LibType {
	return records.LTRuby
}

func (r Ruby) Enhance(opr records.OperationResult) records.EnhancedResult {
	deps := make(map[records.Id]records.DependencyInfo, 0)
	result := records.EnhancedResult{
		OperationResult: opr,
		Deps:            &deps,
	}

	for id, lib := range *opr.Libraries {
		d := records.DependencyInfo{
			ArtifactID:     lib.Artifact,
			Version:        lib.Version,
			SystemPath:     lib.SystemPath,
			DependencyType: lib.LType.String(),
			DependencyFile: opr.ManifestFile,
		}
		d.Sha1 = r.calcRealSha1(lib.SystemPath)
		d.AdditionalSha1 = r.calcAdditionalSha1(lib)
		d.Checksums = &map[string]string{
			hash.Sha1:           d.Sha1,
			hash.AdditionalSha1: d.AdditionalSha1,
		}
		deps[id] = d
	}

	return result
}

func (r Ruby) calcAdditionalSha1(l records.Library) string {
	sha1, err := sha(l, avt{})
	if err != nil {
		fmt.Printf("error calculating additionalSha1 for %s:%s:%s@%s, %v",
			l.GroupId, l.Artifact, l.Version, l.LType, err)
		return ""
	}
	return sha1
}

func (r Ruby) calcRealSha1(path string) string {
	if len(path) == 0 {
		return ""
	}
	realSha1, err := hash.FileSha1(path)
	if err != nil {
		r.logger.Error(err, "failed to calculate sha1", "path", path)
		return ""
	}
	return realSha1
}
