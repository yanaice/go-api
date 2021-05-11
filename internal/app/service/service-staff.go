package service

import (
	"github.com/google/uuid"
	"go-starter-project/internal/app/cache"
	"go-starter-project/internal/app/config"
	"go-starter-project/internal/app/database"
	"go-starter-project/internal/app/model"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/derror"
)

type StaffService interface {
	LoginUser(username, password string, ip string) (auth.AccessRefreshTokenPair, error)
	CreateUser(doer interface{}, username, password, ip string) error
	RenewSessionID(user auth.UserStaff) (auth.UserStaff, error)
	CheckSessionID(user auth.UserStaff) error
	RemoveSessionID(user auth.UserStaff) error
}

type staffServiceImpl struct {
	db    database.StaffDatabase
	cache cache.UserSessionCache
}

func StaffServiceInit(db database.StaffDatabase, cache cache.UserSessionCache) StaffService {
	return &staffServiceImpl{db: db, cache: cache}
}

func (s *staffServiceImpl) LoginUser(username, password string, ip string) (auth.AccessRefreshTokenPair, error) {
	user, err := s.db.GetUserByName(username)
	if err != nil {
		_ = derror.E(err).Log()
	}
	if err := user.CheckPassword(password); err != nil {
		return auth.AccessRefreshTokenPair{}, derror.ErrPasswordMismatch
	}

	if newUser, err := s.RenewSessionID(user); err != nil {
		return auth.AccessRefreshTokenPair{}, derror.E(err)
	} else {
		user = newUser
	}

	// Grant permissions
	user.Permissions = model.GetDefaultPermission(user.RoleLevel)

	tokens, err := auth.IssueAccessRefreshTokens(config.ApplicationName, user)
	if err != nil {
		return auth.AccessRefreshTokenPair{}, derror.E(err).SetDebug("login failed: issue token error").
			SetExtraInfo("username", username).Log()
	}

	return tokens, nil
}

func (s *staffServiceImpl) CreateUser(doer interface{}, username, password, ip string) error {
	roleLevel := auth.Customer
	if auth.HasGlobalPermission(doer, auth.GlobalPermCreateStaffUser) {
		roleLevel = auth.RoleShopStaff
	}

	_, err := s.db.GetUserByName(username)
	if err == nil {
		return derror.Ecode(derror.ErrCodeUserAlreadyExists).SetDebug("create user failed: username already exists").
			SetExtraInfo("username", username).Log()
	} else if err != derror.ErrItemNotFound {
		return derror.E(err).SetDebug("create user failed: database error").Log()
	}

	newUser := auth.UserStaff{Username: username, Password: password}
	newUser, err = newUser.HashPassword()
	if err != nil {
		return derror.E(err).SetDebug("create user failed: cannot hash password").Log()
	}
	_, err = s.db.CreateUser(newUser.Username, newUser.Password, roleLevel)
	if err != nil {
		return derror.E(err).SetDebug("create user failed: database error").Log()
	}
	return nil
}

func (s *staffServiceImpl) RenewSessionID(user auth.UserStaff) (auth.UserStaff, error) {
	user.SessionID = uuid.New().String()
	if err := s.cache.SetSessionID(user.ID, user.SessionID); err != nil {
		return auth.UserStaff{}, derror.E(err).SetDebug("renew session id failed: cannot set new session id")
	}
	return user, nil
}

func (s *staffServiceImpl) CheckSessionID(user auth.UserStaff) error {
	realSessionID, err := s.cache.GetSessionID(user.ID)
	if err != nil && err != derror.ErrItemNotFound {
		return derror.E(err).SetDebug("check session id failed: cache error").Log()
	}
	if err == derror.ErrItemNotFound || realSessionID != user.SessionID {
		return derror.Wcode(derror.ErrCodeUnauthorized).SetDebug("check session id failed: not found or session ID mismatch").Log()
	}
	return nil
}

func (s *staffServiceImpl) RemoveSessionID(user auth.UserStaff) error {
	if err := s.cache.UnsetSessionID(user.ID); err != nil {
		return derror.E(err).SetDebug("remove session id failed: cache error (cannot unset session ID)").Log()
	}
	return nil
}
