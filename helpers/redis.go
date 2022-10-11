package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var redis_host,
	redis_port,
	redis_pass,
	redis_addr, environment string

func init() {
	environment = "development"
}

func RedisConnection() *redis.Client {
	if environment == "test" {
		redis_host = os.Getenv("REDIS_HOST_TEST")
		redis_port = os.Getenv("REDIS_PORT_TEST")
		redis_pass = os.Getenv("REDIS_PASSWORD_TEST")
	} else if environment == "development" {
		redis_host = os.Getenv("REDIS_HOST_DEV")
		redis_port = os.Getenv("REDIS_PORT_DEV")
		redis_pass = os.Getenv("REDIS_PASSWORD_DEV")
	} else {
		redis_host = os.Getenv("REDIS_HOST")
		redis_port = os.Getenv("REDIS_PORT")
		redis_pass = os.Getenv("REDIS_PASSWORD")
	}

	redis_addr = fmt.Sprintf("%s:%s", redis_host, redis_port)
	client := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(pong)
	return client
}

func SetRedisValue(key string, value string) bool {
	client := RedisConnection()
	err := client.Set(key, value, 1*time.Minute).Err()
	if err != nil {
		log.Println(err)
	}
	client.Close()
	return true
}

func GetRedisValue(key string) string {
	client := RedisConnection()
	get, err := client.Get(key).Result()
	if err != nil {
		log.Println(err)
	}
	client.Close()
	return get
}

func DeleteRedisValue(keys []string) bool {
	client := RedisConnection()
	err := client.Del(keys...).Err()
	if err != nil {
		log.Println(err)
	}
	client.Close()
	return true
}

func SearchRedisValue(keys string) []string {
	client := RedisConnection()
	search, err := client.Keys(keys).Result()
	if err != nil {
		log.Println(err)
	}
	client.Close()
	return search
}
