package mfsymlink_test

import (
	"os"
	"testing"
	"time"
)

import (
	"github.com/dhui/mfsymlink"
)

type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
}

func (i *mockFileInfo) Name() string       { return i.name }
func (i *mockFileInfo) Size() int64        { return i.size }
func (i *mockFileInfo) Mode() os.FileMode  { return i.mode }
func (i *mockFileInfo) ModTime() time.Time { return i.modTime }
func (i *mockFileInfo) IsDir() bool        { return i.mode.IsDir() }
func (i *mockFileInfo) Sys() interface{}   { return i.sys }

var _ os.FileInfo = &mockFileInfo{}

const validMFSymlink = `XSym
0026
500e5dbfa9b8c0041e01fe4f7967e287
../XXXX-XX-XX/XXXXXXXXXXXX
																																																																																																																																																																																																																																																									 `

func TestIsPossibleSymlink(t *testing.T) {
	testCases := []struct {
		name     string
		fileInfo os.FileInfo
		expected bool
	}{
		{name: "size too small", fileInfo: &mockFileInfo{size: 10}, expected: false},
		{name: "size too large", fileInfo: &mockFileInfo{size: 10000}, expected: false},
		{name: "success", fileInfo: &mockFileInfo{size: mfsymlink.Size}, expected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := mfsymlink.IsPossibleSymlink(tc.fileInfo); actual != tc.expected {
				t.Error("Failed test")
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name           string
		bytes          []byte
		expectedTarget string
		expectedErr    error
	}{
		{name: "nil bytes", bytes: nil, expectedTarget: "", expectedErr: mfsymlink.ErrNotMFSymlink},
		{name: "empty bytes", bytes: []byte(""), expectedTarget: "", expectedErr: mfsymlink.ErrNotMFSymlink},
		{name: "no marker", bytes: []byte("\n\n\n"), expectedTarget: "", expectedErr: mfsymlink.ErrNotMFSymlink},
		{name: "invalid marker", bytes: []byte("jjj\n\n\n"), expectedTarget: "", expectedErr: mfsymlink.ErrNotMFSymlink},
		{name: "non-numeric length", bytes: []byte("XSym\njj\n\n"), expectedTarget: "", expectedErr: mfsymlink.ErrNotMFSymlink},
		{name: "malformed md5", bytes: []byte("XSym\n10\nx\n"), expectedTarget: "", expectedErr: mfsymlink.ErrNotMFSymlink},
		{name: "invalid md5", bytes: []byte("XSym\n10\n\n"), expectedTarget: "", expectedErr: mfsymlink.ErrMD5Mismatch},
		{name: "no trailing padding", bytes: []byte("XSym\n10\nc59195470191ddf4c0f9e54e33046386\nXXXXXXXXXX"), expectedTarget: "XXXXXXXXXX", expectedErr: nil},
		{name: "success", bytes: []byte(validMFSymlink), expectedTarget: "../XXXX-XX-XX/XXXXXXXXXXXX", expectedErr: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if target, err := mfsymlink.Parse(tc.bytes); err != tc.expectedErr {
				t.Error("Error mismatch:", err, "!=", tc.expectedErr)
			} else if target != tc.expectedTarget {
				t.Error("Target mismatch:", target, "!=", tc.expectedTarget)
			}
		})
	}
}
