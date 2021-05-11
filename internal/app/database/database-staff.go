package database

import "go-starter-project/pkg/auth"

type StaffDatabase interface {
	GetUserByName(username string) (auth.UserStaff, error)
	CreateUser(username, password string, roleLevel int) (auth.UserStaff, error)
}
