// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package cli

import (
	"bufio"
	"fmt"

	"github.com/rootkiwi/screen_share_remote_go/portutil"
)

func readPortNumber(scan *bufio.Scanner) int {
	return readPort(scan, 1, "port", defaultPort)
}

func readWebPortNumber(scan *bufio.Scanner) int {
	return readPort(scan, 2, "web server port", defaultWebPort)
}

func readPort(scan *bufio.Scanner, promptNum int, portName string, defaultPort int) int {
	for {
		fmt.Printf("%d. enter %s number (0-65535) [%d]: ", promptNum, portName, defaultPort)
		scan.Scan()
		portInput := scan.Text()
		if len(portInput) == 0 {
			return defaultPort
		}
		port, err := portutil.Parse(portInput)
		if err != nil {
			fmt.Println(err)
			continue
		}
		return port
	}
}
