package server

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/pressly/warpdrive/proto"
	"github.com/pressly/warpdrive/server/config"
	"github.com/pressly/warpdrive/token"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var tlsTest = struct {
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

func createConfigTest() config.Config {
	return config.Config{
		Server: struct {
			Addr       string `toml:"addr"`
			PublicAddr string `toml:"public_addr"`
			BundlesDir string `toml:"bundles_dir"`
		}{
			PublicAddr: "warpdrive.example.com",
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
			Username: "admin",
			Password: "$2a$10$NRzYde3E6xGwN1eATKGvBeY1DXhWghAjiBFRvxaJLy9AQ0JmTXG2q",
		},
	}
}

func createTempDirTest() (string, func() error, error) {
	path, err := ioutil.TempDir("", "test")
	if err != nil {
		return "", nil, err
	}

	return path, func() error {
		if err := os.RemoveAll(path); err != nil {
			fmt.Println("ERROR", err)
		}

		return err
	}, nil
}

func createTempTLSfilesTest(config *config.Config) error {
	err := ioutil.WriteFile(config.TLS.CA, tlsTest.CA, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.TLS.Private, tlsTest.Private, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.TLS.Public, tlsTest.Public, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func createServerTest(config *config.Config) (*Server, func() error, error) {
	// create a temporary folder
	path, cleanTempDir, err := createTempDirTest()

	config.TLS.CA = filepath.Join(path, config.TLS.CA)
	config.TLS.Public = filepath.Join(path, config.TLS.Public)
	config.TLS.Private = filepath.Join(path, config.TLS.Private)

	err = createTempTLSfilesTest(config)
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

func startServerTest(server *Server) (func(), error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	server.Conf.Server.Addr = ln.Addr().String()
	start, cleanup, err := server.SetupServer()
	if err != nil {
		ln.Close()
		return nil, err
	}

	go func() {
		defer cleanup()
		// fmt.Println("server start at", server.Conf.Server.Addr)
		err = start(ln)
		if err != nil {
			ln.Close()
		}
	}()

	return func() {
		cleanup()
		ln.Close()
	}, nil
}

func createClientTest(cert []byte, addr string, tokValue string) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}

	creds := credentials.NewClientTLSFromCert(cp, "warpdrive")

	dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))

	if tokValue != "" {
		authorization, err := token.NewAuthorizationWithToken(tokValue)
		if err != nil {
			return nil, err
		}

		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(authorization))
	}

	return grpc.Dial(addr, dialOptions...)
}

func TestBasicServer(t *testing.T) {
	conf := createConfigTest()

	server, cleanup, err := createServerTest(&conf)
	assert.Nil(t, err)

	defer cleanup()

	cleanupListener, err := startServerTest(server)
	assert.Nil(t, err)

	defer cleanupListener()
}

