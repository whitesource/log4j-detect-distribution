package supplements

import (
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/whitesource/log4j-detect/records"
	"testing"
)

func TestRealSha(t *testing.T) {
}

func TestJavaLType(t *testing.T) {
	j := Java{}
	if j.LType() != records.LTJava {
		t.Error("failed to assign the correct LType for Java")
	}
}

func TestJavaEnhance1(t *testing.T) {
	lib := records.Library{
		GroupId:    "io.netty",
		Artifact:   "netty-handler",
		Version:    "4.0.56.Final",
		SystemPath: "./testdata/dummy.zip",
		LType:      records.LTJava,
	}
	testJavaSingleLibrary("db39e7ac4886a4221d8e966381305d0638cfc539",
		"bff628cb4b92c865e3d6ccc6eb1c6500e699d61f", lib, t)
}

func TestJavaEnhance2(t *testing.T) {
	lib := records.Library{
		GroupId:  "io.netty",
		Artifact: "netty-handler",
		Version:  "4.0.56.Final",
		LType:    records.LTJava,
	}
	testJavaSingleLibrary("",
		"bff628cb4b92c865e3d6ccc6eb1c6500e699d61f", lib, t)
}

func testJavaSingleLibrary(expSha string, expAddSha string, lib records.Library, t *testing.T) {
	j := Java{}
	lName := records.Id(lib.Artifact)
	opr := records.OperationResult{
		ManifestFile: "abcd",
		Libraries:    &map[records.Id]records.Library{lName: lib},
		LType:        records.LTJava,
	}
	er := j.Enhance(opr)
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
	s, found := (*d.Checksums)["ADDITIONAL_SHA1"]
	if !found {
		t.Errorf("ADDITIONAL_SHA1 was not found in checksums")
	}
	if s != expAddSha {
		t.Errorf("invalid checksum ADDITIONAL_SHA1 recieved expected '%s' actual '%s'", expSha, s)
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

func TestJavaInvalidLibs(t *testing.T) {
	libs := []records.Library{
		{Version: "1.0.2", LType: records.LTJava},
		{GroupId: "io.netty", Version: "1.0.2", LType: records.LTJava},
		{Artifact: "netty-handler", LType: records.LTJava},
		{GroupId: "io.netty", Artifact: "netty-handler", LType: records.LTJava},
		{Artifact: "netty-handler", Version: "1.0.2"},
		{GroupId: "io.netty", Artifact: "netty-handler", Version: "1.0.2"},
		{Artifact: "netty-handler"},
		{GroupId: "io.netty", Artifact: "netty-handler"},
		{Version: "1.0.2"},
		{GroupId: "io.netty"},
		{GroupId: "io.netty", Artifact: "netty-handler"},
	}

	nop := zerolog.Nop()
	lgr := zerologr.New(&nop)
	java := NewJavaEnhancer(lgr)

	for i, l := range libs {
		if java.invalid(&l) == false {
			t.Errorf("lib %d: 'invalid' returned incorrect result", i)
		}
	}
}
