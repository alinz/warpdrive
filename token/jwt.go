package token

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

type Jwt struct {
	token string
}

func (j *Jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j *Jwt) RequireTransportSecurity() bool {
	return true
}

func (j *Jwt) Token(publicKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(j.token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func NewJwtWithToken(token string) (*Jwt, error) {
	return &Jwt{token: strings.Trim(string(token), " \r\t\n")}, nil
}

func NewJwtWithClaims(data map[string]interface{}, privateKey string) (*Jwt, error) {
	keyData, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims(data)
	alg := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(alg, claims)

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	value, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &Jwt{token: value}, nil
}
