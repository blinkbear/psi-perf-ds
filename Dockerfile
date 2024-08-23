FROM golang:1.20 AS builder

WORKDIR /go/src/app
COPY . .

RUN go env -w  GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go build -v ./...
RUN mv ./psi-perf-ds /root/

FROM ubuntu
COPY --from=builder /root/psi-perf-ds /usr/bin/psi-perf-ds
RUN echo 'deb http://security.ubuntu.com/ubuntu jammy-security main' >> /etc/apt/sources.list
RUN apt update && apt install -y libc6
ENTRYPOINT ["psi-perf-ds"]