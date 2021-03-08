#Build from go/src/

#To build:
#sudo docker build ./ -t visual-simd-back

#To run:
#sudo docker run -d --rm --name v-simd-back --network host visual-simd-back

FROM golang:latest as builder

#Update system and install nasm
WORKDIR /
RUN apt-get update
RUN apt-get install nasm

#Need this in order to use minijail
RUN apt-get install -y libcap-dev 
RUN apt-get install -y git

#Need this in order to use microjail
RUN apt-get install -y libseccomp-dev

#Downloading minijail

RUN git clone https://android.googlesource.com/platform/external/minijail ./minijail

#Installing minijail
WORKDIR /minijail
RUN make LIBDIR=/lib64

RUN cp libminijail.so /lib64/
RUN cp libminijailpreload.so /lib64/
RUN cp minijail0 /usr/bin/

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


