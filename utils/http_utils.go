package utils

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"net/http"
)

var HttpClient *http.Client

func SendRequest(r *http.Request) (*http.Response) {
	if HttpClient == nil {
		HttpClient = &http.Client{}
	}
	fmt.Println("client request:")
	fmt.Println(r)
	response, err := HttpClient.Do(r)
	if err != nil {
		logger.Error(err)
		return nil
	}
	fmt.Println("client response:")
	fmt.Println(response)
	return response
}
