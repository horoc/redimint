package database

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/go-redis/redis"
	"os/exec"
	"strings"
)

var Client *redis.Client

func InitRedis() {
	StartRedisServer()
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

func StopRedis() {
	status := Client.Shutdown()
	if status.Err() != nil {
		logger.Error(status.Err())
	}
}

func StartRedisServer() {
	cmd := exec.Command(utils.Config.Redis.RedisBin, utils.Config.Redis.ConfPath)
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
	}
}

func DumpRDBFile() string {
	save := Client.Save()
	return save.Val()
}

func ExecuteCommand(commond string) (string, error) {
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
		return "", err
	}
	return fmt.Sprintf("%v", s), nil
}
