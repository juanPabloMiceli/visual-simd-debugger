FROM golang:latest as builder

RUN apt-get update
RUN apt-get install nasm 

WORKDIR /go/src/app
COPY go.mod .

COPY go.sum .

COPY . .

RUN go mod download


RUN cd main/ && go build -o http

CMD ["./main/http"]


