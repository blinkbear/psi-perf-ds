package main

import (
	"strconv"
	"strings"
	"time"

	perf_collector "github.com/hodgesds/perf-utils"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

func updatePerf(localcache *Cache, perfCollector *PerfCollector, labels map[string][]string) {
	perfCollectorTicker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case <-perfCollectorTicker.C:
			podPid := localcache.GetAllPodPidInfo()
			_queryPerf(perfCollector, podPid, labels)
		}
	}
}

type PerfProfiler struct {
	hwprofiler    perf_collector.HardwareProfiler
	swprofiler    perf_collector.SoftwareProfiler
	cacheprofiler perf_collector.CacheProfiler
}

func startPerfCollector(perfCollector *PerfCollector, podInfo string, containerPids map[string]string, perfLabels map[string][]string) {
	for container, containerPid := range containerPids {
		pid, err := strconv.Atoi(containerPid)
		if err != nil {
			klog.Errorf("failed to convert string to int %v", err)
			continue
		}
		perfProfiler := &PerfProfiler{}
		if len(perfLabels["hw"]) != 0 {
			hwprofiler, err := perf_collector.NewHardwareProfiler(pid, 1)
			if err != nil {
				klog.Errorf("start hardware PerfCollector: %v\n", err)
			}
			if err := hwprofiler.Start(); err != nil {
				klog.Errorf("start hardware perf for %s failed: %v\n", podInfo, err)
				continue
			}
			perfProfiler.hwprofiler = hwprofiler
		}
		if len(perfLabels["sw"]) != 0 {
			swprofiler, err := perf_collector.NewSoftwareProfiler(pid, -1)
			if err != nil {
				klog.Errorf("start software PerfCollector: %v\n", err)
			}
			if err := swprofiler.Start(); err != nil {
				klog.Errorf("start software perf for %s failed: %v\n", podInfo, err)
				continue
			}
			perfProfiler.swprofiler = swprofiler
		}
		if len(perfLabels["cache"]) != 0 {
			cacheprofiler, err := perf_collector.NewCacheProfiler(pid, -1)
			if err != nil {
				klog.Errorf("start cache PerfCollector: %v\n", err)
			}
			if err := cacheprofiler.Start(); err != nil {
				klog.Errorf("start cache perf for %s failed: %v\n", podInfo, err)
				continue
			}
			perfProfiler.cacheprofiler = cacheprofiler
		}

		klog.Infof("start perf for %s success", podInfo)
		container_info := podInfo + "/" + container
		perfCollector.AddNewPerfCollector(container_info, &perfProfiler.hwprofiler, &perfProfiler.swprofiler, &perfProfiler.cacheprofiler)
	}

}

func removePerfCollector(perfCollector *PerfCollector, perfLabels map[string][]string, podInfo string, containerPids map[string]string) {
	podInfos := strings.Split(podInfo, "/")
	podNamespace := podInfos[0]
	podName := podInfos[1]
	for container, pid := range containerPids {
		container_info := podInfo + "/" + container
		hwprofiler, swprofiler, cacheprofiler := perfCollector.DeletePerfCollector(container_info)
		if hwprofiler[container_info] != nil {
			(*hwprofiler[container_info]).Stop()
		}
		if swprofiler[container_info] != nil {
			(*swprofiler[container_info]).Stop()
		}
		if cacheprofiler[container_info] != nil {
			(*cacheprofiler[container_info]).Stop()
		}
		for label, metricTypes := range perfLabels {
			for _, metricType := range metricTypes {
				_deletePerfMetricsInPrometheus(podNamespace, podName, container, pid, label, metricType)
			}
		}
	}
}

