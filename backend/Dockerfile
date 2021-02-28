#Build from go/src/

#To build:
#sudo docker build ./ -t visual-simd-back

#To run:
#sudo docker run -d --rm --name v-simd-back --network host visual-simd-back
FROM golang:latest as builder

#Update system and install nasm
RUN apt-get update
RUN apt-get install nasm 


WORKDIR /go/src
#Copy all dependencies
COPY ./gitlab.com/juampi_miceli/visual-simd-debugger/backend ./gitlab.com/juampi_miceli/visual-simd-debugger/backend


WORKDIR /go/src/gitlab.com/juampi_miceli/visual-simd-debugger/backend/main
#Build go exectutable
RUN go build -o http

#Run go executable
CMD ["./http"]


