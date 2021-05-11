package service

import (
	"go-starter-project/internal/app/cache"
	"go-starter-project/internal/app/config"
	"go-starter-project/internal/app/model"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/derror"
	"go-starter-project/pkg/dtime"
	"time"
)

type TokenService interface {
	CheckToken(token string) error
	RevokeToken(tokensPair model.TokensPair) error
	RefreshToken(tokensPair model.TokensPair) (auth.AccessRefreshTokenPair, error)
}

type tokenServiceImpl struct {
	cache cache.TokenCache

	ss   StaffService
	time dtime.TimeSource
}

func TokenServiceInit(cache cache.TokenCache, ss StaffService) TokenService {
	return &tokenServiceImpl{
		cache: cache,
		ss:    ss,
		time:  &dtime.RealTime{},
	}
}

func (s *tokenServiceImpl) CheckToken(token string) error {
	if token == "" {
		return derror.Ecode(derror.ErrCodeUnauthorized).
			SetDebug("check token failed: blank token").Log()
	}

	claims, err := auth.GetClaims(token)
	if err != nil {
		if err == derror.ErrUnauthorized {
			// token expired
			return derror.Wcode(derror.ErrCodeUnauthorized).SetDebug("check token failed: token validation failed")
		}
		return derror.E(err).SetCode(derror.ErrCodeInputError).SetDebug("check token failed: get claim failed").Log()
	}

	user, err := auth.GetUserFromClaims(claims)
	if err != nil {
		return derror.E(err).SetDebug("check token failed: get user failed").Log()
	}

	if staff, ok := user.(auth.UserStaff); ok {
		if !config.Conf.AllowConcurrentSessions {
			if err := s.ss.CheckSessionID(staff); err != nil {
				return derror.Wrap(err.(derror.Derror)).SetDebug("check token failed: check session failed").Log()
			}
		}
	}

	banned, err := s.cache.IsTokenBanned(token)
	if err != nil {
		return derror.E(err).SetDebug("check token failed: database error").Log()
	}

	if banned {
		return derror.Wcode(derror.ErrCodeUnauthorized).SetDebug("check token failed: token banned").Log()
	}

	return nil
}

func (s *tokenServiceImpl) RefreshToken(tokensPair model.TokensPair) (auth.AccessRefreshTokenPair, error) {

	token := tokensPair.RefreshToken
	if token == "" {
		return auth.AccessRefreshTokenPair{}, derror.Ecode(derror.ErrCodeUnauthorized).
			SetDebug("refresh token failed: blank refresh token").Log()
	}

	claims, err := auth.GetClaims(token)
	if err != nil {
		if err == derror.ErrUnauthorized {
			return auth.AccessRefreshTokenPair{}, derror.Wcode(derror.ErrCodeUnauthorized).
				SetDebug("refresh token failed: token validation failed").SetExtraInfo("token", token).Log()
		}
		return auth.AccessRefreshTokenPair{}, derror.E(err).SetCode(derror.ErrCodeInputError).SetDebug("refresh token failed: get claim failed").Log()
	}

	user, err := auth.GetUserFromClaims(claims)
	if err != nil {
		return auth.AccessRefreshTokenPair{}, derror.E(err).SetDebug("refresh token failed: get user failed").Log()
	}

	if claims.TokenType != auth.TokenTypeRefresh {
		return auth.AccessRefreshTokenPair{}, derror.Ecode(derror.ErrCodeUnauthorized).
			SetDebug("refresh token failed: wrong token type").Log()
	}

	banned, err := s.cache.IsTokenBanned(token)
	if err != nil {
		return auth.AccessRefreshTokenPair{}, derror.E(err).SetDebug("refresh token failed: token ban check error").Log()
	}

	if banned {
		return auth.AccessRefreshTokenPair{}, derror.Wcode(derror.ErrCodeUnauthorized).SetDebug("refresh token failed: token banned").Log()
	}

	if staff, ok := user.(auth.UserStaff); ok {
		if err := s.ss.CheckSessionID(staff); err != nil {
			return auth.AccessRefreshTokenPair{}, derror.Wrap(err.(derror.Derror)).SetDebug("refresh token failed: check session failed").Log()
		}
		if newStaff, err := s.ss.RenewSessionID(staff); err != nil {
			return auth.AccessRefreshTokenPair{}, derror.Wrap(err.(derror.Derror)).SetDebug("refresh token failed: renew session failed").Log()
		} else {
			user = newStaff
		}
	}

	tokens, err := auth.IssueAccessRefreshTokens(config.ApplicationName, user)
	if err != nil {
		return auth.AccessRefreshTokenPair{}, derror.E(err).SetDebug("refresh token failed: cannot issue token").Log()
	}
	err = s.revokeToken(tokensPair.RefreshToken)
	if err != nil {
		return auth.AccessRefreshTokenPair{}, derror.E(err).SetDebug("refresh token failed: cannot revoke token").Log()
	}
	if tokensPair.AccessToken != "" {
		err = s.revokeToken(tokensPair.AccessToken)
		if err != nil {
			_ = derror.W(err).SetDebug("refresh token warning: cannot revoke access token").Log()
		}
	}

	return tokens, nil
}

func (s *tokenServiceImpl) RevokeToken(tokensPair model.TokensPair) error {
	if tokensPair.RefreshToken == "" {
		return derror.Ecode(derror.ErrCodeUnauthorized).
			SetDebug("revoke token failed: blank refresh token").Log()
	}

	claims, err := auth.GetClaims(tokensPair.RefreshToken)
	if err == derror.ErrUnauthorized {
		return derror.W(err).SetCode(derror.ErrCodeUnauthorized).
			SetDebug("revoke token failed: token expired").Log()
	} else if err != nil {
		return derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("revoke token failed: cannot get claims").Log()
	}

	if claims.TokenType != auth.TokenTypeRefresh {
		return derror.Ecode(derror.ErrCodeUnauthorized).
			SetDebug("revoke token failed: wrong token type").Log()
	}

	user, err := auth.GetUserFromClaims(claims)
	if err != nil {
		return derror.E(err).SetDebug("check token failed: get user failed").Log()
	}

	if staff, ok := user.(auth.UserStaff); ok {
		if err := s.ss.CheckSessionID(staff); err != nil {
			return nil
		}
		if err := s.ss.RemoveSessionID(staff); err != nil {
			_ = derror.W(err).SetDebug("revoke token warning: cannot remove session ID").Log()
		}
	}

	if err := s.revokeToken(tokensPair.RefreshToken); err != nil {
		return derror.E(err).SetDebug("revoke token failed: cannot revoke refresh token").Log()
	}

	if tokensPair.AccessToken != "" {
		if err := s.revokeToken(tokensPair.AccessToken); err != nil {
			_ = derror.W(err).SetDebug("revoke token warning: cannot revoke access token").Log()
		}
	}

	return nil
}

func (s *tokenServiceImpl) revokeToken(token string) error {
	claims, err := auth.GetClaims(token)
	if err == derror.ErrUnauthorized {
		return nil
	} else if err != nil {
		return err
	}

	expires := time.Unix(claims.ExpiresAt, 0)

	banDuration := expires.Sub(s.time.Now())
	banDuration += time.Duration(config.Conf.BanTokenGuardThreshold) * time.Second

	if banDuration <= 0 {
		return nil
	}

	err = s.cache.BanToken(token, banDuration)
	if err != nil {
		return err
	}

	return nil
}
