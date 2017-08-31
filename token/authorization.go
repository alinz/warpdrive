package token

import (
	"crypto/rsa"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/pressly/warpdrive/proto"
	"golang.org/x/net/context"
)

type Authorization struct {
	pb.Token
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	token      string
}

// GetRequestMetadata is being used by higher level code in grpc
func (a *Authorization) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	var err error

	token := a.token
	if token == "" {
		token, err = a.GetSignedToken()
		if err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"authorization": token,
	}, nil
}

// RequireTransportSecurity is being used by higher level code in grpc
func (a *Authorization) RequireTransportSecurity() bool {
	return true
}

func (a *Authorization) Valid() error {
	return nil
}

func (a *Authorization) GetSignedToken() (string, error) {
	if a.privateKey == nil {
		return "", fmt.Errorf("private key not given")
	}

	alg := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(alg, a)

	return token.SignedString(a.privateKey)
}

func NewAuthorizationWithToken(token string) (*Authorization, error) {
	return &Authorization{
		token: token,
	}, nil
}

func NewAuthorization(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, value ...string) (*Authorization, error) {
	var authorization *Authorization

	if len(value) > 0 {
		var ok bool

		authorization = &Authorization{}
		token, err := jwt.ParseWithClaims(strings.Trim(value[0], " \r\t\n"), authorization, func(t *jwt.Token) (interface{}, error) {
			if _, ok = t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil {
			return nil, err
		}

		if !ok && !token.Valid {
			return nil, fmt.Errorf("not valid")
		}
	}

	if authorization == nil {
		authorization = &Authorization{}
	}

	authorization.privateKey = privateKey
	authorization.publicKey = publicKey

	return authorization, nil
}
