// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"log"
	"math/big"
	"time"

	"github.com/rootkiwi/screen_share_remote_go/byteutil"
	"github.com/rootkiwi/screen_share_remote_go/sha256"
)

func GenerateForConf() (fingerprint, base64cert string, base64key []byte) {
	key, certBytes, fingerprint := genKeyAndCert()
	base64cert = base64.StdEncoding.EncodeToString(certBytes)
	keyDerBytes := x509.MarshalPKCS1PrivateKey(key)
	base64key = make([]byte, base64.StdEncoding.EncodedLen(len(keyDerBytes)))
	base64.StdEncoding.Encode(base64key, keyDerBytes)
	byteutil.ClearBytes(keyDerBytes)
	return fingerprint, base64cert, base64key
}

func GenerateNoConf() (fingerprint string, cert *tls.Certificate) {
	key, certBytes, fingerprint := genKeyAndCert()
	keyDerBytes := x509.MarshalPKCS1PrivateKey(key)
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDerBytes})
	byteutil.ClearBytes(keyDerBytes)
	newCert, err := tls.X509KeyPair(certPem, keyPem)
	byteutil.ClearBytes(keyPem)
	if err != nil {
		log.Fatalf("error creating certificate: %v\n", err)
	}
	return fingerprint, &newCert
}

func genKeyAndCert() (key *rsa.PrivateKey, certBytes []byte, fingerprint string) {
	key = genKey()
	certBytes = genCert(key)
	fingerprint = sha256.PrettySum(certBytes)
	return key, certBytes, fingerprint
}

func genKey() *rsa.PrivateKey {
	const keySize = 4096
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Panicln(err)
	}
	return key
}

func genCert(key *rsa.PrivateKey) []byte {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"screen_share_remote_go"},
		},
		NotBefore:             time.Now(),
		NotAfter:              endOfTime,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		log.Fatal(err)
	}
	return derBytes
}
