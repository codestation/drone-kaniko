// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package kaniko

import (
	"github.com/drone-plugins/drone-plugin-lib/drone"
)

type pluginImpl struct {
	settings Settings
	pipeline drone.Pipeline
	network  drone.Network
}

// New Plugin from the given Settings, Pipeline, and Network.
func New(settings Settings, pipeline drone.Pipeline, network drone.Network) drone.Plugin {
	return &pluginImpl{
		settings: settings,
		pipeline: pipeline,
		network:  network,
	}
}
