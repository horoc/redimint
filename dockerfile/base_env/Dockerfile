FROM centos:7.6.1810
RUN  yum -y install epel-release && yum -y update && yum -y install redis && yum -y install golang
RUN mkdir /home/DecentralizedRedis
COPY /Users/chenzhou/go/src/github.com/chenzhou9513/DecentralizedRedis/ /home/DecentralizedRedis/
RUN rm /home/DecentralizedRedis/tendermint
RUN rm /home/DecentralizedRedis/tendermint.sock
COPY /Users/chenzhou/tendermint/linux/tendermint /home/DecentralizedRedis/tendermint
WORKDIR /home/DecentralizedRedis
RUN export GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN go build
ENTRYPOINT ["/usr/sbin/init"]