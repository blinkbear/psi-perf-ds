package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"os"
	"time"
)

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
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newPod := obj.(*v1.Pod)
			if newPod.Spec.NodeName == nodeName {
				addFunc(newPod, localcache)
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
				deleteFunc(newPod, localcache)
				klog.Infof("Pod %s deleted", newPod.Name)
			}
			if newPod.Spec.NodeName == nodeName {
				updateFunc(newPod, localcache)
			}

		},
		DeleteFunc: func(obj interface{}) {
			deletePod := obj.(*v1.Pod)
			deleteFunc(deletePod, localcache)
		},
	})
	stopCh := make(chan struct{})
	defer close(stopCh)
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// Start go routine to update PSI values
	go updatePsi(localcache)
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
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}()
	klog.Infof("Start to server")
	<-stopCh
}

func addFunc(newPod *v1.Pod, localcache *Cache) {
	klog.Infof("Pod %s/%s adding", newPod.Namespace, newPod.Name)
	status := newPod.Status.Phase
	if status != v1.PodRunning {
		return
	}
	findPid(localcache, newPod)
	klog.Infof("Pod %s/%s added", newPod.Namespace, newPod.Name)
}

func updateFunc(newPod *v1.Pod, localcache *Cache) {
	status := newPod.Status.Phase
	if status != v1.PodRunning {
		return
	}
	findPid(localcache, newPod)
}

func deleteFunc(deletePod *v1.Pod, localcache *Cache) {
	podName := deletePod.Name
	podNamespace := deletePod.Namespace
	podInfo := podNamespace + "/" + podName
	removePid(localcache, podInfo)
}
