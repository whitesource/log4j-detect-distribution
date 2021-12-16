package hash

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_CalculateStringSha1(t *testing.T) {
	res := StringSha1("cocoon_cocoon-ajax_2.1.11_JAVA")
	if strings.Compare(res, "4a9cb79a26090ddcd04725ac9cf3bd6e804bb4ac") != 0 {
		t.Error("failed to match ", res)
	}
}

func Test_CalculateFileSha1(t *testing.T) {
	tests := []struct {
		file string
		sha1 string
	}{
		{
			file: "./testdata/file_hash/dummy.txt",
			sha1: "fbc8494d99b8017960fe11194b76cf12cf781995",
		},
		{
			file: "./testdata/file_hash/dummy.zip",
			sha1: "db39e7ac4886a4221d8e966381305d0638cfc539",
		},
	}
	for _, tt := range tests {
		got, err := FileSha1(tt.file)
		assert.Nil(t, err, "%s - %v", tt.file, err)
		assert.Equal(t, tt.sha1, got, "%s - expected = %s, got = %s", tt.file, tt.sha1, got)
	}
}
