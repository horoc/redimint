PREFIX?=/usr/local


GIT_REV:=git-$(shell git rev-parse --short HEAD)
GIT_TAG:=$(shell git tag --contains HEAD)
VERSION:=$(if $(GIT_TAG),$(GIT_TAG),$(GIT_REV))

GO_OS:=$(shell go env GOOS)
GO_ARCH:=$(shell go env GOARCH)
GO_ARM:=$(shell go env GOARM)
GO_FLAGS?=-ldflags "-X 'github.com/chenzhou9513/redimint/cmd.Revision=$(GIT_REV)' -X 'github.com/chenzhou9513/redimint/cmd.Version=$(VERSION)'"


ifndef $(output)
	OUTPUT:=./redimint_home
	OUTPUT_BINS:=./redimint_home/bin
	OUTPUT_CONFS:=./redimint_home/conf
	OUTPUT_DATA:=./redimint_home/chain
	OUTPUT_LOG=./redimint_home/log
else
	OUTPUT:=$(output)
	OUTPUT_BINS:=$(OUTPUT)/bin
	OUTPUT_CONFS:=$(OUTPUT)/conf
	OUTPUT_DATA:=$(OUTPUT)/chain
	OUTPUT_LOG=$(OUTPUT)/log
endif


.PHONY: default
default:
	@rm -rf $(OUTPUT)
	@mkdir $(OUTPUT)
	@mkdir $(OUTPUT_BINS)
	@mkdir $(OUTPUT_CONFS)
	@mkdir $(OUTPUT_LOG)
	@go build -o $(OUTPUT_BINS)/redimint ${GO_FLAGS}
	@tendermint init --home=$(OUTPUT_DATA)
	@cp -f ./conf/tendermint/config.toml $(OUTPUT_DATA)/config/config.toml
	@cp -f ./conf/redis/redis.conf $(OUTPUT_CONFS)/redis.conf
	@cp -f ./conf/configuration.yaml $(OUTPUT_CONFS)/configuration.yaml

.PHONY: rebuild
rebuild:
	@go build -o $(OUTPUT_BINS)/redimint ${GO_FLAGS}

.PHONY: clean
clean:
	@rm -rf $(OUTPUT)

.PHONY: reinit
reinit:
	@tendermint init --home=$(OUTPUT_DATA)

.PHONY: reinstall
reinstall:
	@rm -rf $(OUTPUT)
	@mkdir $(OUTPUT)
	@mkdir $(OUTPUT_BINS)
	@mkdir $(OUTPUT_CONFS)
	@mkdir $(OUTPUT_LOG)
	@go build -o $(OUTPUT_BINS)/redimint ${GO_FLAGS}
	@tendermint init --home=$(OUTPUT_DATA)
	@cp -f ./conf/tendermint/config.toml $(OUTPUT_DATA)/config/config.toml
	@cp -f ./conf/redis/redis.conf $(OUTPUT_CONFS)/redis.conf
	@cp -f ./conf/configuration.yaml $(OUTPUT_CONFS)/configuration.yaml