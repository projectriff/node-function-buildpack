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
	"fmt"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	_, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc := libpak.NewDependencyCache(context.Buildpack)
	dc.Logger = b.Logger

	e, ok, err := pr.Resolve("riff-node")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve riff-node plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}

	dep, err := dr.Resolve("invoker", e.Version)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	i := NewInvoker(dep, dc, result.Plan)
	i.Logger = b.Logger
	result.Layers = append(result.Layers, i)

	artifact := ""
	if s, ok := e.Metadata["artifact"].(string); ok {
		artifact = s
	}

	f, err := NewFunction(context.Application.Path, artifact)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create function\n%w", err)
	}
	f.Logger = b.Logger
	result.Layers = append(result.Layers, f)

	command := fmt.Sprintf("streaming-http-adapter node %s/server.js", filepath.Join(context.Layers.Path, i.Name()))
	result.Processes = append(result.Processes,
		libcnb.Process{Type: "node-function", Command: command},
		libcnb.Process{Type: "function", Command: command},
		libcnb.Process{Type: "web", Command: command},
	)

	return result, nil
}
