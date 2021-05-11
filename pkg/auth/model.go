package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	UserTypeStaff = "staff"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

const (
	PetStoreUser         = "PetStoreUser"
	PetStoreAccessToken  = "PetStoreAccessToken"
	PetStoreRefreshToken = "PetStoreRefreshToken"
)

type AccessRefreshTokenPair struct {
	AccessToken  TokenExpires `json:"access_token"`
	RefreshToken TokenExpires `json:"refresh_token"`
}

type TokenExpires struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserClaims struct {
	UserType  string                 `json:"user_type"`
	UserData  map[string]interface{} `json:"user_data"`
	TokenType string                 `json:"token_type"`
	jwt.StandardClaims
}

type User interface {
	HasShopPermission(permission StaffPermission) bool
	HasGlobalPermission(permission StaffGlobalPermission) bool
}

func HasShopPermission(user interface{}, permission StaffShopPermission) bool {
	if u, ok := user.(User); ok {
		return u.HasShopPermission(permission)
	}
	if u, ok := user.(*User); ok {
		return (*u).HasShopPermission(permission)
	}
	return false
}

func HasGlobalPermission(user interface{}, permission StaffGlobalPermission) bool {
	if u, ok := user.(User); ok {
		return u.HasGlobalPermission(permission)
	}
	if u, ok := user.(*User); ok {
		return (*u).HasGlobalPermission(permission)
	}
	return false
}
