package cacheredis

import (
	"github.com/go-redis/redis"
	"go-starter-project/internal/app/config"
	"go-starter-project/pkg/log"
)

var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
