#!/bin/bash

function start(){
    rm -rf tendermint.sock
    go build
    nohup ./redimint > ./log/db.log 2>&1 &
    nohup ./tendermint --home=./chain node --proxy_app=unix://tendermint.sock > ./log/tendermint.log 2>&1 &
}

function stop() {
    PID=$(ps -ef|grep redimint|grep -v grep|awk '{print $2}')
    if [ -z $PID ]; then
      echo "process redimint not exist"
    else
      echo "process id: $PID"
      kill -9 ${PID}
      echo "process redimint killed"
    fi

    PID=$(ps -ef|grep tendermint|grep -v grep|awk '{print $2}')
    if [ -z $PID ]; then
      echo "process tendermint not exist"
    else
      echo "process id: $PID"
      kill -9 ${PID}
      echo "process tendermint killed"
    fi
}

function reinstall() {
    rm -rf ./chain
    rm -rf /tmp/badger
    rm -rf ./tendermint.sock
    ./tendermint init --home=./chain
    cp -f ./conf/tendermint/config.toml ./chain/config/config.toml
}
function stopRedis(){
    PID=$(ps -ef|grep redis-server|grep -v grep|awk '{print $2}')
    if [ -z $PID ]; then
      echo "process redis not exist"
    else
      echo "process id: $PID"
      kill -9 ${PID}
      echo "process redis killed"
    fi
}
function startRedis(){
    nohup redis-server ./conf/redis/redis.conf > ./log/redis.log 2>&1 &
}

case "$1" in
  newDockernode)
    rm -rf ./chain
    mkdir ./chain
    cp -fr ./conf/dockernode/$2/* ./chain
    ;;
  newDocker2node)
    rm -rf ./chain
    mkdir ./chain
    cp -fr ./conf/docker2node/$2/* ./chain
    ;;
  start)
    start
    ;;
  stop)
    stopRedis
    rm -rf dump.rdb
    stop
    ;;
  reinstall)
    stopRedis
    rm -rf dump.rdb
    stop
    reinstall
    ;;
  restart)
    stopRedis
    rm -rf dump.rdb
    stop
    start
    ;;
  newNode)
    rm -rf ./chain
    mkdir ./chain
    cp -fr ./conf/testnet/$2/* ./chain
    ;;
  startRedis)
    startRedis
    ;;
  reinstallRedis)
    stopRedis
    rm -rf dump.rdb
    startRedis
    ;;
  restartRedis)
    stopRedis
    startRedis
    ;;
  start2Docker)
    docker kill $(docker ps -q)
    docker rm $(docker ps -a -q)

    docker run -tid --name node0 --privileged=true -p 30002:30001 base_test_env /sbin/init
    docker run -tid --name node1 --privileged=true  base_test_env /sbin/init

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/redimint node0:/home/redimint
    docker cp /Users/chenzhou/tendermint/linux/tendermint node0:/home/redimint/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/redimint node1:/home/redimint
    docker cp /Users/chenzhou/tendermint/linux/tendermint node1:/home/redimint/tendermint
    ;;
  startAllDocker)
    docker kill $(docker ps -q)
    docker rm $(docker ps -a -q)

    docker run -tid --name node0 --privileged=true -p 30002:30001 base_test_env /sbin/init
    docker run -tid --name node1 --privileged=true  base_test_env /sbin/init
    docker run -tid --name node2 --privileged=true  base_test_env /sbin/init
    docker run -tid --name node3 --privileged=true  base_test_env /sbin/init

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/redimint node0:/home/redimint
    docker cp /Users/chenzhou/tendermint/linux/tendermint node0:/home/redimint/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/redimint node1:/home/redimint
    docker cp /Users/chenzhou/tendermint/linux/tendermint node1:/home/redimint/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/redimint node2:/home/redimint
    docker cp /Users/chenzhou/tendermint/linux/tendermint node2:/home/redimint/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/redimint node3:/home/redimint
    docker cp /Users/chenzhou/tendermint/linux/tendermint node3:/home/redimint/tendermint
    ;;
esac
exit $RETVAL