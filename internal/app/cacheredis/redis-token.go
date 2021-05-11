package cacheredis

import (
	"crypto/sha256"
	"encoding/hex"
	"go-starter-project/internal/app/cache"
	"time"
)

type tokenBlacklistCacheImpl struct{}

func GetTokenBlacklistCache() cache.TokenCache {
	return &tokenBlacklistCacheImpl{}
}

func (c *tokenBlacklistCacheImpl) BanToken(token string, duration time.Duration) error {
	return client.Set(getBlacklistKey(token), token, duration).Err()
}

func getBlacklistKey(token string) string {
	tokenHashBytes := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(tokenHashBytes[:])
	return "blacklist:token:" + tokenHash
}

func (c *tokenBlacklistCacheImpl) IsTokenBanned(token string) (bool, error) {
	key := getBlacklistKey(token)
	exists, err := client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if exists == 0 {
		return false, nil
	}
	bannedToken, err := client.Get(key).Result()
	if err != nil {
		return false, err
	}
	return token == bannedToken, nil
}
