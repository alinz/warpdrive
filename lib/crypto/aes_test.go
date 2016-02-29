package crypto_test

import (
	"testing"

	"github.com/pressly/warpdrive/lib/crypto"
)

func TestAES(t *testing.T) {
	plaintext := "this is an apple."
	key, err := crypto.MakeAESKey(16)
	if err != nil {
		t.Error(err)
	}

	encrypted, err := crypto.AESCrypt([]byte(plaintext), key)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := crypto.AESCrypt(encrypted, key)
	if err != nil {
		t.Error(err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("got %s instead of %s", string(decrypted), plaintext)
	}
}
