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

package node_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/projectriff/node-function-buildpack/node"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestInvoker(t *testing.T) {
	spec.Run(t, "Invoker", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns true if build plan exists", func() {
			f.AddDependency(node.Dependency, filepath.Join("testdata", "stub-invoker.tgz"))
			f.AddPlan(buildpackplan.Plan{Name: node.Dependency})

			_, ok, err := node.NewInvoker(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(ok).To(gomega.BeTrue())
		})

		it("returns false if build plan does not exist", func() {
			_, ok, err := node.NewInvoker(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(ok).To(gomega.BeFalse())
		})

		it("contributes invoker to launch", func() {
			f.AddDependency(node.Dependency, filepath.Join("testdata", "stub-invoker.tgz"))
			f.AddPlan(buildpackplan.Plan{Name: node.Dependency})

			i, _, err := node.NewInvoker(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(i.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer(node.Dependency)
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "fixture-marker")).To(gomega.BeARegularFile())

			streamingCommand := fmt.Sprintf("node %s/server.js", layer.Root)
			command := fmt.Sprintf("streaming-http-adapter %s", streamingCommand)
			g.Expect(f.Build.Layers).To(test.HaveApplicationMetadata(layers.Metadata{
				Processes: []layers.Process{
					{Type: "function", Command: command, Direct: false},
					{Type: "streaming-function", Command: streamingCommand, Direct: false},
					{Type: "web", Command: command, Direct: false},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
