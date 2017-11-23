package process

// includeTopConfig is the configuration for the "top N processes
// filtering" feature
type IncludeTopConfig struct {
	Enabled  bool `config:"enabled"`
	ByCPU    int  `config:"by_cpu"`
	ByMemory int  `config:"by_memory"`
}
