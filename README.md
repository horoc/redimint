
<p align="center"><img src="https://chenzhou-images.oss-cn-shanghai.aliyuncs.com/redimint.png" /></p>

<p align="center"> Redimint 是基于 Redis 和 Tendermint 区块链中间件的去中心化KV数据库系统。</p>

## 概述 

### 目标

- 为`DAPP (Decentralized Application)` 服务的数据库组件
- 为多组织合作场景下数据库可信的需求提供去中心化数据库组件

### Redimint 特性

- 去中心设计，数据库集群不存在主节点，任何节点均可读写数据。
- 通过与区块链相结合，保证任何节点的操作日志不可篡改，可信且可溯源。
- 提供Commit, Private, Async三种数据更新模式，最大化TPS和最小化更新延迟。
- 兼容几乎所有的Redis操作，使用学习成本低。
- 插件化编程，提供插件接口，可定制化区块生成或操作执行中特定步骤的行为。

### Rediminit数据存储逻辑

<p align="center"><img src="https://chenzhou-images.oss-cn-shanghai.aliyuncs.com/%E6%95%B0%E6%8D%AE%E5%AD%98%E5%82%A8%E8%AF%B4%E6%98%8E.png" /></p>

### Redimint 架构

<p align="center"><img src="https://chenzhou-images.oss-cn-shanghai.aliyuncs.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%9E%B6%E6%9E%84%E5%9B%BE%20%283%29.png" /></p>

### TODO

- [X] Redis monitor
- [ ] 使用文档完善
- [ ] 性能测试
- [ ] ... ...

## Getting Started

### 依赖环境安装

Redimint的运行依托于Redis数据库和Tendermint服务, 需要预先安装Tendermint和Redis

#### Redis

