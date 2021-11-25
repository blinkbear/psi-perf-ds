#FROM golang:1.17
#
#WORKDIR /go/src/app
#COPY . .
#
#RUN go env -w  GOPROXY=https://goproxy.cn,direct
#RUN go get -d -v ./...
#RUN go build -v ./...
#RUN mv ./cgroup-sc /usr/bin/
#CMD ["cgroup-sc"]
FROM ubuntu
COPY cgroup-psi-sc /usr/bin/cgroup-psi-sc
CMD ["cgroup-psi-sc"]