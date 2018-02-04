// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package cli

import (
	"fmt"
	"os"
)

func PrintUsage() {
	fmt.Println("Usage: ./screen_share_remote_go-<VERSION>-<PLATFORM> [/path/to/conf | genconf | noconf]")
	fmt.Println()
	fmt.Println("Example: ./screen_share_remote_go-0.1.0-linux-amd64 genconf")
	fmt.Println()
	fmt.Println("Available command-line options:")
	fmt.Println("1. /path/to/screen_share_remote.conf (start screen_share_remote)")
	fmt.Println("2. genconf                           (generate new config file)")
	fmt.Println("3. noconf                            (start without saving config)")
}

func printGenConfInfo() {
	fmt.Println("Generate screen_share_remote.conf file in working directory which is:")
	wd, err := os.Getwd()
	if err != nil {
		wd = fmt.Sprintf("error getting working directory: %v", err)
	}
	fmt.Println(wd)
	fmt.Println()
	fmt.Println("Will overwrite if already exists")
	fmt.Println()
	fmt.Println("The config file will contain these attributes:")
	printConfigItems()
	fmt.Println()
	fmt.Println("Do note that the RSA private key is stored in cleartext, so make sure")
	fmt.Println("to make the config file inaccessible for unauthorized parties.")
	fmt.Println("Or you could run in the 'noconf' mode, which means a new private key will be")
	fmt.Println("generated each time. Without saving to disk.")
	printInputInfo()
}

func printConfigItems() {
	fmt.Println("1. port number                 (port number to enter in screen_share)")
	fmt.Println("2. web server port number      (port the web server will serve on)")
	fmt.Println("3. password                    (password to enter in screen_share)")
	fmt.Println("4. self-signed TLS certificate (whose fingerprint to enter in screen_share)")
	fmt.Println("5. RSA private key             (corresponding to certificate)")
}

func printInputInfo() {
	fmt.Println()
	fmt.Println("Leave empty and press ENTER for the [default] value")
}

func printNoConfInfo() {
	fmt.Println("Start in 'noconf' mode, no config will be saved to disk")
	fmt.Println()
	fmt.Println("These attributes are needed:")
	printConfigItems()
	printInputInfo()
}
