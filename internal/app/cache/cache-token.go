package cache

import (
	"time"
)

type TokenCache interface {
	IsTokenBanned(token string) (bool, error)
	BanToken(token string, duration time.Duration) error
}
