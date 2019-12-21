OUTPUT:=$(output)
OUTPUT_BINS:=$(OUTPUT)/bin
OUTPUT_CONFS:=$(OUTPUT)/conf
OUTPUT_DATA:=$(OUTPUT)/chain
OUTPUT_LOG=$(OUTPUT)/log

default:
	rm -rf $(OUTPUT)
	mkdir $(OUTPUT)
	mkdir $(OUTPUT_BINS)
	mkdir $(OUTPUT_CONFS)
	mkdir $(OUTPUT_LOG)
	go build -o $(OUTPUT_BINS)/redimint
	tendermint init --home=$(OUTPUT_DATA)
	cp -f ./conf/tendermint/config.toml $(OUTPUT_DATA)/config/config.toml
	cp -f ./conf/redis/redis.conf $(OUTPUT_CONFS)/redis.conf
	cp -f ./conf/configuration.yaml $(OUTPUT_CONFS)/configuration.yaml

build:
	rm -rf $(OUTPUT)
	mkdir $(OUTPUT)
	mkdir $(OUTPUT_BINS)
	mkdir $(OUTPUT_CONFS)
	mkdir $(OUTPUT_LOG)
	go build -o $(OUTPUT_BINS)/redimint
	tendermint init --home=$(OUTPUT_DATA)
	cp -f ./conf/tendermint/config.toml $(OUTPUT_DATA)/config/config.toml
	cp -f ./conf/redis/redis.conf $(OUTPUT_CONFS)/redis.conf
	cp -f ./conf/configuration.yaml $(OUTPUT_CONFS)/configuration.yaml



