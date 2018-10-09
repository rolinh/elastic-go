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

deps-dev: deps
	go get -u -v github.com/golang/lint/golint

clean:
	rm -f ${EXEC}

.PHONY: build install test cover check deps-dev clean

