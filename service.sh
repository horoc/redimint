#!/bin/bash

function start(){
    rm -rf tendermint.sock
    go build
    nohup ./DecentralizedRedis > ./log/db.log 2>&1 &
    nohup ./tendermint --home=./chain node --proxy_app=unix://tendermint.sock > ./log/tendermint.log 2>&1 &
}

function stop() {
    PID=$(ps -ef|grep DecentralizedRedis|grep -v grep|awk '{print $2}')
    if [ -z $PID ]; then
      echo "process DecentralizedRedis not exist"
    else
      echo "process id: $PID"
      kill -9 ${PID}
      echo "process DecentralizedRedis killed"
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
function testTPS(){
    curl -s 'http://127.0.0.1:30001/test_tps'
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

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis node0:/home/DecentralizedRedis
    docker cp /Users/chenzhou/tendermint/linux/tendermint node0:/home/DecentralizedRedis/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis node1:/home/DecentralizedRedis
    docker cp /Users/chenzhou/tendermint/linux/tendermint node1:/home/DecentralizedRedis/tendermint
    ;;
  startAllDocker)
    docker kill $(docker ps -q)
    docker rm $(docker ps -a -q)

    docker run -tid --name node0 --privileged=true -p 30002:30001 base_test_env /sbin/init
    docker run -tid --name node1 --privileged=true  base_test_env /sbin/init
    docker run -tid --name node2 --privileged=true  base_test_env /sbin/init
    docker run -tid --name node3 --privileged=true  base_test_env /sbin/init

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis node0:/home/DecentralizedRedis
    docker cp /Users/chenzhou/tendermint/linux/tendermint node0:/home/DecentralizedRedis/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis node1:/home/DecentralizedRedis
    docker cp /Users/chenzhou/tendermint/linux/tendermint node1:/home/DecentralizedRedis/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis node2:/home/DecentralizedRedis
    docker cp /Users/chenzhou/tendermint/linux/tendermint node2:/home/DecentralizedRedis/tendermint

    docker cp -a /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis node3:/home/DecentralizedRedis
    docker cp /Users/chenzhou/tendermint/linux/tendermint node3:/home/DecentralizedRedis/tendermint
    ;;
esac
exit $RETVAL