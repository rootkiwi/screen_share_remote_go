// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package byteutil

func ClearBytes(bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		bytes[i] = 0
	}
}
