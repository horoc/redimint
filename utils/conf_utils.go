package utils

import (
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type SerivceConfig struct {
	Redis      RedisConfig      `yaml:"redis-config"`
	Tendermint TendermintConfig `yaml:"tendermint-config"`
	IPFS IPFSConfig `yaml:"ipfs-config"`
}

type RedisConfig struct {
	Url      string `yaml:"url"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
	RDBPath  string `yaml:"rdb_path"`
}

type TendermintConfig struct {
	Url string `yaml:"url"`
}

type IPFSConfig struct{
	Url string `yaml:"url"`
}

var Config SerivceConfig

func InitConfig() {
	abspath, err := filepath.Abs("./conf/service.yaml")
	if err != nil {
	}
	yamlFile, err := ioutil.ReadFile(abspath)
	if err != nil {
		logger.Error(err)
		return
	}
	yaml.Unmarshal(yamlFile, &Config)

}
