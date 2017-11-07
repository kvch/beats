package instance

import (
	"os"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/common"
	//"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/monitoring"
	sigar "github.com/elastic/gosigar"
)

var lastSample = struct {
	time      time.Time
	procTimes sigar.ProcTime
}{
	time.Now(),
	sigar.ProcTime{},
}
var numCores = runtime.NumCPU()

func init() {
	pid := os.Getpid()
	err := lastSample.procTimes.Get(pid)
	if err != nil {
		panic(err)
	}

	metrics := monitoring.Default.NewRegistry("beat")

	monitoring.NewFunc(metrics, "memstats", reportMemStats, monitoring.Report)
	monitoring.NewFunc(metrics, "cpu", reportCpu, monitoring.Report)
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

func reportCpu(_ monitoring.Mode, V monitoring.Visitor) {
	V.OnRegistryStart()
	defer V.OnRegistryFinished()

	cpuUsage, normalizedCpu := getCpuUsage()
	monitoring.ReportFloat(V, "usage", cpuUsage)
	monitoring.ReportFloat(V, "usage.normalized", normalizedCpu)
}

func getCpuUsage() (float64, float64) {
	pid := os.Getpid()

	sample := struct {
		time      time.Time
		procTimes sigar.ProcTime
	}{
		time.Now(),
		sigar.ProcTime{},
	}

	if err := sample.procTimes.Get(pid); err != nil {
		return 0, 0
	}

	dTime := sample.time.Sub(lastSample.time)
	dMilli := dTime / time.Millisecond
	dCpu := int64(sample.procTimes.Total - lastSample.procTimes.Total)

	usage := float64(dCpu) / float64(dMilli)
	normalized := usage / float64(numCores)

	lastSample = sample
	return common.Round(usage), common.Round(normalized)
}
