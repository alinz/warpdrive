package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

func RSAKeyPair(size int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return "", "", err
	}

	var privateBuffer bytes.Buffer

	pemkey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	pem.Encode(&privateBuffer, pemkey)

	var publicBuffer bytes.Buffer

	bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return "", "", err
	}

	pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bytes,
	}

	pem.Encode(&publicBuffer, pemkey)
	return privateBuffer.String(), publicBuffer.String(), nil
}

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
		Type:  "RSA PRIVATE KEY",
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
