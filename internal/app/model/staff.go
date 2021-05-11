package model

type StaffLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
