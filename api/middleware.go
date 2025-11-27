package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yilinyo/project_bank/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationFormat     = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is empty")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//用于分割 Bearer xxxxxxxx 的token
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != authorizationFormat {
			err := errors.New("authorization header is not valid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken, token.TokenTypeAccessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//设置payload ctx在上下文
		ctx.Set(authorizationPayloadKey, payload)
		//转发到下个handler
		ctx.Next()

	}

}
