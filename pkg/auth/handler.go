package auth

import (
	"github.com/gin-gonic/gin"
	api "go-starter-project/pkg/api/auth"
	"go-starter-project/pkg/derror"
	"go-starter-project/pkg/handler"
	"strings"
)

func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}
	authHeaderSplits := strings.SplitN(authHeader, " ", 2)
	if (len(authHeaderSplits) != 2) || (authHeaderSplits[0] != "Bearer") {
		return "", derror.Ecode(derror.ErrCodeUnauthorized).
			SetDebug("authentication middleware error: invalid header").
			SetExtraInfo("header", authHeader).Log()
	}
	return authHeaderSplits[1], nil
}

func GetMiddlewareAuthToken(a api.AuthTokenAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := GetTokenFromHeader(c)
		if err != nil {
			handler.HandlerReturnError(c, err)
			c.Abort()
			return
		}

		if tokenString == "" {
			return
		}

		claims, err := GetClaims(tokenString)
		if err != nil {
			if err == derror.ErrUnauthorized {
				handler.HandlerReturnError(c, derror.Ecode(derror.ErrCodeUnauthorized))
				c.Abort()
				return
			}
			handler.HandlerReturnError(c, derror.E(err).
				SetDebug("authentication middleware error: get claims error").
				SetExtraInfo("token", tokenString).Log())
			c.Abort()
			return
		}

		user, err := GetUserFromClaims(claims)
		if err != nil {
			handler.HandlerReturnError(c, derror.E(err).
				SetDebug("authentication middleware error: get user from claims error").
				SetExtraInfo("token", tokenString).Log())
			c.Abort()
			return
		}

		tokenType := claims.TokenType

		c.Set(PetStoreUser, user)
		switch tokenType {
		case TokenTypeAccess:
			c.Set(PetStoreAccessToken, true)
		case TokenTypeRefresh:
			c.Set(PetStoreRefreshToken, true)
		}

		if tokenType == TokenTypeAccess {
			// TODO: check session id from token by calling checking token in auth service
			if err = a.CheckToken(tokenString); err != nil {
				handler.HandlerReturnError(c, err)
				c.Abort()
				return
			}
		}
	}
}

func GetMiddlewareAuthCheckAccessToken() gin.HandlerFunc { // TODO: deprecated
	return func(c *gin.Context) {
		_, ok := c.Get(PetStoreUser)
		if ok && !c.GetBool(PetStoreAccessToken) {
			handler.HandlerReturnError(c, derror.Ecode(derror.ErrCodeUnauthorized).
				SetDebug("authentication middleware error: not access token").Log())
			c.Abort()
		}
	}
}

func GetMiddlewareAuthLoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get(PetStoreUser)
		if !ok || !c.GetBool(PetStoreAccessToken) {
			handler.HandlerReturnError(c, derror.Ecode(derror.ErrCodeUnauthorized).
				SetDebug("authentication middleware error: not logged in or wrong token type").Log())
			c.Abort()
		}
	}
}
