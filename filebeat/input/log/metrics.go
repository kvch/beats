package log

import "github.com/elastic/beats/libbeat/monitoring"

var (
	harvesterMetrics = monitoring.Default.NewRegistry("filebeat.harvester")

	harvesterStarted   = monitoring.NewInt(harvesterMetrics, "started")
	harvesterClosed    = monitoring.NewInt(harvesterMetrics, "closed")
	harvesterRunning   = monitoring.NewInt(harvesterMetrics, "running")
	harvesterOpenFiles = monitoring.NewInt(harvesterMetrics, "open_files")
)
