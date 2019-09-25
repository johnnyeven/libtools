FROM registry.profzone.net:5000/profzone/golang:tensorflow

RUN apt-get install -y curl wget unzip \
    && go get -u github.com/kardianos/govendor

COPY ./libtools /go/src/github.com/johnnyeven/libtools
RUN cd /go/src/github.com/johnnyeven/libtools && go install

COPY ./profzone /go/src/github.com/johnnyeven/profzone
RUN cd /go/src/github.com/johnnyeven/profzone && go install