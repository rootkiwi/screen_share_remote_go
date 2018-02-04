// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package cli

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/rootkiwi/screen_share_remote_go/byteutil"
	"github.com/rootkiwi/screen_share_remote_go/password"
	"golang.org/x/crypto/ssh/terminal"
)

func readPasswordHash() (passwordHash string) {
	readPassword := func(prompt string) []byte {
		fmt.Print(prompt)
		pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("error getting terminal: %v\n", err)
		}
		fmt.Println()
		return pass
	}
	for {
		pass := readPassword("3. enter password [random]: ")
		if len(pass) == 0 {
			randomPassword := password.GenRandom()
			fmt.Printf("\npassword:\n%s\n", randomPassword)
			return password.Hash(randomPassword)
		} else {
			passAgain := readPassword("3. enter password again: ")
			if bytes.Equal(pass, passAgain) {
				byteutil.ClearBytes(passAgain)
				return password.Hash(pass)
			} else {
				fmt.Print("password did not match\n\n")
				byteutil.ClearBytes(pass)
				byteutil.ClearBytes(passAgain)
			}
		}
	}
}
