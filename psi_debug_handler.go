package main

import (
	"io"
	"net/http"
	"os"
)

// PsiHandler for debug/*
func PsiHandler(w http.ResponseWriter, r *http.Request) {
	basePsiDir := os.Getenv("PSI_DIR")
	if basePsiDir == "" {
		basePsiDir = "/proc/pressure"
	}
	cpuPsi, err := os.ReadFile(basePsiDir + `/cpu`)
	if check(&err) {
		return
	}
	memPsi, err := os.ReadFile(basePsiDir + `/memory`)
	if check(&err) {
		return
	}
	ioPsi, err := os.ReadFile(basePsiDir + `/io`)
	if check(&err) {
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, `cpu: 
`+string(cpuPsi)+`
mem: 
`+string(memPsi)+`
io: 
`+string(ioPsi))
	if err != nil {
		return
	}
}
