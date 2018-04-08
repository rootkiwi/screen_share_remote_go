// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package portutil

import (
	"fmt"
	"strconv"
)

func Parse(s string) (port int, err error) {
	port, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q (not a number)", s)
	}
	if port < 0 || port > 65535 {
		return 0, fmt.Errorf("%q (0-65535)", s)
	}
	return port, nil
}
