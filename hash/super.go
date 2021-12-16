package hash

import (
	"crypto"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"math"
)

const (
	smallBucketSize       = 1024 * 1.25
	maxSizeForSuperHash   = 2 * 1024 * 1024 * 1024
	minSizeForFullHash    = 512
	minSizeForPartialHash = 2 * 1024
	smallSourceSize       = 3 * 1024
)

var (
	ErrTooLargeForSuperHash   = fmt.Errorf("reader is too large for super hash calculation")
	ErrTooSmallForSuperHash   = fmt.Errorf("reader is too small for super hash calculation")
	ErrTooSmallForSigByteHash = fmt.Errorf("reader is too small for msb and lsb super hash calculations")
)

// SuperHash contains various hashes of a source file.
// All hashes are performed after stripping away all whitespace from the file
type SuperHash struct {
	// FullHash contains a hash of the whole file
	FullHash string

	// MosSigBitsHash contains a hash of the first bits of the file. Amount of bits determined by bucket size.
	MostSigBitsHash string

	// LeastSigBitsHash contains a hash of the last bits of the file. Amount of bits determined by bucket size.
	LeastSigBitsHash string

	// The algorithm to use for calculating hashes
	algo crypto.Hash
}

func CalculateSuperSha1(rs io.ReadSeeker) (*SuperHash, error) {
	return CalculateSuperHash(rs, crypto.SHA1)
}

// CalculateSuperHash calculates the super hash of the provided reader, with the provided hash algorithm
func CalculateSuperHash(rs io.ReadSeeker, algo crypto.Hash) (h *SuperHash, err error) {
	size, err := calcSizeNoWhiteSpace(rs)
	if err != nil {
		return nil, err
	}

	if size > maxSizeForSuperHash {
		err = ErrTooLargeForSuperHash
		return
	}

	if size < minSizeForFullHash {
		err = ErrTooSmallForSuperHash
		return
	}

	h = &SuperHash{algo: algo}

	_, _ = rs.Seek(0, io.SeekStart)
	if h.FullHash, err = h.fullHash(rs); err != nil {
		return
	}

	if size < minSizeForPartialHash {
		err = ErrTooSmallForSigByteHash
		return
	}

	bucketSize := h.bucketSize(float64(size))

	_, _ = rs.Seek(0, io.SeekStart)
	if h.MostSigBitsHash, err = h.msbHash(rs, int64(bucketSize)); err != nil {
		return
	}

	_, _ = rs.Seek(0, io.SeekStart)
	if h.LeastSigBitsHash, err = h.lsbHash(rs, size, bucketSize); err != nil {
		return
	}
	return h, nil
}

func calcSizeNoWhiteSpace(reader io.Reader) (int, error) {
	r := &countingReader{reader: transform.NewReader(reader, noWhitespaceTransformer{})}
	_, err := io.ReadAll(r)
	return r.count, err
}

func (h *SuperHash) fullHash(reader io.Reader) (string, error) {
	return Hash(transform.NewReader(reader, noWs), h.algo)
}

func (h *SuperHash) lsbHash(reader io.Reader, sizeNoWs, offset int) (string, error) {
	r := &offsetReader{
		reader: transform.NewReader(reader, noWs),
		offset: sizeNoWs - offset,
	}
	hash := h.algo.New()
	_, err := io.Copy(hash, r)
	return hashHex(hash), err
}

func (h *SuperHash) msbHash(reader io.Reader, bucketSize int64) (string, error) {
	return Hash(io.LimitReader(transform.NewReader(reader, noWs), bucketSize), h.algo)
}

// bucketSize calculates the amount of bits to hash from the beginning / end of the file
// the algorithm is copied (with some aesthetic improvements) from the WhiteSource agents repository
func (h *SuperHash) bucketSize(size float64) int {
	var base, high, low float64

	if size <= smallSourceSize {
		return smallBucketSize
	}

	base = math.Pow10(int(math.Log10(size)))
	high = math.Ceil((size+1)/base) * base
	low = high - base
	return int(high+low) / 4
}

func OtherPlatformSha1(reader io.Reader) (string, error) {
	return OtherPlatformHash(reader, crypto.SHA1)
}

// OtherPlatformHash calculates the hash of the provided io.Reader contents,
// after replacing the line endings to the other platform type.
// if the line endings are CRLF, they will be switched to LF
// if the line endings are LF, they will be switched to CRLF
func OtherPlatformHash(reader io.Reader, algo crypto.Hash) (string, error) {
	return Hash(transform.NewReader(reader, &lineEndingTransformer{}), algo)
}
