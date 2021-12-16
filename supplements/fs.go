package supplements

import (
	"github.com/whitesource/log4j-detect/hash"
	"github.com/whitesource/log4j-detect/records"
)

type Fs struct{}

func NewFsEnhancer() Fs {
	return Fs{}
}

func (f Fs) LType() records.LibType {
	return records.LTFs
}

func (f Fs) Enhance(opr records.OperationResult) records.EnhancedResult {
	deps := make(map[records.Id]records.DependencyInfo, 0)
	result := records.EnhancedResult{
		OperationResult: opr,
		Deps:            &deps,
	}
	for id, lib := range *opr.Libraries {
		d := records.DependencyInfo{
			GroupID:        lib.GroupId,
			ArtifactID:     lib.Artifact,
			Version:        lib.Version,
			SystemPath:     lib.SystemPath,
			DependencyType: lib.LType.String(),
			DependencyFile: opr.ManifestFile,
			Scope:          lib.LScope.String(),
		}
		sha1 := f.calcRealSha1(lib.SystemPath)
		d.Sha1 = sha1
		d.Checksums = &map[string]string{
			hash.Sha1:           sha1,
			hash.AdditionalSha1: d.AdditionalSha1,
		}
		deps[id] = d
	}

	return result
}

func (f Fs) calcRealSha1(path string) string {
	if len(path) == 0 {
		return ""
	}
	realSha1, _ := hash.FileSha1(path)
	return realSha1
}
