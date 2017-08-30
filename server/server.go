package server

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/pressly/warpdrive/proto"
	"github.com/pressly/warpdrive/server/config"
	"github.com/pressly/warpdrive/token"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	config     *config.Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	db         *storm.DB
}

func (s *Server) SetupApp(ctx context.Context, credential *pb.Credential) (*pb.Certificate, error) {
	adminConfig := s.config.Admin
	if adminConfig.Username != credential.Username ||
		bcrypt.CompareHashAndPassword([]byte(adminConfig.Password), []byte(credential.Password)) != nil {
		return nil, fmt.Errorf("username is not correct")
	}

	var certificate *pb.Certificate

	err := s.transaction(func(tx storm.Node) error {
		app := pb.App{Name: credential.AppName}

		err := s.db.Save(&app)
		if err != nil {
			return err
		}

		certificate, err = s.genCertificate(app.Id, 0, true)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func (s *Server) SetupReleaseAdminCertificate(ctx context.Context, release *pb.Release) (*pb.Certificate, error) {
	auth, err := s.getToken(ctx)
	if err != nil {
		return nil, err
	}

	release.Id = 0
	release.AppId = auth.AppId

	var certificate *pb.Certificate

	err = s.transaction(func(tx storm.Node) error {
		err := tx.Select(q.Eq("Name", release.Name), q.Eq("AppId", release.AppId)).First(release)
		if err != nil {
			err = tx.Save(release)
			if err != nil {
				return err
			}
		}

		certificate, err = s.genCertificate(auth.AppId, release.Id, true)
		return err
	})

	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func (s *Server) SetupReleaseUserCertificate(ctx context.Context, release *pb.Release) (*pb.Certificate, error) {
	auth, err := s.getToken(ctx)
	if err != nil {
		return nil, err
	}

	release.Id = 0
	release.AppId = auth.AppId

	var certificate *pb.Certificate

	err = s.transaction(func(tx storm.Node) error {
		err := tx.Select(q.Eq("Name", release.Name), q.Eq("AppId", release.AppId)).First(release)
		if err != nil {
			return err
		}

		certificate, err = s.genCertificate(auth.AppId, release.Id, false)
		return err
	})

	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func (s *Server) Publish(stream pb.Warpdrive_PublishServer) error {
	return nil
}

func (s *Server) Download(release *pb.Release, stream pb.Warpdrive_DownloadServer) error {
	return nil
}

func (s *Server) Start() error {
	cleanup, err := s.setup()
	if err != nil {
		return err
	}

	defer cleanup()

	tlsConfig := s.config.TLS
	creds, err := credentials.NewServerTLSFromFile(tlsConfig.CA, tlsConfig.Private)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	log.Println("server has started")

	// connected grpc server to server implementation
	pb.RegisterWarpdriveServer(grpcServer, s)

	// start listening to the network
	ln, err := net.Listen("tcp", s.config.Server.Addr)
	if err != nil {
		return err
	}

	return grpcServer.Serve(ln)
}

func (s *Server) loadKeys() error {
	tlsConfig := s.config.TLS

	// load public key
	data, err := ioutil.ReadFile(tlsConfig.Public)
	if err != nil {
		return err
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return err
	}

	s.publicKey = publickey

	// load private key
	data, err = ioutil.ReadFile(tlsConfig.Private)
	if err != nil {
		return err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return err
	}

	s.privateKey = privateKey

	return nil
}

func (s *Server) setup() (func(), error) {
	// load keys
	err := s.loadKeys()
	if err != nil {
		return nil, err
	}

	// load db
	s.db, err = storm.Open(s.config.DB.Path)
	if err != nil {
		return nil, err
	}

	return func() {
		s.db.Close()
	}, nil
}

// getToken extract the Token object from context
func (s *Server) getToken(ctx context.Context) (*pb.Token, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	tokenValue, ok := md["authorization"]
	if !ok || len(tokenValue) != 1 {
		return nil, fmt.Errorf("invalid token")
	}

	authorization, err := token.NewAuthorization(s.privateKey, s.publicKey, tokenValue[0])
	if err != nil {
		return nil, err
	}

	if !authorization.Token.Admin {
		return nil, fmt.Errorf("token is not admin")
	}

	return &authorization.Token, nil
}

func (s *Server) genCertificate(appID, releaseID uint64, admin bool) (*pb.Certificate, error) {
	auth, err := token.NewAuthorization(s.privateKey, s.publicKey)
	if err != nil {
		return nil, err
	}

	auth.Admin = admin
	auth.ReleaseId = releaseID
	auth.AppId = appID

	token, err := auth.GetSignedToken()
	if err != nil {
		return nil, err
	}

	cert, err := ioutil.ReadFile(s.config.TLS.CA)
	if err != nil {
		return nil, err
	}

	return &pb.Certificate{
		Addr:  s.config.Server.PublicAddr,
		Token: token,
		Cert:  string(cert),
	}, nil
}

func (s *Server) transaction(fn func(tx storm.Node) error) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// New creates new Server with given config file
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
