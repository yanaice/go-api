package cache

type UserSessionCache interface {
	GetSessionID(userID string) (string, error)
	SetSessionID(userID, sessionID string) error
	UnsetSessionID(userID string) error
}
