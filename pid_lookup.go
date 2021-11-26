package main

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

func ListPods(client *kubernetes.Clientset, ctx context.Context) (map[string]string, map[string][]string, map[string]string) {
	pods, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	podContainers := make(map[string][]string, 0)
	podNamespaces := make(map[string]string, 0)
	containerIdToName := make(map[string]string, 0)
	for _, pod := range pods.Items {
		containerIds := make([]string, 0)
		if pod.Status.Phase == "Running" {
			for i := 0; i < len(pod.Spec.Containers); i++ {
				containerId := pod.Status.ContainerStatuses[i].ContainerID
				containerName := pod.Status.ContainerStatuses[i].Name
				containerIds = append(containerIds, containerId)
				containerIdToName[containerId] = containerName
			}
		}
		podContainers[pod.Name] = containerIds
		podNamespaces[pod.Name] = pod.Namespace
	}
	return podNamespaces, podContainers, containerIdToName
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

func findPid(pidChan chan map[string]map[string]string, clientset *kubernetes.Clientset, ctx context.Context, done chan bool) {
	podNamespaces, podContainers, containerIdToName := ListPods(clientset, ctx)
	procDir := "/root/proc"
	baseDir := `/var/lib/docker/containers`
	files, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	podPidPath := make(map[string]map[string]string, 0)
	for podName, containerIds := range podContainers {
		containerPid := make(map[string]string, 0)
		for _, containerId := range containerIds {
			for _, file := range files {
				if strings.Contains(containerId, file.Name()) {
					config, err := os.ReadFile(baseDir + "/" + file.Name() + "/config.v2.json")
					if err != nil {
						log.Fatal(err)
					}
					pid, err := getPidFromJson(string(config))
					path := fmt.Sprintf("%s/%s/root/sys/fs/cgroup", procDir, pid)
					containerPid[containerIdToName[containerId]] = path
				}
			}
		}
		podInfo := podNamespaces[podName] + "/" + podName
		podPidPath[podInfo] = containerPid
	}
	if len(podPidPath) > 0 {
		pidChan <- podPidPath
		return
	}
	defer close(done)

}
