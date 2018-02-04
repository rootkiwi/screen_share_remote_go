// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package sha256

import (
	"crypto/sha256"
	"fmt"
)

func PrettySum(data []byte) string {
	return fmt.Sprintf("%064X", sha256.Sum256(data))
}
