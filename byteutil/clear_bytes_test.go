// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package byteutil

import (
	"bytes"
	"testing"
)

func TestClearBytes(t *testing.T) {
	tests := []struct {
		actual, expected []byte
	}{
		{
			actual:   []byte{1},
			expected: []byte{0},
		},
		{
			actual:   []byte{1, 2},
			expected: []byte{0, 0},
		},
		{
			actual:   []byte{0, 12, 37, 255, 111, 1},
			expected: []byte{0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		ClearBytes(tt.actual)
		if !bytes.Equal(tt.actual, tt.expected) {
			t.Error("array not cleared")
		}
	}
}
