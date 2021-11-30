FROM golang:1.16-alpine
WORKDIR /go/src/project/
COPY server.go /go/src/project/
RUN export GOPROXY=https://goproxy.cn && go mod init && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -ldflags '-w -s' -o server

FROM alpine:3.13 
COPY --from=0 /go/src/project/server /
ENTRYPOINT ["/server"]
