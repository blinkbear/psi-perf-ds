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

<h3 align="center">CgroupV2 PSI and Perf Daemonset</h3>

  <p align="center">

[comment]: <> (    CgroupV2 PSI Daemonset can be deployed on any kubernetes pod with access to cgroupv2 PSI metrics.)
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

This is a docker container that can be deployed as a daemonset on any kubernetes pod to monitor PSI metrics. 


### Built With

* [Go Lang](https://golang.org/)
* [Gorilla Mux](https://github.com/gorilla/mux)
* [Prometheus](https://prometheus.io/)


<!-- GETTING STARTED -->
## Getting Started

To deploy a daemonset follow these steps.

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
If you want to run the server locally without a container/kubernetes deployment edit `pid_lookup.go` to resolve the systems cgroup dir.

#### Regular image
1. `bash create_image.sh latest`

please change the image name in `create_image.sh`

#### Port
Set `PORT` env var to specify the metrics port.

<!-- USAGE EXAMPLES -->
## Usage
Assuming all the prerequisites have been met and image built and pushed to your docker repository follow these steps to deploy the daemonset.

In this section I will refer to the monitoring container as the **daemonset** and the container being monitored as the **host** container. 
The daemonset loads `/proc` and `/var/lib/docker` directories. It finds the container pid by searching `/var/lib/docker`. Then it accesses the `/proc` directory with pid to find the container's cgroup information.


The service is used to expose the daemonset webserver where the metrics are hosted.
If you are not using some kind of service mesh make sure your Prometheus deployment is on the same namespace as your daemonset deployment.

**perf collector** is implement based on [perf-utils v0.4.0](https://github.com/hodgesds/perf-utils/tree/v0.4.0). In this project, we exposed the performance information for each container process in one host. If you want to query the total performance of one host, please refer to [node_exporter](https://github.com/prometheus/node_exporter#disabled-by-default).

In order to collect perf, we should keep the daemonset as an administrator, and share the host pid with container. These configuration can be set in pod YAML file, as a `spec.hostPid` and `spec.containers.securityContext`. The example.yaml shows how to configure them.

Because the cost of perf collector is high, so we set an switch to enable and disable perf collector. And you can also set which perf metric should be collected. All of the supported perf metrics can be seen in project [perf-utils](https://github.com/hodgesds/perf-utils/blob/0517eb74ee7dd94e10d33088cac0df2f3342fd86/process_profile.go#L48), as follows

Available HW_PERF_LABELS 
`CPUCycles,Instructions,CacheRefs,CacheMisses,BranchInstr,BranchMisses,BusCycles,StalledCyclesFrontend,StalledCyclesBackend,RefCPUCycles,TimeEnabled,TimeRunning`

Available SW_PERF_LABELS 
`CPUClock,TaskClock,PageFaults,ContextSwitches,CPUMigrations,MinorPageFaults,MajorPageFaults,AlignmentFaults,EmulationFaults,TimeEnabled,TimeRunning`

Available CACHE_PERF_LABELS

`L1DataReadHit,L1DataReadMiss,L1DataWriteHit,L1InstrReadMiss,LastLevelReadHit,LastLevelReadMiss,LastLevelWriteHit,LastLevelWriteMiss,DataTLBReadHit,DataTLBReadMiss,DataTLBWriteHit,DataTLBWriteMiss,InstrTLBReadHit,InstrTLBReadMiss,BPUReadHit,BPUReadMiss,NodeReadHit,NodeReadMiss,NodeWriteHit,NodeWriteMiss,TimeEnabled,TimeRunning`

The fields in the yaml file are as follows:
```yaml
      - name: PERF_COLLECTOR_ENABLED
        value: "true"
      - name: HW_PERF_LABELS
        value: "Instructions"
```


Then just point Prometheus to the `/metrics` endpoint of your pod on the metrics port.

```yaml
- job_name:  'cgroup-monitor'
    kubernetes_sd_configs:
      - role: endpoints
    scheme: http
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      insecure_skip_verify: true
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - source_labels: [__address__]
        separator: ;
        regex: (.+):\d+
        target_label: __address__
        replacement: $1:2333
        action: replace
      - source_labels: [__address__]
        separator: ;
        regex: (.+):\d+
        target_label: instance
        replacement: $1
        action: replace
```

### Example kubernetes yaml
see `example/` directory
### API
There are a few endpoints:
- `/` Homepage
- `/health` K8s health endpoint
- `/psi` Debugging PSI 
- `/metrics` Prom metrics+psi endpoint

## Data Available
The following PSI metrics are reported to Prometheus and are available for querying.
```
# HELP cgroup_monitor_cpu_cycles CPU migration of monitored container
# TYPE cgroup_monitor_cpu_cycles gauge
cgroup_monitor_cpu_cycles{container_name="cgroup-monitor-sc",namespace="monitor",pid="5275",pod_name="cgroup-monitor-sc-nw78c"} 2.6810185e+07
cgroup_monitor_cpu_cycles{container_name="etcd",namespace="kube-system",pid="31708",pod_name="etcd-crack-bedbug"} 5.949181e+06
cgroup_monitor_cpu_cycles{container_name="kube-apiserver",namespace="kube-system",pid="31318",pod_name="kube-apiserver-crack-bedbug"} 0
cgroup_monitor_cpu_cycles{container_name="kube-controller-manager",namespace="kube-system",pid="31809",pod_name="kube-controller-manager-crack-bedbug"} 1.887559e+07
cgroup_monitor_cpu_cycles{container_name="kube-flannel",namespace="kube-system",pid="32265",pod_name="kube-flannel-ds-amd64-cszlv"} 3.4171708e+07
cgroup_monitor_cpu_cycles{container_name="kube-scheduler",namespace="kube-system",pid="32102",pod_name="kube-scheduler-crack-bedbug"} 0
cgroup_monitor_cpu_cycles{container_name="node-exporter",namespace="monitor",pid="22435",pod_name="prometheus-prometheus-node-exporter-jgtv7"} 2.89654836e+08
# HELP cgroup_monitor_instruction instruction of monitored container
# TYPE cgroup_monitor_instruction gauge
cgroup_monitor_instruction{container_name="cgroup-monitor-sc",namespace="monitor",pid="5275",pod_name="cgroup-monitor-sc-nw78c"} 5.0756236e+07
cgroup_monitor_instruction{container_name="etcd",namespace="kube-system",pid="31708",pod_name="etcd-crack-bedbug"} 1.2358213e+07
cgroup_monitor_instruction{container_name="kube-apiserver",namespace="kube-system",pid="31318",pod_name="kube-apiserver-crack-bedbug"} 0
cgroup_monitor_instruction{container_name="kube-controller-manager",namespace="kube-system",pid="31809",pod_name="kube-controller-manager-crack-bedbug"} 1.5420931e+07
cgroup_monitor_instruction{container_name="kube-flannel",namespace="kube-system",pid="32265",pod_name="kube-flannel-ds-amd64-cszlv"} 5.9731916e+07
cgroup_monitor_instruction{container_name="kube-scheduler",namespace="kube-system",pid="32102",pod_name="kube-scheduler-crack-bedbug"} 0
cgroup_monitor_instruction{container_name="node-exporter",namespace="monitor",pid="22435",pod_name="prometheus-prometheus-node-exporter-jgtv7"} 1.89660562e+08
# HELP cgroup_monitor_sc_monitored_cpu_psi CPU PSI of monitored container
# TYPE cgroup_monitor_sc_monitored_cpu_psi gauge
cgroup_monitor_sc_monitored_cpu_psi{container_name="carts",instance="172.169.8.219",job="cgroup-monitor",pod_name="carts-677b598f6f-lb9zn",type="some",window="10s"} 0
cgroup_monitor_sc_monitored_cpu_psi{container_name="carts",instance="172.169.8.219",job="cgroup-monitor",pod_name="carts-677b598f6f-lb9zn",type="some",window="300s"} 0
cgroup_monitor_sc_monitored_cpu_psi{container_name="carts",instance="172.169.8.219",job="cgroup-monitor",pod_name="carts-677b598f6f-lb9zn",type="some",window="60s"} 0
cgroup_monitor_sc_monitored_cpu_psi{container_name="carts",instance="172.169.8.219",job="cgroup-monitor",pod_name="carts-677b598f6f-lb9zn",type="some",window="total"} 328157795

# HELP cgroup_monitor_sc_monitored_io_psi IO PSI of monitored container
# TYPE cgroup_monitor_sc_monitored_io_psi gauge

cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="10s"} 0
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="300s"} 0
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="60s"} 0
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="total"} 69165
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="10s"} 0
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="300s"} 0
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="60s"} 0
cgroup_monitor_sc_monitored_io_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="total"} 69210

# HELP cgroup_monitor_sc_monitored_mem_psi Mem PSI of monitored container
# TYPE cgroup_monitor_sc_monitored_mem_psi gauge
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="10s"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="300s"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="60s"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="full", window="total"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="10s"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="300s"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="60s"} 0
cgroup_monitor_sc_monitored_mem_psi{container_name="carts", instance="172.169.8.219", job="cgroup-monitor", pod_name="carts-677b598f6f-lb9zn", type="some", window="total"} 0
```


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