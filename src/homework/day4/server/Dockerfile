FROM golang:1.16-alpine AS build
RUN apk add --no-cache git && go get github.com/golang/glog
WORKDIR /go/src/project/
COPY . /go/src/project/
RUN go build -o /bin/httpserver

FROM scratch
LABEL multi.name="httpserver" multi.version="v1.0" author="jeromeji"
COPY --from=build /bin/httpserver /bin/httpserver
EXPOSE 80
ENTRYPOINT /httpserver