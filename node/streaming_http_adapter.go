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
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// StreamingDependency is a key identifying the streaming HTTP adapter dependency.
const StreamingDependency = "streaming-http-adapter"

// StreamingHTTPAdapter represents the streaming HTTP adapter contribute by the buildpack.
type StreamingHTTPAdapter struct {
	layer layers.DependencyLayer
}

// Contributes makes the contribution to the launch layer.
func (s StreamingHTTPAdapter) Contribute() error {
	return s.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.Body("Expanding to %s", layer.Root)
		return helper.ExtractTarGz(artifact, filepath.Join(layer.Root, "bin"), 0)
	}, layers.Launch)
}

// NewStreamingHTTPAdapter creates a new instance.
func NewStreamingHTTPAdapter(build build.Build) (StreamingHTTPAdapter, error) {
	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return StreamingHTTPAdapter{}, err
	}

	dep, err := deps.Best(StreamingDependency, "", build.Stack)
	if err != nil {
		return StreamingHTTPAdapter{}, err
	}

	return StreamingHTTPAdapter{build.Layers.DependencyLayer(dep)}, nil
}
