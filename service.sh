#!/bin/bash

function start(){
    rm -rf tendermint.sock
    go build
    nohup ./DecentralizedRedis > db.log 2>&1 &
    nohup ./tendermint --home=./chain node --proxy_app=unix://tendermint.sock > tendermint.log 2>&1 &
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
    ./tendermint init --home=./chain
    cp -f ./conf/tendermint/config.toml ./chain/config/config.toml
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
    cp -fr ./conf/testnet/$2/* ./chain
    ;;
esac
exit $RETVAL