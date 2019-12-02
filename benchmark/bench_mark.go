package benchmark

import (
	"errors"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/service"
	"strings"
	"sync"
	"time"
)

type BenchMark struct {
	mutex          *sync.Mutex
	wg             *sync.WaitGroup
	connectionNums int
	txPerSecond    int
	txNums         int
	testBegin      time.Time
	testEnd        time.Time
	totalLatency   int64
	mode           string
	cmd            string
}

func (b *BenchMark) SendTestRequest(req *models.CommandRequest, mode string) {
	t1 := time.Now()
	if strings.EqualFold(mode, "sync") {
		service.AppService.ExecuteAsync(req)
	} else if strings.EqualFold(mode, "commit") {
		service.AppService.Execute(req)
	}
	t2 := time.Now()
	b.mutex.Lock()
	b.totalLatency += t2.Sub(t1).Nanoseconds() / 1e6
	b.mutex.Unlock()

	b.wg.Done()
}

func NewBenchMark(b *models.BenchMarkRequest) (*BenchMark, error) {
	if !strings.EqualFold(b.Mode, "sync") && !strings.EqualFold(b.Mode, "commit") {
		return nil, errors.New("Invalid mode")
	}
	return &BenchMark{
		mutex:          &sync.Mutex{},
		wg:             &sync.WaitGroup{},
		connectionNums: b.Connections,
		txNums:         b.TxNums,
		txPerSecond:    b.TxSendPerSec,
		mode:           b.Mode,
		cmd:            b.Cmd,
		totalLatency:   0,
	}, nil
}

func (b *BenchMark) StartTest() *models.BenchMarkResponse {
	b.wg.Add(b.txNums)
	//d := time.Duration(60.0/b.txPerSecond) * time.Second
	b.testBegin = time.Now()
	for i := 0; i < b.txNums; i++ {
		fmt.Printf("test %d transaction\n", i)
		go b.SendTestRequest(&models.CommandRequest{Cmd: b.cmd}, b.mode)
		//time.Sleep(d)
	}
	b.wg.Wait()
	b.testEnd = time.Now()
	duration := (b.testEnd.Sub(b.testBegin)).Nanoseconds() / 1e6
	return &models.BenchMarkResponse{
		Latency: &models.BenchMarkDetail{
			Avg:   fmt.Sprintf("%f ms", float64(b.totalLatency)/float64(b.txNums)),
			Max:   "",
			Stdev: "",
		},
		Tps: &models.BenchMarkDetail{
			Avg:   fmt.Sprintf("%f tps", float64(b.txNums)/(float64(duration)/1e3)),
			Max:   "",
			Stdev: "",
		},
	}

}
