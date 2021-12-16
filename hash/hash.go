package hash

import (
	"bytes"
	"crypto"
	_ "crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"strings"
)

const (
	Sha1           = "SHA1"
	AdditionalSha1 = "ADDITIONAL_SHA1"
)

func BytesSha1(bs []byte) string {
	return BytesHash(bs, crypto.SHA1)
}

func BytesHash(bs []byte, algo crypto.Hash) string {
	h, _ := Hash(bytes.NewReader(bs), algo)
	return h
}

func StringSha1(str string) string {
	return StringHash(str, crypto.SHA1)
}

func StringHash(str string, algo crypto.Hash) string {
	h, _ := Hash(strings.NewReader(str), algo)
	return h
}

func FileSha1(path string) (string, error) {
	return FileHash(path, crypto.SHA1)
}

func FileHash(path string, algo crypto.Hash) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}
	return Hash(file, algo)
}

func Hash(reader io.Reader, algo crypto.Hash) (string, error) {
	h := algo.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}
	return hashHex(h), nil
}

func hashHex(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}
