/*
 * Copyright 2018-2020 the original author or authors.
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

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/projectriff/libfnbuildpack"
)

type Function struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger
	Path             string
}

func NewFunction(applicationPath string, artifactPath string) (Function, error) {
	return Function{
		LayerContributor: libpak.NewLayerContributor(libfnbuildpack.FormatFunction("NodeJS", artifactPath),
			map[string]interface{}{"artifact": artifactPath}),
		Path: filepath.Join(applicationPath, artifactPath),
	}, nil
}

func (f Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.LayerContributor.Logger = f.Logger

	return f.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		layer.LaunchEnvironment.Default("FUNCTION_URI", f.Path)

		return layer, nil
	}, libpak.LaunchLayer)
}

func (Function) Name() string {
	return "function"
}
