package database

import (
	"os"

	"github.com/go-redis/redis"
)

// ConnectRedis connects with Redis based on an instance in GCP
func ConnectRedis() *redis.Client {
	REDIS_IP_PORT := os.Getenv("REDIS_IP_PORT")
	client := redis.NewClient(&redis.Options{
		Addr: REDIS_IP_PORT,
		DB:   0,
	})

	return client
}
