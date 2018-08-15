VERSION=1.53
GOBIN=go
GOGET=$(GOBIN) get
GOINSTALL=$(GOBIN) install
PATH:=/usr/local/bin:$(GOPATH)/bin:$(PATH)
BINARY_NAME=tokenbalance
XGO=GOPATH=$(GOPATH) $(GOPATH)/bin/xgo -go 1.10.x --dest=build
BUILDVERSION=-ldflags="-X main.VERSION=$(VERSION)"

all: deps install clean

release: deps build-all compress

build: deps
	$(GOBIN) build -ldflags="-X main.VERSION=$(VERSION)" -o $(BINARY_NAME) -v ./cmd

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	$(GOPATH)/bin/$(BINARY_NAME) version

run: build
	./$(BINARY_NAME)

test:
	$(GOBIN) test -p 1 -ldflags="-X main.VERSION=$(VERSION)" -coverprofile=coverage.out -v ./...

coverage:
	$(GOPATH)/bin/goveralls -coverprofile=coverage.out -service=travis -repotoken $(COVERALLS)

build-all: clean
	mkdir build
	$(XGO) $(BUILDVERSION) --targets=linux/amd64 ./cmd
	$(XGO) $(BUILDVERSION) --targets=linux/386 ./cmd
	$(XGO) $(BUILDVERSION) --targets=linux/arm-7 ./cmd
	$(XGO) $(BUILDVERSION) --targets=linux/arm64 ./cmd
	$(XGO) $(BUILDVERSION) --targets=darwin/amd64 ./cmd
	$(XGO) $(BUILDVERSION) --targets=darwin/386 ./cmd
	$(XGO) $(BUILDVERSION) --targets=windows-6.0/amd64 ./cmd
	$(XGO) $(BUILDVERSION) --targets=windows-6.0/386 ./cmd
	$(XGO) --targets=linux/amd64 -ldflags="-X main.VERSION=$VERSION -linkmode external -extldflags -static" -out alpine ./cmd

docker:
	$(DOCKER) build -t hunterlong/tokenbalance:latest ./cmd

deps:
	$(GOGET) github.com/stretchr/testify/assert
	$(GOGET) golang.org/x/tools/cmd/cover
	$(GOGET) github.com/mattn/goveralls
	$(GOINSTALL) github.com/mattn/goveralls
	$(GOGET) github.com/rendon/testcli
	$(GOGET) github.com/karalabe/xgo
	$(GOGET) github.com/ethereum/go-ethereum
	$(GOGET) -d ./...

clean:
	rm -rf build
	rm -f coverage.out

tag:
	git tag "v$(VERSION)" --force

compress:
	mv build/alpine-linux-amd64 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-linux-alpine.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)
	mv build/cmd-darwin-10.6-amd64 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-osx-x64.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)
	mv build/cmd-darwin-10.6-386 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-osx-x32.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)
	mv build/cmd-linux-amd64 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-linux-x64.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)
	mv build/cmd-linux-386 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-linux-x32.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)
	mv build/cmd-windows-6.0-amd64.exe build/$(BINARY_NAME).exe
	zip build/$(BINARY_NAME)-windows-x64.zip build/$(BINARY_NAME).exe  && rm -f build/$(BINARY_NAME).exe
	mv build/cmd-windows-6.0-386.exe build/$(BINARY_NAME).exe
	zip build/$(BINARY_NAME)-windows-x32.zip build/$(BINARY_NAME).exe  && rm -f build/$(BINARY_NAME).exe
	mv build/cmd-linux-arm-7 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-linux-arm7.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)
	mv build/cmd-linux-arm64 build/$(BINARY_NAME)
	tar -czvf build/$(BINARY_NAME)-linux-arm64.tar.gz build/$(BINARY_NAME) && rm -f build/$(BINARY_NAME)

.PHONY: build test