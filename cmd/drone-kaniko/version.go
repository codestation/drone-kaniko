package main

import (
	"strconv"
	"time"
)

var (
	// Version indicates the application version
	Version = "dev"
	// Commit indicates the git commit of the build
	Commit = "unknown"
	// BuildTime indicates the date when the binary was built (set by -ldflags)
	BuildTime = "unknown"
)

func init() {
	if BuildTime != "unknown" {
		i, err := strconv.ParseInt(BuildTime, 10, 64)
		if err == nil {
			tm := time.Unix(i, 0)
			BuildTime = tm.Format("Mon Jan _2 15:04:05 2006")
		} else {
			BuildTime = "unknown"
		}
	}
}
