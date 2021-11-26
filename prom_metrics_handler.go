package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func updatePsi(pidChan chan map[string]map[string]string, ticker *time.Ticker, done chan bool) {
	podPidPath := <-pidChan
	for {
		select {
		case <-done:
			klog.Infof("done")
			return
		case <-ticker.C:
			_queryNodePsi()
			_queryPsi(podPidPath)
		}
	}
}

func _queryNodePsi() {
	basePsiDir := "/root/cgroup/"
	//basePsiDir := nodePidPath["basePsiDir"]
	cpuPsi, err := os.ReadFile(basePsiDir + `/cpu.pressure`)
	if err != nil {
		klog.Errorf("Failed to read cpu.pressure: %v", err)
		return
	}
	memPsi, err := os.ReadFile(basePsiDir + `/memory.pressure`)
	if err != nil {
		klog.Errorf("Failed to read mem.pressure: %v", err)
		return
	}
	ioPsi, err := os.ReadFile(basePsiDir + `/io.pressure`)
	if err != nil {
		klog.Errorf("Failed to read io.pressure: %v", err)
		return
	}

	FLOAT_BIT_SIZE := 64

	reSomeMatch, _ := regexp.Compile(`some avg10=(\d+.\d+) avg60=(\d+.\d+) avg300=(\d+.\d+) total=(\d+)`)
	reFullMatch, _ := regexp.Compile(`full avg10=(\d+.\d+) avg60=(\d+.\d+) avg300=(\d+.\d+) total=(\d+)`)

	cpuSomeMatches := reSomeMatch.FindAllStringSubmatch(string(cpuPsi), -1)
	cpuSome10, _ := strconv.ParseFloat(cpuSomeMatches[0][1], FLOAT_BIT_SIZE)
	cpuSome60, _ := strconv.ParseFloat(cpuSomeMatches[0][2], FLOAT_BIT_SIZE)
	cpuSome300, _ := strconv.ParseFloat(cpuSomeMatches[0][3], FLOAT_BIT_SIZE)
	cpuSomeTotal, _ := strconv.ParseFloat(cpuSomeMatches[0][4], FLOAT_BIT_SIZE)
	nodeCpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "10s"}).Set(cpuSome10)
	nodeCpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "60s"}).Set(cpuSome60)
	nodeCpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "300s"}).Set(cpuSome300)
	nodeCpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "total"}).Set(cpuSomeTotal)

	memSomeMatches := reSomeMatch.FindAllStringSubmatch(string(memPsi), -1)
	memSome10, _ := strconv.ParseFloat(memSomeMatches[0][1], FLOAT_BIT_SIZE)
	memSome60, _ := strconv.ParseFloat(memSomeMatches[0][2], FLOAT_BIT_SIZE)
	memSome300, _ := strconv.ParseFloat(memSomeMatches[0][3], FLOAT_BIT_SIZE)
	memSomeTotal, _ := strconv.ParseFloat(memSomeMatches[0][4], FLOAT_BIT_SIZE)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "some", "window": "10s"}).Set(memSome10)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "some", "window": "60s"}).Set(memSome60)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "some", "window": "300s"}).Set(memSome300)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "some", "window": "total"}).Set(memSomeTotal)

	memFullMatches := reFullMatch.FindAllStringSubmatch(string(memPsi), -1)
	memFull10, _ := strconv.ParseFloat(memFullMatches[0][1], FLOAT_BIT_SIZE)
	memFull60, _ := strconv.ParseFloat(memFullMatches[0][2], FLOAT_BIT_SIZE)
	memFull300, _ := strconv.ParseFloat(memFullMatches[0][3], FLOAT_BIT_SIZE)
	memFullTotal, _ := strconv.ParseFloat(memFullMatches[0][4], FLOAT_BIT_SIZE)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "full", "window": "10s"}).Set(memFull10)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "full", "window": "60s"}).Set(memFull60)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "full", "window": "300s"}).Set(memFull300)
	nodeMemPsiGauge.With(prometheus.Labels{"type": "full", "window": "total"}).Set(memFullTotal)

	ioSomeMatches := reSomeMatch.FindAllStringSubmatch(string(ioPsi), -1)
	ioSome10, _ := strconv.ParseFloat(ioSomeMatches[0][1], FLOAT_BIT_SIZE)
	ioSome60, _ := strconv.ParseFloat(ioSomeMatches[0][2], FLOAT_BIT_SIZE)
	ioSome300, _ := strconv.ParseFloat(ioSomeMatches[0][3], FLOAT_BIT_SIZE)
	ioSomeTotal, _ := strconv.ParseFloat(ioSomeMatches[0][4], FLOAT_BIT_SIZE)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "some", "window": "10s"}).Set(ioSome10)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "some", "window": "60s"}).Set(ioSome60)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "some", "window": "300s"}).Set(ioSome300)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "some", "window": "total"}).Set(ioSomeTotal)

	ioFullMatches := reFullMatch.FindAllStringSubmatch(string(ioPsi), -1)
	ioFull10, _ := strconv.ParseFloat(ioFullMatches[0][1], FLOAT_BIT_SIZE)
	ioFull60, _ := strconv.ParseFloat(ioFullMatches[0][2], FLOAT_BIT_SIZE)
	ioFull300, _ := strconv.ParseFloat(ioFullMatches[0][3], FLOAT_BIT_SIZE)
	ioFullTotal, _ := strconv.ParseFloat(ioFullMatches[0][4], FLOAT_BIT_SIZE)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "full", "window": "10s"}).Set(ioFull10)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "full", "window": "60s"}).Set(ioFull60)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "full", "window": "300s"}).Set(ioFull300)
	nodeIoPsiGauge.With(prometheus.Labels{"type": "full", "window": "total"}).Set(ioFullTotal)

}

