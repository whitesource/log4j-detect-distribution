package supplements

import (
	"github.com/whitesource/log4j-detect/records"
	"testing"
)

func TestGolangLType(t *testing.T) {
	g := Golang{}
	if g.LType() != records.LTGolang {
		t.Error("failed to assign the correct LType for golang")
	}
}

func TestGolangEnhance1(t *testing.T) {
	lib := records.Library{
		Artifact: "github.com/karrick/godirwalk",
		Version:  "v1.16.1",
		LType:    records.LTGolang,
	}
	testGolangSingleLibrary("ab6e9c1f0b207f8782f594b79807275d6349753b", lib, t)
}

func TestGolangEnhance2(t *testing.T) {
	lib := records.Library{
		Artifact: "github.com/karrick/godirwalk",
		Commit:   "aed5e4c7ecf9",
		LType:    records.LTGolang,
	}
	testGolangSingleLibrary("4e4a226070034452c1518840acd82d06f995ac7e", lib, t)
}

func TestGolangEnhance3(t *testing.T) {
	lib := records.Library{
		Artifact: "github.com/karrick/godirwalk",
		Version:  "v1.16.1",
		Commit:   "aed5e4c7ecf9",
		LType:    records.LTGolang,
	}
	testGolangSingleLibrary("4e4a226070034452c1518840acd82d06f995ac7e", lib, t)
}

func TestGolangEnhance4(t *testing.T) {
	libs := []records.Library{
		{
			Artifact: "github.com/olekukonko/tablewriter",
			Version:  "v0.0.5",
			LType:    records.LTGolang,
		},
		{
			Artifact: "github.com/spf13/cobra",
			Commit:   "eb3b6397b1b5",
			LType:    records.LTGolang,
		},
		{
			Artifact: "github.com/karrick/godirwalk",
			Version:  "v1.16.1",
			Commit:   "aed5e4c7ecf9",
			LType:    records.LTGolang,
		},
	}

	expSha := map[string]string{
		"github.com/karrick/godirwalk":      "4e4a226070034452c1518840acd82d06f995ac7e",
		"github.com/olekukonko/tablewriter": "54d6b791ea3f57636fd65079a51ec0aeced7beed",
		"github.com/spf13/cobra":            "6d97c2e6f7161add56d02f5f333c7b44595f0fea",
	}

	testGolangMultipleLibraries(expSha, libs, t)
}

func testGolangSingleLibrary(expSha string, lib records.Library, t *testing.T) {
	g := Golang{}
	lName := records.Id(lib.Artifact)
	opr := records.OperationResult{
		ManifestFile: "abcd",
		Libraries:    &map[records.Id]records.Library{lName: lib},
		LType:        records.LTGolang,
	}
	er := g.Enhance(opr)
	if er.Deps == nil {
		t.Errorf("inavlid nil result recieved")
		t.FailNow()
	}

	deps := *(er.Deps)
	if l := len(deps); l != 1 {
		t.Errorf("expected 1 library actual %d", l)
	}

	d, found := deps[lName]
	if !found {
		t.Errorf("library '%s' was not found", lName)
	}
	if d.Sha1 != expSha {
		t.Errorf("invalid sha1 recieved expected '%s' actual '%s'", expSha, d.Sha1)
	}
	s, found := (*d.Checksums)["SHA1"]
	if !found {
		t.Errorf("sha1 was not found in checksums")
	}
	if s != expSha {
		t.Errorf("invalid checksum sha1 recieved expected '%s' actual '%s'", expSha, s)
	}
	if d.Version != lib.Version {
		t.Errorf("invalid version expected '%s' actual '%s'", lib.Version, d.Version)
	}
	if d.ArtifactID != lib.Artifact {
		t.Errorf("invalid artifact expected '%s' actual '%s'", lib.Artifact, d.ArtifactID)
	}
	if d.Commit != lib.Commit {
		t.Errorf("invalid commit expected '%s' actual '%s'", lib.Commit, d.Commit)
	}
	if d.DependencyType != lib.LType.String() {
		t.Errorf("invalid type expected '%s' actual '%s'", lib.LType, d.Type)
	}

}

func testGolangMultipleLibraries(expSha map[string]string, libs []records.Library, t *testing.T) {
	g := Golang{}
	opr := records.OperationResult{
		ManifestFile: "abcd",
		Libraries:    &map[records.Id]records.Library{},
		LType:        records.LTGolang,
	}

	for _, lib := range libs {
		(*opr.Libraries)[records.Id(lib.Artifact)] = lib
	}
	er := g.Enhance(opr)
	if er.Deps == nil {
		t.Errorf("invalid nil result recieved")
		t.FailNow()
	}

	deps := *(er.Deps)
	if l := len(deps); l != len(libs) {
		t.Errorf("expected %d libraries actual %d", len(libs), l)
	}
	for _, dep := range libs {
		d, exists := deps[records.Id(dep.Artifact)]
		if !exists {
			t.Errorf("library '%s' was not found", dep.Artifact)
		}
		if d.Sha1 != expSha[d.ArtifactID] {
			t.Errorf("expected '%s' sha1 for '%s' lib but got '%s' ", expSha[d.ArtifactID], d.ArtifactID, d.Sha1)
		}
		s, found := (*d.Checksums)["SHA1"]
		if !found {
			t.Errorf("sha1 was not found in checksums")
		}
		if s != expSha[d.ArtifactID] {
			t.Errorf("invalid checksum sha1 recieved expected '%s' actual '%s'", expSha, s)
		}
	}
}

func TestGolangLibraryNegativeScenarios(t *testing.T) {
	g := Golang{}
	libs := map[records.Id]records.Library{
		"a": {
			Version: "v1.16.1",
			Commit:  "aed5e4c7ecf9",
			LType:   records.LTGolang,
		},
		"b": {
			Artifact: "abcd",
			Version:  "v1.16.1",
			Commit:   "aed5e4c7ecf9",
		},
		"c": {
			Artifact: "abcd",
			LType:    records.LTGolang,
		},
		"d": {
			Artifact: "abcd",
			LType:    records.LTGolang,
			Version:  "v0.0.0",
		},
	}

	opr := records.OperationResult{
		ManifestFile:      "abcd",
		Direct:            nil,
		LibraryToChildren: nil,
		Libraries:         &libs,
		LType:             records.LTGolang,
		Err:               nil,
	}

	er := g.Enhance(opr)
	if er.Deps == nil {
		t.Errorf("inavlid nil result recieved")
		t.FailNow()
	}

	deps := *(er.Deps)
	if l := len(deps); l != 0 {
		t.Errorf("expected 1 library actual %d", l)
	}
}
