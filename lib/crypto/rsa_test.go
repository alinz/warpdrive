package crypto_test

import (
	"testing"

	"github.com/pressly/warpdrive/lib/crypto"
)

func TestParsePublicKey(t *testing.T) {
	_, public, err := crypto.RSAKeyPair(1024)

	if err != nil {
		t.Error(err)
	}

	_, err = crypto.PublicKey(public)

	if err != nil {
		t.Error(err)
	}
}

func TestParsePrivateKey(t *testing.T) {
	private, _, err := crypto.RSAKeyPair(1024)

	if err != nil {
		t.Error(err)
	}

	_, err = crypto.PrivateKey(private)

	if err != nil {
		t.Error(err)
	}
}

func TestRSA(t *testing.T) {
	private, public, err := crypto.RSAKeyPair(1024)

	if err != nil {
		t.Error(err)
	}

	message := "Hello This is Warpdrive."

	encrypted, err := crypto.EncryptByPublicRSA([]byte(message), public, "sha256")
	if err != nil {
		t.Error(err)
	}

	decrypted, err := crypto.DecryptByPrivateRSA(encrypted, private, "sha256")
	if err != nil {
		t.Error(err)
	}

	if string(decrypted) != message {
		t.Errorf("got '%s' instead of '%s'", decrypted, message)
	}
}
