package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcConfig struct {
	certPool    *x509.CertPool
	certificate tls.Certificate
}

func (g *GrpcConfig) CreateServer() (*grpc.Server, error) {
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{g.certificate},
		ClientCAs:    g.certPool,
	}

	serverOption := grpc.Creds(credentials.NewTLS(tlsConfig))
	server := grpc.NewServer(serverOption)

	return server, nil
}

func (g *GrpcConfig) CreateClient(name, addr string) (*grpc.ClientConn, error) {
	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   name,
		Certificates: []tls.Certificate{g.certificate},
		RootCAs:      g.certPool,
	})

	dialOption := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial(addr, dialOption)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %s", err)
	}

	return conn, nil
}

func NewGrpcConfig(ca, crt, key string) (*GrpcConfig, error) {
	certificate, err := tls.LoadX509KeyPair(crt, key)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	return &GrpcConfig{certPool, certificate}, nil
}
