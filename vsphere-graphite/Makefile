GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
RELEASE_DIR=releases
SRC_FILES=$(wildcard *.go)
MUSL_BUILD_FLAGS=-ldflags '-linkmode external -s -w -extldflags "-static"' -a
BUILD_FLAGS=-ldflags -s -a
MUSL_CC=musl-gcc
MUSL_CCGLAGS="-static"

deps:
	go get github.com/takama/daemon
	go get golang.org/x/net/context
	go get github.com/vmware/govmomi
	go get github.com/marpaia/graphite-golang
	go get github.com/influxdata/influxdb/client/v2
	go get github.com/pquerna/ffjson/fflib/v1
	go get code.cloudfoundry.org/bytefmt

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

dist-windows-amd64:
	@$(MAKE) dist GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-linux-amd64:
	CC=$(MUSL_CC) CCGLAGS=$(MUSL_CCGLAGS) go build $(MUSL_BUILD_FLAGS) -o $(RELEASE_DIR)/linux/amd64//vsphere-graphite .
	upx -qq --best $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite
	cp vsphere-graphite-example.json $(RELEASE_DIR)/linux/amd64/vsphere-graphite.json

dist-linux-amd64:
	@$(MAKE) dist GOOS=linux GOARCH=amd64

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

dist-darwin-amd64:
	@$(MAKE) dist GOOS=darwin GOARCH=amd64
    
build-linux-arm:
	@$(MAKE) build GOOS=linux GOARCH=arm GOARM=5

dist-linux-arm:
	@$(MAKE) dist GOOS=linux GOARCH=arm GOARM=5

docker: $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite
	cp $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/* docker/
	mkdir -p docker/etc
	cp vsphere-graphite-example.json docker/etc/vsphere-graphite.json
	docker build -f docker/Dockerfile -t cblomart/$(PREFIX)vsphere-graphite docker

docker-linux-amd64:
	@$(MAKE) docker GOOS=linux GOARCH=amd64

docker-linux-arm:
	@$(MAKE) docker GOOS=linux GOARCH=arm PREFIX=rpi-

docker-darwin-amd64: ;

docker-windows-amd64: ;

checks:
	go get honnef.co/go/tools/cmd/gosimple
	go get github.com/golang/lint/golint
	go get github.com/gordonklaus/ineffassign
	go get github.com/GoASTScanner/gas
	gosimple ./...
	gofmt -s -d .
	go vet ./...
	golint ./...
	ineffassign ./
	gas ./...
	go tool vet ./..

$(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite$(SUFFIX): $(SRC_FILES)
	go build $(BUILD_FLAGS) -o $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite$(SUFFIX) .
	upx -qq --best $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite$(SUFFIX)
	cp vsphere-graphite-example.json $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite.json

$(RELEASE_DIR)/vsphere-graphite_$(GOOS)_$(GOARCH).tgz: $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite$(SUFFIX)
	cd $(RELEASE_DIR)/$(GOOS)/$(GOARCH); tar czf /tmp/vsphere-graphite_$(GOOS)_$(GOARCH).tgz ./vsphere-graphite$(SUFFIX) ./vsphere-graphite.json

dist: $(RELEASE_DIR)/vsphere-graphite_$(GOOS)_$(GOARCH).tgz

build: $(RELEASE_DIR)/$(GOOS)/$(GOARCH)/vsphere-graphite$(SUFFIX)

clean:
	rm -rf $(RELEASE_DIR)
	
all:
	@$(MAKE) dist-windows-amd64 
	@$(MAKE) dist-linux-amd64
	@$(MAKE) dist-darwin-amd64
	@$(MAKE) dist-linux-arm
