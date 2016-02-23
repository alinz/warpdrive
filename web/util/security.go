package util

import (
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/warpdrive/web/constant"
	"golang.org/x/net/context"
)

func UserIsRoot(ctx context.Context) bool {
	return ctx.Value(constant.CtxIsRoot) != nil
}

func LoggedInUserID(ctx context.Context) int64 {
	token := ctx.Value(constant.CtxJWT).(*jwt.Token)
	userIDStr := token.Claims["user_id"].(string)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	return userID
}

func UUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
