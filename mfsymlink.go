// Package mfsymlink provides methods for parsing Minshall-French symlinks
//
// See: https://wiki.samba.org/index.php/UNIX_Extensions#Minshall.2BFrench_symlinks
package mfsymlink

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
)

const (
	// Size is the size of mfsymlink files in bytes
	Size = 1067
	// Marker is the marker used at the front of mfsymlink files
	Marker      = "XSym"
	emptyString = ""
)

var (
	// ErrNotMFSymlink signifies that the given bytes/file is not a mfsymlink
	ErrNotMFSymlink = errors.New("not a mfsymlink")
	// ErrMD5Mismatch signifies a corrupt mfsymlink due to the md5 checksum not matching
	ErrMD5Mismatch  = errors.New("corrupt mfsymlink: md5 checksum mismatch")
	mfsymlinkMarker = []byte(Marker)
	newline         = []byte("\n")
)

// IsPossibleSymlink determines if the given os.FileInfo is a possible mfsymlink
func IsPossibleSymlink(i os.FileInfo) bool {
	return i.Size() == Size
}

// Parse parses the given bytes and returns the mfsymlink target
func Parse(b []byte) (string, error) {
	lines := bytes.SplitN(b, newline, 4)
	if len(lines) != 4 {
		return emptyString, ErrNotMFSymlink
	}
	// Check file byte marker
	if !bytes.Equal(lines[0], mfsymlinkMarker) {
		return emptyString, ErrNotMFSymlink
	}
	l, err := strconv.Atoi(string(lines[1]))
	if err != nil {
		return emptyString, ErrNotMFSymlink
	}
	expectedMD5 := make([]byte, hex.DecodedLen(len(lines[2])))
	if _, err := hex.Decode(expectedMD5, lines[2]); err != nil {
		return emptyString, ErrNotMFSymlink
	}
	targetBytes := lines[3]
	if len(targetBytes) > l {
		targetBytes = targetBytes[:l]
	}
	actualMD5 := md5.Sum(targetBytes)
	if !bytes.Equal(expectedMD5, actualMD5[:]) {
		return emptyString, ErrMD5Mismatch
	}
	return string(targetBytes), nil
}
