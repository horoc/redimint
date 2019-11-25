package logrus

import (
	"fmt"
	"os"
)
import  "github.com/sirupsen/logrus"

var logger = logrus.New()


func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func InitLogger(){
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	customFormatter.DisableLevelTruncation = false

	logger.SetFormatter(customFormatter)

	f, err := os.OpenFile("./log/server.log", os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	logger.SetOutput(f)
	logger.SetLevel(logrus.InfoLevel)
}