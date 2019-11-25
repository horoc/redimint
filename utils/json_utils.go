package utils

import (
	"encoding/json"
	"fmt"
)

func StructToJson(v interface{}) []byte{
	json, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return json
}


