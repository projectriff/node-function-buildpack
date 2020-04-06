buildpacks = [
  { id = "io.projectriff.node",            uri = "../../artifactory/io/projectriff/node/io.projectriff.node/latest" },
  { id = "paketo-buildpacks/node-engine",  uri = "https://github.com/cloudfoundry/node-engine-cnb/releases/download/{{ go mod download -json | jq -r 'select(.Path == "github.com/cloudfoundry/node-engine-cnb").Version' }}/node-engine-cnb-{{ go mod download -json | jq -r 'select(.Path == "github.com/cloudfoundry/node-engine-cnb").Version' | sed -e 's/^v//g' }}.tgz" },
  { id = "paketo-buildpacks/yarn-install", uri = "https://github.com/cloudfoundry/yarn-install-cnb/releases/download/{{ go mod download -json | jq -r 'select(.Path == "github.com/cloudfoundry/yarn-install-cnb").Version' }}/yarn-install-cnb-{{ go mod download -json | jq -r 'select(.Path == "github.com/cloudfoundry/yarn-install-cnb").Version' | sed -e 's/^v//g' }}.tgz" },
  { id = "paketo-buildpacks/npm",          uri = "https://github.com/cloudfoundry/npm-cnb/releases/download/{{ go mod download -json | jq -r 'select(.Path == "github.com/cloudfoundry/npm-cnb").Version' }}/npm-cnb-{{ go mod download -json | jq -r 'select(.Path == "github.com/cloudfoundry/npm-cnb").Version' | sed -e 's/^v//g' }}.tgz" },
]

[[order]]
group = [
  { id = "paketo-buildpacks/node-engine" },
  { id = "paketo-buildpacks/yarn-install" },
  { id = "io.projectriff.node" },
]

[[order]]
group = [
  { id = "paketo-buildpacks/node-engine" },
  { id = "paketo-buildpacks/npm" },
  { id = "io.projectriff.node" },
]

[[order]]
group = [
  { id = "paketo-buildpacks/node-engine" },
  { id = "io.projectriff.node" },
]

[stack]
id          = "io.buildpacks.stacks.bionic"
build-image = "cloudfoundry/build:base-cnb"
run-image   = "cloudfoundry/run:base-cnb"
