package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
)

func StructToJson(v interface{}) []byte {
	json, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return json
}

func JsonToStruct(j []byte, s interface{}) {
	err := json.Unmarshal(j, s)
	if err != nil {
		logger.Error(err)
		return
	}
}

func ByteToHex(b []byte) string{
	return fmt.Sprintf("%X",b)
}

func HexToByte(s string) []byte{
	bytes, err := hex.DecodeString(s)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return bytes
}

func SignToHex(b []byte) string{
	return fmt.Sprintf("%X",b)
}


