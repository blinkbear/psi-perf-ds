package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// check whether the node support cgroup v2
func checkSupportCGV2OrNot() bool {
	CGV2 := false
	files, _ := ioutil.ReadDir("/root/cgroup")
	for _, f := range files {
		if strings.Contains(f.Name(), "pressure") {
			CGV2 = true
		}
	}
	return CGV2
}

type BasicConfig struct {
	CgroupBaseDir        string
	ContainerRuntimePath string
	ContainerRuntime     string
	ProcBaseDir          string
	DockerBaseDir        string
	PerfBaseDir          string
	PerfCollectorEnabled bool
	PsiCollectorEnabled  bool
	PerfLabels           map[string][]string
	PsiInterval          int
	PerfInterval         int
}

func startLoader() *BasicConfig {
	basicConfig := &BasicConfig{
		CgroupBaseDir:        "/root/cgroup/",
		ContainerRuntimePath: "unix://run/containerd/containerd.sock",
		ContainerRuntime:     "containerd",
		ProcBaseDir:          "/root/proc",
		DockerBaseDir:        "/root/docker",
		PerfBaseDir:          "/sys/",
		PerfCollectorEnabled: false,
		PsiCollectorEnabled:  true,
		PsiInterval:          5,
		PerfInterval:         5,
		PerfLabels:           make(map[string][]string),
	}
	perfCollectorEnabled := basicConfig.PerfCollectorEnabled
	psiCollectorEnabled := basicConfig.PsiCollectorEnabled
	availableHwPerfLabels := GetPerfKeys(hwPerfPromGaugeMap)
	availableSwPerfLabels := GetPerfKeys(swPerfPromGaugeMap)
	availableCachePerfLabels := GetPerfKeys(cachePerfPromGaugeMap)
	perfCollectorEnabledStr := os.Getenv("PERF_COLLECTOR_ENABLED")
	psiCollectorEnabledStr := os.Getenv("PSI_COLLECTOR_ENABLED")
	psiIntervalStr := os.Getenv("PSI_INTERVAL")
	psiInterval, err := strconv.Atoi(psiIntervalStr)
	if err != nil {
		klog.Infof("Failed to get PSI interval")
	} else {
		basicConfig.PsiInterval = psiInterval
	}
	perfIntervalStr := os.Getenv("PERF_INTERVAL")
	perfInterval, err := strconv.Atoi(perfIntervalStr)
	if err != nil {
		klog.Infof("Failed to get Perf interval")
	} else {
		basicConfig.PerfInterval = perfInterval
	}
	basicConfig.CgroupBaseDir = os.Getenv("CGROUP_BASE_DIR")
	basicConfig.ContainerRuntime = os.Getenv("CONTAINER_RUNTIME")
	basicConfig.ContainerRuntimePath = os.Getenv("CONTAINER_RUNTIME_PATH")
	basicConfig.ProcBaseDir = os.Getenv("PROC_BASE_DIR")
	basicConfig.DockerBaseDir = os.Getenv("DOCKER_BASE_DIR")
	basicConfig.PerfBaseDir = os.Getenv("PERF_BASE_DIR")
	hwPerfLabelString := os.Getenv("HW_PERF_LABELS")
	swPerfLabelString := os.Getenv("SW_PERF_LABELS")
	cachePerfLabelString := os.Getenv("CACHE_PERF_LABELS")
	hwPerfLabels := strings.Split(hwPerfLabelString, ",")
	swPerfLabels := strings.Split(swPerfLabelString, ",")
	cachePerfLabels := strings.Split(cachePerfLabelString, ",")
	truelyHwPerfLabels := make([]string, 0)
	truelySwPerfLabels := make([]string, 0)
	truelyCachePerfLabels := make([]string, 0)
	perfLabels := make(map[string][]string)

	if perfCollectorEnabledStr == "false" {
		klog.Info("Perf collector is disabled")
		perfCollectorEnabled = false
		// return nil, perfCollectorEnabled
	} else if perfCollectorEnabledStr == "true" {
		perfCollectorEnabled = true
	}
	if !checkSupportCGV2OrNot() {
		klog.Infof("The node doesn't support CGroup V2 and PSI!")
		psiCollectorEnabled = false
	} else if psiCollectorEnabledStr == "false" {
		klog.Info("PSI collector is disabled")
		psiCollectorEnabled = false
	} else if psiCollectorEnabledStr == "true" {
		psiCollectorEnabled = true
	}
	basicConfig.PerfCollectorEnabled = perfCollectorEnabled
	basicConfig.PsiCollectorEnabled = psiCollectorEnabled
	klog.Infof("Strat perf collector %v and psi collector %v", basicConfig.PerfCollectorEnabled, basicConfig.PsiCollectorEnabled)
	if basicConfig.PerfCollectorEnabled {
		for _, label := range hwPerfLabels {
			if !stringInSlice(label, availableHwPerfLabels) {
				klog.Warningf("HW perf label %s is not available", label)
				continue
			}
			truelyHwPerfLabels = append(truelyHwPerfLabels, label)
		}

		for _, label := range swPerfLabels {
			if !stringInSlice(label, availableSwPerfLabels) {
				klog.Warningf("SW perf label %s is not available", label)
				continue
			}
			truelySwPerfLabels = append(truelySwPerfLabels, label)
		}
		for _, label := range cachePerfLabels {
			if !stringInSlice(label, availableCachePerfLabels) {
				klog.Warningf("Cache perf label %s is not available", label)
				continue
			}
			truelyCachePerfLabels = append(truelyCachePerfLabels, label)
		}
		perfLabels["hw"] = truelyHwPerfLabels
		perfLabels["sw"] = truelySwPerfLabels
		perfLabels["cache"] = truelyCachePerfLabels
	}
	basicConfig.PerfLabels = perfLabels

	return basicConfig
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	client := kubernetes.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactory(client, time.Duration(5)*time.Second)
	podInformer := factory.Core().V1().Pods().Informer()

	localcache := NewCache()
	nodeName := os.Getenv("NODE_NAME")
	basicConfig := startLoader()
	localPerfCollector := NewPerfCollector()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newPod := obj.(*v1.Pod)
			if newPod.Spec.NodeName == nodeName {
				addFunc(newPod, localcache, localPerfCollector, basicConfig)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)
			if oldPod.ResourceVersion == newPod.ResourceVersion {
				return
			}
			if newPod.DeletionTimestamp != nil {
				klog.Infof("Pod %s deleting", newPod.Name)
				deleteFunc(newPod, localcache, localPerfCollector, basicConfig)
				klog.Infof("Pod %s deleted", newPod.Name)
			}
			if newPod.Spec.NodeName == nodeName {
				updateFunc(newPod, localcache, localPerfCollector, basicConfig)
			}

		},
		DeleteFunc: func(obj interface{}) {
			deletePod := obj.(*v1.Pod)
			deleteFunc(deletePod, localcache, localPerfCollector, basicConfig)
		},
	})
	stopCh := make(chan struct{})
	defer close(stopCh)
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// Start go routine to update PSI values
	if basicConfig.PsiCollectorEnabled {
		go updatePsi(localcache, basicConfig.CgroupBaseDir, basicConfig.PsiInterval)
	}
	if basicConfig.PerfCollectorEnabled {
		go updatePerf(localcache, localPerfCollector, basicConfig.PerfLabels, basicConfig.PerfInterval, basicConfig.ProcBaseDir)
	}
	if !basicConfig.PsiCollectorEnabled && !basicConfig.PerfCollectorEnabled {
		klog.Infof("All collectors are disabled")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fmt.Printf("Listening on port %s\n", port)
	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/", HomeHandler)
		r.HandleFunc("/health", HealthHandler)
		// r.HandleFunc("/psi", PsiHandler)
		r.Handle("/metrics", promhttp.Handler())
		http.Handle("/", r)
		// Register metrics
		if basicConfig.PsiCollectorEnabled {
			prometheus.MustRegister(cpuPsiGauge)
			prometheus.MustRegister(memPsiGauge)
			prometheus.MustRegister(ioPsiGauge)
		}
		if basicConfig.PerfCollectorEnabled {
			prometheus.MustRegister(nodeCpuPsiGauge)
			prometheus.MustRegister(nodeMemPsiGauge)
			prometheus.MustRegister(nodeIoPsiGauge)
			for label, promGauge := range hwPerfPromGaugeMap {
				klog.V(5).Infof("Registering prometheus gauge for %s", label)
				prometheus.MustRegister(promGauge)
			}
			for label, promGauge := range swPerfPromGaugeMap {
				klog.V(5).Infof("Registering prometheus gauge for %s", label)
				prometheus.MustRegister(promGauge)
			}
			for label, promGauge := range cachePerfPromGaugeMap {
				klog.V(5).Infof("Registering prometheus gauge for %s", label)
				prometheus.MustRegister(promGauge)
			}
		}
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}()
	klog.Infof("Start to server")
	<-stopCh
}

