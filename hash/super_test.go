//go:build !race
// +build !race

package hash

import (
	"crypto"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

var superHashFilesDir = path.Join("testdata", "super_hash_files")
var otherPlatformFileDir = path.Join("testdata", "other_platform")

func Test_calculateBucketSize(t *testing.T) {
	tests := []struct {
		size               float64
		expectedBucketSize int
	}{
		{
			size:               47038576,
			expectedBucketSize: 2.25e7,
		},
		{
			size:               3855366,
			expectedBucketSize: 1750000,
		},
		{
			size:               3156,
			expectedBucketSize: 1750,
		},
		{
			size:               130725,
			expectedBucketSize: 75000,
		},
		{
			size:               18039443,
			expectedBucketSize: 7500000,
		},
	}

	h := &SuperHash{algo: crypto.SHA1}
	for _, tt := range tests {
		got := h.bucketSize(tt.size)
		assert.Equal(t, tt.expectedBucketSize, got)
	}
}

func Test_CalculateSuperHashSha1(t *testing.T) {
	tests := []struct {
		path     string
		expected SuperHash
	}{
		{
			path: "plus.pdf",
			expected: SuperHash{
				FullHash:         "7fbad47886605db0c4d8818813300596c3f91af1",
				MostSigBitsHash:  "9d9f8db8145c404e9a19c6a70807514a76e9a448",
				LeastSigBitsHash: "13b59bdef20a3fc512dbb23e538389012b31779b",
			},
		},
		{
			path: "post.js",
			expected: SuperHash{
				FullHash:         "b114665429a97c137bdf5694d7f51d2211f616c2",
				MostSigBitsHash:  "8869a209f04a6c0ee04c80b6900de78494da9347",
				LeastSigBitsHash: "937318d2591e2831e9600d6d58294a8d0a311315",
			},
		},
		{
			path: "tsserver.js",
			expected: SuperHash{
				FullHash:         "8fd21684e9c4b3a659360b2c9d63e6bf7f56c380",
				MostSigBitsHash:  "0341ea8d75d00896d0e808b19b83a31a217f93dd",
				LeastSigBitsHash: "e038603a2c1f964bf97a7ecf60ce6c8a8ddfb7d3",
			},
		},
		{
			path: "update-request.txt",
			expected: SuperHash{
				FullHash:         "0070cc6e7351550148e42ac7e8e8583586596d5e",
				MostSigBitsHash:  "721bb0f1c61c6739dcfdc98dabd7c0f1e0c53514",
				LeastSigBitsHash: "00f532f4a15bda98756f8da1d5c4a9b3084ac69a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			p := path.Join(superHashFilesDir, tt.path)
			f, err := os.Open(p)
			assert.Nilf(t, err, "expected open %s to be successful, got %v", p, err)
			h, err := CalculateSuperSha1(f)
			assert.Nilf(t, err, "expected SuperHash calculation for %s to complete without err, got %v", p, err)
			assert.Equal(t, tt.expected.FullHash, h.FullHash)
			assert.Equal(t, tt.expected.MostSigBitsHash, h.MostSigBitsHash)
			assert.Equal(t, tt.expected.LeastSigBitsHash, h.LeastSigBitsHash)
		})
	}
}

func Test_OtherPlatformSha1(t *testing.T) {
	tests := []struct {
		file     string
		expected string
	}{
		{
			file:     "README_lf.md",
			expected: "36be27a95c647cd70ee5b4d4d58622ce2ae0ff1f",
		},
		{
			file:     "README_crlf.md",
			expected: "01f253d1e123025268e0ba806c9d1dce413d0075",
		},
		{
			file:     "gamesList.txt",
			expected: "408cb3fd0985ff362ee4f9a484227854ecee0c69",
		},
		{
			file:     "go1_crlf.txt",
			expected: "02c3e857243f2feec3d5da8b2909707cbfa0dbcb",
		},
		{
			file:     "go1_lf.txt",
			expected: "aea15f18d082f7517a2520b2baebd4c5d5baabaf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			p := path.Join(otherPlatformFileDir, tt.file)
			f, err := os.Open(p)
			assert.Nilf(t, err, "expected open %s to be successful, got %v", p, err)
			hash, err := OtherPlatformSha1(f)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, hash)
		})
	}
}
