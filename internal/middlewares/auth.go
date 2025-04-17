package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	adminServices "server/internal/app/admin/services"
	kanboardServices "server/internal/app/kanboard/services"
	"server/internal/common"
	"server/internal/constant"
	"server/internal/global"
	"server/pkg/crypto"

	"github.com/gin-gonic/gin"
)

var (
	AUTH_KEY             = "Authorization"
	AUTH_QUERY           = "token"
	TOKEN_KEY            = "Token"
	RENEW_TOKEN_DURATION = 15 * time.Minute
)

func tokenError(ctx *gin.Context, msg string) {
	common.ResponseWithOmitempty(ctx, common.BaseRspWithOmitempty{
		Status: http.StatusUnauthorized,
		Code:   constant.FAIL,
		Msg:    msg,
	})
}

func Auth(namespace string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AUTH_KEY)
		authQuery := ctx.Query(AUTH_QUERY)

		if authHeader == "" && authQuery == "" {
			tokenError(ctx, "Authorization header is required")
			return
		}

		var token string
		if authHeader != "" {
			parts := strings.Split(authHeader, ":")
			if len(parts) != 2 || parts[0] != TOKEN_KEY {
				tokenError(ctx, "Invalid token format")
				return
			}
			token = parts[1]
		} else {
			token = authQuery
		}

		jwtClaims, err := crypto.ParseToken(token)
		if err != nil || jwtClaims.ID == 0 {
			tokenError(ctx, err.Error())
			return
		}

		userId := strconv.Itoa(int(jwtClaims.ID))
		result := global.Redis.Get(namespace, userId)
		if result != token {
			tokenError(ctx, "Invalid token")
			return
		}

		ttl := global.Redis.GetTTL(namespace, userId)
		if ttl <= 0 {
			tokenError(ctx, "Token expired")
			return
		}

		if ttl.Minutes() <= RENEW_TOKEN_DURATION.Minutes() {
			var token string
			var err error
			if namespace == constant.KANBOARD_TOKEN {
				token, err = kanboardServices.GenerateToken(jwtClaims.ID, jwtClaims.Username)
			} else {
				token, err = adminServices.GenerateToken(jwtClaims.ID, jwtClaims.Username)
			}
			if err != nil {
				tokenError(ctx, err.Error())
				return
			}
			ctx.Header("Token", token)
		}

		ctx.Next()
	}
}
