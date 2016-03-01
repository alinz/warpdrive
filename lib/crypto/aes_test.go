package crypto_test

import (
	"os"
	"testing"

	"github.com/pressly/warpdrive/lib/crypto"
)

func TestAES(t *testing.T) {
	plaintext := "this is an apple."
	key, err := crypto.MakeAESKey(16)
	if err != nil {
		t.Error(err)
	}

	encrypted, err := crypto.AESEncrypt([]byte(plaintext), key)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := crypto.AESDecrypt(encrypted, key)
	if err != nil {
		t.Error(err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("got %s instead of %s", string(decrypted), plaintext)
	}
}

func TestAESAsFile(t *testing.T) {
	plaintext := "this is an apple"
	key := "presslywarpdrive"

	encrypted, err := crypto.AESEncrypt([]byte(plaintext), []byte(key))
	if err != nil {
		t.Error(err)
	}

	out, err := os.Create("/Users/ali/temp/output.enc")
	if err != nil {
		t.Error(err)
	}

	out.Write(encrypted)
	out.Close()
}
