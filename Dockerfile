FROM golang:1.12-alpine3.9 AS builder
RUN apk update && apk add --no-cache git
WORKDIR /workdir
COPY . .
RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o elastic elastic.go 

FROM alpine:3.9
RUN apk update && \
    apk add --no-cache ca-certificates && \
    addgroup -S elastic && \
    adduser -S elastic -G elastic
USER elastic:elastic
COPY --from=builder /workdir/elastic /usr/local/bin/elastic
ENTRYPOINT ["/usr/local/bin/elastic"]
CMD ["--help"]
