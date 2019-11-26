package database

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/go-redis/redis"
	"strings"
)

var Client *redis.Client

func InitRedisClient() {
	Client = NewRedisClient()
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     utils.Config.Redis.Url,
		Password: utils.Config.Redis.Password, // no password set
		DB:       utils.Config.Redis.Db,       // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

func ExecuteCommand(commond string) string {
	split := strings.Split(commond, " ")
	slice := make([]interface{}, len(split))
	for i := 0; i < len(split); i++ {
		slice[i] = split[i]
	}
	cmd := redis.NewCmd(slice...)
	Client.Process(cmd)
	s, err := cmd.Result()
	if err != nil {
		logger.Error(err)
		return err.Error()
	}
	return fmt.Sprintf("%v", s)
}
