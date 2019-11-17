package consensus

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
)

func GetBroadcastTxCommitRequest(operation string)(*http.Request, error) {
	str := "http://127.0.0.1:26657/broadcast_tx_commit"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("tx", "\""+operation+"\"")
	u.RawQuery = q.Encode()
	fmt.Println(u)
	request, e := http.NewRequest("GET",fmt.Sprint(u),nil)

	if e!=nil{
		fmt.Println(e)
	}
	return request,nil
}

func GetQueryRequest(operation string)(*http.Request, error) {
	str := "http://127.0.0.1:26657/abci_query"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("data", "\""+operation+"\"")
	u.RawQuery = q.Encode()
	fmt.Println(u)
	request, e := http.NewRequest("GET",fmt.Sprint(u),nil)

	if e!=nil{
		fmt.Println(e)
	}
	return request,nil
}

func SendBroadcastTxCommitRequest(operation string)(string, error) {
	fmt.Println(operation)
	fmt.Println("localhost:26657/broadcast_tx_commit?tx=\""+operation+"\"")
	bytes, e := exec.Command("curl", "-s", "localhost:26657/broadcast_tx_commit?tx=\""+operation+"\"").Output()
	if e != nil {
		//log.Fatal(err)
	}
	fmt.Print(string(bytes))
	return string(bytes),e
}