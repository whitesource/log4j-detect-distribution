package supplements

import (
	"github.com/stretchr/testify/assert"
	"github.com/whitesource/log4j-detect/records"
	"strings"
	"testing"
)

func TestShaMethodGAVTLC(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		GroupId:  "TestGroup",
		LType:    records.LTJava,
	}
	res, err := sha(lib, gavtlc{})
	if err != nil {
		t.Errorf("Failed to calculate GAVTLC sha %v", err)
	}
	if strings.Compare(res, "1f901287ab64af4f454c7abdebf4cc1dc5adeee3") != 0 {
		t.Errorf("Failed to calculate the right SHA Expected: %s actual %s",
			"1f901287ab64af4f454c7abdebf4cc1dc5adeee3", res)
	}

}

func TestShaMethodACT(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Commit:   "de6cabda18fe637bc2a9a796963b76598464adeb",
		LType:    records.LTJava,
	}
	res, err := sha(lib, act{})
	if err != nil {
		t.Errorf("Failed to calculate GACTLC sha %v", err)
	}
	if strings.Compare(res, "86813b876315a1cda86ae0c818b8560918289b4e") != 0 {
		t.Errorf("Failed to calculate the right SHA Expected: %s actual %s",
			"86813b876315a1cda86ae0c818b8560918289b4e", res)
	}

}

func TestShaMethodAVT(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	res, err := sha(lib, avt{})
	if err != nil {
		t.Errorf("Failed to calculate GACTLC sha %v", err)
	}
	if strings.Compare(res, "e8ca41d5262b5c0bf2c56c8cce96205aec85c925") != 0 {
		t.Errorf("Failed to calculate the right SHA Expected: %s actual %s",
			"e8ca41d5262b5c0bf2c56c8cce96205aec85c925", res)
	}

}

func TestShaMethodAVTLC(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	res, err := sha(lib, avtlc{})
	if err != nil {
		t.Errorf("Failed to calculate GACTLC sha %v", err)
	}
	if strings.Compare(res, "d4312820a14d8c23b77629f57175949d5d28fbc2") != 0 {
		t.Errorf("Failed to calculate the right SHA Expected: %s actual %s",
			"d4312820a14d8c23b77629f57175949d5d28fbc2", res)
	}

}

func TestShaMethodAVTLCA(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1-RC",
		LType:    records.LTJava,
	}
	res, err := sha(lib, avtlca{})
	if err != nil {
		t.Errorf("Failed to calculate GACTLC sha %v", err)
	}
	if strings.Compare(res, "62ff56fbce148300b50bd7eb868a29265bc4116d") != 0 {
		t.Errorf("Failed to calculate the right SHA Expected: %s actual %s",
			"62ff56fbce148300b50bd7eb868a29265bc4116d", res)
	}

}

func TestShaMethodGAVT(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		GroupId:  "TestGroup",
		LType:    records.LTJava,
	}
	res, err := sha(lib, gavt{})
	if err != nil {
		t.Errorf("Failed to calculate GACTLC sha %v", err)
	}
	if strings.Compare(res, "fc2762755136cfc27272bf63a42d03b5db81425e") != 0 {
		t.Errorf("Failed to calculate the right SHA Expected: %s actual %s",
			"fc2762755136cfc27272bf63a42d03b5db81425e", res)
	}

}

func TestShaMethodEmptyValueGAVT(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		GroupId:  "",
		LType:    records.LTJava,
	}
	_, err := sha(lib, gavt{})
	if err == nil {
		t.Errorf("Failed to validate sha compontents for GAVT sha %v", err)
	}
}

func TestShaMethodEmptyValueGAVTLC(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		GroupId:  "",
		LType:    records.LTJava,
	}
	_, err := sha(lib, gavtlc{})
	if err == nil {
		t.Errorf("Failed to validate sha compontents for GAVT sha %v", err)
	}
}

func TestShaMethodEmptyValueAVT(t *testing.T) {
	lib := records.Library{
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	_, err := sha(lib, avt{})
	if err == nil {
		t.Errorf("Failed to validate sha compontents for GAVT sha %v", err)
	}
}

func TestShaMethodEmptyValueAVTLC(t *testing.T) {
	lib := records.Library{
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	_, err := sha(lib, avtlc{})
	if err == nil {
		t.Errorf("Failed to validate sha compontents for GAVT sha %v", err)
	}
}

func TestShaMethodEmptyValueAVTLCA(t *testing.T) {
	lib := records.Library{
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	_, err := sha(lib, avtlca{})
	if err == nil {
		t.Errorf("Failed to validate sha compontents for GAVT sha %v", err)
	}
}

func TestShaMethodEmptyValueACT(t *testing.T) {
	lib := records.Library{
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	_, err := sha(lib, act{})
	if err == nil {
		t.Errorf("Failed to validate sha compontents for GAVT sha %v", err)
	}
}

func TestValidateGAVPositive(t *testing.T) {
	lib := records.Library{
		GroupId:  "TestGroup",
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	err := validateGAVT(lib)
	if err != nil {
		t.Errorf("GAV validation failed as not expected %v", err)
	}
}

func TestValidateGAVNegative(t *testing.T) {
	lib := records.Library{
		GroupId:  "",
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	err := validateGAVT(lib)
	if err == nil {
		t.Errorf("GAV validation passed as not expected")
	}
}

func TestValidateAVTPositive(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	err := validateAVT(lib)
	if err != nil {
		t.Errorf("AV validation failed as not expected")
	}
}

func TestValidateAVTNegative(t *testing.T) {
	lib := records.Library{
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	err := validateAVT(lib)
	if err == nil {
		t.Errorf("AV validation passed as not expected")
	}
}

func TestValidateACTPositive(t *testing.T) {
	lib := records.Library{
		Artifact: "TestArtifact",
		Version:  "1.1.1",
		Commit:   "de6cabda18fe637bc2a9a796963b76598464adeb",
		LType:    records.LTJava,
	}
	err := validateACT(lib)
	if err != nil {
		t.Errorf("AC validation failed as not expected")
	}
}

func TestValidateACTNegative(t *testing.T) {
	lib := records.Library{
		Artifact: "",
		Version:  "1.1.1",
		LType:    records.LTJava,
	}
	err := validateACT(lib)
	if err == nil {
		t.Errorf("AC validation passed as not expected")
	}
}

func Test_underscoreHash(t *testing.T) {
	tests := []struct {
		elems    []string
		expected string
	}{
		{
			elems:    []string{"a", "b", "c"},
			expected: "1680e5e1c8d101fbbefa27b901a7a20e63a9bea6",
		},
		{
			elems:    []string{"mystr"},
			expected: "900dc1f727411f49ed49bb4a08f9220afd6ef0f4",
		},
		{
			elems:    []string{"", "string", ""},
			expected: "c88967feb5cdee67db8bb31cf8d2bbdb3dab85ec",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, underscoreSha1(tt.elems...))
	}
}
