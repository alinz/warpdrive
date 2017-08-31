package section

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	pb "github.com/pressly/warpdrive/proto"
	"github.com/pressly/warpdrive/token"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func stringHasValue(values ...*string) bool {
	for _, value := range values {
		if value == nil || len(strings.Trim(*value, " \t\r\n")) == 0 {
			return false
		}
	}

	return true
}

func loadCert() ([]byte, error) {
	data, err := ioutil.ReadFile("./warpdrive.crt")
	if err != nil {
		return nil, fmt.Errorf("warpdrive.crt not found")
	}

	return data, nil
}

func loadCertificate() (*pb.Certificate, error) {
	file, err := os.Open("./warpdrive.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var certificate pb.Certificate
	err = json.NewDecoder(file).Decode(&certificate)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func grpcConnection(addr string, tokValue string) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption

	cert, err := loadCert()
	if err != nil {
		return nil, err
	}

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

func saveAdminCertificate(certificate *pb.Certificate) error {
	file, err := os.OpenFile("./warpdrive.admin.json", os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(certificate)
	if err != nil {
		return err
	}

	return nil
}

func saveUserCertificate(certificate *pb.Certificate) error {
	file, err := os.OpenFile("./warpdrive.json", os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(certificate)
	if err != nil {
		return err
	}

	return nil
}

var root *cobra.Command

func init() {
	root = &cobra.Command{
		Use:   "warp",
		Short: "In-App upgrade service for React-Native! Supporting iOS and Android apps",
		Long: `A Fast and Flexible upgrade service for React-Native apps!
Created by Ali Najafizadeh (alinz) at Pressly Inc.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("please run 'warp -h' for usage")
		},
	}
}
