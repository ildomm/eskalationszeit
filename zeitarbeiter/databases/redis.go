package databases

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/config"
	"log"
)

func Redis() *redis.Client {
	return conn.Redis
}

func redisClient() *redis.Client {
	redisInfo := fmt.Sprintf("%s:%s",
		config.App.Database.Host,
		config.App.Database.Port )

	log.Println("Redis: " + redisInfo)

	conn := redis.NewClient(&redis.Options{
		Addr:     redisInfo,
		Password: "",
		DB:       0,
	})

	_, err := conn.Ping().Result()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return conn
}

func ResetRedisDB() {
	Redis().FlushDB()
}