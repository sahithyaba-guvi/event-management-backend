package ConnectRedis

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	var err = godotenv.Load("env/local/.env")
	if err != nil {
		fmt.Println(err)
	}
	return os.Getenv(key)
}

func ConnectToRedis() (*redis.Client, error) {
	// Create a new Redis client
	conn := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL") + ":6379",
		Password: "",
		DB:       0,
	})
	redisURl := getEnv("REDIS_URL") + ":6379"
	fmt.Println("redis url ", redisURl)
	pong, err := conn.Ping().Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil, err
	}
	fmt.Println("Redis connection successful:", pong)

	return conn, nil
}
