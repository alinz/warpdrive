package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

var (
	iv = []byte("PresslyWarpdrive") //it needs to be 16 bytes
)

//MakeAESKey generates random key for aes
func MakeAESKey(size int) ([]byte, error) {
	switch size {
	case 16, 24, 32:
	default:
		return nil, errors.New("size must be 16, 24 or 32")
	}

	key := make([]byte, size)
	n, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read new random key: %s", err)
	}
	if n < size {
		return nil, fmt.Errorf("failed to read entire key, only read %d out of %d", n, size)
	}
	return key, nil
}

//AESCrypt encrypt or decrypt data
func AESCrypt(data, key []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(blockCipher, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
