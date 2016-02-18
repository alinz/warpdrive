package util

import (
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
