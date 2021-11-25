# create cgroup psi image, please replace the tag and the image name
rm cgroup-psi-sc
TAG=$1
go build -v ./...
docker build -f ./Dockerfile . -t t.harbor.siat.ac.cn:100/library/cgroup-sc:$TAG
docker push t.harbor.siat.ac.cn:100/library/cgroup-sc:$TAG