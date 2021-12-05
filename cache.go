package main

import (
	"sync"

	perf_collector "github.com/hodgesds/perf-utils"
)

type Cache struct {
	podPidPath      map[string]map[string]string
	podContainerPid map[string]map[string]string
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		podPidPath:      make(map[string]map[string]string, 0),
		podContainerPid: make(map[string]map[string]string, 0),
	}
}

func (c *Cache) AddNewPodInfo(podInfo string, containerPid, containerPidPath map[string]string) {
	c.Lock()
	defer c.Unlock()
	c.podPidPath[podInfo] = containerPidPath
	c.podContainerPid[podInfo] = containerPid
}

func (c *Cache) GetPodPidPathFromPodInfo(podInfo string) map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.podPidPath[podInfo]
}

func (c *Cache) DeletePodInfo(podInfo string) (map[string]string, map[string]string) {
	c.Lock()
	defer c.Unlock()
	deletePathItem := make(map[string]string)
	deletePidItem := make(map[string]string)
	for container, path := range c.podPidPath[podInfo] {
		deletePathItem[container] = path
	}
	for container, pid := range c.podContainerPid[podInfo] {
		deletePidItem[container] = pid
	}
	delete(c.podPidPath, podInfo)
	return deletePathItem, deletePidItem
}

func (c *Cache) GetAllPodPathInfo() map[string]map[string]string {
	c.RLock()
	defer c.RUnlock()
	podPathInfoMap := make(map[string]map[string]string, 0)
	for podInfo, containerPid := range c.podPidPath {
		podPathInfoMap[podInfo] = containerPid
	}
	return podPathInfoMap
}

func (c *Cache) GetAllPodPidInfo() map[string]map[string]string {
	c.RLock()
	defer c.RUnlock()
	podPidInfoMap := make(map[string]map[string]string, 0)
	for podInfo, containerPid := range c.podContainerPid {
		podPidInfoMap[podInfo] = containerPid
	}
	return podPidInfoMap
}

func (c *Cache) GetPodPidInfoFromPodInfo(podInfo string) map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.podContainerPid[podInfo]
}

type PerfCollector struct {
	podHwPerfCollector    map[string]*perf_collector.HardwareProfiler
	podSwPerfCollector    map[string]*perf_collector.SoftwareProfiler
	podCachePerfCollector map[string]*perf_collector.CacheProfiler
	sync.RWMutex
}

func NewPerfCollector() *PerfCollector {
	return &PerfCollector{
		podHwPerfCollector:    make(map[string]*perf_collector.HardwareProfiler, 0),
		podSwPerfCollector:    make(map[string]*perf_collector.SoftwareProfiler, 0),
		podCachePerfCollector: make(map[string]*perf_collector.CacheProfiler, 0),
	}
}

func (p *PerfCollector) AddNewPerfCollector(containerInfo string, hwprofiler *perf_collector.HardwareProfiler, swprofiler *perf_collector.SoftwareProfiler, cacheprofiler *perf_collector.CacheProfiler) {
	p.Lock()
	defer p.Unlock()
	p.podHwPerfCollector[containerInfo] = hwprofiler
	p.podSwPerfCollector[containerInfo] = swprofiler
	p.podCachePerfCollector[containerInfo] = cacheprofiler
}

func (p *PerfCollector) GetHwPerfCollectorFromContainerInfo(containerInfo string) *perf_collector.HardwareProfiler {
	p.RLock()
	defer p.RUnlock()
	return p.podHwPerfCollector[containerInfo]
}

func (p *PerfCollector) GetSwPerfCollectorFromContainerInfo(containerInfo string) *perf_collector.SoftwareProfiler {
	p.RLock()
	defer p.RUnlock()
	return p.podSwPerfCollector[containerInfo]
}

func (p *PerfCollector) DeletePerfCollector(containerInfo string) (map[string]*perf_collector.HardwareProfiler, map[string]*perf_collector.SoftwareProfiler, map[string]*perf_collector.CacheProfiler) {
	p.Lock()
	defer p.Unlock()
	deleteHwItem := make(map[string]*perf_collector.HardwareProfiler)
	deleteSwItem := make(map[string]*perf_collector.SoftwareProfiler)
	deleteCacheItem := make(map[string]*perf_collector.CacheProfiler)
	deleteHwItem[containerInfo] = p.podHwPerfCollector[containerInfo]
	deleteSwItem[containerInfo] = p.podSwPerfCollector[containerInfo]
	deleteCacheItem[containerInfo] = p.podCachePerfCollector[containerInfo]
	delete(p.podHwPerfCollector, containerInfo)
	delete(p.podSwPerfCollector, containerInfo)
	delete(p.podCachePerfCollector, containerInfo)
	return deleteHwItem, deleteSwItem, deleteCacheItem
}

func (p *PerfCollector) GetAllHwPerfCollector() map[string]*perf_collector.HardwareProfiler {
	p.RLock()
	defer p.RUnlock()
	hwPerfCollectorMap := make(map[string]*perf_collector.HardwareProfiler, 0)
	for cotnainerInfo, hwprofiler := range p.podHwPerfCollector {
		hwPerfCollectorMap[cotnainerInfo] = hwprofiler
	}
	return hwPerfCollectorMap
}

func (p *PerfCollector) GetAllSwPerfCollector() map[string]*perf_collector.SoftwareProfiler {
	p.RLock()
	defer p.RUnlock()
	swPerfCollectorMap := make(map[string]*perf_collector.SoftwareProfiler, 0)
	for containerInfo, swprofiler := range p.podSwPerfCollector {
		swPerfCollectorMap[containerInfo] = swprofiler
	}
	return swPerfCollectorMap
}
func (p *PerfCollector) GetAllCachePerfCollector() map[string]*perf_collector.CacheProfiler {
	p.RLock()
	defer p.RUnlock()
	cachePerfCollectorMap := make(map[string]*perf_collector.CacheProfiler, 0)
	for containerInfo, cacheprofiler := range p.podCachePerfCollector {
		cachePerfCollectorMap[containerInfo] = cacheprofiler
	}
	return cachePerfCollectorMap
}