1. 安装 Redis要求Redis版本4.0以上。安装方法 [参考这里](https://redis.io/download) 
2. 要求redis-server的环境变量。
   
```bash
$ redis-server --help
Usage: ./redis-server [/path/to/redis.conf] [options]
       ./redis-server - (read config from stdin)
       ./redis-server -v or --version
       ./redis-server -h or --help
       ./redis-server --test-memory <megabytes>

Examples:
       ./redis-server (run the server with default conf)
       ./redis-server /etc/redis/6379.conf
       ./redis-server --port 7777
       ./redis-server --port 7777 --slaveof 127.0.0.1 8888
       ./redis-server /etc/myredis.conf --loglevel verbose

Sentinel mode:
       ./redis-server /etc/sentinel.conf --sentinel
```

#### Tendermint

`Redimint` 基于 [Tendermint](https://github.com/tendermint/tendermint) 0.32.7版本 

Linux:

```bash
sudo yum install -y unzip
wget https://github.com/tendermint/tendermint/releases/download/v0.32.7/tendermint_v0.32.7_linux_amd64.zip
unzip tendermint_v0.32.7_linux_amd64.zip
rm tendermint_v0.32.7_linux_amd64.zip
sudo mv tendermint /usr/local/bin
```

MacOs:

```bash
sudo brew install -y unzip
wget https://github.com/tendermint/tendermint/releases/download/v0.32.7/tendermint_v0.32.7_darwin_amd64.zip
unzip tendermint_v0.32.7_darwin_amd64.zip
rm tendermint_v0.32.7_darwin_amd64.zip
sudo mv tendermint /usr/local/bin
```

查看版本：

```bash
$ tendermint
Tendermint Core (BFT Consensus) in Go

Usage:
  tendermint [command]

Available Commands:
  gen_node_key                Generate a node key for this node and print its ID
  gen_validator               Generate new validator keypair
  help                        Help about any command
  init                        Initialize Tendermint
  lite                        Run lite-client proxy server, verifying tendermint rpc
  node                        Run the tendermint node
  probe_upnp                  Test UPnP functionality
  replay                      Replay messages from WAL
  replay_console              Replay messages from WAL in a console
  show_node_id                Show this node's ID
  show_validator              Show this node's validator info
  testnet                     Initialize files for a Tendermint testnet
  unsafe_reset_all            (unsafe) Remove all the data and WAL, reset this node's validator to genesis state
  unsafe_reset_priv_validator (unsafe) Reset this node's validator to genesis state
  version                     Show version info

Flags:
  -h, --help               help for tendermint
      --home string        directory for config and data (default "/Users/chenzhou/.tendermint")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors

Use "tendermint [command] --help" for more information about a command.
```

### Redimint 安装

#### 1. 安装 Go 运行环境 [参考这里](https://golang.org/doc/install)
要求Go版本>=1.12

```bash
$ go version
go version go1.13.4 darwin/amd64
```
#### 2. 设置编译环境

添加 `$GOPATH/bin` 到 `$PATH`，例如：`PATH=$PATH:$GOPATH/bin`。

```bash
$ go env GOPATH
/home/go
```
#### 3. 下载 Redimint 源代码

```bash
$ mkdir -p $GOPATH/src/github.com/chenzhou9513
$ cd $GOPATH/src/github.com/chenzhou9513 && git clone https://github.com/chenzhou9513/redimint.git -b master
```
#### 4. 编译 Redimint 源代码

```bash
$ cd $GOPATH/src/github.com/chenzhou9513/redimint
$ make
```
make 默认会在当前目录创建redimint_home目录, 目录结构：

    ├── 默认会在当前目录创建redimint_home目录
        │   ├── bin
        │   │   ├── redimint
        │   ├── conf
        │   │   ├── redis.conf
        │   │   ├── configuration.yaml
        │   ├── chain
        │   │   ├── config
        │   │   │   ├── genesis.json
        │   │   │   ├── config.toml
        │   │   │   ├── ... ...
        │   │   ├── data
        │   │   │   ├── ... ...
        

通过output参数自定义redimint home目录地址

```bash
make output=my_home_path
```

#### 4. 启动 Redimint

验证Redimint编译是否成功

```bash
$ cd ./redimint_home/bin
$ ./redimint

Description:
  Decentralized K-V database based on Redis and Tendermint Blockchain.

Usage:
  redimint [command]

Available Commands:
  help        Help about any command
  init        Initialization redimint service
  restart     Restart redimint server
  start       Start redimint server
  stop        Stop redimint server
  version     Get redimint version

Flags:
  -h, --help   help for redimint

Use "redimint [command] --help" for more information about a command.

```

单节点启动Redimint服务，后台启动使用 -d

```bash
$ cd ./redimint_home/bin
$ ./redimint start

Redis daemon process started
badger 2019/12/26 20:42:40 INFO: All 0 tables opened in 0s
Tendermint daemon process started
I[2019-12-26|20:42:40.952] Starting ABCIServer                          impl=ABCIServer
I[2019-12-26|20:42:40.957] Waiting for new connection...
I[2019-12-26|20:42:41.076] Accepted a new connection
I[2019-12-26|20:42:41.076] Waiting for new connection...
I[2019-12-26|20:42:41.076] Accepted a new connection
I[2019-12-26|20:42:41.076] Waiting for new connection...
I[2019-12-26|20:42:41.076] Accepted a new connection
I[2019-12-26|20:42:41.076] Waiting for new connection...
```

多节点启动Redimint服务:

1.获取每个节点的Node address 

```bash
$ tendermint show_node_id --home ./redimint_home/chain

12c1fb57614c43761e8bbe65c4454be11e756267
```

2.修改./redimint_home/chain/config/config.toml文件, 写入每个节点的address和ip

```bash
persistent_peers = "12c1fb57614c43761e8bbe65c4454be11e756267@IP0:26656,8a223d1493fa45496a4fa1b054d2a7dd6116b50c@IP1:26656,d46f4422738b543bda1dfae06973896d290385c4@IP2:26656,744a8d89611d584dd88055e6eddf625212705c20@IP3:26656"
```

3.修改genesis.json文件, 添加validator (power可自定义，建议不小于10, pub_key由./redimint_home/chain/config/priv_validator_key.json文件里获取)，示例：

```bash
"validators": [
    {
      "address": "540B94741F1A64B787A6D885A0E382A54249B659",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "R1lIn7piIdxfoB/hUSN6kBukTiD27GI2EuiPH6isn4Y="
      },
      "power": "10",
      "name": "node0"
    },
    {
      "address": "B8F3744A0E5DB5411932B2848F0A2C1F47A0B2AD",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "ECfG6qfbiS+1cNDJlVqEzmFkYKRCHpqneCfyjeeFsPg="
      },
      "power": "10",
      "name": "node1"
    },
    {
      "address": "6B760E33FF91A5376F754EFB5E56481ADE93222E",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "no5aL13aED+rlVEvmEp+UjS6JRrsUtXuQRyCQtjeFlg="
      },
      "power": "10",
      "name": "node2"
    },
    {
      "address": "B4B9A08F2965E4E8BA07570FB6CED4024252AC0B",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "vfaA1//UmE76BukRvtG1KftlPQ38SAipdbjPpNKDmq8="
      },
      "power": "10",
      "name": "node3"
    }
  ],
```

4.如果仅用于测试，可以使用test_home目录下的四个已配置好的节点配置文件。

5.启动各个节点，启动方式与单节点一致

#### 5. 测试

数据库执行 `set k v ` 指令

```bash
$ curl -X POST http://localhost:30001/db/execute -H 'Content-Type: application/json' -d '{"cmd":"set k v","mode":"commit"}'

{
    "code":0,
    "code_info":"Success!",
    "data":{
        "command":"set k v",
        "execute_result":"Result:OK",
        "signature":"F869FC98EFF68069760A70D4F0C39E6D103AFB004623FC97C5D89D58C82484E9595DCAF62E5866C93A2429E07EDE56969541993748B11045C591398ECC8D1803",
        "sequence":"03AB90AE27E011EABCB588E9FE67FE42",
        "time_stamp":"1577365367595",
        "hash":"D69821214F6AED6BBBDFE374916575ABBCF73A3CDED6790FF18C587DC15E0468",
        "height":5
    }
}
```

#### 6. 退出Redimint

```bash
$ ./redimint stop   
```

## 使用文档

> TODO

## Author

* **Chen Zhou** - *Initial work* - [ChenZhou](https://github.com/chenzhou9513)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

