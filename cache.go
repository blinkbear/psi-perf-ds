package main

import (
	"sync"

	perf_collector "github.com/hodgesds/perf-utils"
)

type Cache struct {
	podPidPathes      map[string]map[string]map[string]string
	podContainerPids map[string]map[string][]string
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		podPidPathes:      make(map[string]map[string]map[string]string, 0),
		podContainerPids: make(map[string]map[string][]string, 0),
	}
}

func (c *Cache) AddNewPodInfo(podInfo string, containerPids map[string][]string, containerPidPathes map[string]map[string]string) {
	c.Lock()
	defer c.Unlock()
	c.podPidPathes[podInfo] = containerPidPathes
	c.podContainerPids[podInfo] = containerPids
}

// Not used
// func (c *Cache) GetPodPidPathFromPodInfo(podInfo string) map[string]string {
// 	c.RLock()
// 	defer c.RUnlock()
// 	return c.podPidPathes[podInfo]
// }

func (c *Cache) DeletePodInfo(podInfo string) (map[string]map[string]string, map[string][]string) {
	c.Lock()
	defer c.Unlock()
	deletePathesItem := make(map[string]map[string]string)
	deletePidsItem := make(map[string][]string)
	for container, pathes := range c.podPidPathes[podInfo] {
		deletePathesItem[container] = pathes
	}
	for container, pids := range c.podContainerPids[podInfo] {
		deletePidsItem[container] = pids
	}
	delete(c.podPidPathes, podInfo)
	return deletePathesItem, deletePidsItem
}

func (c *Cache) GetAllPodPathInfo() map[string]map[string]map[string]string {
	c.RLock()
	defer c.RUnlock()
	podPathInfoMap := make(map[string]map[string]map[string]string, 0)
	for podInfo, containerPids := range c.podPidPathes {
		podPathInfoMap[podInfo] = containerPids
	}
	return podPathInfoMap
}

func (c *Cache) GetAllPodPidInfo() map[string]map[string][]string {
	c.RLock()
	defer c.RUnlock()
	podPidInfoMap := make(map[string]map[string][]string, 0)
	for podInfo, containerPids := range c.podContainerPids {
		podPidInfoMap[podInfo] = containerPids
	}
	return podPidInfoMap
}

func (c *Cache) GetPodPidInfoFromPodInfo(podInfo string) map[string][]string {
	c.RLock()
	defer c.RUnlock()
	return c.podContainerPids[podInfo]
}

type PerfCollector struct {
	podHwPerfCollector    map[string]map[string]*perf_collector.HardwareProfiler
	podSwPerfCollector    map[string]map[string]*perf_collector.SoftwareProfiler
	podCachePerfCollector map[string]map[string]*perf_collector.CacheProfiler
	sync.RWMutex
}

func NewPerfCollector() *PerfCollector {
	return &PerfCollector{
		podHwPerfCollector:    make(map[string]map[string]*perf_collector.HardwareProfiler, 0),
		podSwPerfCollector:    make(map[string]map[string]*perf_collector.SoftwareProfiler, 0),
		podCachePerfCollector: make(map[string]map[string]*perf_collector.CacheProfiler, 0),
	}
}

func (p *PerfCollector) AddNewPerfCollector(containerInfo string, pid string, hwprofiler *perf_collector.HardwareProfiler, swprofiler *perf_collector.SoftwareProfiler, cacheprofiler *perf_collector.CacheProfiler) {
	p.Lock()
	defer p.Unlock()
	var exist bool
	_, exist = p.podHwPerfCollector[containerInfo]
	if !exist {
		p.podHwPerfCollector[containerInfo] = map[string]*perf_collector.HardwareProfiler{}
	}
	_, exist = p.podHwPerfCollector[containerInfo][pid]
	if !exist {
		p.podHwPerfCollector[containerInfo][pid] = hwprofiler
	}

	_, exist = p.podSwPerfCollector[containerInfo]
	if !exist {
		p.podSwPerfCollector[containerInfo] = map[string]*perf_collector.SoftwareProfiler{}
	}
	_, exist = p.podSwPerfCollector[containerInfo][pid]
	if !exist {
		p.podSwPerfCollector[containerInfo][pid] = swprofiler
	}

	_, exist = p.podCachePerfCollector[containerInfo]
	if !exist {
		p.podCachePerfCollector[containerInfo] = map[string]*perf_collector.CacheProfiler{}
	}
	_, exist = p.podCachePerfCollector[containerInfo][pid]
	if !exist{
		p.podCachePerfCollector[containerInfo][pid] = cacheprofiler
	}
}

