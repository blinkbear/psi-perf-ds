package main

import "github.com/prometheus/client_golang/prometheus"

// Gauge declarations
var (
	cpuPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_cpu_psi",
			Help:      "CPU PSI of monitored container",
		},
		[]string{"namespace", "pod_name", "container_name", "type", "window"})

	memPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_mem_psi",
			Help:      "Mem PSI of monitored container",
		},
		[]string{"namespace", "pod_name", "container_name", "type", "window"})

	ioPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_io_psi",
			Help:      "IO PSI of monitored container",
		},
		[]string{"namespace", "pod_name", "container_name", "type", "window"})

	nodeCpuPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_node_cpu_psi",
			Help:      "CPU PSI of monitored container",
		},
		[]string{"type", "window"})

	nodeMemPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_node_mem_psi",
			Help:      "Mem PSI of monitored container",
		},
		[]string{"type", "window"})

	nodeIoPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_node_io_psi",
			Help:      "IO PSI of monitored container",
		},
		[]string{"type", "window"})
)
