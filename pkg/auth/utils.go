package auth

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-starter-project/pkg/derror"
	"time"

	"go-starter-project/pkg/config"
)

func getTokenTypeFromJWTToken(token *jwt.Token) string {
	switch token.Claims.(type) {
	case UserClaims:
		return token.Claims.(UserClaims).TokenType
	case *UserClaims:
		return token.Claims.(*UserClaims).TokenType
	}

	return ""
}

func remarshalUserDataAsMap(userData interface{}) (map[string]interface{}, error) {
	j, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func getTokenTimeoutDuration(tokenType string) time.Duration {
	switch tokenType {
	case TokenTypeAccess:
		duration, err := time.ParseDuration(config.Conf.JWT.Timeout)
		if err != nil {
			return time.Duration(10 * time.Minute)
		}
		return duration

	case TokenTypeRefresh:
		duration, err := time.ParseDuration(config.Conf.JWT.RefreshTimeout)
		if err != nil {
			return time.Duration(1 * time.Hour)
		}
		return duration
	}
	return time.Duration(1 * time.Hour)
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	tokenType := getTokenTypeFromJWTToken(token)
	switch tokenType {
	case TokenTypeRefresh:
		return hex.DecodeString(config.Conf.JWT.RefreshKey)
	case TokenTypeAccess:
		return jwt.ParseRSAPublicKeyFromPEM([]byte(config.Conf.JWT.Key))
	}

	return nil, errors.New("invalid token type")
}

func jwtKeyFuncSigning(token *jwt.Token) (interface{}, error) {
	tokenType := getTokenTypeFromJWTToken(token)
	switch tokenType {
	case TokenTypeAccess:
		return jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Conf.JWT.SignKey))
	}

	return jwtKeyFunc(token)
}

func getTokenSigningAlg(tokenType string) jwt.SigningMethod {
	switch tokenType {
	case TokenTypeAccess:
		return jwt.SigningMethodRS256
	}

	return jwt.SigningMethodHS256
}

func remarshalUserData(userType string, userData map[string]interface{}) (interface{}, error) {
	j, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	unmarshal := func(user interface{}) error {
		return json.Unmarshal(j, user)
	}

	var retUser interface{}

	switch userType {
	case UserTypeStaff:
		var user UserStaff
		err = unmarshal(&user)
		retUser = user
	default:
		return nil, derror.ErrInvalidUserType
	}

	if err != nil {
		return nil, err
	}

	return retUser, nil
}

func GetUserFromClaims(claims UserClaims) (interface{}, error) {
	return remarshalUserData(claims.UserType, claims.UserData)
}

func GetClaims(tokenString string) (UserClaims, error) {
	var claims UserClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, jwtKeyFunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return UserClaims{}, derror.ErrUnauthorized
		}
		return UserClaims{}, err
	}

	if !token.Valid {
		return UserClaims{}, derror.ErrUnauthorized
	}

	return claims, nil
}
