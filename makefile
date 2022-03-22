INFO = github.com/mmbros/quotes/cmd

VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
# GOXVER := $(shell go version | awk '{print $$3}')
GOXVER := $(shell go version)
GITCOMMIT := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date '+%F %T %z')
OSARCH := $(shell uname -s -m)

TIMESTAMP := $(shell date +%Y%m%dT%H%M%S)
# go build -v -ldflags="-X 'main.Version=v1.0.0' -X 'app/build.User=$(id -u -n)' -X 'app/build.Time=$(date)'"


COMMON_LDFLAGS = -X '${INFO}.BuildTime=${BUILDTIME}' -X '${INFO}.GitCommit=${GITCOMMIT}' -X '${INFO}.GoVersion=${GOXVER}' -X '${INFO}.OsArch=${OSARCH}'
PROD_LDFLAGS = -ldflags "-X '${INFO}.Version=${VERSION}' ${COMMON_LDFLAGS}"
DEV_LDFLAGS = -ldflags "-X '${INFO}.Version=dev-${TIMESTAMP}' ${COMMON_LDFLAGS}"

BIN=./bin/quotes

all: prod


dev:
	go build ${DEV_LDFLAGS} -o ${BIN} *.go

prod:
	go build ${PROD_LDFLAGS} -o ${BIN} *.go
