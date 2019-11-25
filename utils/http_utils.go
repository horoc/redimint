package utils

import (
	"fmt"
	"net/http"
)

var HttpClient *http.Client

func SendRequest(r *http.Request) (*http.Response) {
	if HttpClient == nil{
		HttpClient = &http.Client{}
	}
	fmt.Println("client request:")
	fmt.Println(r)
	response, e := HttpClient.Do(r)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("client response:")
	fmt.Println(response)
	return response
}