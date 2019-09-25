FROM golang:1.10-alpine

COPY proc_id.go.patch /proc_id.go

RUN cd $(go env GOROOT)/src/runtime \
    && mv /proc_id.go . \
    && go install