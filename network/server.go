package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/consensus"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	url  string
	port string
	//node *Node
}

//func NewServer(nodeID string) *Server {
//	node := NewNode(nodeID)
//	server := &Server{node.NodeTable[nodeID], node}
//
//	server.setRoute()
//
//	return server
//}


func NewServer(host string, port string) *Server {
	url := host + ":" + port
	//node := NewNode(nodeID, url);
	server := &Server{url, port}
	server.setRoute()
	return server
}

//
//func (server *Server) JoinNode(nodeTable map[string]string) {
//	server.node.JoinNode(nodeTable)
//	fmt.Println(server.node.NodeTable)
//}

func (server *Server) Start() {
	fmt.Printf("Server will be started at %s...\n", server.url)
	if err := http.ListenAndServe(server.url, nil); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) setRoute() {
	http.HandleFunc("/execute", server.getReq)
	//http.HandleFunc("/req", server.getReq)
	//http.HandleFunc("/preprepare", server.getPrePrepare)
	//http.HandleFunc("/prepare", server.getPrepare)
	//http.HandleFunc("/commit", server.getCommit)
	//http.HandleFunc("/reply", server.getReply)
	//http.HandleFunc("/restore", server.doRestore)
}

//
//func (server *Server) doRestore(writer http.ResponseWriter, request *http.Request) {
//	//TODO 重新加载rdb和日志文件
//}
//
//
//
//func (server *Server) SendConsensus(request *pbft.RequestMsg) *pbft.ComfirmedMsg {
//	server.node.routeMsg(request)
//	sequenceID := request.SequenceID
//
//	select {
//	case <-server.node.ComfirmedEvent:
//		return &pbft.ComfirmedMsg{
//			Msg:        request.Operation,
//			SequenceID: sequenceID,
//		}
//	case <-time.After(100 * time.Second):
//		return nil
//	}
//	return nil
//
//}

func (server *Server) getReq(writer http.ResponseWriter, request *http.Request) {
	var msg RequestMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := consensus.GetBroadcastTxCommitRequest(msg.Operation)
	if err != nil {
	}
	response := sendRequest(req)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	command := strings.Split(msg.Operation, "=")[1]
	res := database.ExecuteCommand(command)
	fmt.Println("res:"+res)
}

//func (server *Server) getPrePrepare(writer http.ResponseWriter, request *http.Request) {
//	var msg pbft.PrePrepareMsg
//	err := json.NewDecoder(request.Body).Decode(&msg)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	server.node.MsgEntrance <- &msg
//}
//
//func (server *Server) getPrepare(writer http.ResponseWriter, request *http.Request) {
//	var msg pbft.VoteMsg
//	err := json.NewDecoder(request.Body).Decode(&msg)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	server.node.MsgEntrance <- &msg
//}
//
//func (server *Server) getCommit(writer http.ResponseWriter, request *http.Request) {
//	var msg pbft.VoteMsg
//	err := json.NewDecoder(request.Body).Decode(&msg)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	server.node.MsgEntrance <- &msg
//}
//
//func (server *Server) getReply(writer http.ResponseWriter, request *http.Request) {
//	var msg pbft.ReplyMsg
//	err := json.NewDecoder(request.Body).Decode(&msg)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	server.node.GetReply(&msg)
//}

func send(url string, msg []byte) {
	buff := bytes.NewBuffer(msg)
	fmt.Println("send : " + url)
	http.Post("http://"+url, "application/json", buff)
}
func sendRequest(r *http.Request) (*http.Response) {
	client := http.Client{}
	fmt.Println("client request:")
	fmt.Println(r)
	response, e := client.Do(r)
	if e != nil {
		fmt.Println(e)
	}
	//response.Body.Close()
	fmt.Println("client response:")
	fmt.Println(response)
	return response

}
