package server

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/asdine/storm"
	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/pressly/warpdrive/proto"
	"github.com/pressly/warpdrive/server/config"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	config       *config.Config
	jwtPublicKey *rsa.PublicKey
	db           *storm.DB
}

func (s *Server) SetupApp(ctx context.Context, credential *pb.Credential) (*pb.Certificate, error) {
	return nil, nil
}

func (s *Server) SetupReleaseAdminCertificate(ctx context.Context, release *pb.Release) (*pb.Certificate, error) {
	return nil, nil
}

func (s *Server) SetupReleaseUserCertificate(ctx context.Context, relese *pb.Release) (*pb.Certificate, error) {
	return nil, nil
}

func (s *Server) Publish(stream pb.Warpdrive_PublishServer) error {
	return nil
}

func (s *Server) Download(release *pb.Release, stream pb.Warpdrive_DownloadServer) error {
	return nil
}

func (s *Server) Start() error {
	tlsConfig := s.config.TLS
	creds, err := credentials.NewServerTLSFromFile(tlsConfig.CA, tlsConfig.Private)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// load public key for validating JWT
	pubKeyData, err := ioutil.ReadFile(tlsConfig.Public)
	if err != nil {
		return err
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return err
	}

	s.jwtPublicKey = publickey

	// load db
	s.db, err = storm.Open(s.config.DB.Path)
	if err != nil {
		return err
	}

	defer s.db.Close()

	log.Println("server has started")

	// connected grpc server to server implementation
	pb.RegisterWarpdriveServer(grpcServer, s)

	// start listening to the network
	ln, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return err
	}

	return grpcServer.Serve(ln)
}

func New(configPath string) (*Server, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}

	config := config.Config{}
	_, err = toml.DecodeFile(configPath, &config)
	if err != nil {
		return nil, err
	}

	server := Server{
		config: &config,
	}

	return &server, nil
}
