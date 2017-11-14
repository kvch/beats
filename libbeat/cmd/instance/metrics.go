// +build darwin linux freebsd windows

package instance

import (
	"os"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
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
		time: time.Now(),
	}
)

func init() {
	pid := os.Getpid()
	err := lastSample.procTimes.Get(pid)
	if err != nil {
		logp.Err("Error getting process ID of the beat: %v", err, "CPU usage might be wrong.")
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

	cpuUsage, normalizedCPU, err := getProcessCPUUsage()
	if err != nil {
		logp.Err("Error retrieving CPU usage of the Beat: %v", err)
	}

	monitoring.ReportFloat(V, "usage", cpuUsage)
	monitoring.ReportFloat(V, "usage.normalized", normalizedCPU)
}

// getProcessCPUUsage return the CPU usage of the Beat
// during the period between the given samples.
// The values are between 0 and 1 in case of normalized values and
// 0 and number of cores in case of unnormalized values
func getProcessCPUUsage() (float64, float64, error) {
	pid := os.Getpid()

	sample := cpuSample{
		time: time.Now(),
	}

	if err := sample.procTimes.Get(pid); err != nil {
		return 0, 0, err
	}

	dTime := sample.time.Sub(lastSample.time)
	dMilli := dTime / time.Millisecond
	dCPU := int64(sample.procTimes.Total - lastSample.procTimes.Total)

	usage := float64(dCPU) / float64(dMilli)
	normalized := usage / float64(numCores)

	lastSample = sample
	return common.Round(usage, 4), common.Round(normalized, 4), nil
}
