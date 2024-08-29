FROM golang:1.20 AS builder

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go env -w  GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...

COPY . .
RUN go build -v ./...
RUN mv ./psi-perf-ds /root/

FROM ubuntu
RUN echo 'deb http://security.ubuntu.com/ubuntu jammy-security main' >> /etc/apt/sources.list
RUN apt update && apt install -y libc6
COPY --from=builder /root/psi-perf-ds /usr/bin/psi-perf-ds
ENTRYPOINT ["psi-perf-ds"]