// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package manifest

import (
	"fmt"
	"log/slog"

	"github.com/estesp/manifest-tool/v2/pkg/registry"
	"github.com/estesp/manifest-tool/v2/pkg/types"
)

type Config struct {
	Username      string
	Password      string
	IgnoreMissing bool
	Insecure      bool
	PlainHTTP     bool
	ConfigDir     string
}

func Push(target string, tags []string, srcImages []types.ManifestEntry, config Config) error {
	yamlInput := types.YAMLInput{
		Image:     target,
		Tags:      tags,
		Manifests: srcImages,
	}

	manifestType := types.Docker

	digest, length, err := registry.PushManifestList(
		config.Username,
		config.Password,
		yamlInput,
		config.IgnoreMissing,
		config.Insecure,
		config.PlainHTTP,
		manifestType,
		config.ConfigDir,
	)
	if err != nil {
		return fmt.Errorf("failed to push manifest list: %w", err)
	}

	slog.Info("Manifest pushed to registry", "digest", digest, "length", length)

	return nil
}
