package crypto

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

const filechunk = 8192

//HashFile accepting a filename and convert it into secure hash value
func HashFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash, err := Hash(file)

	return fmt.Sprintf("%x", hash), err
}

//Hash hashes any io.Reader into a 20 byte hash
func Hash(src io.Reader) ([]byte, error) {
	var result []byte

	hash := sha1.New()
	if _, err := io.Copy(hash, src); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}
