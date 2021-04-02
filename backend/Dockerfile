#Build from go/src/

#To build:
#sudo docker build ./ -t visual-simd-back

#To run:
#sudo docker run -d --rm --name v-simd-back --network host visual-simd-back

FROM golang:latest as builder

#Update system and install nasm
#Then install git and libcap-dev because minijail needs it
#Then install libseccomp-dev for using microjail
#Then download minijail
WORKDIR /
RUN \
apt-get update && \
apt-get install nasm && \
apt-get install -y libcap-dev && \ 
apt-get install -y git && \
apt-get install -y libseccomp-dev && \
git clone https://android.googlesource.com/platform/external/minijail ./minijail

#Installing minijail
WORKDIR /minijail

RUN \
make LIBDIR=/lib64 && \
cp libminijail.so /lib64/ && \
cp libminijailpreload.so /lib64/ && \
cp minijail0 /usr/bin/

WORKDIR /go/src
#Copy all dependencies
COPY ./gitlab.com/juampi_miceli/visual-simd-debugger/backend ./gitlab.com/juampi_miceli/visual-simd-debugger/backend

#Installing microjail
WORKDIR /go/src/gitlab.com/juampi_miceli/visual-simd-debugger/backend/microjail
RUN make docker-install


WORKDIR /go/src/gitlab.com/juampi_miceli/visual-simd-debugger/backend/main
#Build go exectutable
RUN go build -o http

#Run go executable
CMD ["./http"]


