package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

var (
	iv = []byte("PresslyWarpdrive") //it needs to be 16 bytes
)

// MakeAESKey generates random key for aes
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

func padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func AESEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext = padding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func AESDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return unpadding(ciphertext), nil
}

func AESEncryptStream2(key []byte, input io.Reader, output io.Writer) error {
	aes, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	enc := cipher.NewCBCEncrypter(aes, iv)
	buf := make([]byte, enc.BlockSize())

	for {
		_, err = io.ReadFull(input, buf)
		if err != nil {
			if err == io.EOF {
				break
			} else if err == io.ErrUnexpectedEOF {
				// nothing
			} else {
				return err
			}
		}
		enc.CryptBlocks(buf, buf)
		if _, err = output.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

func AESDecryptStream2(key []byte, input io.Reader, output io.Writer) error {
	aes, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	dec := cipher.NewCBCDecrypter(aes, iv)
	buf := make([]byte, dec.BlockSize())

	for {
		_, err := io.ReadFull(input, buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		dec.CryptBlocks(buf, buf)
		if _, err = output.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

func AESEncryptStream(key []byte, input io.Reader, output io.Writer) error {
	block, err := aes.NewCipher(key)

	if err != nil {
		return err
	}

	// beacuse every bundle in this packge uses different key, then this is
	// going to be ok if we use zero IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	writer := &cipher.StreamWriter{S: stream, W: output}
	if _, err := io.Copy(writer, input); err != nil {
		return err
	}

	return nil
}

func AESDecryptStream(key []byte, input io.Reader, output io.Writer) error {
	block, err := aes.NewCipher(key)

	if err != nil {
		return err
	}

	// beacuse every bundle in this packge uses different key, then this is
	// going to be ok if we use zero IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: input}
	if _, err := io.Copy(output, reader); err != nil {
		return err
	}

	return nil
}

func NewAESEncryptWriter(key []byte, output io.Writer) io.Writer {
	// whatever I write into `w`, will be avilable in `r`
	r, w := io.Pipe()

	go func() {
		AESEncryptStream(key, r, output)
	}()

	return w
}

func NewAESDecryptWriter(key []byte, output io.Writer) io.Writer {
	// whatever I write into `w`, will be avilable in `r`
	r, w := io.Pipe()

	go func() {
		AESDecryptStream(key, r, output)
	}()

	return w
}
