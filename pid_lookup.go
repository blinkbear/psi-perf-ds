package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func removePid(localcache *Cache, podInfo string) map[string][]string {
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
	containerPidPath := make(map[string]map[string]string, 0)
	containerPid := make(map[string][]string, 0)
	context := context.Background()
	for _, containerId := range containerIds {
		containerPid[containerIdToName[containerId]] = []string{}
		containerPidPath[containerIdToName[containerId]] = map[string]string{}
		status, err := r.ContainerStatus(context, containerId, true)
		if err != nil {
			klog.Errorf("Failed to get container info: %v", err)
		}
		info := status.GetInfo()["info"]
		re := regexp.MustCompile(`"pid":(\d+)`)
		match := re.FindStringSubmatch(info)
		if len(match) >= 2 {
			pid := match[1]
			childrenPid, err := getChildrenPid(pid)
			if err != nil {
				klog.Errorf("Failed to get children PID for container %v", containerId)
			}
			for _, childPid := range childrenPid {
				path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procBaseDir, childPid)
				containerPidPath[containerIdToName[containerId]][childPid] = path
				containerPid[containerIdToName[containerId]] = append(containerPid[containerIdToName[containerId]], childPid)
			}
		} else {
			fmt.Println("Pid not found in the string.")
		}
	}
	localCache.AddNewPodInfo(podInfo, containerPid, containerPidPath)
}

func updatePids(localcache *Cache, podInfo string, procBaseDir string) {
	containerPids := localcache.GetPodPidInfoFromPodInfo(podInfo)
	newContainerPids := map[string][]string{}
	newContainerPidPathes := map[string]map[string]string{}
	for containerName, pids := range containerPids {
		newContainerPids[containerName] = []string{}
		newContainerPidPathes[containerName] = map[string]string{}
		for _, pid := range pids {
			newPids, err := getChildrenPid(pid)
			if err != nil {
				klog.Errorf("updatePids failed: %+v", err)
				return
			}
			newContainerPids[containerName] = append(newContainerPids[containerName], newPids...)
		}
		encountered := make(map[string]bool)
		result := []string{}
		for _, pid := range newContainerPids[containerName] {
			if !encountered[pid] {
				encountered[pid] = true
				result = append(result, pid)
				path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procBaseDir, pid)
				newContainerPidPathes[containerName][pid] = path
			}
		}
		klog.Infof("%s %s PIDs are : %s", podInfo, containerName, result)
		newContainerPids[containerName] = result
	}
	localcache.AddNewPodInfo(podInfo, newContainerPids, nil)
}

func findPids(localcache *Cache, pod *v1.Pod, procBaseDir string, dockerBaseDir string) {
	containerIds, containerIdToName := ListPods(pod)
	// procBaseDir := "/root/proc"
	containerBaseDir := dockerBaseDir + "/containers"
	files, err := os.ReadDir(containerBaseDir)
	if err != nil {
		klog.Errorf("Fail to read dir %s", containerBaseDir)
		return
	}
	podInfo := pod.Namespace + "/" + pod.Name
	containerPidPath := make(map[string]map[string]string, 0)
	containerPid := make(map[string][]string, 0)
	for _, containerId := range containerIds {
		containerPid[containerIdToName[containerId]] = []string{}
		containerPidPath[containerIdToName[containerId]] = map[string]string{}
		for _, file := range files {
			if strings.Contains(containerId, file.Name()) {
				klog.Infof("Investigating: %s", file.Name())
				config, err := os.ReadFile(containerBaseDir + "/" + file.Name() + "/config.v2.json")
				if err != nil {
					klog.Errorf("Failed to open file for container %v", containerId)
					continue
				}
				pid, err := getPidFromJson(string(config))
				if err != nil {
					klog.Infof("Failed to get PID from container config json file")
					continue
				}
				klog.Infof("%s %s PID is : %s", podInfo, containerIdToName[containerId], pid)
				childrenPid, err := getChildrenPid(pid)
				klog.Infof("%s %s PIDs are : %v", podInfo, containerIdToName[containerId], childrenPid)
				if err != nil {
					klog.Errorf("Failed to get children PID for container %v", containerId)
					klog.Errorf("Caused by: %+v", err)
					continue
				}
				for _, childPid := range childrenPid {
					path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procBaseDir, childPid)
					containerPidPath[containerIdToName[containerId]][childPid] = path
					containerPid[containerIdToName[containerId]] = append(containerPid[containerIdToName[containerId]], childPid)
				}
			}
		}
	}
	localcache.AddNewPodInfo(podInfo, containerPid, containerPidPath)
}

func getChildrenPid(pid string) ([]string, error) {
	if pid == "0" {
		// When deleting pythonpi driver, it will get PID == 0
		return nil, errors.New("Container is completed")
	}
	allChildPids := []string{pid}

	cmd := exec.Command("pgrep", "-P", pid)
	output, err := cmd.Output()
	if err != nil {
		// If pgrep of a pid doesn't get children
		// It will raise an error
		// Simply ignore it
		return []string{pid}, nil
	}

	childPids := strings.Fields(string(output))
	for _, childPid := range childPids {
		allChildPids = append(allChildPids, childPid)
		grandChildPids, err := getChildrenPid(childPid)
		if err != nil {
			// If pgrep of a pid doesn't get children
			// It will raise an error
			// Simply ignore it
			continue
		}

		allChildPids = append(allChildPids, grandChildPids...)
	}

	return allChildPids, nil
}
