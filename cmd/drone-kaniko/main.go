package main

import (
	"fmt"
	"log"
	"os"

	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/urfave/cli/v2"
	"megpoid.xyz/go/drone-kaniko/pkg/kaniko"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "kaniko plugin"
	app.Usage = "kaniko plugin"
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)

	// Run the application
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
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

	// Run the plugin
	if err := plugin.Execute(); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}
