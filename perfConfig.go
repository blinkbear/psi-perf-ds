package main

type PerfConfig struct {
	StartPerfCollector string `json:"startPerfCollector"`
	HwMetrics          string `json:"hwMetrics"`
	SwMetrics          string `json:"swMetrics"`
	CachedMetrics      string `json:"cacheMetrics"`
}
