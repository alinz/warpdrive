package crypto_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pressly/warpdrive/lib/crypto"
)

func TestHash(t *testing.T) {
	f, err := os.Open("plaintext")
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	hash, err := crypto.Hash(f)
	if err != nil {
		t.Error(err)
	}

	actualHash := "22596363b3de40b06f981fb85d82312e8c0ed511"
	result := fmt.Sprintf("%x", hash)

	if result != actualHash {
		t.Errorf("got %s instead of %s", result, actualHash)
	}
}
