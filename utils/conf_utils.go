package utils

import (
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Configuration struct {
	Redis      RedisConfig      `yaml:"redis-config"`
	Tendermint TendermintConfig `yaml:"tendermint-config"`
	IPFS       IPFSConfig       `yaml:"ipfs-config"`
	HttpServer HttpServerConfig `yaml:"http-server-config"`
	Rpc        RpcConfig        `yaml:"rpc-config"`
}

type RedisConfig struct {
	Url      string `yaml:"url"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
	RDBPath  string `yaml:"rdb_path"`
	RedisBin string `yaml:"redis_bin"`
	ConfPath string `yaml:"conf_path"`
}

type TendermintConfig struct {
	Url string `yaml:"url"`
}

type IPFSConfig struct {
	Url string `yaml:"url"`
}

type HttpServerConfig struct {
	RunMode      string        `yaml:"run_mode"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type RpcConfig struct {
	Port int `yaml:"port"`
}

var Config Configuration

func InitConfig() {
	abspath, err := filepath.Abs("./conf/configuration.yaml")
	if err != nil {
	}
	yamlFile, err := ioutil.ReadFile(abspath)
	if err != nil {
		logger.Error(err)
		return
	}
	yaml.Unmarshal(yamlFile, &Config)

}
