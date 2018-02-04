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

	"github.com/rootkiwi/screen_share_remote_go/cert"
	"github.com/rootkiwi/screen_share_remote_go/conf"
)

func NoConf() *conf.Config {
	printNoConfInfo()
	scan := bufio.NewScanner(os.Stdin)
	port := readPortNumber(scan)
	webPort := readWebPortNumber(scan)
	passwordHash := readPasswordHash()
	fmt.Print("\nGenerating a 4096-bit RSA key pair and a self-signed certificate...")
	fingerprint, cert := cert.GenerateNoConf()
	fmt.Println(" done.")
	fmt.Println()
	fmt.Println("certificate fingerprint:")
	fmt.Println(fingerprint)
	fmt.Println()
	return &conf.Config{
		Port:         port,
		WebPort:      webPort,
		PasswordHash: passwordHash,
		Cert:         cert,
	}
}
