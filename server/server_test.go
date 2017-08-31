package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pressly/warpdrive/server/config"
)

var testTLS = struct {
	CA      []byte
	Private []byte
	Public  []byte
}{
	CA: []byte(`-----BEGIN CERTIFICATE-----
MIIBnzCCAQgCCQDWXptzyyuH6DANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDEwl3
YXJwZHJpdmUwHhcNMTcwODMwMTk0ODQ4WhcNMjcwODI4MTk0ODQ4WjAUMRIwEAYD
VQQDEwl3YXJwZHJpdmUwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMME5Ppp
jc1xMLEgb/8K+VCk/lHCYF5WQJPbsmypPI6ZXx8Iqxb212dc//JAmIXbF9fF25Mx
aq+98USPZfgxA/ihH3ZC9jaQuPaSxATZ6XykwaEkj/2dYV80EmVL5BCoJ0KBUQT2
F8/imAyZpV3JmCSVcNEQp4snvjFNeWBQt0mNAgMBAAEwDQYJKoZIhvcNAQELBQAD
gYEAKf/1qXn6vzrJlqIX5pCWFwzyH1VCucTn3fMpgtyNu+6XoVwDUoHdV/x13vnS
aPLMfpMKR4RbRyi3n1+WTqndVw7/7gLkqj8A0MA8UbsTXD8AgvxiVS5/8P2TPunD
wyWYMOfea+1o6nYurX2S4BzCsK45iqISpGrSFvM4B42Ziy0=
-----END CERTIFICATE-----`),
	Private: []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDDBOT6aY3NcTCxIG//CvlQpP5RwmBeVkCT27JsqTyOmV8fCKsW
9tdnXP/yQJiF2xfXxduTMWqvvfFEj2X4MQP4oR92QvY2kLj2ksQE2el8pMGhJI/9
nWFfNBJlS+QQqCdCgVEE9hfP4pgMmaVdyZgklXDREKeLJ74xTXlgULdJjQIDAQAB
AoGAQnHTdkIqbznGhkLwBax+f2yHveGFJf8rJ3VuGDmdCVTWJOO2Ly/Q+kWkaqx5
ivm36OtfwYnPuKr1ng9hhatll2JbVAUA/2cGDmA9FOWjjxLKzR4BigqvRvu/I7to
wC+hlNRImf/OTaU8KB2UeQzN7MAUhqi+IwAEf3KmJxp5uH0CQQD5PgZTx2pr8bs0
L5bOvHIjv//784dAJKfVsnm/mTnV+S206mhX4bL0rLcTDXOwNo9dS76D3Mwf4Edf
9fkBe/6vAkEAyE6ClXssCTMrYkz3ruDQay25gDfv/i2NUKlHUX6nWhZ1T8W1wYIk
o5siDSOMwpeThb9kpall0gt7QVxbXgjqgwJBAPBHRzpFKOdfZyXsKuqq6S5lzpZK
M702mUZ+hLidMxCA4/thb64pO6h9SRDpCvp53sQGXWgp1+9y+9wa+S7hJqkCQQDA
r+WmbmqaHwMo+Ol67QERWVcNJMJVPPSoF29n0fKjEt+e8Y46rDsat202PnB18OIU
01y6kA5G1Iyo/3NVLjaJAkAHG5ZpGiws8wA3Ub9z0tq/ZjXuul2OTpqCfvuaQWfb
1zWDOshMvt/kasdjRPIJP2ij3gnQhSDChVbBhokKRdFg
-----END RSA PRIVATE KEY-----`),
	Public: []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDBOT6aY3NcTCxIG//CvlQpP5R
wmBeVkCT27JsqTyOmV8fCKsW9tdnXP/yQJiF2xfXxduTMWqvvfFEj2X4MQP4oR92
QvY2kLj2ksQE2el8pMGhJI/9nWFfNBJlS+QQqCdCgVEE9hfP4pgMmaVdyZgklXDR
EKeLJ74xTXlgULdJjQIDAQAB
-----END PUBLIC KEY-----`),
}

func createTempDir() (string, func() error, error) {
	path, err := ioutil.TempDir("", "test")
	if err != nil {
		return "", nil, err
	}

	return path, func() error {
		return os.Remove(path)
	}, nil
}

func createTempTLSfiles(config *config.Config) error {
	err := ioutil.WriteFile(config.TLS.CA, testTLS.CA, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.TLS.Private, testTLS.Private, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.TLS.Public, testTLS.Public, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func createTestServer(config *config.Config) (*Server, func() error, error) {
	// create a temporary folder
	path, cleanTempDir, err := createTempDir()

	config.TLS.CA = filepath.Join(path, config.TLS.CA)
	config.TLS.Public = filepath.Join(path, config.TLS.Public)
	config.TLS.Private = filepath.Join(path, config.TLS.Private)

	err = createTempTLSfiles(config)
	if err != nil {
		cleanTempDir()
		return nil, nil, err
	}

	config.DB.Path = filepath.Join(path, config.DB.Path)
	config.Server.BundlesDir = filepath.Join(path, config.Server.BundlesDir)

	return &Server{
			Conf: config,
		}, func() error {
			cleanTempDir()
			return nil
		}, nil
}

func startTestServer(server *Server) (func(), error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	server.Conf.Server.Addr = ln.Addr().String()
	start, err := server.SetupServer()
	if err != nil {
		ln.Close()
		return nil, err
	}

	go func() {
		fmt.Println("server start at", server.Conf.Server.Addr)
		err = start(ln)
		if err != nil {
			fmt.Println("server closed")
			ln.Close()
		}
	}()

	return func() {
		fmt.Println("server closed")
		ln.Close()
	}, nil
}

func TestBasicServer(t *testing.T) {
	conf := config.Config{
		Server: struct {
			Addr       string `toml:"addr"`
			PublicAddr string `toml:"public_addr"`
			BundlesDir string `toml:"bundles_dir"`
		}{
			BundlesDir: "/bundles",
		},

		DB: struct {
			Path string `toml:"path"`
		}{
			Path: "/db",
		},

		TLS: struct {
			CA      string `toml:"ca"`
			Private string `toml:"private"`
			Public  string `toml:"public"`
		}{
			CA:      "ca",
			Private: "private",
			Public:  "public",
		},

		Admin: struct {
			Username string `toml:"username"`
			Password string `toml:"password"`
		}{
			Username: "",
			Password: "",
		},
	}

	server, cleanup, err := createTestServer(&conf)
	if err != nil {
		t.Fatal(err)
	}

	defer cleanup()

	cleanupListener, err := startTestServer(server)
	if err != nil {
		t.Fatal(err)
	}

	defer cleanupListener()

	<-time.After(time.Second * 5)
}
