#!/bin/bash

function start(){
    rm -rf tendermint.sock
    go build
    nohup ./DecentralizedRedis > db.log 2>&1 &
    nohup ./tendermint --home=/tmp/example node --proxy_app=unix://tendermint.sock > tendermint.log 2>&1 &
}

function stop() {
    PID=$(ps -ef|grep DecentralizedRedis|grep -v grep|awk '{print $2}')
    if [ -z $PID ]; then
      echo "process DecentralizedRedis not exist"
      exit
    else
      echo "process id: $PID"
      kill -9 ${PID}
      echo "process DecentralizedRedis killed"
    fi

    PID=$(ps -ef|grep tendermint|grep -v grep|awk '{print $2}')
    if [ -z $PID ]; then
      echo "process tendermint not exist"
      exit
    else
      echo "process id: $PID"
      kill -9 ${PID}
      echo "process tendermint killed"
    fi
}

function reinstall() {
    rm -rf /tmp/example
    ./tendermint init --home=/tmp/example
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  reinstall)
    reinstall
    ;;
  restart)
    stop
    start
esac
exit $RETVAL