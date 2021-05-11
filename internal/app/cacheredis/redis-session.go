package cacheredis

import (
	"github.com/go-redis/redis"
	"go-starter-project/internal/app/cache"
	"go-starter-project/pkg/derror"
)

type userSessionCacheImpl struct{}

func GetUserSessionsCache() cache.UserSessionCache {
	return &userSessionCacheImpl{}
}

func (c *userSessionCacheImpl) GetSessionID(userID string) (string, error) {
	res, err := client.Get(getSessionIDKey(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", derror.ErrItemNotFound
		}
	}
	return res, nil
}

func (c *userSessionCacheImpl) SetSessionID(userID, sessionID string) error {
	return client.Set(getSessionIDKey(userID), sessionID, 0).Err()
}

func (c *userSessionCacheImpl) UnsetSessionID(userID string) error {
	return client.Del(getSessionIDKey(userID)).Err()
}

func getSessionIDKey(userID string) string {
	return "session:" + userID
}
