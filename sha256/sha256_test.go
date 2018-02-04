// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package sha256

import "testing"

func TestPrettySum(t *testing.T) {
	tests := []struct {
		args []byte
		want string
	}{
		{[]byte{0, 1, 2, 3, 4, 5}, "17E88DB187AFD62C16E5DEBF3E6527CD006BC012BC90B51A810CD80C2D511F43"},
		{[]byte{123, 125, 1, 0, 32, 255, 182, 192}, "38B7A0D5A3F66D3CC8CE861162751FDF67EBAF1E999D26088185D1C76F4C9129"},
	}
	for _, tt := range tests {
		if got := PrettySum(tt.args); got != tt.want {
			t.Errorf("PrettySum() = %v, want %v", got, tt.want)
		}
	}
}
