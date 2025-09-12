# $ go env -w GO111MODULE=on
#$ go env -w GOPROXY=https://goproxy.cn,direct
MODULE = $(shell go list -m)
CONFIG_FILE ?= ./config/local.yml

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: build
build:
	go build  -o ./runtime/web ./bin/web/server.go
run:
	./runtime/web --conf=runtime/conf/base.json