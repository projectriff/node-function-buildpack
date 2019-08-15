// +build acceptance

/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package acceptance

import (
	"testing"

	fntesting "github.com/projectriff/libfnbuildpack/testing"
)

func TestBuilder(t *testing.T) {
	tcs := &fntesting.Testcases{
		Common: fntesting.Testcase{
			Repo:        "https://github.com/projectriff/fats",
			Refspec:     "c64950508fac705837205503ed4c8611acccc9c0", // master as of 2018-08-14
			Input:       "builder",
			ContentType: "text/plain",
			Output:      "BUILDER",
		},
		Testcases: []fntesting.Testcase{
			{
				Name:    "node",
				SubPath: "functions/uppercase/node",
			},
			{
				Name:        "npm",
				SubPath:     "functions/uppercase/npm",
				SkipRebuild: true,
			},
		},
	}

	tcs.Run(t)
}
