package security

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

// Jwt is a base structure for creating signed token
type Jwt struct {
	token string
}

// GetRequestMetadata is being used by higher level code in grpc
func (j *Jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

// RequireTransportSecurity is being used by higher level code in grpc
func (j *Jwt) RequireTransportSecurity() bool {
	return true
}

// Token accepts public key to check the signature of token
func (j *Jwt) Token(publicKey []byte) (*jwt.Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(j.token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// NewJwtToken create a Jwt token based on given token
func NewJwtToken(token string) (*Jwt, error) {
	// data, err := ioutil.ReadFile(tokenFile)
	// if err != nil {
	// 	return nil, err
	// }
	return &Jwt{token: strings.Trim(token, " \r\t\n")}, nil
}

// NewJwtData accepts data and privateKey to sign the jwt
func NewJwtData(data map[string]interface{}, privateKey []byte) (*Jwt, error) {
	// keyData, err := ioutil.ReadFile(privateKeyPath)
	// if err != nil {
	// 	return nil, err
	// }

	claims := jwt.MapClaims(data)
	alg := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(alg, claims)

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	value, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &Jwt{token: value}, nil
}
