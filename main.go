package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
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

	dirChan := make(chan string)
	psiQueryTicker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	defer close(dirChan)
	defer close(done)

	// Lookup mounted file and verify contents match defined hash
	go findPidDir(dirChan, done)
	// Start go routine to update PSI values
	go updatePsi(dirChan, psiQueryTicker, done)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
