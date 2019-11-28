package network

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/consensus"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/service"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Server struct {
	url  string
	port string
	service service.Service
}

func NewServer(host string, port string) *Server {
	service := service.ServiceImpl{}
	url := host + ":" + port
	//node := NewNode(nodeID, url);
	server := &Server{url, port,service}
	server.setRoute()
	return server
}

func (server *Server) Start() {
	fmt.Printf("Server will be started at %s...\n", server.url)
	if err := http.ListenAndServe(server.url, nil); err != nil {
		logger.Error(err)
		return
	}
}

func (server *Server) setRoute() {

	//tendermint
	http.HandleFunc("/chain/height", server.getChainHeightBlock)
	http.HandleFunc("/chain/info", server.getChainInfo)
	http.HandleFunc("/chain/state", server.getChainState)
	http.HandleFunc("/chain/abci_query", server.getACBIQuery)
	http.HandleFunc("/chain/tx", server.getTxFromHash)
	http.HandleFunc("/chain/search_tx", server.getSearchTx)

	//db execute
	http.HandleFunc("/execute", server.executeReq)
	http.HandleFunc("/db/execute", server.execute)


	//db query
	http.HandleFunc("/query", server.getQuery)

	//redis
	http.HandleFunc("/db/dump", server.dump)

	http.HandleFunc("/logs", server.getLogsFromHeight)
	http.HandleFunc("/test_tps", server.testTps)

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

func (server *Server) getChainInfo(writer http.ResponseWriter, request *http.Request) {

	min := request.URL.Query().Get("min")
	max := request.URL.Query().Get("max")

	minH, err := strconv.Atoi(min)
	if err != nil {
		logger.Error(err)
		return
	}

	maxH, err := strconv.Atoi(max)
	if err != nil {
		logger.Error(err)
		return
	}

	info := consensus.GetChainInfo(minH, maxH)
	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(info))
}

func (server *Server) getChainHeightBlock(writer http.ResponseWriter, request *http.Request) {

	h := request.URL.Query().Get("h")
	height, err := strconv.Atoi(h)
	if err != nil {
		logger.Error(err)
		return
	}
	block := consensus.GetBlockFromHeight(height)

	var txHashList = make([]string, 0)
	for i := 0; i < len(block.Block.Data.Txs); i++ {
		txHashList = append(txHashList, fmt.Sprintf("%x", block.Block.Data.Txs[i].Hash()))
	}
	var obj map[string]interface{}
	err = json.Unmarshal(utils.StructToJson(block), &obj)
	if err != nil {
		logger.Error(err)
		return
	}
	obj["tx_hash"] = txHashList
	output, err := json.Marshal(obj)
	if err != nil {
		logger.Error(err)
		return
	}
	writer.Header().Set("Content-type", "application/json")
	writer.Write(output)
}

func (server *Server) getACBIQuery(writer http.ResponseWriter, request *http.Request) {
	data := request.URL.Query().Get("data")
	query := consensus.ABCIDataQuery("", []byte(data))
	//TODO state里面有字段是Byte显示的
	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(query))
}

func (server *Server) getChainState(writer http.ResponseWriter, request *http.Request) {

	state := consensus.GetChainState()
	//TODO state里面有字段是Byte显示的
	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(state))
}