// func (p *PerfCollector) GetHwPerfCollectorFromContainerInfo(containerInfo string) *perf_collector.HardwareProfiler {
// 	p.RLock()
// 	defer p.RUnlock()
// 	return p.podHwPerfCollector[containerInfo]
// }

// func (p *PerfCollector) GetSwPerfCollectorFromContainerInfo(containerInfo string) *perf_collector.SoftwareProfiler {
// 	p.RLock()
// 	defer p.RUnlock()
// 	return p.podSwPerfCollector[containerInfo]
// }

func (p *PerfCollector) DeletePerfCollector(containerInfo string) (map[string]map[string]*perf_collector.HardwareProfiler, map[string]map[string]*perf_collector.SoftwareProfiler, map[string]map[string]*perf_collector.CacheProfiler) {
	p.Lock()
	defer p.Unlock()
	deleteHwItem := make(map[string]map[string]*perf_collector.HardwareProfiler)
	deleteSwItem := make(map[string]map[string]*perf_collector.SoftwareProfiler)
	deleteCacheItem := make(map[string]map[string]*perf_collector.CacheProfiler)
	deleteHwItem[containerInfo] = p.podHwPerfCollector[containerInfo]
	deleteSwItem[containerInfo] = p.podSwPerfCollector[containerInfo]
	deleteCacheItem[containerInfo] = p.podCachePerfCollector[containerInfo]
	delete(p.podHwPerfCollector, containerInfo)
	delete(p.podSwPerfCollector, containerInfo)
	delete(p.podCachePerfCollector, containerInfo)
	return deleteHwItem, deleteSwItem, deleteCacheItem
}

func (p *PerfCollector) GetAllHwPerfCollector() map[string]map[string]*perf_collector.HardwareProfiler {
	p.RLock()
	defer p.RUnlock()
	hwPerfCollectorMap := make(map[string]map[string]*perf_collector.HardwareProfiler, 0)
	for cotnainerInfo, hwprofilerMap := range p.podHwPerfCollector {
		hwPerfCollectorMap[cotnainerInfo] = map[string]*perf_collector.HardwareProfiler{}
		for pid, hwprofiler := range hwprofilerMap {
			hwPerfCollectorMap[cotnainerInfo][pid] = hwprofiler
		}
	}
	return hwPerfCollectorMap
}

func (p *PerfCollector) GetAllSwPerfCollector() map[string]map[string]*perf_collector.SoftwareProfiler {
	p.RLock()
	defer p.RUnlock()
	swPerfCollectorMap := make(map[string]map[string]*perf_collector.SoftwareProfiler, 0)
	for containerInfo, swprofilerMap := range p.podSwPerfCollector {
		swPerfCollectorMap[containerInfo] = map[string]*perf_collector.SoftwareProfiler{}
		for pid, swprofiler := range swprofilerMap {
			swPerfCollectorMap[containerInfo][pid] = swprofiler
		}
	}
	return swPerfCollectorMap
}
func (p *PerfCollector) GetAllCachePerfCollector() map[string]map[string]*perf_collector.CacheProfiler {
	p.RLock()
	defer p.RUnlock()
	cachePerfCollectorMap := make(map[string]map[string]*perf_collector.CacheProfiler, 0)
	for containerInfo, cacheprofilerMap := range p.podCachePerfCollector {
		cachePerfCollectorMap[containerInfo] = map[string]*perf_collector.CacheProfiler{}
		for pid, cacheprofiler := range cacheprofilerMap {
			cachePerfCollectorMap[containerInfo][pid] = cacheprofiler
		}
	}
	return cachePerfCollectorMap
}
