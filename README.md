<div id="top"></div>

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[comment]: <> ([![Contributors][contributors-shield]][contributors-url])

[comment]: <> ([![Forks][forks-shield]][forks-url])

[comment]: <> ([![Stargazers][stars-shield]][stars-url])

[comment]: <> ([![Issues][issues-shield]][issues-url])

[comment]: <> ([![MIT License][license-shield]][license-url])

[comment]: <> ([![LinkedIn][linkedin-shield]][linkedin-url])



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/github_username/repo_name">
    <img src="https://golang.org/lib/godoc/images/go-logo-blue.svg" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">CgroupV2 PSI Sidecar</h3>

  <p align="center">

[comment]: <> (    CgroupV2 PSI Sidecar can be deployed on any kubernetes pod with access to cgroupv2 PSI metrics.)
  </p>
</div>



[comment]: <> (<!-- TABLE OF CONTENTS -->)

[comment]: <> (  <summary>Table of Contents</summary>)

[comment]: <> (  <ol>)

[comment]: <> (    <li>)

[comment]: <> (      <a href="#about-the-project">About The Project</a>)

[comment]: <> (      <ul>)

[comment]: <> (        <li><a href="#built-with">Built With</a></li>)

[comment]: <> (      </ul>)

[comment]: <> (    </li>)

[comment]: <> (    <li>)

[comment]: <> (      <a href="#getting-started">Getting Started</a>)

[comment]: <> (      <ul>)

[comment]: <> (        <li><a href="#prerequisites">Prerequisites</a></li>)

[comment]: <> (        <li><a href="#installation">Installation</a></li>)

[comment]: <> (      </ul>)

[comment]: <> (    </li>)

[comment]: <> (    <li><a href="#usage">Usage</a></li>)

[comment]: <> (    <li><a href="#roadmap">Roadmap</a></li>)

[comment]: <> (    <li><a href="#contributing">Contributing</a></li>)

[comment]: <> (    <li><a href="#license">License</a></li>)

[comment]: <> (    <li><a href="#contact">Contact</a></li>)

[comment]: <> (    <li><a href="#acknowledgments">Acknowledgments</a></li>)

[comment]: <> (  </ol>)



<!-- ABOUT THE PROJECT -->
## About

This is a docker container that can be deployed as a sidecar on any kubernetes pod to monitor PSI metrics. 


### Built With

