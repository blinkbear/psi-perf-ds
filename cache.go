package main

import (
	"sync"
)

type Cache struct {
	podPidPath map[string]map[string]string
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		podPidPath: make(map[string]map[string]string, 0),
	}
}

func (c *Cache) AddNewPod(podInfo string, containerPid map[string]string) {
	c.Lock()
	defer c.Unlock()
	c.podPidPath[podInfo] = containerPid
}

func (c *Cache) GetPodPidPathFromPodInfo(podInfo string) map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.podPidPath[podInfo]
}
func getKeys1(m map[string]map[string]string) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}
func (c *Cache) DeletePodInfo(podInfo string) map[string]string {
	c.Lock()
	defer c.Unlock()
	deleteItem := make(map[string]string)
	for container, path := range c.podPidPath[podInfo] {
		deleteItem[container] = path
	}
	delete(c.podPidPath, podInfo)
	return deleteItem
}

func (c *Cache) GetAllPodInfo() map[string]map[string]string {
	c.RLock()
	defer c.RUnlock()
	podInfoMap := make(map[string]map[string]string, 0)
	for podInfo, containerPid := range c.podPidPath {
		podInfoMap[podInfo] = containerPid
	}
	return podInfoMap
}
