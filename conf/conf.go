// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package conf

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rootkiwi/screen_share_remote_go/byteutil"
	"github.com/rootkiwi/screen_share_remote_go/password"
	"github.com/rootkiwi/screen_share_remote_go/portutil"
	"github.com/rootkiwi/screen_share_remote_go/sha256"
)

func CreateConfigFile(port, webPort int, fingerprint, passwordHash, base64Cert string, base64Key []byte) {
	createSimple(port, webPort, fingerprint, passwordHash, base64Cert, base64Key)
}

func createSimple(port, webPort int, fingerprint, passwordHash, base64Cert string, base64Key []byte) {
	file, err := os.Create("screen_share_remote.conf")
	if err != nil {
		log.Fatalf("error creating screen_share_remote.conf: %v\n", err)
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("port    (changeable): %d\n", port))
	buf.WriteString(fmt.Sprintf("webPort (changeable): %d\n", webPort))
	buf.WriteString(fmt.Sprintf("fingerprint: %s\n", fingerprint))
	buf.WriteString(fmt.Sprintln(passwordHash))
	buf.WriteString(fmt.Sprintln(base64Cert))
	buf.Write(base64Key)
	buf.WriteRune('\n')
	byteutil.ClearBytes(base64Key)
	_, err = io.Copy(file, bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Fatalf("error writing screen_share_remote.conf: %v\n", err)
	}
}

type Config struct {
	Port         int
	WebPort      int
	PasswordHash string
	Cert         *tls.Certificate
}

func ParseConfigFile(path string) (*Config, error) {
	return parseSimple(path)
}

func parseSimple(path string) (*Config, error) {
	fail := func(err error) (*Config, error) { return nil, err }

	file, err := os.Open(path)
	if err != nil {
		return fail(fmt.Errorf("error opening config: %v", err))
	}
	defer file.Close()

	var lines [6][]byte
	scanner := bufio.NewScanner(file)
	for i := 0; i < 6; i++ {
		if scanner.Scan() {
			lines[i] = append([]byte{}, scanner.Bytes()...)
		} else {
			return fail(errors.New("invalid config: too few lines"))
		}
	}
	portLine := string(lines[0])
	webPortLine := string(lines[1])
	fingerprintLine := string(lines[2])
	passwordHashLine := string(lines[3])
	certificateLine := string(lines[4])
	privateKeyLine := lines[5]

	portLineSplit := strings.Split(portLine, ":")
	if len(portLineSplit) != 2 {
		return fail(fmt.Errorf("invalid config port line: %q", portLine))
	}
	port, err := portutil.Parse(strings.TrimSpace(portLineSplit[1]))
	if err != nil {
		return fail(fmt.Errorf("invalid config port: %v", err))
	}

	webPortLineSplit := strings.Split(webPortLine, ":")
	if len(webPortLineSplit) != 2 {
		return fail(fmt.Errorf("invalid config webPort line: %q", webPortLine))
	}
	webPort, err := portutil.Parse(strings.TrimSpace(webPortLineSplit[1]))
	if err != nil {
		return fail(fmt.Errorf("invalid webPort config: %v", err))
	}

	if !password.ValidHash(passwordHashLine) {
		return fail(errors.New("invalid config: invalid password hash"))
	}

	certBytes, err := base64.StdEncoding.DecodeString(certificateLine)
	if err != nil {
		return fail(fmt.Errorf("invalid config certificate: %v", err))
	}

	fingerprintLineSplit := strings.Split(fingerprintLine, ":")
	if len(fingerprintLineSplit) != 2 {
		return fail(fmt.Errorf("invalid config fingerprint line: %q", fingerprintLine))
	}
	if sha256.PrettySum(certBytes) != strings.TrimSpace(fingerprintLineSplit[1]) {
		return fail(errors.New("invalid config: fingerprint not matching certificate"))
	}

	privKeyBytes := make([]byte, base64.StdEncoding.DecodedLen(len(privateKeyLine)))
	n, err := base64.StdEncoding.Decode(privKeyBytes, privateKeyLine)
	if err != nil {
		return fail(fmt.Errorf("invalid config private key: %v", err))
	}
	privKeyBytes = privKeyBytes[:n]

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKeyBytes})
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return fail(fmt.Errorf("invalid config certificate or key: %v", err))
	}
	byteutil.ClearBytes(privateKeyLine)
	byteutil.ClearBytes(privKeyBytes)
	byteutil.ClearBytes(keyPem)

	return &Config{
		Port:         port,
		WebPort:      webPort,
		PasswordHash: passwordHashLine,
		Cert:         &cert,
	}, nil
}
