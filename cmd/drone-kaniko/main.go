// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/urfave/cli/v2"
	"go.megpoid.dev/drone-kaniko/pkg/kaniko"
)

func printVersion(c *cli.Context) {
	slog.Info("Kaniko plugin",
		slog.String("version", Tag),
		slog.String("commit", Revision),
		slog.Time("date", LastCommit),
		slog.Bool("clean_build", !Modified),
	)
}

func main() {
	app := cli.NewApp()
	app.Name = "Kaniko plugin"
	app.Usage = "Kaniko plugin"
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)
	app.Version = Tag
	cli.VersionPrinter = printVersion

	// Run the application
	if err := app.Run(os.Args); err != nil {
		slog.Error("Kaniko plugin failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	printVersion(ctx)

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
