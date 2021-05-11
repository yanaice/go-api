package auth

import (
	"github.com/dgrijalva/jwt-go"
	"go-starter-project/pkg/derror"
)

var TimeFunc = &jwt.TimeFunc

func issueToken(tokenType string, issuer string, userDataI interface{}) (TokenExpires, error) {
	userType := ""

	currentTime := (*TimeFunc)()
	expireAt := currentTime.Add(getTokenTimeoutDuration(tokenType))

	switch userDataI.(type) {
	case UserStaff:
		userType = UserTypeStaff
		userDataI = userDataI.(UserStaff).StripPassword()
	default:
		return TokenExpires{}, derror.ErrInvalidUserType
	}

	userData, err := remarshalUserDataAsMap(userDataI)
	if err != nil {
		return TokenExpires{}, nil
	}

	claims := UserClaims{}
	claims.UserType = userType
	claims.UserData = userData
	claims.TokenType = tokenType
	claims.Issuer = issuer
	claims.IssuedAt = currentTime.UTC().Unix() - 600
	claims.NotBefore = currentTime.UTC().Unix() - 600
	claims.ExpiresAt = expireAt.UTC().Unix()

	token := jwt.NewWithClaims(getTokenSigningAlg(tokenType), claims)
	hexKey, err := jwtKeyFuncSigning(token)
	if err != nil {
		return TokenExpires{}, err
	}
	jwtStr, err := token.SignedString(hexKey)
	if err != nil {
		return TokenExpires{}, err
	}
	return TokenExpires{Token: jwtStr, ExpiresAt: expireAt}, nil
}

func IssueAccessRefreshTokens(issuer string, userDataI interface{}) (AccessRefreshTokenPair, error) {
	access, err := issueToken(TokenTypeAccess, issuer, userDataI)
	if err != nil {
		return AccessRefreshTokenPair{}, err
	}
	refresh, err := issueToken(TokenTypeRefresh, issuer, userDataI)
	if err != nil {
		return AccessRefreshTokenPair{}, err
	}
	return AccessRefreshTokenPair{AccessToken: access, RefreshToken: refresh}, nil
}
