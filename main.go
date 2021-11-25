package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

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

	pidChan := make(chan map[string]map[string]string)
	done := make(chan bool)
	defer close(pidChan)
	defer close(done)
	ctx, _ := context.WithCancel(context.Background())
	psiQueryTicker := time.NewTicker(1 * time.Second)
	// Lookup mounted file and verify contents match defined hash
	go findPid(pidChan, clientset, ctx, done)
	// Start go routine to update PSI values
	go updatePsi(pidChan, psiQueryTicker, done)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
