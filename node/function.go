/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package node

import (
	"fmt"
	"path/filepath"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// functionArtifact is a key identifying the path to the function entrypoint in the build plan.
const FunctionArtifact = "fn"

// Function represents the function to be executed.
type Function struct {
	application application.Application
	functionJS  string
	layer       layers.Layer
}

// Contributes makes the contribution to the launch layer.
func (f Function) Contribute() error {
	path := filepath.Join(f.application.Root, f.functionJS)

	return f.layer.Contribute(marker{"NodeJS", path}, func(layer layers.Layer) error {
		return layer.OverrideLaunchEnv("FUNCTION_URI", path)
	}, layers.Launch)
}

// NewFunction creates a new instance returning true if the riff-invoker-node plan exists.
func NewFunction(build build.Build) (Function, bool, error) {
	p, ok, err := build.Plans.GetShallowMerged(Dependency)
	if err != nil {
		return Function{}, false, err
	}
	if !ok {
		return Function{}, false, nil
	}

	fa, ok := p.Metadata[FunctionArtifact]
	if !ok {
		fa = ""
	}

	exec, ok := fa.(string)
	if !ok {
		return Function{}, false, fmt.Errorf("node metadata of incorrect type: %v", p.Metadata[FunctionArtifact])
	}

	return Function{
		build.Application,
		exec,
		build.Layers.Layer("node-function"),
	}, true, nil
}

type marker struct {
	Type     string `toml:"type"`
	Function string `toml:"function"`
}

func (m marker) Identity() (string, string) {
	return m.Type, m.Function
}
