FROM golang:onbuild

RUN ln -vs /go/bin/app /usr/local/bin/elastic

ENTRYPOINT [ "elastic" ]
