package main

import "github.com/prometheus/client_golang/prometheus"

// Gauge declarations
var (
	cpuPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor_sc",
			Name:      "monitored_cpu_psi",
			Help:      "CPU PSI of monitored container",
		},
		[]string{"type", "window"})

	memPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor_sc",
			Name:      "monitored_mem_psi",
			Help:      "Mem PSI of monitored container",
		},
		[]string{"type", "window"})

	ioPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor_sc",
			Name:      "monitored_io_psi",
			Help:      "IO PSI of monitored container",
		},
		[]string{"type", "window"})
)
