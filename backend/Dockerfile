#Build from go/src/

#To build:
#sudo docker build ./ -t visual-simd-backend

#To run:
#sudo docker run -d --rm --name v-simd-back --network host visual-simd-back

FROM golang:1.15.10 as builder

#Update system and install nasm
#Then install git and libcap-dev because minijail needs it
#Then install libseccomp-dev for using microjail
#Then download minijail
WORKDIR /
RUN \
    apt-get update && \
    apt-get install -y apt-utils && \
    apt-get install nasm && \
    apt-get install -y libcap-dev && \ 
    apt-get install -y git && \
    apt-get install -y libseccomp-dev && \
    git clone https://android.googlesource.com/platform/external/minijail ./minijail && \
    mkdir clients && \

    #Installing minijail
    WORKDIR /minijail

RUN \
    make LIBDIR=/lib64 && \
    cp libminijail.so /lib64/ && \
    cp libminijailpreload.so /lib64/ && \
    cp minijail0 /usr/bin/

WORKDIR /go/src
#Copy all dependencies
COPY ./github.com/juanPabloMiceli/visual-simd-debugger/backend ./github.com/juanPabloMiceli/visual-simd-debugger/backend

#Installing microjail
WORKDIR /go/src/github.com/juanPabloMiceli/visual-simd-debugger/backend/microjail
RUN make docker-install


WORKDIR /go/src/github.com/juanPabloMiceli/visual-simd-debugger/backend/main
#Build go exectutable
RUN go build -o http

#Run go executable
CMD ["./http"]


