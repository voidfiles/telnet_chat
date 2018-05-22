PROJECT         :=telnet_chat
CW              :=$(shell pwd)
GOFILES         :=$(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES      :=$(shell go list ./... | grep -v /vendor/| grep -v /checkers)
OS              := $(shell go env GOOS)
ARCH            := $(shell go env GOARCH)
CACHE           :=download-cache

BIN             := $(CW)/bin

GITHASH         :=$(shell git rev-parse --short HEAD)
GITBRANCH       :=$(shell git rev-parse --abbrev-ref HEAD)
BUILDDATE      	:=$(shell date -u +%Y%m%d%H%M)
GO_LDFLAGS		  ?= -s -w
GO_BUILD_FLAGS  :=-ldflags "${GOLDFLAGS} -X main.BuildVersion=${GITHASH} -X main.GitHash=${GITHASH} -X main.GitBranch=${GITBRANCH} -X main.BuildDate=${BUILDDATE}"
WORKDIR 	      :=$(CW)/_work


build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(WORKDIR)/$(PROJECT)_linux_amd64 $(GO_BUILD_FLAGS)

build:
	CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS)


dependencies:
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/alecthomas/gometalinter
	go get github.com/golang/dep/cmd/dep
	go get github.com/stretchr/testify
	go get github.com/jstemmer/go-junit-report
	dep ensure
	gometalinter --install

lint:
	echo "metalinter..."
	gometalinter --enable=goimports --enable=unparam --enable=unused --disable=golint --disable=govet $(GOPACKAGES)
	echo "megacheck..."
	megacheck $(GOPACKAGES)
	echo "golint..."
	golint -set_exit_status $(GOPACKAGES)
	echo "go vet..."
	go vet --all $(GOPACKAGES)

init: dependencies

clean:
	rm -fR $(WORKDIR)

test:
	CGO_ENABLED=0 go test $(GOPACKAGES)

coverage:
	CGO_ENABLED=0 go test -cover $(GOPACKAGES)

test-race:
	CGO_ENABLED=1 go test -race $(GOPACKAGES)

docs:
	godoc -http=":6060" && \
	open "http://localhost:6060/pkg/github.com/voidfiles/telnet_chat/"

run:
	_work/telnet_chat_darwin_amd64 -configpath ./examples/config.toml

zipit:
	zip -r /tmp/archive.zip ../
