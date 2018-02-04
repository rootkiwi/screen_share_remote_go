// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rootkiwi/screen_share_remote_go/cert"
	"github.com/rootkiwi/screen_share_remote_go/conf"
)

const (
	defaultPort    = 50000
	defaultWebPort = 8081
)

func GenConf() {
	printGenConfInfo()
	scan := bufio.NewScanner(os.Stdin)
	port := readPortNumber(scan)
	webPort := readWebPortNumber(scan)
	passwordHash := readPasswordHash()
	fmt.Print("\nGenerating a 4096-bit RSA key pair and a self-signed certificate...")
	fingerprint, base64cert, base64key := cert.GenerateForConf()
	fmt.Println(" done.")
	fmt.Println()
	fmt.Println("certificate fingerprint:")
	fmt.Println(fingerprint)
	conf.CreateConfigFile(port, webPort, fingerprint, passwordHash, base64cert, base64key)
	fmt.Println()
	fmt.Println("config file created:")
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("/???/screen_share_remote.conf")
		fmt.Printf("error getting working directory: %v\n", err)
	} else {
		fmt.Println(filepath.Join(wd, "screen_share_remote.conf"))
	}
	fmt.Println()
	fmt.Println("The settings 'port' and 'webPort' is changeable, the rest is not")
	fmt.Println("If you need to change the password/certificate run genconf again")
}
