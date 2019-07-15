package databases

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
)

var conn *DATABASE
var once sync.Once
var env string

type DATABASE struct {
	Redis    *redis.Client
}

func Setup() *DATABASE {
	once.Do(func() {
		conn = instances()
	})
	return conn
}

func instances() *DATABASE {
	log.Println("Connecting database")
	return &DATABASE{
		Redis: redisClient(),
	}
}

func ResetDatabases() {
	ResetRedisDB()
}

func CleanDatabases() {
	ResetRedisDB()
}