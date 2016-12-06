package crypto_test

import (
	"bytes"
	"fmt"
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

func TestAESStream(t *testing.T) {
	var given bytes.Buffer
	var input bytes.Buffer
	var output bytes.Buffer
	var final bytes.Buffer

	message := "1"

	given.WriteString(message)
	input.WriteString(message)

	key, err := crypto.MakeAESKey(32)
	if err != nil {
		t.Error(err.Error())
	}

	err = crypto.AESEncryptStream(key, &input, &output)
	if err != nil {
		t.Error(err.Error())
	}

	err = crypto.AESDecryptStream(key, &output, &final)
	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(input.Bytes(), final.Bytes()) {
		t.Log("given:", given.Len(), given.String())
		t.Log("result:", final.Len(), final.String())
		t.Error(fmt.Errorf("final decrypted message is not the same as given input"))

	}
}
