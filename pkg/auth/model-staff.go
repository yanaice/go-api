package auth

import (
	"crypto/sha256"
	"go-starter-project/pkg/derror"
	"golang.org/x/crypto/bcrypt"
)

const RoleSuperAdmin = 256
const RoleShopStaff = 4
const Customer = 1

type UserStaff struct {
	ID          string              `bson:"-" json:"id"`
	Username    string              `bson:"username" json:"username"`
	Password    string              `bson:"password" json:"password,omitempty"`
	SessionID   string              `bson:"-" json:"session_id,omitempty"`
	Permissions StaffPermissionsMap `bson:"-" json:"permissions,omitempty"`
	RoleLevel   int32               `bson:"role_level" json:"role_level"`
}

func (u UserStaff) CheckPassword(password string) error {
	h := sha256.New()
	_, err := h.Write([]byte(password))
	if err != nil {
		return err
	}

	shaHashed := h.Sum(nil)
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), shaHashed)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return derror.ErrPasswordMismatch
		}
		return err
	}

	return nil
}

func (u UserStaff) HashPassword() (UserStaff, error) {
	h := sha256.New()
	_, err := h.Write([]byte(u.Password))
	if err != nil {
		return UserStaff{}, err
	}

	shaHashed := h.Sum(nil)
	hashed, err := bcrypt.GenerateFromPassword(shaHashed, 14)
	if err != nil {
		return UserStaff{}, err
	}
	u.Password = string(hashed)
	return u, nil
}

func (u UserStaff) StripPassword() UserStaff {
	u.Password = ""
	return u
}

func (u UserStaff) HasShopPermission(permission StaffPermission) bool {
	// TODO: Other conditions
	return u.Permissions[permission.GetPermissionKey()]
}

func (u UserStaff) HasGlobalPermission(permission StaffGlobalPermission) bool {
	return u.Permissions[permission.GetPermissionKey()]
}
