package main

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"os"
	"strconv"
	"strings"
)

func ListPods(pod *v1.Pod) ([]string, map[string]string) {
	containerIdToName := make(map[string]string, 0)
	containerIds := make([]string, 0)
	if pod.Status.Phase == "Running" {
		for i := 0; i < len(pod.Spec.Containers); i++ {
			containerId := pod.Status.ContainerStatuses[i].ContainerID
			containerName := pod.Status.ContainerStatuses[i].Name
			containerIds = append(containerIds, containerId)
			containerIdToName[containerId] = containerName
		}
	}

	return containerIds, containerIdToName
}

func getPidFromJson(config string) (string, error) {
	var pid string
	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(config), &configMap)
	if err != nil {
		klog.Errorf("Failed to parse json")
		return "", err
	}
	pid = strconv.FormatFloat(configMap["State"].(map[string]interface{})["Pid"].(float64), 'f', -1, 64)
	return pid, nil
}

func removePid(localcache *Cache, podInfo string) {
	deleteItem := localcache.DeletePodInfo(podInfo)
	deletePsi(podInfo, deleteItem)
}

func findPid(localcache *Cache, pod *v1.Pod) {
	containerIds, containerIdToName := ListPods(pod)
	procDir := "/root/proc"
	baseDir := `/var/lib/docker/containers`
	files, err := os.ReadDir(baseDir)
	if err != nil {
		klog.Errorf("Fail to read dir %s", baseDir)
		return
	}
	podInfo := pod.Namespace + "/" + pod.Name
	containerPid := make(map[string]string, 0)
	for _, containerId := range containerIds {
		for _, file := range files {
			if strings.Contains(containerId, file.Name()) {
				config, err := os.ReadFile(baseDir + "/" + file.Name() + "/config.v2.json")
				if err != nil {
					klog.Errorf("Failed to open file for container %v", containerId)
					continue
				}
				pid, err := getPidFromJson(string(config))
				path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procDir, pid)
				containerPid[containerIdToName[containerId]] = path
			}
		}
	}
	localcache.AddNewPod(podInfo, containerPid)
	return
}
