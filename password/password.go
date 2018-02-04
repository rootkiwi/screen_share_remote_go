// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package password

import (
	"bytes"
	"crypto/rand"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/rootkiwi/screen_share_remote_go/byteutil"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password []byte) (hash string) {
	hash = bcryptGenerate(password)
	byteutil.ClearBytes(password)
	return hash
}

func Validate(hash string, password []byte) (valid bool) {
	valid = bcryptValidate(hash, password)
	byteutil.ClearBytes(password)
	return valid
}

func GenRandom() (random []byte) {
	return simpleRandom()
}

func ValidHash(hash string) bool {
	// $2a$10$AOFSBPq1QVsqU.oAlN5/L.ExuBEZzWMADdcwOqvxWwt7Osu9yfCZW
	split := strings.Split(hash, "$")
	if len(split) != 4 {
		return false
	}
	switch split[1] {
	case "2a", "2b", "2x", "2y":
	default:
		return false
	}
	cost, err := strconv.Atoi(split[2])
	if err != nil {
		return false
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return false
	}
	// check if valid secret hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), nil)
	if err != bcrypt.ErrMismatchedHashAndPassword {
		return false
	}

	return true
}

func bcryptGenerate(password []byte) (hash string) {
	bcryptHash, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Fatal(err)
	}
	return string(bcryptHash)
}

func bcryptValidate(hash string, password []byte) (valid bool) {
	valid = bcrypt.CompareHashAndPassword([]byte(hash), password) == nil
	return valid
}

func simpleRandom() (random []byte) {
	const length = 40
	alphabet := []rune("123456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz")
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		buffer.WriteRune(alphabet[randomIndex.Int64()])
	}
	return buffer.Bytes()
}
