package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"path/filepath"
	"strings"
)

var NodeKey *p2p.NodeKey
var ValidatorKey *privval.FilePVKey

func InitKey(){
	InitNodeKey()
	InitValidatorKey()
}

func InitNodeKey() {
	key, err := p2p.LoadNodeKey("./chain/config/node_key.json")
	if err != nil {
		logger.Error(err)
		return
	}
	NodeKey = key
	fmt.Println(strings.ToUpper(hex.EncodeToString(NodeKey.PrivKey.PubKey().Address())))
}

func InitValidatorKey(){
	keyFile := filepath.Join("./chain", "config", "priv_validator_key.json")
	stateFile := filepath.Join("./chain", "data", "priv_validator_state.json")
	fpv := privval.LoadFilePV(keyFile, stateFile)
	ValidatorKey = &fpv.Key
}

func GetNodeID() string{
	return string(p2p.PubKeyToID(NodeKey.PubKey()))
}

func NodeSign(msg []byte) []byte {
	bytes, err := NodeKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return bytes
}

func NodeStringSign(msg []byte) string {
	bytes, err := NodeKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return SignToHex(bytes)
}

func ValidatorSign(msg []byte) []byte {
	bytes, err := ValidatorKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return bytes
}

func ValidatorStringSign(msg []byte) string {
	bytes, err := ValidatorKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return SignToHex(bytes)
}


