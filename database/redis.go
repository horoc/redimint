package database

import (
	"errors"
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

func ExecuteCommand(command string) (string, error) {
	split := strings.Split(command, " ")
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

func IsValidCmd(command string) bool {
	return true
}

func IsQueryCmd(command string) bool {
	return false
}

func GetKey(command string) (string, error) {
	if !IsValidCmd(command) || IsQueryCmd(command) {
		return "", errors.New("Not a update command : " + command)
	}
	split := strings.Split(command, " ")
	if len(split) < 2 {
		return "", errors.New("Invalid command : " + command)
	}
	return split[1], nil
}

func IsPrivateKey(command string) bool {
	if !IsValidCmd(command) {
		logger.Error("Invalid command : " + command)
		return false
	}
	split := strings.Split(command, " ")
	if len(split) < 2 {
		logger.Error("Invalid command : " + command)
		return false
	}
	return strings.HasPrefix(split[1], utils.ValidatorKey.Address.String())
}

func ReplaceKey(command string, key string) (string, error) {
	if !IsValidCmd(command) {
		return "", errors.New("Invalid command : " + command)
	}
	split := strings.Split(command, " ")
	if len(split) < 2 {
		return "", errors.New("Invalid command : " + command)
	}
	split[1] = key
	cmd := strings.Join(split, " ")
	return cmd, nil
}
