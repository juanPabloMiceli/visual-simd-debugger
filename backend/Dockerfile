#Build from go/src/

FROM golang:latest as builder

#Update system and install nasm
RUN apt-get update
RUN apt-get install nasm 


WORKDIR /go/src
#Copy all dependencies
COPY ./github.com/zchee/kube-timeleap/pkg/ptrace ./github.com/zchee/kube-timeleap/pkg/ptrace
COPY ./gitlab.com/juampi_miceli/visual-simd-debugger/backend ./gitlab.com/juampi_miceli/visual-simd-debugger/backend
COPY ./golang.org/x/sys/unix ./golang.org/x/sys/unix
COPY ./golang.org/x/sys/internal/unsafeheader ./golang.org/x/sys/internal/unsafeheader
COPY ./github.com/go-logr/logr ./github.com/go-logr/logr
COPY ./sigs.k8s.io/controller-runtime/pkg/log ./sigs.k8s.io/controller-runtime/pkg/log


WORKDIR /go/src/gitlab.com/juampi_miceli/visual-simd-debugger/backend/main
#Build go exectutable
RUN go build -o http

#Run go executable
CMD ["./http"]


