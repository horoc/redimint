package utils

import (
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/tendermint/tendermint/p2p"
)

var NodeKey *p2p.NodeKey

func InitNodeKey() {
	key, err := p2p.LoadNodeKey("./chain/config/node_key.json")
	if err != nil {
		logger.Error(err)
		return
	}
	NodeKey = key
}

func Sign(msg []byte) []byte {
	bytes, err := NodeKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return bytes
}


