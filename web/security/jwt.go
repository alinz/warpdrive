package security

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/web/constant"
)

type jwtAuth struct {
	signKey   []byte
	verifyKey []byte
	signer    jwt.SigningMethod
}

var (
	tokenAuth *jwtAuth
)

//JwtEncode converts map to string
func JwtEncode(claims map[string]interface{}) (string, error) {
	t := jwt.New(tokenAuth.signer)
	t.Claims = claims
	tokenString, err := t.SignedString(tokenAuth.signKey)
	t.Raw = tokenString
	return tokenString, err
}

//JwtDecode convert jwt string token to jwt object
func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if tokenAuth.verifyKey != nil && len(tokenAuth.verifyKey) > 0 {
			return tokenAuth.verifyKey, nil
		}
		return tokenAuth.signKey, nil
	})
}

//TryFindJwt tries to find jwt token in query string, header and cookies
func TryFindJwt(r *http.Request) (*jwt.Token, error) {
	var tokenStr string
	var err error

	// Get token from query params
	tokenStr = r.URL.Query().Get("jwt")

	// Get token from authorization header
	if tokenStr == "" {
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenStr = bearer[7:]
		}
	}

	// Get token from cookie
	if tokenStr == "" {
		cookie, err := r.Cookie("jwt")
		if err == nil {
			tokenStr = cookie.Value
		}
	}

	// Token is required, cya
	if tokenStr == "" {
		return nil, constant.ErrUnauthorized
	}

	// Verify the token
	token, err := JwtDecode(tokenStr)
	if err != nil || !token.Valid || token.Method != tokenAuth.signer {
		return nil, constant.ErrUnauthorized
	}

	return token, nil
}

func SetJwtCookie(token string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "jwt",
		Domain:   warpdrive.Config.JWT.Domain,
		Path:     warpdrive.Config.JWT.Path,
		Secure:   warpdrive.Config.JWT.Secure,
		MaxAge:   warpdrive.Config.JWT.MaxAge,
		HttpOnly: true,
		Value:    token,
	}
	http.SetCookie(w, &cookie)
}

func RemoveJwtCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "jwt",
		Domain:   warpdrive.Config.JWT.Domain,
		Path:     warpdrive.Config.JWT.Path,
		Secure:   warpdrive.Config.JWT.Secure,
		MaxAge:   -1,
		HttpOnly: true,
		Value:    "",
	}
	http.SetCookie(w, &cookie)
}

func UserIDFromJwt(ctx context.Context) (int64, error) {
	token := ctx.Value(constant.CtxJWT).(*jwt.Token)
	content := token.Claims["user_id"].(string)
	id, err := strconv.ParseInt(content, 10, 64)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func SetupWebSecurity() {
	//it has been initialized before, no need to initialize again.
	if tokenAuth != nil {
		return
	}

	tokenAuth = &jwtAuth{
		signKey:   []byte(warpdrive.Config.JWT.SecretKey),
		verifyKey: nil,
		signer:    jwt.GetSigningMethod("HS256"),
	}
}
