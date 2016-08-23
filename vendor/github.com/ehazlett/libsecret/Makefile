CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
COMMIT=`git rev-parse --short HEAD`
APP=secret
REPO?=ehazlett/libsecret
TAG?=latest
export GO15VENDOREXPERIMENT=1

all: build

add-deps:
	@godep save ./...

build:
	@cd cmd/$(APP) && godep go build -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT)" .

build-static:
	@cd cmd/$(APP) && godep go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT)" .

test: build
	@godep go test -v ./...

clean:
	@rm cmd/$(APP)/$(APP)

.PHONY: add-deps build build-static image release test clean
