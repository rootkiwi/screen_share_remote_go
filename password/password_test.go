// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package password

import (
	"bytes"
	"testing"
)

func TestGenRandom(t *testing.T) {
	random1 := GenRandom()
	random2 := GenRandom()
	if bytes.Equal(random1, random2) {
		t.Errorf("%s == %s, want random", random1, random2)
	}
}

func TestValidateTrue(t *testing.T) {
	tests := []struct {
		hash     string
		password []byte
	}{
		{Hash([]byte(`password`)), []byte(`password`)},
		{Hash([]byte(`ïØÎ#àÿH½½¿õcÒóYÅ#KY¡SqÙY#·¾,åøµ]ì)N´Îû4ª`)), []byte(`ïØÎ#àÿH½½¿õcÒóYÅ#KY¡SqÙY#·¾,åøµ]ì)N´Îû4ª`)},
		{Hash([]byte(` `)), []byte(` `)},
		{Hash([]byte(``)), []byte(``)},
	}
	for i, tt := range tests {
		if !Validate(tt.hash, tt.password) {
			t.Errorf("%d: Validate() = false, want true", i)
		}
	}
}

func TestValidateTrueRandom(t *testing.T) {
	for i := 0; i < 5; i++ {
		random := GenRandom()
		randomString := string(random)
		hash := Hash([]byte(randomString))
		if !Validate(hash, []byte(randomString)) {
			t.Errorf("Validate(%q, %q) = false, want true", hash, randomString)
		}
	}
}

func TestValidateFalse(t *testing.T) {
	tests := []struct {
		hash     string
		password []byte
	}{
		{Hash([]byte(` password`)), []byte(`password`)},
		{Hash([]byte(`password`)), []byte(`password `)},
		{Hash([]byte(`password `)), []byte(`password`)},
		{Hash([]byte(`ïØÎ#àÿH½½¿õcÒóYÅ#KY¡SqÙY#·¾,åøµ]ì)N´Îû4ª`)), []byte(`ØÎ#àÿH½½¿õcÒóYÅ#KY¡SqÙY#·¾,åøµ]ì)N´Îû4ª`)},
		{Hash([]byte(` `)), []byte(``)},
		{Hash([]byte(`123`)), []byte(`12`)},
	}
	for i, tt := range tests {
		if Validate(tt.hash, tt.password) {
			t.Errorf("%d: Validate() = true, want false", i)
		}
	}
}
