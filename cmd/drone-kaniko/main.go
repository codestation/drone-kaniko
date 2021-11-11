package main

import (
	"fmt"
	"os"

	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"megpoid.xyz/go/drone-kaniko/pkg/kaniko"
)

const versionFormatter = `Kaniko plugin version: %s, commit: %s, built at: %s`

func printVersion(c *cli.Context) {
	_, _ = fmt.Fprintf(c.App.Writer, versionFormatter, Version, Commit, BuildTime)
}

func main() {
	app := cli.NewApp()
	app.Name = "Kaniko plugin"
	app.Usage = "Kaniko plugin"
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)
	app.Version = Version
	cli.VersionPrinter = printVersion

	// Run the application
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	logrus.Infof(versionFormatter, Version, Commit, BuildTime)

	plugin := kaniko.New(
		settingsFromContext(ctx),
		urfave.PipelineFromContext(ctx),
		urfave.NetworkFromContext(ctx),
	)

	// Validate the settings
	if err := plugin.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Run the plugin
	if err := plugin.Execute(); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}
