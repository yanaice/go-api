package derror

import (
	"errors"
)

var ErrEndpointNotFound = errors.New("endpoint not found")
var ErrItemNotFound = errors.New("item not found")
var ErrInputValidationFailed = errors.New("input validation failed")
var ErrPasswordMismatch = errors.New("password mismatch")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrInvalidUserType = errors.New("invalid user type")
var ErrInvalidColumn = errors.New("invalid column")