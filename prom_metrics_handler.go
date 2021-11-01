package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"regexp"
	"strconv"
	"time"
)

func updatePsi(dirChan chan string, ticker *time.Ticker, done chan bool) {
	fileDir := <-dirChan
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			_queryPsi(fileDir)
		}
	}
}

func _queryPsi(fileDir string) {
	basePsiDir := fileDir
	cpuPsi, err := os.ReadFile(basePsiDir + `/cpu.pressure`)
	if check(&err) {
		return
	}
	memPsi, err := os.ReadFile(basePsiDir + `/memory.pressure`)
	if check(&err) {
		return
	}
	ioPsi, err := os.ReadFile(basePsiDir + `/io.pressure`)
	if check(&err) {
		return
	}

	FLOAT_BIT_SIZE := 64

	reSomeMatch, _ := regexp.Compile(`some avg10=(\d+.\d+) avg60=(\d+.\d+) avg300=(\d+.\d+) total=(\d+)`)
	reFullMatch, _ := regexp.Compile(`full avg10=(\d+.\d+) avg60=(\d+.\d+) avg300=(\d+.\d+) total=(\d+)`)

	cpuSomeMatches := reSomeMatch.FindAllStringSubmatch(string(cpuPsi), -1)
	cpuSome10, _ := strconv.ParseFloat(cpuSomeMatches[0][1], FLOAT_BIT_SIZE)
	cpuSome60, _ := strconv.ParseFloat(cpuSomeMatches[0][2], FLOAT_BIT_SIZE)
	cpuSome300, _ := strconv.ParseFloat(cpuSomeMatches[0][3], FLOAT_BIT_SIZE)
	cpuSomeTotal, _ := strconv.ParseFloat(cpuSomeMatches[0][4], FLOAT_BIT_SIZE)
	cpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "10s"}).Set(cpuSome10)
	cpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "60s"}).Set(cpuSome60)
	cpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "300s"}).Set(cpuSome300)
	cpuPsiGauge.With(prometheus.Labels{"type": "some", "window": "total"}).Set(cpuSomeTotal)

	memSomeMatches := reSomeMatch.FindAllStringSubmatch(string(memPsi), -1)
	memSome10, _ := strconv.ParseFloat(memSomeMatches[0][1], FLOAT_BIT_SIZE)
	memSome60, _ := strconv.ParseFloat(memSomeMatches[0][2], FLOAT_BIT_SIZE)
	memSome300, _ := strconv.ParseFloat(memSomeMatches[0][3], FLOAT_BIT_SIZE)
	memSomeTotal, _ := strconv.ParseFloat(memSomeMatches[0][4], FLOAT_BIT_SIZE)
	memPsiGauge.With(prometheus.Labels{"type": "some", "window": "10s"}).Set(memSome10)
	memPsiGauge.With(prometheus.Labels{"type": "some", "window": "60s"}).Set(memSome60)
	memPsiGauge.With(prometheus.Labels{"type": "some", "window": "300s"}).Set(memSome300)
	memPsiGauge.With(prometheus.Labels{"type": "some", "window": "total"}).Set(memSomeTotal)

	memFullMatches := reFullMatch.FindAllStringSubmatch(string(memPsi), -1)
	memFull10, _ := strconv.ParseFloat(memFullMatches[0][1], FLOAT_BIT_SIZE)
	memFull60, _ := strconv.ParseFloat(memFullMatches[0][2], FLOAT_BIT_SIZE)
	memFull300, _ := strconv.ParseFloat(memFullMatches[0][3], FLOAT_BIT_SIZE)
	memFullTotal, _ := strconv.ParseFloat(memFullMatches[0][4], FLOAT_BIT_SIZE)
	memPsiGauge.With(prometheus.Labels{"type": "full", "window": "10s"}).Set(memFull10)
	memPsiGauge.With(prometheus.Labels{"type": "full", "window": "60s"}).Set(memFull60)
	memPsiGauge.With(prometheus.Labels{"type": "full", "window": "300s"}).Set(memFull300)
	memPsiGauge.With(prometheus.Labels{"type": "full", "window": "total"}).Set(memFullTotal)

	ioSomeMatches := reSomeMatch.FindAllStringSubmatch(string(ioPsi), -1)
	ioSome10, _ := strconv.ParseFloat(ioSomeMatches[0][1], FLOAT_BIT_SIZE)
	ioSome60, _ := strconv.ParseFloat(ioSomeMatches[0][2], FLOAT_BIT_SIZE)
	ioSome300, _ := strconv.ParseFloat(ioSomeMatches[0][3], FLOAT_BIT_SIZE)
	ioSomeTotal, _ := strconv.ParseFloat(ioSomeMatches[0][4], FLOAT_BIT_SIZE)
	ioPsiGauge.With(prometheus.Labels{"type": "some", "window": "10s"}).Set(ioSome10)
	ioPsiGauge.With(prometheus.Labels{"type": "some", "window": "60s"}).Set(ioSome60)
	ioPsiGauge.With(prometheus.Labels{"type": "some", "window": "300s"}).Set(ioSome300)
	ioPsiGauge.With(prometheus.Labels{"type": "some", "window": "total"}).Set(ioSomeTotal)

	ioFullMatches := reFullMatch.FindAllStringSubmatch(string(ioPsi), -1)
	ioFull10, _ := strconv.ParseFloat(ioFullMatches[0][1], FLOAT_BIT_SIZE)
	ioFull60, _ := strconv.ParseFloat(ioFullMatches[0][2], FLOAT_BIT_SIZE)
	ioFull300, _ := strconv.ParseFloat(ioFullMatches[0][3], FLOAT_BIT_SIZE)
	ioFullTotal, _ := strconv.ParseFloat(ioFullMatches[0][4], FLOAT_BIT_SIZE)
	ioPsiGauge.With(prometheus.Labels{"type": "full", "window": "10s"}).Set(ioFull10)
	ioPsiGauge.With(prometheus.Labels{"type": "full", "window": "60s"}).Set(ioFull60)
	ioPsiGauge.With(prometheus.Labels{"type": "full", "window": "300s"}).Set(ioFull300)
	ioPsiGauge.With(prometheus.Labels{"type": "full", "window": "total"}).Set(ioFullTotal)
}
