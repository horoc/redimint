package plugins

import (
	"github.com/chenzhou9513/redimint/utils"
)

var pluginsMap map[string]TransactionPlugin

type TransactionPlugin interface {
	CustomTxValidationCheck(tx []byte) (bool, string)
	CustomTransactionDeliverLog(tx []byte, result string) string
}

func register(name string, plugin interface{}) {
	if pluginsMap == nil {
		pluginsMap = make(map[string]TransactionPlugin, 0)
	}
	pluginsMap[name] = plugin.(TransactionPlugin)
}

func GetConfigPlugin() TransactionPlugin {
	var plugin TransactionPlugin
	customPlugin := utils.Config.App.Plugin
	if pluginsMap == nil || !containsPlugin(customPlugin) {
		plugin = &DefaultTransactionPlugin{}
	} else {
		plugin = pluginsMap[customPlugin]
	}
	return plugin
}

func containsPlugin(name string) bool {
	_, ok := pluginsMap[name]
	return ok
}
