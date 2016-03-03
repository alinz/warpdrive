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

	// fmt.Println(private)
	// fmt.Println(public)

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

func TestRSA2(t *testing.T) {
	privateKey := `-----BEGIN PRIVATE KEY-----
MIICXQIBAAKBgQC2aAHeEafO9hZGm/ZVIDKVrmRM3eU/b5Z6aDWY1amQIynGW2Es
qJlyDt1cAG6FsDT8iD9aj4KwgbzAl2bPNZX2uY11extm+XKxCwe0E0bI1Zm1K8Wy
vGDwbQI4k3qWmmrG5xjYS/G9FvH5TJ3zL+2EZ1NEZFrtLvFM3TzbaxOm4wIDAQAB
AoGABhTXoxjBmIPZ4EbI4rOtHBJxY6KuRvwobzJUPyE4gwa5GNTpG30PiJ74QF3/
UVO7oIPGYPWR7OKWcXFayyPFOSMymeoHN9b9jQh0GfNxm0PQgflhshfgN/r+bDOO
XkxjG9vEI3sOnWtFocaEejUFI4NhebK+UA4zamI/WIIufMECQQDdtPdQ+prls/zd
z0zwbPeUD6cmakxXOOleVpqAo8PQUtFLPd0fa+i43JJquCvAQRqoGnbrKNaunSCu
Zfy/uhDDAkEA0p7V9nGnnXtzJqS7E+fpOZlPhxg10qquOSczdU5YnhwmAvVuUlrk
bpj4Q3bbw6NyPSxBykVTe5ZNBycnmYqvYQJBAMrK/vWZZRn7Gq8hMTUx1vwdnTzs
OkwGCKB8AvLr2O6y8jIqshpNsB930o2/THWcl29wVZogTs6FdyFOtHQDE9UCQAw+
k63KGbZ8EMu0U/PqTZK9qPPvomFm7s3/y2wMa/Z1KHiPkCRViGYtmnFBnbEX9XI8
+m4p7ZqHuF6sFg9FEsECQQCTo1oDn6myvyDXZUs5iPjfg9LGHSDlP4IgBYk1Sh7q
eA0BgWNUJaOmLt5tIMOhWsk+vnAHBy44St4lNNoZRtxQ
-----END PRIVATE KEY-----`

	publicKey := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC2aAHeEafO9hZGm/ZVIDKVrmRM
3eU/b5Z6aDWY1amQIynGW2EsqJlyDt1cAG6FsDT8iD9aj4KwgbzAl2bPNZX2uY11
extm+XKxCwe0E0bI1Zm1K8WyvGDwbQI4k3qWmmrG5xjYS/G9FvH5TJ3zL+2EZ1NE
ZFrtLvFM3TzbaxOm4wIDAQAB
-----END PUBLIC KEY-----`

	message := "Hello This is Warpdrive."

	encrypted, err := crypto.EncryptByPublicRSA([]byte(message), publicKey, "sha256")
	if err != nil {
		t.Error(err)
	}

	decrypted, err := crypto.DecryptByPrivateRSA(encrypted, privateKey, "sha256")
	if err != nil {
		t.Error(err)
	}

	if string(decrypted) != message {
		t.Errorf("got '%s' instead of '%s'", decrypted, message)
	}
}
