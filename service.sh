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
  start)
    start
    ;;
  stop)
    stop
    ;;
  reinstall)
    stop
    reinstall
    ;;
  restart)
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
esac
exit $RETVAL