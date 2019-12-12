package utils

import (
	"bufio"
	"fmt"
	"os"
)

const DbTxLogFilePath = "./log/db_transactions"

func InitFiles() {
	os.Remove(DbTxLogFilePath)
}

func AppendToDBLogFile(lines []string) error {
	dbTxLogFile, err := os.OpenFile(DbTxLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer dbTxLogFile.Close()

	w := bufio.NewWriter(dbTxLogFile)
	for _, v := range lines {
		w.WriteString(v + "\n")
	}
	return w.Flush()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		return false
	}
}
