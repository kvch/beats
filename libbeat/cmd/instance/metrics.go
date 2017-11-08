package instance

import (
	"os"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/monitoring"
	sigar "github.com/elastic/gosigar"
)

type cpuSample struct {
	time      time.Time
	procTimes sigar.ProcTime
}

var (
	numCores   = runtime.NumCPU()
	lastSample = cpuSample{
		time:      time.Now(),
		procTimes: sigar.ProcTime{},
	}
)

func init() {
	pid := os.Getpid()
	err := lastSample.procTimes.Get(pid)
	if err != nil {
		panic(err)
	}

	metrics := monitoring.Default.NewRegistry("beat")

	monitoring.NewFunc(metrics, "memstats", reportMemStats, monitoring.Report)
	monitoring.NewFunc(metrics, "cpu", reportCPU, monitoring.Report)
	monitoring.NewFunc(metrics, "info", reportInfo, monitoring.Report)
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

func reportCPU(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	cpuUsage, normalizedCPU := getCPUUsage()
	monitoring.ReportFloat(V, "usage", cpuUsage)
	monitoring.ReportFloat(V, "usage.normalized", normalizedCPU)
}

func getCPUUsage() (float64, float64) {
	pid := os.Getpid()

	sample := cpuSample{
		time:      time.Now(),
		procTimes: sigar.ProcTime{},
	}

	if err := sample.procTimes.Get(pid); err != nil {
		return 0, 0
	}

	dTime := sample.time.Sub(lastSample.time)
	dMilli := dTime / time.Millisecond
	dCPU := int64(sample.procTimes.Total - lastSample.procTimes.Total)

	usage := float64(dCPU) / float64(dMilli)
	normalized := usage / float64(numCores)

	lastSample = sample
	return common.Round(usage), common.Round(normalized)
}
