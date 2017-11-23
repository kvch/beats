// +build darwin linux freebsd windows

package instance

import (
	"fmt"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/metrics/process"
	"github.com/elastic/beats/libbeat/metrics/system"
	"github.com/elastic/beats/libbeat/monitoring"
)

var (
	cpuMonitor       *system.CPUMonitor
	beatProcessStats *process.ProcStats
)

func init() {
	beatMetrics := monitoring.Default.NewRegistry("beat")
	monitoring.NewFunc(beatMetrics, "memstats", reportMemStats, monitoring.Report)
	monitoring.NewFunc(beatMetrics, "cpu", reportBeatCPU, monitoring.Report)
	monitoring.NewFunc(beatMetrics, "info", reportInfo, monitoring.Report)

	hostMetrics := monitoring.Default.NewRegistry("host")
	monitoring.NewFunc(hostMetrics, "load_average", reportSystemLoadAverage, monitoring.Report)
	monitoring.NewFunc(hostMetrics, "cpu", reportSystemCPUUsage, monitoring.Report)
}

func setupMetrics(name string) error {
	cpuMonitor = new(system.CPUMonitor)

	logp.Info("beat name: %v", name)
	beatProcessStats = &process.ProcStats{
		Procs:        []string{name},
		EnvWhitelist: nil,
		CpuTicks:     false,
		CacheCmdLine: true,
		IncludeTop:   process.IncludeTopConfig{},
	}
	err := beatProcessStats.InitProcStats()

	return err
}

func reportMemStats(m monitoring.Mode, V monitoring.Visitor) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	monitoring.ReportInt(V, "memory_total", int64(stats.TotalAlloc))
	if m == monitoring.Full {
		monitoring.ReportInt(V, "memory_alloc", int64(stats.Alloc))
		monitoring.ReportInt(V, "gc_next", int64(stats.NextGC))
	}
}

func reportInfo(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	delta := time.Since(startTime)
	uptime := int64(delta / time.Millisecond)
	monitoring.ReportInt(V, "uptime.ms", uptime)
}

func reportBeatCPU(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	cpuUsage, cpuUsageNorm, totalCpuUsage, err := getCPUPercentages()
	if err != nil {
		logp.Err("Error retrieving CPU percentages: %v", err)
		return
	}
	monitoring.ReportFloat(V, "usage", cpuUsage)
	monitoring.ReportFloat(V, "normalized", cpuUsageNorm)
	monitoring.ReportFloat(V, "usage.total", totalCpuUsage)
}

func getCPUPercentages() (float64, float64, float64, error) {
	state, err := beatProcessStats.GetProcStats()
	if err != nil {
		return 0.0, 0.0, 0.0, fmt.Errorf("error retrieving process stats")
	}

	if len(state) != 1 {
		return 0.0, 0.0, 0.0, fmt.Errorf("more or less than one processes: %v", len(state))
	}

	beatState := state[0]
	iCpuUsage, err := beatState.GetValue("cpu.total.pct")
	if err != nil {
		return 0.0, 0.0, 0.0, fmt.Errorf("error getting total CPU usage: %v", err)
	}
	iCpuUsageNorm, err := beatState.GetValue("cpu.total.norm.pct")
	if err != nil {
		return 0.0, 0.0, 0.0, fmt.Errorf("error getting normalized CPU percentage: %v", err)
	}

	iTotalCpuUsage, err := beatState.GetValue("cpu.total.total_pct")
	if err != nil {
		return 0.0, 0.0, 0.0, fmt.Errorf("error getting total CPU: %v", err)
	}

	cpuUsage, ok := iCpuUsage.(float64)
	if !ok {
		return 0.0, 0.0, 0.0, fmt.Errorf("error converting value of CPU usage")
	}

	cpuUsageNorm, ok := iCpuUsageNorm.(float64)
	if !ok {
		return 0.0, 0.0, 0.0, fmt.Errorf("error converting value of normalized CPU usage")
	}

	totalCpuUsage, ok := iTotalCpuUsage.(float64)
	if !ok {
		return 0.0, 0.0, 0.0, fmt.Errorf("error converting value of CPU usage")
	}

	return cpuUsage, cpuUsageNorm, totalCpuUsage, nil
}

func reportSystemLoadAverage(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	load, err := system.Load()
	if err != nil {
		logp.Err("Error retrieving load average: %v", err)
		return
	}
	avgs := load.Averages()
	monitoring.ReportFloat(V, "1m", avgs.OneMinute)
	monitoring.ReportFloat(V, "5m", avgs.FiveMinute)
	monitoring.ReportFloat(V, "15m", avgs.FifteenMinute)

	normAvgs := load.NormalizedAverages()
	monitoring.ReportFloat(V, "norm.1m", normAvgs.OneMinute)
	monitoring.ReportFloat(V, "norm.5m", normAvgs.FiveMinute)
	monitoring.ReportFloat(V, "norm.15m", normAvgs.FifteenMinute)
}

func reportSystemCPUUsage(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	sample, err := cpuMonitor.Sample()
	if err != nil {
		logp.Err("Error retrieving CPU usage of the system: %v", err)
		return
	}

	pct := sample.Percentages()
	monitoring.ReportFloat(V, "usage", pct.Total)

	normalizedPct := sample.NormalizedPercentages()
	monitoring.ReportFloat(V, "usage.normalized", normalizedPct.Total)
}
