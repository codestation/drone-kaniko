package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"megpoid.xyz/go/drone-kaniko/pkg/kaniko"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime string
)

func main() {
	app := cli.NewApp()
	app.Name = "Kaniko plugin"
	app.Usage = "Kaniko plugin"
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)
	app.Version = version

	i, err := strconv.ParseInt(buildTime, 10, 64)
	if err == nil {
		tm := time.Unix(i, 0)
		buildTime = tm.Format("Mon Jan _2 15:04:05 2006")
	} else {
		buildTime = "unknown"
	}

	// Run the application
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	plugin := kaniko.New(
		settingsFromContext(ctx),
		urfave.PipelineFromContext(ctx),
		urfave.NetworkFromContext(ctx),
	)

	// Validate the settings
	if err := plugin.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	logrus.Infof("Starting Kaniko plugin version: %s, commit: %s, built at: %s", version, commit, buildTime)

	// Run the plugin
	if err := plugin.Execute(); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}
