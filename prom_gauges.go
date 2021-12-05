package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge declarations
var (
	// Gauge map for the psi of running pods
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
			Help:      "CPU PSI of monitored node",
		},
		[]string{"type", "window"})

	nodeMemPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_node_mem_psi",
			Help:      "Mem PSI of monitored node",
		},
		[]string{"type", "window"})

	nodeIoPsiGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup_monitor",
			Name:      "monitored_node_io_psi",
			Help:      "IO PSI of monitored node",
		},
		[]string{"type", "window"})
)

var (
	hwPerfPromGaugeMap = map[string]*prometheus.GaugeVec{
		"CPUCycles": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "cpu_cycles",
				Help:      "CPU migration of monitored container",
			},
			[]string{"namespace", "pod_name", "container_name"}),
		"Instructions": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "instruction",
				Help:      "instruction of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"CacheRefs": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "cache_refs",
				Help:      "cache refs of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"CacheMisses": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "cache_misses",
				Help:      "cache misses of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"BranchInstr": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "branch_instructions",
				Help:      "branch instructions of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"BranchMisses": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "branch_misses",
				Help:      "branch misses of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"BusCycles": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "bus_cycles",
				Help:      "bus cycles of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"StalledCyclesFrontend": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "stalled_cycles_frontend",
				Help:      "stalled cycles frontend of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"StalledCyclesBackend": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "stalled_cycles_backend",
				Help:      "stalled cycles backend of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"RefCpuCycles": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "ref_cpu_cycles",
				Help:      "ref cpu cycles of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
	}
	swPerfPromGaugeMap = map[string]*prometheus.GaugeVec{
		"CPUClock": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "cpu_clock",
				Help:      "CPU clock of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"TaskClock": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "task_clock",
				Help:      "task clock of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"PageFaults": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "page_faults",
				Help:      "page faults of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"MajorPageFaults": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "major_page_faults",
				Help:      "major page faults of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"MinorPageFaults": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "minor_page_faults",
				Help:      "minor page faults of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"ContextSwitches": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "context_switches",
				Help:      "context switches of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"CPUMigrations": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "cpu_migrations",
				Help:      "cpu migrations of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"AlignmentFaults": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "alignment_faults",
				Help:      "alignment faults of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"EmulationFaults": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "emulation_faults",
				Help:      "emulation faults of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
	}
	cachePerfPromGaugeMap = map[string]*prometheus.GaugeVec{
		"L1DataReadHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "l1_data_read_hit",
				Help:      "l1 data read hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"L1DataReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "l1_data_read_miss",
				Help:      "l1 data read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"L1DataWriteHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "l1_data_write_hit",
				Help:      "l1 data write hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"L1InstrReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "l1_instr_read_miss",
				Help:      "l1 instr read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"LastLevelReadHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "last_level_read_hit",
				Help:      "last level read hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"LastLevelReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "last_level_read_miss",
				Help:      "last level read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"LastLevelWriteHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "last_level_write_hit",
				Help:      "last level write hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"LastLevelWriteMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "last_level_write_miss",
				Help:      "last level write miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"DataTLBReadHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "data_tlb_read_hit",
				Help:      "data tlb read hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"DataTLBReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "data_tlb_read_miss",
				Help:      "data tlb read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"DataTLBWriteHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "data_tlb_write_hit",
				Help:      "data tlb write hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"DataTLBWriteMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "data_tlb_write_miss",
				Help:      "data tlb write miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"InstrTLBReadHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "instr_tlb_read_hit",
				Help:      "instr tlb read hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"InstrTLBReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "instr_tlb_read_miss",
				Help:      "instr tlb read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"BPUReadHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "bpu_read_hit",
				Help:      "bpu read hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"BPUReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "bpu_read_miss",
				Help:      "bpu read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"NodeReadHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "node_read_hit",
				Help:      "node read hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"NodeReadMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "node_read_miss",
				Help:      "node read miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"NodeWriteHit": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "node_write_hit",
				Help:      "node write hit of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
		"NodeWriteMiss": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup_monitor",
				Name:      "node_write_miss",
				Help:      "node write miss of monitored container",
			}, []string{"namespace", "pod_name", "container_name"}),
	}
)