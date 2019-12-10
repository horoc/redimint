package ipfs

import (
	"bufio"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"os"
)

var ipfsClient *shell.Shell

/**
	普通场景，日常同步：
		1. 间隔height由各节点自己的配置文件里，各节点可动态￿修改，但是，必须多数节点保持一致才能保证系统正常运作（替代方案，程序里写死）
		2. 到了EndBlcok的时候，发现height到了，封锁全局锁，lock住Endblock
		3. 生成rdb, 开启全局锁，起线程，上传ipfs，拿到hash，签名，braodcommittx
		4. 收集满签名，如果与自身的一致，什么都不做
		5. 不一致，封锁全局锁，lock住Endblock，拉ipfs, 收集rdb height到目前锁的height的所有请求，重启redis，执行所有请求，解锁全局锁，

 */

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

func UpdateRDBFile(hash string){
	GetFile(hash, utils.Config.Redis.RDBPath)
}

func GetFile(hash string, path string) {
	err := ipfsClient.Get(hash, path)
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}
}
