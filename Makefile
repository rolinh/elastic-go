EXEC = elastic
PKG  = github.com/Rolinh/elastic-go

all: check test build

build: clean
	go build -o ${EXEC} ${PKG}

install: clean
	go install ${PKG}

test:
	go test -v ${PKG}/...

cover:
	go test -cover ${PKG}/...

check:
	go vet ${PKG}/...
	golint

deps:
	go get -u github.com/codegangsta/cli
	go get -u github.com/gilliek/go-xterm256/xterm256
	go get -u github.com/hokaccha/go-prettyjson

deps-dev: deps
	go get -u -v github.com/golang/lint/golint

clean:
	rm -f ${EXEC}

.PHONY: build install test cover check deps deps-dev clean

