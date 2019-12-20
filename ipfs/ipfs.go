package ipfs

import (
	"bufio"
	"github.com/chenzhou9513/redimint/logger"
	"github.com/chenzhou9513/redimint/utils"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"os"
)

var ipfsClient *shell.Shell

func InitIPFS() {
	ipfsClient = shell.NewShell(utils.Config.IPFS.Url)
}

func AddFile(r io.Reader) string {
	hash, err := ipfsClient.Add(r)
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}
	return hash
}

func UploadRDB() string {
	file, err := os.Open(utils.Config.Redis.RDBPath)
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	hash := AddFile(r)
	return hash
}

func UpdateRDBFile(hash string) {
	GetFile(hash, utils.Config.Redis.RDBPath)
}

func GetFile(hash string, path string) {
	err := ipfsClient.Get(hash, path)
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}
}
