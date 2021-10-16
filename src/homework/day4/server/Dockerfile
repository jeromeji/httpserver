FROM golang:1.16-alpine AS build
WORKDIR /go/src/project/
COPY server.go /go/src/project/
RUN export GOPROXY=https://goproxy.cn && go mod init && go mod tidy
RUN go build -o /bin/httpserver

FROM scratch
LABEL multi.name="httpserver" multi.version="v1.0" author="jeromeji"
COPY --from=build /bin/httpserver /bin/httpserver
EXPOSE 80
ENTRYPOINT ["/bin/httpserver"]