* [Go Lang](https://golang.org/)
* [Gorilla Mux](https://github.com/gorilla/mux)
* [Prometheus](https://prometheus.io/)


<!-- GETTING STARTED -->
## Getting Started

To deploy a sidecar follow these steps.

### Prerequisites

#### Minimum versions:
* `Docker 20.10`
* `Linux 5.2`
* `Kubernetes 1.17`

The host machine for all the nodes on the cluster must be using cgroupv2.

#### Check CgroupV2 Availability
Ensure that your machine has cgroupv2 available:

```sh
$ grep cgroup /proc/filesystems
nodev	cgroup
nodev	cgroup2
```

Just because you have cgroupv2 it doesn't mean you are using it. 
Check that the unified cgroup is enabled by checking the hierarchy.
```shell
$ ll /sys/fs/cgroup/
total 0
dr-xr-xr-x   5 root root 0 Oct 31 14:52 ./
drwxr-xr-x  10 root root 0 Oct 31 14:52 ../
-r--r--r--   1 root root 0 Nov  1 08:45 cgroup.controllers
-rw-r--r--   1 root root 0 Nov  1 08:45 cgroup.max.depth
-rw-r--r--   1 root root 0 Nov  1 08:45 cgroup.max.descendants
-rw-r--r--   1 root root 0 Nov  1 08:45 cgroup.procs
-r--r--r--   1 root root 0 Nov  1 08:45 cgroup.stat
-rw-r--r--   1 root root 0 Oct 31 14:52 cgroup.subtree_control
-rw-r--r--   1 root root 0 Nov  1 08:45 cgroup.threads
-rw-r--r--   1 root root 0 Nov  1 08:45 cpu.pressure
-r--r--r--   1 root root 0 Nov  1 08:45 cpuset.cpus.effective
-r--r--r--   1 root root 0 Nov  1 08:45 cpuset.mems.effective
drwxr-xr-x   2 root root 0 Nov  1 08:45 init.scope/
-rw-r--r--   1 root root 0 Nov  1 08:45 io.cost.model
-rw-r--r--   1 root root 0 Nov  1 08:45 io.cost.qos
-rw-r--r--   1 root root 0 Nov  1 08:45 io.pressure
-rw-r--r--   1 root root 0 Nov  1 08:45 memory.pressure
drwxr-xr-x 106 root root 0 Nov  1 08:45 system.slice/
drwxr-xr-x   3 root root 0 Oct 31 14:52 user.slice/
```
_Note the slice dirs._

If you have cgroupv2 but it isn't enabled the above structure will be available in `/sys/fs/cgroup/unified`.

#### Enable cgroupv2
Edit `/etc/default/grub` and add `systemd.unified_cgroup_hierarchy=1` to `GRUB_CMDLINE_LINUX`

Run `sudo update-grub` and reboot the system.

_If cgroupv2 is not available on the system you will have to update the kernel version to meet the prerequisites above._

### Build Image
There are two docker files one for regular deployment and the other for debugging.
If you want to run the server locally without a container/kubernetes deployment edit `sidecar_pid_lookup.go` to resolve the systems cgroup dir.

#### Regular image
1. `docker build -f ./Dockerfile . -t evankrul/cgroup-sc:v.1.2`
2. `docker push evankrul/cgroup-sc:v.1.2`
#### Debug image
1. `docker build -f ./Dockerfile.debug . -t evankrul/cgroup-sc:v.1.2-debug`
2. `docker push evankrul/cgroup-sc:v.1.2-dubug`

#### Port
Set `PORT` env var to specify the metrics port.

<!-- USAGE EXAMPLES -->
## Usage
Assuming all the prerequisites have been met and image built and pushed to your docker repository follow these steps to deploy the sidecar.

In this section I will refer to the monitoring container as the **sidecar** and the container being monitored as the **host** container. 
The sidecar makes use of the `shareProcessNamespace` option to access the host cgroup metrics.
The sidecar has access to process dirs in `/proc`. The sidecar finds the pid dir of the host by searching the dirs in `/proc`.

For each dir the sidecar looks at the contents of `/proc/{id}/root/etc/pid_flag` and checks that it exists and matches the contents of `/etc/pid_flag_sc`.
If a match is found then this is the host container. The pid_flag and pid_flag_sc are mounted in the deployment configuration as a ConfigMap using a VolumeMount.

The service is used to expose the sidecar webserver where the metrics are hosted.
If you are not using some kind of service mesh make sure your Prometheus deployment is on the same namespace as your sidecar deployment.
Then just point Prometheus to the `/metrics` endpoint of your pod on the metrics port.

```yaml
- job_name: 'cgroup_monitor_sc'
        scrape_interval: 1s
        static_configs:
          - targets: ['cgroup-monitor-sc:2333']
```

### Example kubernetes yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: stress-ng
  namespace: default
spec:
  selector:
    matchLabels:
      app: stress-ng
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: stress-ng
    spec:
      terminationGracePeriodSeconds: 5
      shareProcessNamespace: true
      containers:
        - name: CONTAINER_TO_BE_MONITORED
          ...
          volumeMounts:
            - name: pid-flag-volume
              mountPath: /etc/pid_flag
        - name: cgroup-monitor-sc
          image: evankrul/cgroup-sc:prom.v.1.2
          imagePullPolicy: Always
          ports:
            - containerPort: 2333
              name: metrics
          securityContext:
            capabilities:
              add:
                - SYS_PTRACE
          env:
            - name: PORT
              value: "2333"
          resources:
            requests:
              cpu: 1
              memory: "500Mi"
            limits:
              cpu: 1
              memory: "500Mi"
          volumeMounts:
            - name: pid-flag-volume
              mountPath: /etc/pid_flag_sc
      volumes:
        - name: pid-flag-volume
          configMap:
            name: pid-flag-config-map
---
#Cgroup config map
kind: ConfigMap
apiVersion: v1
metadata:
  name: pid-flag-config-map
data:
  pid_flag: stess-ng-1
---
#Cgroup Monitor SC Service
apiVersion: v1
kind: Service
metadata:
  name: cgroup-monitor-sc #this will be the Domain name
  namespace: default
spec:
  selector:
    app: stress-ng
  ports:
    - name: stress
      port: 2335
      targetPort: 2335
    - name: metrics
      port: 2333
      targetPort: 2333
  type: LoadBalancer
```
### API
There are a few endpoints:
- `/` Homepage
- `/health` K8s health endpoint
- `/psi` Debugging PSI 
- `/metrics` Prom metrics+psi endpoint

## Data Available
The following PSI metrics are reported to Prometheus and are available for querying.
```
# HELP cgroup_monitor_sc_monitored_cpu_psi CPU PSI of monitored container
# TYPE cgroup_monitor_sc_monitored_cpu_psi gauge
cgroup_monitor_sc_monitored_cpu_psi{type="some",window="10s"} 0
cgroup_monitor_sc_monitored_cpu_psi{type="some",window="300s"} 0
cgroup_monitor_sc_monitored_cpu_psi{type="some",window="60s"} 0
cgroup_monitor_sc_monitored_cpu_psi{type="some",window="total"} 385

# HELP cgroup_monitor_sc_monitored_io_psi IO PSI of monitored container
# TYPE cgroup_monitor_sc_monitored_io_psi gauge
cgroup_monitor_sc_monitored_io_psi{type="full",window="10s"} 0
cgroup_monitor_sc_monitored_io_psi{type="full",window="300s"} 0
cgroup_monitor_sc_monitored_io_psi{type="full",window="60s"} 0
cgroup_monitor_sc_monitored_io_psi{type="full",window="total"} 330809
cgroup_monitor_sc_monitored_io_psi{type="some",window="10s"} 0
cgroup_monitor_sc_monitored_io_psi{type="some",window="300s"} 0
cgroup_monitor_sc_monitored_io_psi{type="some",window="60s"} 0
cgroup_monitor_sc_monitored_io_psi{type="some",window="total"} 330815

# HELP cgroup_monitor_sc_monitored_mem_psi Mem PSI of monitored container
# TYPE cgroup_monitor_sc_monitored_mem_psi gauge
cgroup_monitor_sc_monitored_mem_psi{type="full",window="10s"} 0
cgroup_monitor_sc_monitored_mem_psi{type="full",window="300s"} 0
cgroup_monitor_sc_monitored_mem_psi{type="full",window="60s"} 0
cgroup_monitor_sc_monitored_mem_psi{type="full",window="total"} 0
cgroup_monitor_sc_monitored_mem_psi{type="some",window="10s"} 0
cgroup_monitor_sc_monitored_mem_psi{type="some",window="300s"} 0
cgroup_monitor_sc_monitored_mem_psi{type="some",window="60s"} 0
cgroup_monitor_sc_monitored_mem_psi{type="some",window="total"} 0
```


<!-- TODO -->
## TODO
- Sidecars may not be needed, and it may be worthwhile to replace it with a dameonset for the cluster.
- The PID matching algorithm should be improved.
- The use of go channels is not entirely correct.
<!-- CONTACT -->
## Contact

Evan Krul - [Website](https://krul.ca)


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/github_username/repo_name.svg?style=for-the-badge
[contributors-url]: https://github.com/github_username/repo_name/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/github_username/repo_name.svg?style=for-the-badge
[forks-url]: https://github.com/github_username/repo_name/network/members
[stars-shield]: https://img.shields.io/github/stars/github_username/repo_name.svg?style=for-the-badge
[stars-url]: https://github.com/github_username/repo_name/stargazers
[issues-shield]: https://img.shields.io/github/issues/github_username/repo_name.svg?style=for-the-badge
[issues-url]: https://github.com/github_username/repo_name/issues
[license-shield]: https://img.shields.io/github/license/github_username/repo_name.svg?style=for-the-badge
[license-url]: https://github.com/github_username/repo_name/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/linkedin_username
[product-screenshot]: https://golang.org/lib/godoc/images/go-logo-blue.svg