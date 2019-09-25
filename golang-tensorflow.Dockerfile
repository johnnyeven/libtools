FROM tensorflow/tensorflow

RUN apt-get install -y curl git

ENV GOLANG_VERSION 1.13
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz && \
    tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR "/go"

COPY proc_id.go.patch /proc_id.go

RUN cd $(go env GOROOT)/src/runtime \
    && mv /proc_id.go . \
    && go install

ENV TENSORFLOW_LIB_GZIP libtensorflow-cpu-linux-x86_64-1.13.1.tar.gz
ENV TARGET_DIRECTORY /usr/local
RUN  curl -fsSL "https://storage.googleapis.com/tensorflow/libtensorflow/$TENSORFLOW_LIB_GZIP" -o $TENSORFLOW_LIB_GZIP && \
     tar -C $TARGET_DIRECTORY -xzf $TENSORFLOW_LIB_GZIP && \
     rm -Rf $TENSORFLOW_LIB_GZIP

ENV LD_LIBRARY_PATH=$TARGET_DIRECTORY/lib
ENV LIBRARY_PATH=$TARGET_DIRECTORY/lib

#RUN go get -v github.com/tensorflow/tensorflow/tensorflow/go

#RUN mkdir -p $GOPATH/src/github.com/tensorflow/ && cd $GOPATH/src/github.com/tensorflow/ && \
#	git clone https://github.com/tensorflow/tensorflow.git && \
#	cd tensorflow/go && \
#	go install

ADD ./tensorflow.tar.gz /go/src/github.com/
RUN cd /go/src/github.com/tensorflow/tensorflow/tensorflow/go && go install