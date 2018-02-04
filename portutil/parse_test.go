// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package portutil

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		argPort  string
		wantPort int
		wantErr  bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"4568", 4568, false},
		{"65535", 65535, false},
		{"-1", 0, true},
		{"65536", 0, true},
		{"", 0, true},
		{"a", 0, true},
	}
	for _, tt := range tests {
		gotPort, err := Parse(tt.argPort)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) error = %v, wantErr %v", tt.argPort, err, tt.wantErr)
		}
		if gotPort != tt.wantPort {
			t.Errorf("Parse(%q) = %v, want %v", tt.argPort, gotPort, tt.wantPort)
		}
	}
}