func _queryPsi(podPidPath map[string]map[string]string) {
	for podInfo, containerPath := range podPidPath {
		//klog.Infof("query psi for pod %s", podName)
		podInfos := strings.Split(podInfo, "/")
		podNamespace := podInfos[0]
		podName := podInfos[1]
		for container, path := range containerPath {
			basePsiDir := path
			cpuPsi, err := os.ReadFile(basePsiDir + `/cpu.pressure`)
			if err != nil {
				klog.V(3).Infof("Failed to read cpu psi files %v,error is %v", basePsiDir, err)
				return
			}
			memPsi, err := os.ReadFile(basePsiDir + `/memory.pressure`)
			if err != nil {
				klog.V(3).Infof("Failed to read mem psi files %v,error is %v", basePsiDir, err)
				return
			}
			ioPsi, err := os.ReadFile(basePsiDir + `/io.pressure`)
			if err != nil {
				klog.V(3).Infof("Failed to read io psi files %v,error is %v", basePsiDir, err)
				return
			}

			FLOAT_BIT_SIZE := 64

			reSomeMatch, _ := regexp.Compile(`some avg10=(\d+.\d+) avg60=(\d+.\d+) avg300=(\d+.\d+) total=(\d+)`)
			reFullMatch, _ := regexp.Compile(`full avg10=(\d+.\d+) avg60=(\d+.\d+) avg300=(\d+.\d+) total=(\d+)`)

			cpuSomeMatches := reSomeMatch.FindAllStringSubmatch(string(cpuPsi), -1)
			cpuSome10, _ := strconv.ParseFloat(cpuSomeMatches[0][1], FLOAT_BIT_SIZE)
			cpuSome60, _ := strconv.ParseFloat(cpuSomeMatches[0][2], FLOAT_BIT_SIZE)
			cpuSome300, _ := strconv.ParseFloat(cpuSomeMatches[0][3], FLOAT_BIT_SIZE)
			cpuSomeTotal, _ := strconv.ParseFloat(cpuSomeMatches[0][4], FLOAT_BIT_SIZE)
			cpuPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "10s"}).Set(cpuSome10)
			cpuPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "60s"}).Set(cpuSome60)
			cpuPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "300s"}).Set(cpuSome300)
			cpuPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "total"}).Set(cpuSomeTotal)

			memSomeMatches := reSomeMatch.FindAllStringSubmatch(string(memPsi), -1)
			memSome10, _ := strconv.ParseFloat(memSomeMatches[0][1], FLOAT_BIT_SIZE)
			memSome60, _ := strconv.ParseFloat(memSomeMatches[0][2], FLOAT_BIT_SIZE)
			memSome300, _ := strconv.ParseFloat(memSomeMatches[0][3], FLOAT_BIT_SIZE)
			memSomeTotal, _ := strconv.ParseFloat(memSomeMatches[0][4], FLOAT_BIT_SIZE)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "10s"}).Set(memSome10)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "60s"}).Set(memSome60)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "300s"}).Set(memSome300)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "total"}).Set(memSomeTotal)

			memFullMatches := reFullMatch.FindAllStringSubmatch(string(memPsi), -1)
			memFull10, _ := strconv.ParseFloat(memFullMatches[0][1], FLOAT_BIT_SIZE)
			memFull60, _ := strconv.ParseFloat(memFullMatches[0][2], FLOAT_BIT_SIZE)
			memFull300, _ := strconv.ParseFloat(memFullMatches[0][3], FLOAT_BIT_SIZE)
			memFullTotal, _ := strconv.ParseFloat(memFullMatches[0][4], FLOAT_BIT_SIZE)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "10s"}).Set(memFull10)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "60s"}).Set(memFull60)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "300s"}).Set(memFull300)
			memPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "total"}).Set(memFullTotal)

			ioSomeMatches := reSomeMatch.FindAllStringSubmatch(string(ioPsi), -1)
			ioSome10, _ := strconv.ParseFloat(ioSomeMatches[0][1], FLOAT_BIT_SIZE)
			ioSome60, _ := strconv.ParseFloat(ioSomeMatches[0][2], FLOAT_BIT_SIZE)
			ioSome300, _ := strconv.ParseFloat(ioSomeMatches[0][3], FLOAT_BIT_SIZE)
			ioSomeTotal, _ := strconv.ParseFloat(ioSomeMatches[0][4], FLOAT_BIT_SIZE)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "10s"}).Set(ioSome10)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "60s"}).Set(ioSome60)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "300s"}).Set(ioSome300)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "some", "window": "total"}).Set(ioSomeTotal)

			ioFullMatches := reFullMatch.FindAllStringSubmatch(string(ioPsi), -1)
			ioFull10, _ := strconv.ParseFloat(ioFullMatches[0][1], FLOAT_BIT_SIZE)
			ioFull60, _ := strconv.ParseFloat(ioFullMatches[0][2], FLOAT_BIT_SIZE)
			ioFull300, _ := strconv.ParseFloat(ioFullMatches[0][3], FLOAT_BIT_SIZE)
			ioFullTotal, _ := strconv.ParseFloat(ioFullMatches[0][4], FLOAT_BIT_SIZE)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "10s"}).Set(ioFull10)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "60s"}).Set(ioFull60)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "300s"}).Set(ioFull300)
			ioPsiGauge.With(prometheus.Labels{"namespace": podNamespace, "pod_name": podName, "container_name": container, "type": "full", "window": "total"}).Set(ioFullTotal)
		}
	}

}
