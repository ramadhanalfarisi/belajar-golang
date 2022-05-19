package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func RedisConnection() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(pong)
	return client
}

func SetRedisValue(key string, value string) bool {
	client := RedisConnection()
	err := client.Set(key, value, 1  * time.Minute).Err()
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
	return true
}

func GetRedisValue(key string) string {
	client := RedisConnection()
	get, err := client.Get(key).Result()
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
	return get
}

func DeleteRedisValue(keys []string) bool {
	client := RedisConnection()
	err := client.Del(keys...).Err()
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
	return true
}

func SearchRedisValue(keys string) []string{
	client := RedisConnection()
	search,err := client.Keys(keys).Result();
	if err != nil{
		log.Fatal(err)
	}
	client.Close()
	return search
}