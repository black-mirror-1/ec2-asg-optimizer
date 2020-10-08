package asgoptimizer

import "io"

// Config for all the configurations
type Config struct {
	LogFile        io.Writer
	LogFlag        int
	Region         string `default:"us-west-2"`
	EnableDebug    bool
	EnabledRegions string `defalut:"us-east-1,us-west-2"`
}
