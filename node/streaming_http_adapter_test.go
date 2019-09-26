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
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/projectriff/node-function-buildpack/node"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestStreamingHTTPAdapter(t *testing.T) {
	spec.Run(t, "Invoker", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("contributes invoker to launch", func() {
			f.AddDependency(node.StreamingDependency, filepath.Join("testdata", "stub-streaming-http-adapter.tgz"))

			s, err := node.NewStreamingHTTPAdapter(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(s.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer(node.StreamingDependency)
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "bin", "fixture-marker")).To(gomega.BeARegularFile())
		})
	}, spec.Report(report.Terminal{}))
}
