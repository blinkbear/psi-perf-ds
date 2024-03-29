package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
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

func removePid(localcache *Cache, podInfo string) map[string]string {
	deletePathItem, deletePidItem := localcache.DeletePodInfo(podInfo)
	deletePsi(podInfo, deletePathItem)
	return deletePidItem
}

// func findPid(localcache *Cache, pod *v1.Pod, procBaseDir string)
func findPidInContainerd(localCache *Cache, pod *v1.Pod, procBaseDir string, containerRuntimePath string) {
	sockpath := "unix:///run/containerd/containerd.sock"
	r, err := remote.NewRemoteRuntimeService(sockpath, time.Duration(10)*time.Second, nil)
	if err != nil {
		klog.Errorf("Failed to connect to containerd: %v", err)
		return
	}
	containerIds, containerIdToName := ListPods(pod)
	podInfo := pod.Namespace + "/" + pod.Name
	containerPidPath := make(map[string]string, 0)
	containerPid := make(map[string]string, 0)
	context := context.Background()
	for _, containerId := range containerIds {
		containerPidPath[containerIdToName[containerId]] = ""
		containerPid[containerIdToName[containerId]] = ""
		status, err := r.ContainerStatus(context, containerId, true)
		if err != nil {
			klog.Errorf("Failed to get container info: %v", err)
		}
		info := status.GetInfo()["info"]
		re := regexp.MustCompile(`"pid":(\d+)`)
		match := re.FindStringSubmatch(info)
		if len(match) >= 2 {
			pid := match[1]
			path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procBaseDir, pid)
			containerPidPath[containerIdToName[containerId]] = path
			containerPid[containerIdToName[containerId]] = pid
		} else {
			fmt.Println("Pid not found in the string.")
		}
	}
	localCache.AddNewPodInfo(podInfo, containerPid, containerPidPath)
}

func findPid(localcache *Cache, pod *v1.Pod, procBaseDir string, dockerBaseDir string) {
	containerIds, containerIdToName := ListPods(pod)
	// procBaseDir := "/root/proc"
	containerBaseDir := dockerBaseDir + "/containers"
	files, err := os.ReadDir(containerBaseDir)
	if err != nil {
		klog.Errorf("Fail to read dir %s", containerBaseDir)
		return
	}
	podInfo := pod.Namespace + "/" + pod.Name
	containerPidPath := make(map[string]string, 0)
	containerPid := make(map[string]string, 0)
	for _, containerId := range containerIds {
		for _, file := range files {
			if strings.Contains(containerId, file.Name()) {
				config, err := os.ReadFile(containerBaseDir + "/" + file.Name() + "/config.v2.json")
				if err != nil {
					klog.Errorf("Failed to open file for container %v", containerId)
					continue
				}
				pid, err := getPidFromJson(string(config))
				if err != nil {
					klog.Infof("Failed to get PID from container config json file")
				}
				path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procBaseDir, pid)
				containerPidPath[containerIdToName[containerId]] = path
				containerPid[containerIdToName[containerId]] = pid
			}
		}
	}
	localcache.AddNewPodInfo(podInfo, containerPid, containerPidPath)
}
