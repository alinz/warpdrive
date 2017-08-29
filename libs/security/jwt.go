package security

import (
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

func NewJwtToken(tokenFile string) (*Jwt, error) {
	data, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	return &Jwt{token: strings.Trim(string(data), " \r\t\n")}, nil
}

func NewJwtData(data map[string]interface{}, privateKeyPath string) (*Jwt, error) {
	keyData, err := ioutil.ReadFile(privateKeyPath)
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
