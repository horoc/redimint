package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"

)

type SerivceConfig struct {
	Redis RedisConfig `yaml:"redis-config"`
	Tendermint TendermintConfig `yaml:"tendermint-config"`
}

type RedisConfig struct {
	Url string `yaml:"url"`
	Db int `yaml:"db"`
	Password string `yaml:"password"`
}

type TendermintConfig struct{
	Url string `yaml:"url"`
}

var Config SerivceConfig

func InitConfig()(){
	abspath, err := filepath.Abs("./conf/service.yaml")
	if err != nil{
		fmt.Println(err)
	}
	yamlFile, err := ioutil.ReadFile(abspath)
	if err != nil{
		fmt.Print(err)
	}
	yaml.Unmarshal(yamlFile,&Config)

}