func (server *Server) testTps(writer http.ResponseWriter, request *http.Request) {

	//url := "http://127.0.0.1:30001/execute"
	//payload := strings.NewReader("{\n    \"operation\": \"lpush xx 10\"\n}")
	//req, _ := http.NewRequest("POST", url, payload)
	//req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("cache-control", "no-cache")



	response := &TpsTest{Details: make([]TestDetail, 0)}

	//for i := 0; i < 10; i++ {
	//	detail := &TestDetail{}
	//	ti0 := time.Now()
	//	res, _ := http.DefaultClient.Do(req)
	//	body, _ := ioutil.ReadAll(res.Body)
	//	ti1 := time.Now()
	//	detail.Time = fmt.Sprint(ti1.Sub(ti0))
	//	detail.Response = string(body)
	//	response.Details = append(response.Details, *detail)
	//}

	t0 := time.Now()
	consensus.TestTendermintTime([]byte("asdasdasdasdasdasdasd"))
	t1 := time.Now()
	fmt.Printf("time 1 %s",fmt.Sprint(t1.Sub(t0)))


	t0 = time.Now()
	str := "http://" + utils.Config.Tendermint.Url + "/broadcast_tx_commit"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("tx", "\""+"asdasdasdasdasdasdasd"+"\"")
	u.RawQuery = q.Encode()
	request, _ = http.NewRequest("GET", fmt.Sprint(u), nil)
	res := utils.SendRequest(request)
	fmt.Println(res)
	t1 = time.Now()
	fmt.Printf("time 2 %s",fmt.Sprint(t1.Sub(t0)))


	response.TotalTime = fmt.Sprint(t1.Sub(t0))
	response.TotalTx = 10
	j, _ := json.Marshal(response)
	writer.Header().Set("Content-type", "application/json")
	writer.Write(j)
}

func (server *Server) getLogsFromHeight(writer http.ResponseWriter, request *http.Request) {
	h, err := strconv.Atoi(request.URL.Query().Get("height"))
	if err != nil {
		logger.Error(err)
		return
	}
	logs := consensus.LogStoreApp.GetLogFromHeight(h)
	writer.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(&QueryLogResponse{h, logs})
	if err != nil {
		logger.Error(err)
	}
	writer.Write(j)
}

func (server *Server) getQuery(writer http.ResponseWriter, request *http.Request) {
	var msg ExecutionRequest
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		logger.Error(err)
		return
	}
	res := database.ExecuteCommand(msg.Operation)
	resBody := &QueryResponse{}
	resBody.Operation = msg.Operation
	resBody.Result = res
	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(resBody))
}

func (server *Server) getTxFromHash(writer http.ResponseWriter, request *http.Request) {
	hash := request.URL.Query().Get("hash")
	bytes := utils.HexToByte(hash)
	tx := consensus.GetTx(bytes)

	var obj map[string]interface{}
	err := json.Unmarshal(utils.StructToJson(tx), &obj)
	if err != nil {
		logger.Error(err)
		return
	}
	obj["tx_decode"] = string(tx.Tx)

	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(obj))
}

func (server *Server) getSearchTx(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("query")
	tx := consensus.SearchTx(query, 1, 30)
	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(tx))
}


func (server *Server) execute(writer http.ResponseWriter, request *http.Request) {
	var msg ExecutionRequest
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		logger.Error(err)
		return
	}
	res := server.service.Execute(&models.ExecuteRequest{msg.Operation})

	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(res))
}

func (server *Server) executeReq(writer http.ResponseWriter, request *http.Request) {
	var msg ExecutionRequest
	err := json.NewDecoder(request.Body).Decode(&msg)

	if err != nil {
		logger.Error(err)
		return
	}
	op := &models.TxCommitBody{}

	u := uuid.NewV4()
	u1 := binary.BigEndian.Uint64(u[0:8])
	u2 := binary.BigEndian.Uint64(u[8:16])

	op.Sequence = fmt.Sprintf("%x%x", u1, u2)
	sign, err := utils.NodeKey.PrivKey.Sign([]byte(msg.Operation))
	if err != nil {
		logger.Error(err)
		return
	}

	op.Signature = utils.SignToHex(sign)
	op.Operation = msg.Operation
	commitMsg := consensus.BroadcastTxCommit(op)
	writer.Header().Set("Content-type", "application/json")
	writer.Write(utils.StructToJson(commitMsg))
}

func (server *Server) dump(writer http.ResponseWriter, request *http.Request) {
	val := database.DumpRDBFile()
	writer.Header().Set("Content-type", "application/json")
	writer.Write([]byte(val))
}