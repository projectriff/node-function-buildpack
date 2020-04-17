/*
 * Copyright 2018 the original author or authors.
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

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	nodeEngine "github.com/paketo-buildpacks/node-engine/node"
	"github.com/projectriff/libfnbuildpack/function"
	"github.com/projectriff/node-function-buildpack/node"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			b node.Buildpack
			f *test.DetectFactory
		)

		it.Before(func() {
			b = node.Buildpack{}
			f = test.NewDetectFactory(t)
		})

		when("id", func() {

			it("returns id", func() {
				g.Expect(b.Id()).To(gomega.Equal("node"))
			})
		})

		when("detect", func() {

			it("fails with non-js handler", func() {
				test.TouchFile(t, f.Detect.Application.Root, "test-file")

				code, err := b.Detect(f.Detect, function.Metadata{Artifact: "test-file"})
				g.Expect(code).To(gomega.Equal(function.Error_ComponentInternal))
				g.Expect(err).To(gomega.MatchError(fmt.Sprintf("artifact is not a javascript file: %s", filepath.Join(f.Detect.Application.Root, "test-file"))))
			})

			it("passes without handler", func() {
				g.Expect(b.Detect(f.Detect, function.Metadata{})).To(gomega.Equal(detect.PassStatusCode))
				g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
					Provides: []buildplan.Provided{
						{Name: node.Dependency},
					},
					Requires: []buildplan.Required{
						{
							Name: nodeEngine.Dependency,
							Metadata: map[string]interface{}{
								"build":  true,
								"launch": true,
							},
						},
						{Name: node.ModulesDependency},
						{Name: node.Dependency},
					},
				}))
			})

			it("passes with handler", func() {
				test.TouchFile(t, f.Detect.Application.Root, "test-file.js")

				g.Expect(b.Detect(f.Detect, function.Metadata{Artifact: "test-file.js"})).To(gomega.Equal(detect.PassStatusCode))
				g.Expect(f.Plans).To(test.HavePlans(
					buildplan.Plan{
						Provides: []buildplan.Provided{
							{Name: node.Dependency},
						},
						Requires: []buildplan.Required{
							{
								Name: nodeEngine.Dependency,
								Metadata: map[string]interface{}{
									"build":  true,
									"launch": true,
								},
							},
							{Name: node.ModulesDependency},
							{Name: node.Dependency},
						},
					},
					buildplan.Plan{
						Provides: []buildplan.Provided{
							{Name: node.Dependency},
						},
						Requires: []buildplan.Required{
							{
								Name: nodeEngine.Dependency,
								Metadata: map[string]interface{}{
									"build":  true,
									"launch": true,
								},
							},
							{
								Name: node.Dependency,
								Metadata: map[string]interface{}{
									node.FunctionArtifact: "test-file.js",
								},
							},
						},
					},
				))
			})

			it("passes with non-js handler and override", func() {
				test.TouchFile(t, f.Detect.Application.Root, "test-file")

				g.Expect(b.Detect(f.Detect, function.Metadata{Artifact: "test-file", Override: "node"})).To(gomega.Equal(detect.PassStatusCode))
				g.Expect(f.Plans).To(test.HavePlans(
					buildplan.Plan{
						Provides: []buildplan.Provided{
							{Name: node.Dependency},
						},
						Requires: []buildplan.Required{
							{
								Name: nodeEngine.Dependency,
								Metadata: map[string]interface{}{
									"build":  true,
									"launch": true,
								},
							},
							{Name: node.ModulesDependency},
							{Name: node.Dependency},
						},
					},
					buildplan.Plan{
						Provides: []buildplan.Provided{
							{Name: node.Dependency},
						},
						Requires: []buildplan.Required{
							{
								Name: nodeEngine.Dependency,
								Metadata: map[string]interface{}{
									"build":  true,
									"launch": true,
								},
							},
							{
								Name: node.Dependency,
								Metadata: map[string]interface{}{
									node.FunctionArtifact: "test-file",
								},
							},
						},
					},
				))
			})
		})
	}, spec.Report(report.Terminal{}))
}
