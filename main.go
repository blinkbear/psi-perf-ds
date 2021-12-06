package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func startLoader() (map[string][]string, bool) {
	availableHwPerfLabels := GetPerfKeys(hwPerfPromGaugeMap)
	availableSwPerfLabels := GetPerfKeys(swPerfPromGaugeMap)
	availableCachePerfLabels := GetPerfKeys(cachePerfPromGaugeMap)
	perfCollectorEnabled := os.Getenv("PERF_COLLECTOR_ENABLED")
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

	if perfCollectorEnabled == "false" {
		return nil, false
	}
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
	return perfLabels, true
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
	perfLabels, perfCollectorEnabled := startLoader()
	if !perfCollectorEnabled {
		klog.Info("Perf collector is disabled")
		return
	}
	localPerfCollector := NewPerfCollector()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newPod := obj.(*v1.Pod)
			if newPod.Spec.NodeName == nodeName {
				addFunc(newPod, localcache, localPerfCollector, perfLabels)
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
				deleteFunc(newPod, localcache, localPerfCollector, perfLabels)
				klog.Infof("Pod %s deleted", newPod.Name)
			}
			if newPod.Spec.NodeName == nodeName {
				updateFunc(newPod, localcache, localPerfCollector, perfLabels)
			}

		},
		DeleteFunc: func(obj interface{}) {
			deletePod := obj.(*v1.Pod)
			deleteFunc(deletePod, localcache, localPerfCollector, perfLabels)
		},
	})
	stopCh := make(chan struct{})
	defer close(stopCh)
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// Start go routine to update PSI values
	go updatePsi(localcache)
	go updatePerf(localcache, localPerfCollector, perfLabels)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fmt.Printf("Listening on port %s\n", port)
	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/", HomeHandler)
		r.HandleFunc("/health", HealthHandler)
		r.HandleFunc("/psi", PsiHandler)
		r.Handle("/metrics", promhttp.Handler())
		http.Handle("/", r)
		// Register metrics
		prometheus.MustRegister(cpuPsiGauge)
		prometheus.MustRegister(memPsiGauge)
		prometheus.MustRegister(ioPsiGauge)
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
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}()
	klog.Infof("Start to server")
	<-stopCh
}

func addFunc(newPod *v1.Pod, localcache *Cache, localPerfCollector *PerfCollector, perfLabels map[string][]string) {
	klog.Infof("Pod %s/%s adding", newPod.Namespace, newPod.Name)
	status := newPod.Status.Phase
	if status != v1.PodRunning {
		return
	}
	podInfo := newPod.GetNamespace() + "/" + newPod.GetName()
	findPid(localcache, newPod)
	podPidInfo := localcache.GetPodPidInfoFromPodInfo(podInfo)
	startPerfCollector(localPerfCollector, podInfo, podPidInfo, perfLabels)
	klog.Infof("Pod %s/%s added", newPod.Namespace, newPod.Name)
}

func updateFunc(newPod *v1.Pod, localcache *Cache, localPerfCollector *PerfCollector, perfLabels map[string][]string) {
	status := newPod.Status.Phase
	if status != v1.PodRunning {
		return
	}
	podInfo := newPod.GetNamespace() + "/" + newPod.GetName()
	findPid(localcache, newPod)
	podPidInfo := localcache.GetPodPidInfoFromPodInfo(podInfo)
	startPerfCollector(localPerfCollector, podInfo, podPidInfo, perfLabels)
}

func deleteFunc(deletePod *v1.Pod, localcache *Cache, localPerfCollector *PerfCollector, perfLabels map[string][]string) {
	podName := deletePod.Name
	podNamespace := deletePod.Namespace
	podInfo := podNamespace + "/" + podName
	containerPids := removePid(localcache, podInfo)
	removePerfCollector(localPerfCollector, perfLabels, podInfo, containerPids)
}
