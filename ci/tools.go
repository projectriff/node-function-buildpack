// +build tools

// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import (
	_ "github.com/buildpacks/pack"
	_ "github.com/cloudfoundry/node-engine-cnb/node"
	_ "github.com/cloudfoundry/npm-cnb/npm"
)
