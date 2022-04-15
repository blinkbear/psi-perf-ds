package main

import (
	"io"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, "PSI Perf Monitor DaemonSet")
	if err != nil {
		return
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, `{"alive": true}`)
	if err != nil {
		return
	}
}
