package asgoptimizer

import "io"

// Config for all the configurations
type Config struct {
	LogFile        io.Writer
	LogFlag        int
	Region         string
	EnableDebug    bool
	EnabledRegions string
}
