package pkg

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisConnect() *redis.Client {
	redistHost := os.Getenv("RDSHOST")
	redisPort := os.Getenv("RDSPORT")

	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redistHost, redisPort),
	})
}
