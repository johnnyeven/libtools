FROM registry.profzone.net:5000/profzone/golang:latest

ENV CGO_ENABLED 0

RUN sed -i "s|http://dl-cdn.alpinelinux.org|http://mirrors.aliyun.com|g" /etc/apk/repositories

RUN apk add --no-cache curl git openssh wget unzip \
    && go get -u github.com/kardianos/govendor

COPY libtools /go/src/github.com/johnnyeven/libtools
RUN cd /go/src/github.com/johnnyeven/libtools && go install

COPY profzone /go/src/github.com/johnnyeven/profzone
RUN cd /go/src/github.com/johnnyeven/profzone && go install