func TestGenCertificate(t *testing.T) {
	conf := createConfigTest()

	server, serverCleanup, err := createServerTest(&conf)
	assert.Nil(t, err)

	defer serverCleanup()

	setupCleanup, err := server.setup()
	assert.Nil(t, err)

	defer setupCleanup()

	certificate, err := server.genCertificate(1, 1, true)
	assert.Nil(t, err)

	expected := &pb.Certificate{
		Addr:  "warpdrive.example.com",
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXBwSWQiOjEsInJlbGVhc2VJZCI6MX0.Hk_q2NQB5IZNSrw_xXrYSr80N2awhW42RETOpk4nlEH06PQJSWSlgz-hgdP1eZO9_FTvaIBlWE4R6cK6P5vDX_9RBkVR6CzAja4Q36m53XAVAmtSWFJFlWSjkIXq8dAal-guF2UIJeiNjdkom7VtvzpzqmBDJLdq9-RBWrjZlxI",
		Cert:  "-----BEGIN CERTIFICATE-----\nMIIBnzCCAQgCCQDWXptzyyuH6DANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDEwl3\nYXJwZHJpdmUwHhcNMTcwODMwMTk0ODQ4WhcNMjcwODI4MTk0ODQ4WjAUMRIwEAYD\nVQQDEwl3YXJwZHJpdmUwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMME5Ppp\njc1xMLEgb/8K+VCk/lHCYF5WQJPbsmypPI6ZXx8Iqxb212dc//JAmIXbF9fF25Mx\naq+98USPZfgxA/ihH3ZC9jaQuPaSxATZ6XykwaEkj/2dYV80EmVL5BCoJ0KBUQT2\nF8/imAyZpV3JmCSVcNEQp4snvjFNeWBQt0mNAgMBAAEwDQYJKoZIhvcNAQELBQAD\ngYEAKf/1qXn6vzrJlqIX5pCWFwzyH1VCucTn3fMpgtyNu+6XoVwDUoHdV/x13vnS\naPLMfpMKR4RbRyi3n1+WTqndVw7/7gLkqj8A0MA8UbsTXD8AgvxiVS5/8P2TPunD\nwyWYMOfea+1o6nYurX2S4BzCsK45iqISpGrSFvM4B42Ziy0=\n-----END CERTIFICATE-----",
	}

	assert.Equal(t, expected, certificate)
}

func TestSetupApp(t *testing.T) {
	conf := createConfigTest()

	server, cleanup, err := createServerTest(&conf)
	assert.Nil(t, err)

	defer cleanup()

	cleanupListener, err := startServerTest(server)
	assert.Nil(t, err)

	defer cleanupListener()

	cert, err := ioutil.ReadFile(conf.TLS.CA)
	assert.Nil(t, err)

	conn, err := createClientTest(cert, conf.Server.Addr, "")
	assert.Nil(t, err)

	defer conn.Close()

	client := pb.NewWarpdriveClient(conn)

	_, err = client.SetupApp(context.Background(), &pb.Credential{Username: "admin", Password: "admin", AppName: "My Awesome App"})
	assert.Nil(t, err)
}

func TestRelease(t *testing.T) {
	conf := createConfigTest()

	server, cleanup, err := createServerTest(&conf)
	assert.Nil(t, err)

	defer cleanup()

	cleanupListener, err := startServerTest(server)
	assert.Nil(t, err)

	defer cleanupListener()

	var adminCert *pb.Certificate
	var clientCert *pb.Certificate

	// we need to call the create app first with admin credential
	// this will create admin certificate and token.
	// however, release id has not been set yet, so we can use this
	// certificare to create a release object
	{
		cert, err := ioutil.ReadFile(conf.TLS.CA)
		assert.Nil(t, err)

		conn, err := createClientTest(cert, conf.Server.Addr, "")
		assert.Nil(t, err)

		defer conn.Close()

		client := pb.NewWarpdriveClient(conn)

		adminCert, err = client.SetupApp(context.Background(), &pb.Credential{Username: "admin", Password: "admin", AppName: "My Awesome App"})
		assert.Nil(t, err)
	}

	// we are going to create a release for current app
	// this will generate a new certificate which we need to store to `warpdrive.admin.json`
	{
		conn, err := createClientTest([]byte(adminCert.Cert), conf.Server.Addr, adminCert.Token)
		assert.Nil(t, err)

		defer conn.Close()

		client := pb.NewWarpdriveClient(conn)

		adminCert, err = client.SetupReleaseAdminCertificate(context.Background(), &pb.Release{Name: "prod"})
		assert.Nil(t, err)
	}

	// next we need another certificate for bundling with mobile app
	// this call will generate mobile app certificate
	{
		conn, err := createClientTest([]byte(adminCert.Cert), conf.Server.Addr, adminCert.Token)
		assert.Nil(t, err)

		defer conn.Close()

		client := pb.NewWarpdriveClient(conn)

		clientCert, err = client.SetupReleaseAdminCertificate(context.Background(), &pb.Release{Name: "prod"})
		assert.Nil(t, err)
	}

	assert.NotNil(t, clientCert)
}
