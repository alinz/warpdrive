package crypto

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

//PublicKey gets the encoded version of public key and returns a rsa.PublicKey
//struct
func PublicKey(publicKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKey))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Value returned from ParsePKIXPublicKey was not an RSA public key")
	}

	return rsaPub, nil
}

//PrivateKey gets the encoded version of private key and returns a
//rsa.PrivateKey struct
func PrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	rsaPriv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsaPriv, nil
}

func hashAlgorithm(algo string) (hash.Hash, crypto.Hash) {
	switch algo {
	case "md5":
		return md5.New(), crypto.MD5
	case "sha1":
		return sha1.New(), crypto.SHA1
	case "sha256":
		return sha256.New(), crypto.SHA256
	case "sha512":
		return sha512.New(), crypto.SHA512
	default:
		log.Fatalf("%s is not a valid hash algorithm. Must be one of md5, sha1, sha256, sha512", algo)
	}
	panic("something is wrong")
}

//EncryptByPublicRSA encrypts using Public key and message
func EncryptByPublicRSA(message []byte, publicKey, algo string) ([]byte, error) {
	public, err := PublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	//func EncryptPKCS1v15(rand io.Reader, pub *PublicKey, msg []byte) (out []byte, err error)

	//h, _ := hashAlgorithm(algo)
	//encrypted, err := rsa.EncryptOAEP(h, rand.Reader, public, message, nil)
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, public, message)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

//DecryptByPrivateRSA decrypt message using Private key
func DecryptByPrivateRSA(message []byte, privateKey, algo string) ([]byte, error) {
	key, err := PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	//func DecryptPKCS1v15(rand io.Reader, priv *PrivateKey, ciphertext []byte) (out []byte, err error)

	//h, _ := hashAlgorithm(algo)
	//plaintext, err := rsa.DecryptOAEP(h, rand.Reader, key, message, nil)
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, key, message)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return plaintext, nil
}

//RSAKeyPair returns a pair of encoded string version of public/private keys
func RSAKeyPair(size int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return "", "", err
	}

	var privateBuffer bytes.Buffer

	pemkey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	pem.Encode(&privateBuffer, pemkey)

	var publicBuffer bytes.Buffer

	bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return "", "", err
	}

	pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes,
	}

	pem.Encode(&publicBuffer, pemkey)
	return privateBuffer.String(), publicBuffer.String(), nil
}

//SSHKeyPair returns pair of public and private key in format of ssh
func SSHKeyPair(size int, pubKeyPath, privateKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return err
	}

	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}