func addFunc(newPod *v1.Pod, localcache *Cache, localPerfCollector *PerfCollector, basicConfig *BasicConfig) {
	klog.Infof("Pod %s/%s adding", newPod.Namespace, newPod.Name)
	status := newPod.Status.Phase
	if status != v1.PodRunning {
		return
	}
	podInfo := newPod.GetNamespace() + "/" + newPod.GetName()
	if basicConfig.ContainerRuntime == "docker" {
		findPids(localcache, newPod, basicConfig.ProcBaseDir, basicConfig.DockerBaseDir)
	} else if basicConfig.ContainerRuntime == "containerd" {
		findPidInContainerd(localcache, newPod, basicConfig.ProcBaseDir, basicConfig.ContainerRuntimePath)
	}
	podPidInfo := localcache.GetPodPidInfoFromPodInfo(podInfo)
	klog.Infof("%s got pid infos %v", podInfo, podPidInfo)
	if basicConfig.PerfCollectorEnabled {
		startPerfCollector(localPerfCollector, podInfo, podPidInfo, basicConfig.PerfLabels)
	}
	klog.Infof("Pod %s/%s added", newPod.Namespace, newPod.Name)
}

func updateFunc(newPod *v1.Pod, localcache *Cache, localPerfCollector *PerfCollector, basicConfig *BasicConfig) {
	status := newPod.Status.Phase
	if status != v1.PodRunning {
		return
	}
	podInfo := newPod.GetNamespace() + "/" + newPod.GetName()
	if basicConfig.ContainerRuntime == "docker" {
		findPids(localcache, newPod, basicConfig.ProcBaseDir, basicConfig.DockerBaseDir)
	} else if basicConfig.ContainerRuntime == "containerd" {
		findPidInContainerd(localcache, newPod, basicConfig.ProcBaseDir, basicConfig.ContainerRuntimePath)
	}
	// findPids(localcache, newPod, basicConfig.ProcBaseDir, basicConfig.DockerBaseDir)
	podPidInfo := localcache.GetPodPidInfoFromPodInfo(podInfo)
	if basicConfig.PerfCollectorEnabled {
		startPerfCollector(localPerfCollector, podInfo, podPidInfo, basicConfig.PerfLabels)
	}
}

func deleteFunc(deletePod *v1.Pod, localcache *Cache, localPerfCollector *PerfCollector, basicConfig *BasicConfig) {
	podName := deletePod.Name
	podNamespace := deletePod.Namespace
	podInfo := podNamespace + "/" + podName
	containerPids := removePid(localcache, podInfo)
	if basicConfig.PerfCollectorEnabled {
		removePerfCollector(localPerfCollector, basicConfig.PerfLabels, podInfo, containerPids)
	}
}
