package plugins

import "github.com/chenzhou9513/redimint/logger"

type DefaultTransactionPlugin struct {
}

func (d DefaultTransactionPlugin) CustomTxValidationCheck(tx []byte) (bool, string) {
	logger.Log.Info("Default transaction plugin CustomTxValidationCheck: do nothing")
	return true, ""
}

func (d DefaultTransactionPlugin) CustomTransactionDeliverLog(tx []byte, result string) string {
	logger.Log.Info("Default transaction plugin CustomTransactionDeliverLog: do nothing")
	return ""
}

func init() {
	register("default", &DefaultTransactionPlugin{})
}
