# create cgroup psi image, please replace the tag and the image name
rm cgroup-psi-sc
TAG=$1
# go build -v ./...
docker build -f ./Dockerfile . -t 10.119.46.41:30003/library/psi-perf-ds:$TAG
docker push 10.119.46.41:30003/library/psi-perf-ds:$TAG