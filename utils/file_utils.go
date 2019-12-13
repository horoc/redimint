package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

//command | addr | sign | seq | height | time
func ReadTxFromDBLogFile(beginHeight int, endHeight int) ([]string, error) {
	f, err := os.Open(DbTxLogFilePath)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	strList := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		split := strings.Split(line, " | ")
		if len(split) < 6 {
			break
		}
		height, err := strconv.Atoi(split[4])
		if err != nil {
			return nil, err
		}
		if height > endHeight {
			break
		} else if height < beginHeight {
			continue
		}
		strList = append(strList, split[0])
	}
	return strList, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		return false
	}
}
