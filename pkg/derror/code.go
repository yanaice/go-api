package derror

type DerrorCode string
type DerrorLevel string

const (
	// specific error code
	ErrCodeUserAlreadyExists DerrorCode = "USER_ALREADY_EXISTS"

	// common error code
	ErrCodeUnauthorized DerrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    DerrorCode = "FORBIDDEN"
	ErrCodeInputError   DerrorCode = "INPUT_ERROR"
	ErrCodeNotFound     DerrorCode = "NOT_FOUND"
	ErrCodeServerError  DerrorCode = "SERVER_ERROR"
)

const (
	ErrLevelWarn  DerrorLevel = "warn"
	ErrLevelError DerrorLevel = "error"
	ErrLevelFatal DerrorLevel = "fatal"
)