func _queryPerf(perfCollector *PerfCollector, podPid map[string]map[string]string, labels map[string][]string) {
	podHwPerfCollector := perfCollector.GetAllHwPerfCollector()
	podSwPerfCollector := perfCollector.GetAllSwPerfCollector()
	podCachePerfCollector := perfCollector.GetAllCachePerfCollector()
	for podInfo, containerPid := range podPid {
		podInfos := strings.Split(podInfo, "/")
		podNamespace := podInfos[0]
		podName := podInfos[1]
		klog.V(4).Infof("query perf for %s/%s", podNamespace, podName)
		hwPerfMetrics := make(map[string]*uint64)
		swPerfMetrics := make(map[string]*uint64)
		cachePerfMetrics := make(map[string]*uint64)
		for container, pid := range containerPid {
			klog.V(4).Infof("query perf for %s/%s/%s/%s", podNamespace, podName, container, pid)
			container_info := podInfo + "/" + container
			if len(labels["hw"]) != 0 {
				if hwprofiler, ok := podHwPerfCollector[container_info]; ok {
					hwProfle, err := (*hwprofiler).Profile()
					if err != nil {
						klog.V(4).Infof("query perf for %s failed: %v", container_info, err)
						continue
					}
					hwPerfMetrics = Struct2Map(*hwProfle)
				}
				for _, label := range labels["hw"] {
					if hwPerfMetric, ok := hwPerfMetrics[label]; ok {
						klog.V(4).Infof("query perf for %s/%s/%s/%s/%s", podNamespace, podName, container, pid, label)
						_updatePerfMetricsInPrometheus(podNamespace, podName, container, pid, label, hwPerfMetric, "hw")
					}
				}
			}
			if len(labels["sw"]) != 0 {
				if swprofiler, ok := podSwPerfCollector[container_info]; ok {
					swProfle, err := (*swprofiler).Profile()
					if err != nil {
						klog.V(4).Infof("query perf for %s failed: %v", container_info, err)
						continue
					}
					swPerfMetrics = Struct2Map(*swProfle)
				}
				for _, label := range labels["hw"] {
					if swPerfMetric, ok := swPerfMetrics[label]; ok {
						klog.V(4).Infof("query perf for %s/%s/%s/%s/%s", podNamespace, podName, container, pid, label)
						_updatePerfMetricsInPrometheus(podNamespace, podName, container, pid, label, swPerfMetric, "hw")
					}
				}
			}
			if len(labels["cache"]) != 0 {
				if cacheprofiler, ok := podCachePerfCollector[container_info]; ok {
					cacheProfle, err := (*cacheprofiler).Profile()
					if err != nil {
						klog.V(4).Infof("query perf for %s failed: %v", container_info, err)
						continue
					}
					cachePerfMetrics = Struct2Map(*cacheProfle)
				}
				for _, label := range labels["cache"] {
					if cachePerfMetrics, ok := cachePerfMetrics[label]; ok {
						klog.V(4).Infof("query perf for %s/%s/%s/%s/%s", podNamespace, podName, container, pid, label)
						_updatePerfMetricsInPrometheus(podNamespace, podName, container, pid, label, cachePerfMetrics, "cache")
					}
				}
			}

		}

	}
}

func _updatePerfMetricsInPrometheus(podNamespace, podName, container, pid, label string, metricValue *uint64, metricType string) {
	if metricType == "hw" {
		hwPerfPromGaugeMap[label].With(prometheus.Labels{"pod_namespace": podNamespace, "pod_name": podName, "container": container, "pid": pid}).Set(float64(*metricValue))
	}
	if metricType == "sw" {
		swPerfPromGaugeMap[label].With(prometheus.Labels{"pod_namespace": podNamespace, "pod_name": podName, "container": container, "pid": pid}).Set(float64(*metricValue))
	}
}

func _deletePerfMetricsInPrometheus(podNamespace, podName, container, pid, label, metricType string) {
	if metricType == "hw" {
		hwPerfPromGaugeMap[label].Delete(prometheus.Labels{"pod_namespace": podNamespace, "pod_name": podName, "container": container, "pid": pid})
	}
	if metricType == "sw" {
		swPerfPromGaugeMap[label].Delete(prometheus.Labels{"pod_namespace": podNamespace, "pod_name": podName, "container": container, "pid": pid})
	}
	if metricType == "cache" {
		cachePerfPromGaugeMap[label].Delete(prometheus.Labels{"pod_namespace": podNamespace, "pod_name": podName, "container": container, "pid": pid})
	}
}
